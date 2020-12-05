package commands_test

import (
	"context"
	"fmt"
	"hellper/internal/app"
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

	ctx context.Context
	app *app.App

	channelID   string
	userID      string
	triggerID   string
	mockDetails bot.DialogSubmission

	mockIncident model.Incident
}

func (f *cancelCommandFixture) setup(t *testing.T) {
	var (
		loggerMock     = log.NewLoggerMock()
		clientMock     = bot.NewClientMock()
		repositoryMock = model.NewIncidentRepositoryMock()
	)
	f.ctx = context.Background()

	//LoggerMock
	loggerMock.On(
		"Info",
		f.ctx,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]log.Value"),
	).Return()

	//Repository Mock
	repositoryMock.On(
		"GetIncident",
		f.channelID,
	).Return(f.mockIncident, nil)
	repositoryMock.On(
		"CancelIncident",
		f.ctx,
		f.channelID,
		mock.AnythingOfType("string"),
	).Return(nil)

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
	clientMock.On(
		"PostMessage",
		mock.AnythingOfType("string"),            //channel
		mock.AnythingOfType("[]slack.MsgOption"), //options
	).Return("", "", nil)
	clientMock.On(
		"AddPin",
		mock.AnythingOfType("string"),        //channel
		mock.AnythingOfType("slack.ItemRef"), //item
	).Return(nil)
	clientMock.On(
		"ArchiveConversationContext",
		f.ctx,
		f.channelID,
	).Return(nil)

	f.app = &app.App{
		Logger:             loggerMock,
		Client:             clientMock,
		IncidentRepository: repositoryMock,
	}
}
func TestOpenCancelIncidentDialog(t *testing.T) {
	table := []cancelCommandFixture{
		{
			testName:     "Check error if incident is not open",
			expectError:  true,
			errorMessage: "Incident is not open for cancel. The current incident status is resolved",
			channelID:    "ABCD",
			userID:       "ABCD",
			triggerID:    "ABCD",
			mockIncident: buildCancelIncidentMock(model.StatusResolved),
		},
		{
			testName:     "If incident is open, return nil error",
			expectError:  false,
			channelID:    "XYZ",
			userID:       "XYZ",
			triggerID:    "XYZ",
			mockIncident: buildCancelIncidentMock(model.StatusOpen),
		},
	}
	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := commands.OpenCancelIncidentDialog(
				f.ctx,
				f.app,
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

func TestCancelIncidentByDialog(t *testing.T) {
	table := []cancelCommandFixture{
		{
			testName:    "Incident canceled properly",
			expectError: false,
			channelID:   "CT50JJGP5",
			userID:      "U0G9QF9C6",
			mockDetails: buildCancelIncidentDetails(),
		},
	}
	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := commands.CancelIncidentByDialog(
				f.ctx,
				f.app,
				f.mockDetails,
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

func buildCancelIncidentMock(status string) model.Incident {
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
		Team:                    "shield",
		RootCause:               "PR #00",
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

func buildCancelIncidentDetails() bot.DialogSubmission {
	return bot.DialogSubmission{
		Channel: bot.Channel{
			ID: "CT50JJGP5",
		},
		User: bot.User{
			ID: "U0G9QF9C6",
		},
		Submission: bot.Submission{
			IncidentDescription: "Incident Canceled!",
		},
	}
}
