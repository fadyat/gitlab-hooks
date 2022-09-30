package tests

import (
	"github.com/fadyat/gitlab-hooks/app/entities"
	"github.com/fadyat/gitlab-hooks/app/helpers"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type asanaRegexTestModel struct {
	message  string
	expected []entities.AsanaURL
}

var asanaRegexTests = []asanaRegexTestModel{
	{
		"test",
		[]entities.AsanaURL{},
	},
	{
		"https://app.asana.com/0/1/2",
		[]entities.AsanaURL{},
	},
	{
		"ref|https://app.asana.com/0/1/2 ref|https://app.asana.com/0/3/4",
		[]entities.AsanaURL{
			{
				Option:    "",
				ProjectId: "1",
				TaskId:    "2",
			},
			{
				Option:    "",
				ProjectId: "3",
				TaskId:    "4",
			},
		},
	},
	{
		"ref|Added feature https://app.asana.com/0/1/2",
		[]entities.AsanaURL{},
	},
	{
		"complete|ref|https://app.asana.com/0/1/2",
		[]entities.AsanaURL{
			{
				Option:    "complete",
				ProjectId: "1",
				TaskId:    "2",
			},
		},
	},
	{
		"complete|ref|https://app.asana.com/0/1/2 close|ref|https://app.asana.com/0/2/3",
		[]entities.AsanaURL{
			{
				Option:    "complete",
				ProjectId: "1",
				TaskId:    "2",
			},
			{
				Option:    "close",
				ProjectId: "2",
				TaskId:    "3",
			},
		},
	},
	{
		"completed|https://app.asana.com/0/1/2",
		[]entities.AsanaURL{},
	},
}

func TestAsanaURLRegex(t *testing.T) {
	for _, test := range asanaRegexTests {
		actual := helpers.GetAsanaURLS(test.message)

		if len(actual) != len(test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}

		for i := 0; i < len(actual); i++ {
			if !cmp.Equal(actual[i], test.expected[i]) {
				t.Errorf("Expected %v, got %v", test.expected, actual)
			}
		}
	}
}
