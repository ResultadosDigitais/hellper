package commands

import (
	"fmt"
	"hellper/internal/config"
	"hellper/internal/model"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
)

type reminderCommandFixture struct {
	testName           string
	expectedStatus     string
	expectedRecurrence time.Duration
	mockIncident       model.Incident
}

func TestStatusNotify(t *testing.T) {
	test := []reminderCommandFixture{
		{
			testName:       "When status is resolved",
			expectedStatus: config.Env.ReminderResolvedNotifyMsg,
			mockIncident: model.Incident{
				Status: "resolved",
			},
		},
		{
			testName:       "When status is open",
			expectedStatus: config.Env.ReminderOpenNotifyMsg,
			mockIncident: model.Incident{
				Status: "open",
			},
		},
		{
			testName:       "When there's nothing",
			expectedStatus: "",
			mockIncident: model.Incident{
				Status: "",
			},
		},
	}

	for index, f := range test {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			result := statusNotify(f.mockIncident)
			assert.Equal(t, f.expectedStatus, result, "The return should be the same")
		})
	}
}

func TestSetRecurrence(t *testing.T) {
	test := []reminderCommandFixture{
		{
			testName:           "When status is resolved",
			expectedRecurrence: time.Duration(config.Env.ReminderResolvedStatusSeconds) * time.Second,
			mockIncident: model.Incident{
				Status: "resolved",
			},
		},
		{
			testName:           "When status is open",
			expectedRecurrence: time.Duration(config.Env.ReminderOpenStatusSeconds) * time.Second,
			mockIncident: model.Incident{
				Status: "open",
			},
		},
		{
			testName:           "When there's nothing",
			expectedRecurrence: 0,
			mockIncident: model.Incident{
				Status: "",
			},
		},
	}

	for index, f := range test {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			result := setRecurrence(f.mockIncident)
			assert.Equal(t, f.expectedRecurrence, result, "The return should be the same")
		})
	}
}
