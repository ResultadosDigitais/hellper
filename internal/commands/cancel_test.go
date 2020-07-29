package commands_test

import (
	"context"
	"fmt"
	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/log"
	"hellper/internal/model"
	"testing"
	"time"

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

	channelID string
	userID    string
	triggerID string

	mockIncident model.Incident
}

func (f *cancelCommandFixture) setup(t *testing.T) {
	var (
		loggerMock     = log.NewLoggerMock()
		clientMock     = bot.NewClientMock()
		repositoryMock = model.NewRepositoryMock()
	)
	f.ctx = context.Background()

	//Repository Mock
	repositoryMock.On(
		"GetIncident",
		f.channelID,
	).Return(f.mockIncident, nil)

	//Client Mock
	clientMock.On(
		"PostEphemeralContext",
		f.ctx,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]slack.MsgOption"),
	).Return("", nil)
	clientMock.On(
		"OpenDialog",
		f.triggerID,                         //triggerID
		mock.AnythingOfType("slack.Dialog"), //dialog
	).Return(nil)

	f.mockLogger = loggerMock
	f.mockClient = clientMock
	f.mockRepository = repositoryMock

}
func TestOpenCancelIncidentDiolog(t *testing.T) {
	table := []cancelCommandFixture{
		{
			testName:     "Check error if incident is not open",
			expectError:  true,
			errorMessage: "Incident is not open for cancel. The current incident status is resolved",
			channelID:    "ABCD",
			userID:       "ABCD",
			triggerID:    "ABCD",
			mockIncident: buildIncidentMock(model.StatusResolved),
		},
		{
			testName:     "If incident is open, return nil error",
			expectError:  false,
			channelID:    "XYZ",
			userID:       "XYZ",
			triggerID:    "XYZ",
			mockIncident: buildIncidentMock(model.StatusOpen),
		},
	}
	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := commands.OpenCancelIncidentDialog(
				f.ctx, f.mockLogger,
				f.mockClient,
				f.mockRepository,
				f.channelID,
				f.userID,
				f.triggerID,
			)

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

func buildIncidentMock(status string) model.Incident {
	var (
		startDate          = time.Date(2020, time.March, 19, 12, 00, 00, 00, time.UTC)
		identificationDate = time.Date(2020, time.March, 19, 14, 20, 00, 00, time.UTC)
		endDate            = time.Date(2020, time.March, 19, 22, 30, 00, 00, time.UTC)
	)

	return model.Incident{
		Id:                      0,
		Title:                   "Incident Dates Command",
		StartTimestamp:          &startDate,
		IdentificationTimestamp: &identificationDate,
		EndTimestamp:            &endDate,
		Responsibility:          "Product",
		Team:                    "shield",
		Functionality:           "hellper",
		RootCause:               "PR #00",
		CustomerImpact:          2300,
		StatusPageUrl:           "status.io",
		PostMortemUrl:           "google.com",
		Status:                  status,
		Product:                 "RDSM",
		SeverityLevel:           3,
		ChannelName:             "inc-dates-command",
		UpdatedAt:               &endDate,
		DescriptionStarted:      "An incident ocurred with the dates command",
		DescriptionCancelled:    "",
		DescriptionResolved:     "PR was reverted",
		ChannelId:               "CT50JJGP5",
	}
}
