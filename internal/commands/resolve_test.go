package commands_test

import (
	"context"
	"hellper/internal/bot"
	calendar "hellper/internal/calendar"
	"hellper/internal/log"
	"hellper/internal/model"
	"testing"

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

	channelID    string
	mockIncident model.Incident
	mockEvent    *model.Event
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
		mock.AnythingOfType("string"),       //triggerID
		mock.AnythingOfType("slack.Dialog"), //dialog
	).Return(nil)
	clientMock.On(
		"AddPin",
		mock.AnythingOfType("string"),        //channel
		mock.AnythingOfType("slack.ItemRef"), //item
	).Return(nil)
	clientMock.On(
		"PostMessage",
		f.channelID,                              //channel
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
