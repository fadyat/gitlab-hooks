package app

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// APIConfig is the configuration for the API
type APIConfig struct {
	// Asana Personal Access Token
	AsanaAPIKey string `envconfig:"ASANA_API_KEY" required:"true"`

	// Gitlab Secret Tokens
	GitlabSecretTokens []string `envconfig:"GITLAB_SECRET_TOKENS" required:"true"`

	// Asana last commit field name in task
	LastCommitFieldName string `envconfig:"LAST_COMMIT_FIELD_NAME" required:"true" default:"Last Commit"`

	// Asana message field name in task
	MessageCommitFieldName string `envconfig:"MESSAGE_COMMIT_FIELD_NAME" required:"true" default:"Message"`
}

// LoadConfig loads the configuration from the environment
func LoadConfig() (*APIConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var cfg APIConfig
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
