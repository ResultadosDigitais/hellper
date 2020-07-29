package commands_test

import (
	"context"
	"fmt"
	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/log"
	"hellper/internal/model"
	"testing"

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
	).Return(model.Incident{}, nil)

	//Client Mock
	clientMock.On(
		"PostEphemeralContext",
		f.ctx,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]slack.MsgOption"),
	).Return("", nil)

	f.mockLogger = loggerMock
	f.mockClient = clientMock
	f.mockRepository = repositoryMock

}
func TestOpenCancelIncidentDiolog(t *testing.T) {
	table := []cancelCommandFixture{
		{
			testName:    "Check error if incident is not open",
			expectError: false,
			channelID:   "ABCD",
			userID:      "ABCD",
			triggerID:   "ABCD",
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
