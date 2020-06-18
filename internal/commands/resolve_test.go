package commands_test

import (
	"context"
	"fmt"
	"hellper/internal/bot"
	calendar "hellper/internal/calendar"
	"hellper/internal/commands"
	"hellper/internal/log"
	"hellper/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type resolveCommandFixture struct {
	testName     string
	expectError  bool
	errorMessage string

	ctx            context.Context
	mockLogger     log.Logger
	mockClient     bot.Client
	mockRepository model.Repository
	mockCalendar   calendar.Calendar

	triggerID       string
	channelID       string
	incidentDetails bot.DialogSubmission
	mockIncident    model.Incident
	mockEvent       *model.Event
}

func (f *resolveCommandFixture) setup(t *testing.T) {
	var (
		loggerMock     = log.NewLoggerMock()
		clientMock     = bot.NewClientMock()
		repositoryMock = model.NewRepositoryMock()
		calendarMock   = calendar.NewCalendarMock()
	)

	f.ctx = context.Background()

	//Logger Mock
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

	//Client Mock
	clientMock.On(
		"OpenDialog",
		f.triggerID,                         //triggerID
		mock.AnythingOfType("slack.Dialog"), //dialog
	).Return(nil)
	clientMock.On(
		"AddPin",
		mock.AnythingOfType("string"),        //channel
		mock.AnythingOfType("slack.ItemRef"), //item
	).Return(nil)
	clientMock.On(
		"PostMessage",
		mock.AnythingOfType("string"),            //channel
		mock.AnythingOfType("[]slack.MsgOption"), //options
	).Return("", "", nil)

	//Repository Mock
	repositoryMock.On(
		"ResolveIncident",
		f.ctx,                                  //ctx
		mock.AnythingOfType("*model.Incident"), //inc
	).Return(nil)
	repositoryMock.On(
		"GetIncident",
		f.channelID, //channelID
	).Return(f.mockIncident, nil)

	//Calendar Mock
	calendarMock.On(
		"CreateCalendarEvent",
		f.ctx,
		mock.AnythingOfType("string"),   //start
		mock.AnythingOfType("string"),   //end
		mock.AnythingOfType("string"),   //summary
		mock.AnythingOfType("string"),   //commander
		mock.AnythingOfType("[]string"), //emails
	).Return(f.mockEvent, nil)

	f.mockLogger = loggerMock
	f.mockClient = clientMock
	f.mockRepository = repositoryMock
	f.mockCalendar = calendarMock
}

func TestResolveIncidentDialog(t *testing.T) {
	table := []resolveCommandFixture{
		{
			testName:    "Dialog created properly",
			expectError: false,
		},
	}

	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := commands.ResolveIncidentDialog(f.mockClient, f.triggerID)

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

func TestResolveIncidentByDialog(t *testing.T) {
	table := []resolveCommandFixture{
		{
			testName:        "Incident Resolved properly, with PM Meeting",
			expectError:     false,
			incidentDetails: buildSubmissionMock("true"),
			channelID:       "CT50JJGP5",
			mockEvent:       buildEventMock(),
		},
		// {
		// 	testName:        "Incident Resolved properly, without PM Meeting",
		// 	expectError:     false,
		// 	incidentDetails: buildSubmissionMock("false"),
		// 	channelID:       "CT50JJGP5",
		// },
	}

	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := commands.ResolveIncidentByDialog(f.ctx, f.mockClient, f.mockLogger, f.mockRepository, f.incidentDetails, f.mockCalendar)

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

func buildGetIncidentMock() model.Incident {
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

func buildEventMock() *model.Event {
	var (
		startDate = time.Date(2020, time.March, 19, 12, 00, 00, 00, time.UTC)
		endDate   = time.Date(2020, time.March, 19, 22, 30, 00, 00, time.UTC)
	)

	return &model.Event{
		EventURL: "www.xunda.com",
		Start:    &startDate,
		End:      &endDate,
		Summary:  "Incident Post Mortem Meeting",
	}
}

func buildSubmissionMock(postMortemMeeting string) bot.DialogSubmission {
	return bot.DialogSubmission{
		Channel: bot.Channel{
			ID:   "CT50JJGP5",
			Name: "inc-xunda",
		},
		User: bot.User{
			ID:   "U0G9QF9C6",
			Name: "Guilherme Fonseca",
		},
		Submission: bot.Submission{
			IncidentDescription: "Incident Resolved!",
			StatusIO:            "status.io",
			PostMortemMeeting:   postMortemMeeting,
		},
	}
}
