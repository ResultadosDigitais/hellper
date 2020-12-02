package commands_test

import (
	"context"
	"fmt"
	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/log"
	"hellper/internal/model"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type statusCommandFixture struct {
	testName     string
	expectError  bool
	errorMessage string

	ctx            context.Context
	mockClient     bot.Client
	mockLogger     log.Logger
	mockRepository model.IncidentRepository

	channelID string
	userID    string
}

func (f *statusCommandFixture) setup(t *testing.T) {
	var (
		loggerMock     = log.NewLoggerMock()
		clientMock     = bot.NewClientMock()
		repositoryMock = model.NewRepositoryMock()
	)
	f.ctx = context.Background()

	//Logger Mock
	loggerMock.On(
		"Info",
		f.ctx,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]log.Value"),
	).Return()
	loggerMock.On(
		"Error",
		f.ctx,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]log.Value"),
	).Return()

	//Client Mock
	clientMock.On(
		"ListPins",
		f.channelID,
	).Return([]slack.Item{}, new(slack.Paging), nil)
	clientMock.On(
		"PostMessage",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]slack.MsgOption"),
	).Return("", "", nil)

	//Repository Mock
	repositoryMock.On(
		"GetIncident",
		f.channelID, //channelID
	).Return(model.Incident{}, nil)

	f.mockLogger = loggerMock
	f.mockClient = clientMock
	f.mockRepository = repositoryMock
}

func TestShowStatus(t *testing.T) {
	table := []statusCommandFixture{
		{
			testName:    "Dialog created properly",
			expectError: false,
		},
	}

	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := commands.ShowStatus(f.ctx, f.mockClient, f.mockLogger, f.mockRepository, f.channelID, f.userID)

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
