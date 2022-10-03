package gitlab

import (
	"bitbucket.org/mikehouston/asana-go"
	"errors"
	"fmt"
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/entities"
	"github.com/fadyat/hooks/api/helpers"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"net/http"
	"strings"
)

// MergeRequestAsana godoc
// @Summary     Gitlab merge request hook
// @Description Endpoint to set last commit url to custom field in asana task, passed via commit message
// @Tags        gitlab
// @Accept      json
// @Produce     json
// @Param       X-Gitlab-Token header   string                          true "Gitlab token"
// @Param       body           body     entities.GitlabMergeRequestHook true "Gitlab merge request"
// @Success     200            {object} gitlab.SuccessResponse
// @Failure     400            {object} gitlab.ErrorResponse
// @Failure     401            {object} gitlab.ErrorResponse
// @Failure     500            {object} gitlab.ErrorResponse
// @Router      /api/v1/asana/merge [post]
func MergeRequestAsana(c *gin.Context) {
	icfg, exists := c.Get("HTTPAPI")
	if !exists {
		helpers.EndWithError(c, errors.New("apiConfig not found"), http.StatusInternalServerError, &log.Logger)
		return
	}

	cfg := icfg.(*api.HTTPAPI)
	var gitlabRequest entities.GitlabMergeRequestHook
	if err := c.BindJSON(&gitlabRequest); err != nil {
		helpers.EndWithError(c, err, http.StatusBadRequest, &log.Logger)
		return
	}

	logger := log.Logger.With().Str("pr", gitlabRequest.ObjectAttributes.URL).Logger()
	if !slices.Contains(cfg.GitlabSecretTokens, c.GetHeader("X-Gitlab-Token")) {
		helpers.EndWithError(c, errors.New("invalid gitlab token"), http.StatusUnauthorized, &log.Logger)
		return
	}

	const cutset string = "\f\t\r\n "
	lastCommit := gitlabRequest.ObjectAttributes.LastCommit
	lastCommitURL := strings.Trim(lastCommit.URL, cutset)

	urls := helpers.GetAsanaURLS(lastCommit.Message)
	if len(urls) == 0 {
		logger.Info().Msg("No asana URLS found")
	}

	client := asana.NewClientWithAccessToken(cfg.AsanaAPIKey)
	for _, asanaURL := range urls {
		p := &asana.Project{ID: asanaURL.ProjectID}

		err := p.Fetch(client)
		if err != nil {
			e := err.(*asana.Error)
			logger.Info().Msg(fmt.Sprintf("Failed to fetch asana project %s, %s", asanaURL.ProjectID, e.Message))
			continue
		}

		t := &asana.Task{ID: asanaURL.TaskID}

		lastCommitField, asanaErr := helpers.GetCustomField(p, cfg.LastCommitFieldName)

		if asanaErr != nil {
			logger.Info().Msg(fmt.Sprintf("Failed to get custom field %s, %s", cfg.LastCommitFieldName, asanaErr.Message))
			comment := fmt.Sprintf("%s\n\n %s", lastCommit.URL, lastCommit.Message)
			helpers.CreateTaskCommentWithLogs(t, client, &comment, &logger)
			continue
		}

		err = t.Update(client, &asana.UpdateTaskRequest{
			CustomFields: map[string]interface{}{
				lastCommitField.ID: lastCommitURL,
			},
		})

		if err != nil {
			e := err.(*asana.Error)
			logger.Info().Msg(fmt.Sprintf("Failed to update asana task %s, %s", asanaURL.TaskID, e.Message))
			comment := fmt.Sprintf("%s\n\n %s", lastCommit.URL, lastCommit.Message)
			helpers.CreateTaskCommentWithLogs(t, client, &comment, &logger)
			continue
		}

		logger.Debug().Msg(fmt.Sprintf("Updated asana task %s", asanaURL.TaskID))
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
