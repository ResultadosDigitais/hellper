package commands

import (
	"context"
	"fmt"
	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type cancelCommandFixture struct {
	testName     string
	expectError  bool
	errorMessage string

	ctx            context.Context
	mockLogger     log.Logger
	mockClient     bot.Client
	mockRepository model.Repository

	triggerID       string
	mockUser        slack.User
	channelID       string
	incidentDetails bot.DialogSubmission
	mockIncident    model.Incident
}

func (f *cancelCommandFixture) setup(t *testing.T) {
	var (
		loggerMock     = log.NewLoggerMock()
		clientMock     = bot.NewClientMock()
		repositoryMock = model.NewRepositoryMock()
	)

	f.mockUser.ID = "ABC123"
	f.ctx = context.Background()

	// Logger Mock
	loggerMock.On(
		"Info",
		f.ctx,                              //ctx
		mock.AnythingOfType("string"),      //msg
		mock.AnythingOfType("[]log.Value"), //values
	).Return()
	loggerMock.On(
		"Error",
		f.ctx,                              //ctx
		mock.AnythingOfType("string"),      //msg
		mock.AnythingOfType("[]log.Value"), //values
	).Return()

	// Client Mock
	clientMock.On(
		"OpenDialog",
		f.triggerID,                         //triggerID
		mock.AnythingOfType("slack.Dialog"), //dialog
	).Return(nil)
	clientMock.On(
		"OpenDialog",
		f.channelID,                         //channelID
		mock.AnythingOfType("slack.Dialog"), //dialog
	).Return(nil)

	//Repository Mock
	repositoryMock.On(
		"CancelIncident",
		f.ctx,                                  //ctx
		mock.AnythingOfType("*model.Incident"), //inc
	).Return(nil)
	repositoryMock.On(
		"GetIncident",
		f.channelID, //channelID
	).Return(f.mockIncident, nil)

	f.mockLogger = loggerMock
	f.mockClient = clientMock

}

func TestOpenCancelIncidentDialog(t *testing.T) {
	table := []cancelCommandFixture{
		{
			testName:    "Dialog created properly",
			expectError: false,
		},
	}

	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := OpenCancelIncidentDialog(f.ctx, f.mockLogger, f.mockClient, f.mockRepository, f.channelID, f.mockUser.ID, f.triggerID)

			if f.expectError {
				if err == nil {
					t.Fatal("an error was expected, but not occurred")
				}

				assert.EqualError(t, err, f.errorMessage)
			} else {
				if err != nil {
					t.Fatal(
						"an error occurred, but was not expected\n",
						"error: ",
						err,
					)
				}
			}
		})
	}
}

func TestCancelIncidentByDialog(t *testing.T) {
	table := []cancelCommandFixture{
		{
			testName:        "Incident canceled and description not null",
			expectError:     false,
			incidentDetails: buildSubmissionCancelMock("False alarm"),
			channelID:       "ABC123",
		},
		{
			testName:        "Incident canceled without description",
			expectError:     true,
			incidentDetails: buildSubmissionCancelMock(""),
			channelID:       "ABC123",
		},
	}

	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := CancelIncidentByDialog(f.ctx, f.mockClient, f.mockLogger, f.mockRepository, f.incidentDetails)

			if f.expectError {
				if err == nil {
					t.Fatal("an error was expected, but not occurred")
				}

				assert.EqualError(t, err, f.errorMessage)
			} else {
				if err != nil {
					t.Fatal(
						"an error occurred, but was not expected\n",
						"error: ",
						err,
					)
				}
			}
		})
	}
}

func buildSubmissionCancelMock(description string) bot.DialogSubmission {
	return bot.DialogSubmission{
		Channel: bot.Channel{
			ID:   "ABC123",
			Name: "inc-miojo",
		},
		User: bot.User{
			ID:   "DEF456",
			Name: "Lucas Feijo",
		},
		Submission: bot.Submission{
			IncidentDescription: description,
		},
	}
}
