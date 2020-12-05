package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"hellper/internal/app"
	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"context"
	"testing"
)

type datesCommandFixture struct {
	testName     string
	expectError  bool
	errorMessage string

	ctx                    context.Context
	mockLogger             log.Logger
	mockClient             bot.Client
	mockIncidentRepository model.IncidentRepository

	mockIncident             model.Incident
	getIncidentError         error
	incidentDetails          bot.DialogSubmission
	updateIncidentDatesError error
	channelID                string
	userID                   string
	triggerID                string
	timeZone                 string
}

func (f *datesCommandFixture) setup(t *testing.T) {
	var (
		loggerMock     = log.NewLoggerMock()
		clientMock     = bot.NewClientMock()
		repositoryMock = model.NewIncidentRepositoryMock()
	)

	f.ctx = context.Background()

	loggerMock.On("Info", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()
	loggerMock.On("Error", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()

	clientMock.On("PostEphemeralContext", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("[]slack.MsgOption")).Return("", nil)
	clientMock.On("PostMessage", f.channelID, mock.AnythingOfType("[]slack.MsgOption")).Return("", "", nil)
	clientMock.On("OpenDialog", mock.AnythingOfType("string"), mock.AnythingOfType("slack.Dialog")).Return(nil)

	repositoryMock.On("GetIncident", f.channelID).Return(f.mockIncident, f.getIncidentError)
	repositoryMock.On("UpdateIncidentDates", f.ctx, mock.AnythingOfType("*model.Incident")).Return(f.updateIncidentDatesError)

	f.mockLogger = loggerMock
	f.mockClient = clientMock
	f.mockIncidentRepository = repositoryMock
}

func TestUpdateDatesDialog(t *testing.T) {
	table := []datesCommandFixture{
		{
			testName:     "Dialog is created properly",
			expectError:  false,
			channelID:    "CT50JJGP5",
			userID:       "U0G9QF9C6",
			mockIncident: buildDatesIncidentMock(),
		},
		{
			testName:         "GetIncident error",
			expectError:      true,
			errorMessage:     "Incient not found",
			channelID:        "xunda",
			userID:           "U0G9QF9C6",
			mockIncident:     model.Incident{},
			getIncidentError: errors.New("Incient not found"),
		},
	}

	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := UpdateDatesDialog(
				f.ctx,
				&app.App{
					Logger:             f.mockLogger,
					Client:             f.mockClient,
					IncidentRepository: f.mockIncidentRepository,
				},
				f.channelID, f.userID, f.triggerID,
			)

			if f.expectError {
				if err == nil {
					t.Fatal("an error was expected, but not occurred")
				}

				assert.EqualError(t, err, f.errorMessage)
			}

			if !f.expectError && err != nil {
				t.Fatal("an error occurred, but was not expected")
			}
		})
	}
}

func TestUpdateDatesByDialog(t *testing.T) {
	table := []datesCommandFixture{
		{
			testName:        "Dates updated by dialog",
			expectError:     false,
			channelID:       "CT50JJGP5",
			userID:          "U0G9QF9C6",
			incidentDetails: buildSubmissionMock("", "0"),
		},
		{
			testName:        "Dates updated by dates not in UTC",
			expectError:     false,
			channelID:       "CT50JJGP5",
			userID:          "U0G9QF9C6",
			incidentDetails: buildSubmissionMock("", "-3"),
		},
		{
			testName:        "Error timeZone blank",
			expectError:     true,
			errorMessage:    `strconv.Atoi: parsing "": invalid syntax`,
			channelID:       "CT50JJGP5",
			userID:          "U0G9QF9C6",
			incidentDetails: buildSubmissionMock("", ""),
		},
		{
			testName:        "Error initDate out of format",
			expectError:     true,
			errorMessage:    `parsing time "2020-03-09 00:00:00" as "02/01/2006 15:04:05": cannot parse "20-03-09 00:00:00" as "/"`,
			channelID:       "CT50JJGP5",
			userID:          "U0G9QF9C6",
			incidentDetails: buildSubmissionMock("initDate", "0"),
		},
		{
			testName:        "Error identificationDate out of format",
			expectError:     true,
			errorMessage:    `parsing time "2020-03-09 00:00:00" as "02/01/2006 15:04:05": cannot parse "20-03-09 00:00:00" as "/"`,
			channelID:       "CT50JJGP5",
			userID:          "U0G9QF9C6",
			incidentDetails: buildSubmissionMock("identificationDate", "0"),
		},
		{
			testName:        "Error endDate out of format",
			expectError:     true,
			errorMessage:    `parsing time "2020-03-09 00:00:00" as "02/01/2006 15:04:05": cannot parse "20-03-09 00:00:00" as "/"`,
			channelID:       "CT50JJGP5",
			userID:          "U0G9QF9C6",
			incidentDetails: buildSubmissionMock("endDate", "0"),
		},
		{
			testName:                 "Error incident not found",
			expectError:              true,
			errorMessage:             `incident not found`,
			channelID:                "xunda",
			userID:                   "U0G9QF9C6",
			updateIncidentDatesError: errors.New("incident not found"),
			incidentDetails:          buildSubmissionMock("", "0"),
		},
	}

	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := UpdateDatesByDialog(
				f.ctx,
				&app.App{
					Logger:             f.mockLogger,
					Client:             f.mockClient,
					IncidentRepository: f.mockIncidentRepository,
				},
				f.incidentDetails,
			)

			if f.expectError {
				if err == nil {
					t.Fatal("an error was expected, but not occurred")
				}

				assert.EqualError(t, err, f.errorMessage)
			}

			if !f.expectError && err != nil {
				t.Fatal("an error occurred, but was not expected\n", "Error: ", err)
			}
		})
	}
}

func buildDatesIncidentMock() model.Incident {
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
		CustomerImpact:          sql.NullInt64{Int64: 2300, Valid: true},
		StatusPageUrl:           "status.io",
		PostMortemUrl:           "google.com",
		Status:                  "closed",
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

func buildSubmissionMock(wrongFormat string, timeZone string) bot.DialogSubmission {
	var (
		initDate           = "19/03/2020 12:00:00"
		identificationDate = "19/03/2020 14:20:00"
		endDate            = "19/03/2020 22:30:00"
	)

	switch wrongFormat {
	case "initDate":
		initDate = "2020-03-09 00:00:00"
	case "identificationDate":
		identificationDate = "2020-03-09 00:00:00"
	case "endDate":
		endDate = "2020-03-09 00:00:00"
	}

	return bot.DialogSubmission{
		Channel: bot.Channel{
			ID: "CT50JJGP5",
		},
		User: bot.User{
			ID:   "U0G9QF9C6",
			Name: "Guilherme Fonseca",
		},
		Submission: bot.Submission{
			TimeZone:           timeZone,
			InitDate:           initDate,
			IdentificationDate: identificationDate,
			EndDate:            endDate,
		},
	}
}
