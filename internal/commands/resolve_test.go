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

	loggerMock.On("Info", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()
	loggerMock.On("Error", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()

	clientMock.On("OpenDialog", mock.AnythingOfType("string"), mock.AnythingOfType("slack.Dialog")).Return(nil)
	clientMock.On("AddPin", mock.AnythingOfType("string"), mock.AnythingOfType("slack.ItemRef")).Return(nil)
	clientMock.On("PostMessage", f.channelID, mock.AnythingOfType("[]slack.MsgOption")).Return("", "", nil)

	repositoryMock.On("ResolveIncident", f.ctx, mock.AnythingOfType("*model.Incident")).Return(nil)
	repositoryMock.On("GetIncident", f.channelID).Return(f.mockIncident, nil)

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
