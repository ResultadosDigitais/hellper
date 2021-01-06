package commands_test

import (
	"context"
	"database/sql"
	"errors"
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

type closeCommandFixture struct {
	testName     string
	expectError  bool
	errorMessage string

	ctx            context.Context
	mockLogger     log.Logger
	mockClient     bot.Client
	mockRepository model.Repository

	channelID                       string
	userID                          string
	triggerID                       string
	mockDetails                     bot.DialogSubmission
	getIncidentError                error
	postMessageError                error
	closeIncidentError              error
	archiveConversationContextError error

	mockIncident model.Incident
}

func (f *closeCommandFixture) setup(t *testing.T) {
	var (
		loggerMock     = log.NewLoggerMock()
		clientMock     = bot.NewClientMock()
		repositoryMock = model.NewRepositoryMock()
	)
	f.ctx = context.Background()

	//LoggerMock
	loggerMock.On("Error", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()
	loggerMock.On("Info", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()

	//Repository Mock
	repositoryMock.On("GetIncident", f.channelID).Return(f.mockIncident, f.getIncidentError)
	repositoryMock.On("CloseIncident", mock.AnythingOfType("*model.Incident")).Return(f.closeIncidentError)

	//Client Mock
	clientMock.On("OpenDialog", f.triggerID, mock.AnythingOfType("slack.Dialog")).Return(nil)
	clientMock.On("PostMessage", mock.AnythingOfType("string"), mock.AnythingOfType("[]slack.MsgOption")).Return("", "", f.postMessageError)
	clientMock.On("AddPin", mock.AnythingOfType("string"), mock.AnythingOfType("slack.ItemRef")).Return(nil)
	clientMock.On("ArchiveConversationContext", f.ctx, f.channelID).Return(f.archiveConversationContextError)
	clientMock.On("PostEphemeralContext", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("[]slack.MsgOption")).Return("", nil)

	f.mockLogger = loggerMock
	f.mockClient = clientMock
	f.mockRepository = repositoryMock

}
func TestOpenCloseIncidentDialog(t *testing.T) {
	table := []closeCommandFixture{
		{
			testName:     "When incident status is resolved",
			expectError:  false,
			mockIncident: buildCloseIncidentMock(model.StatusResolved),
		},
		{
			testName:     "When incident status is open",
			expectError:  false,
			mockIncident: buildCloseIncidentMock(model.StatusOpen),
		},
		{
			testName:     "When incident status is closed",
			expectError:  false,
			mockIncident: buildCloseIncidentMock(model.StatusClosed),
		},
		{
			testName:     "When incident status is canceled",
			expectError:  false,
			mockIncident: buildCloseIncidentMock(model.StatusCancel),
		},
		{
			testName:     "When incident data is blank",
			expectError:  false,
			mockIncident: model.Incident{},
		},
		{
			testName:         "PostMessage Error",
			expectError:      true,
			errorMessage:     "Please, call the command `/hellper_update_dates` to receive the current dates and update each one.",
			postMessageError: errors.New("Please, call the command `/hellper_update_dates` to receive the current dates and update each one."),
		},
		{
			testName:         "GetIncident Error",
			expectError:      true,
			getIncidentError: errors.New("Ops!"),
			errorMessage:     "Ops!",
		},
	}
	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := commands.CloseIncidentDialog(
				f.ctx,
				f.mockLogger,
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
					t.Fatalf("an error occurred, but was not expected\n%#v", err)
				}
			}
		})
	}
}

func TestCloseIncidentByDialog(t *testing.T) {
	table := []closeCommandFixture{
		{
			testName:    "Incident close properly 1",
			expectError: false,
			mockDetails: bot.DialogSubmission{
				Submission: bot.Submission{
					Impact:         "1",
					SeverityLevel:  "1",
					Responsibility: "1",
				},
			},
		},
		{
			testName:    "Incident close properly 2",
			expectError: false,
			mockDetails: bot.DialogSubmission{
				Submission: bot.Submission{
					Impact:         "1",
					SeverityLevel:  "1",
					Responsibility: "0",
				},
			},
		},
		{
			testName:    "Incident close properly 3",
			expectError: false,
			mockDetails: bot.DialogSubmission{
				Submission: bot.Submission{
					Impact:         "1",
					SeverityLevel:  "1",
					Responsibility: "",
				},
			},
		},
		{
			testName:           "CloseIncident Error",
			expectError:        true,
			errorMessage:       "Ops!",
			closeIncidentError: errors.New("Ops!"),
			mockDetails: bot.DialogSubmission{
				Submission: bot.Submission{
					Impact:         "1",
					SeverityLevel:  "1",
					Responsibility: "1",
				},
			},
		},
		{
			testName:                        "ArchiveConversationContext Error",
			expectError:                     true,
			errorMessage:                    "Ops!",
			archiveConversationContextError: errors.New("Ops!"),
			mockDetails: bot.DialogSubmission{
				Submission: bot.Submission{
					Impact:         "1",
					SeverityLevel:  "1",
					Responsibility: "1",
				},
			},
		},
		{
			testName:         "GetIncident Error",
			expectError:      true,
			errorMessage:     "Ops!",
			getIncidentError: errors.New("Ops!"),
			mockDetails: bot.DialogSubmission{
				Submission: bot.Submission{
					Impact:         "1",
					SeverityLevel:  "1",
					Responsibility: "1",
				},
			},
		},
		{
			testName:     "When Impact is invalid data",
			expectError:  true,
			errorMessage: "strconv.ParseInt: parsing \"A\": invalid syntax",
			mockDetails: bot.DialogSubmission{
				Submission: bot.Submission{
					Impact:         "A",
					SeverityLevel:  "1",
					Responsibility: "1",
				},
			},
		},
		{
			testName:     "When SeverityLevel is invalid data",
			expectError:  true,
			errorMessage: "strconv.ParseInt: parsing \"B\": invalid syntax",
			mockDetails: bot.DialogSubmission{
				Submission: bot.Submission{
					Impact:        "1",
					SeverityLevel: "B",
				},
			},
		},
	}
	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := commands.CloseIncidentByDialog(
				f.ctx,
				f.mockClient,
				f.mockLogger,
				f.mockRepository,
				f.mockDetails,
			)

			if f.expectError {
				if err == nil {
					t.Fatal("an error was expected, but not occurred")
				}

				assert.EqualError(t, err, f.errorMessage)
			} else {
				if err != nil {
					t.Fatalf("an error occurred, but was not expected: %#v\n", err)
				}
			}
		})
	}
}

func buildCloseIncidentMock(status string) model.Incident {
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
		Responsibility:          "0",
		Team:                    "shield",
		Functionality:           "hellper",
		RootCause:               "PR #00",
		CustomerImpact:          sql.NullInt64{Int64: 2300},
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
