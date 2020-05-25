package commands

import (
	"errors"
	"fmt"
	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"context"
	"testing"
)

type userCommandFixture struct {
	testName                string
	expectError             bool
	errorMessage            string
	getUserInfoContextError error
	mockUser                *model.User
	mockSlackUser           *slack.User

	ctx        context.Context
	mockLogger log.Logger
	mockClient bot.Client

	userID string
}

func (f *userCommandFixture) setup(t *testing.T) {
	var (
		loggerMock = log.NewLoggerMock()
		clientMock = bot.NewClientMock()
	)

	f.ctx = context.Background()

	loggerMock.On("Info", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()
	loggerMock.On("Error", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()
	clientMock.On("GetUserInfoContext", f.ctx, mock.AnythingOfType("string")).Return(f.mockSlackUser, f.getUserInfoContextError)

	f.mockLogger = loggerMock
	f.mockClient = clientMock
}

func TestGetSlackUserInfo(t *testing.T) {
	users := []userCommandFixture{
		{
			testName:      "Get slack user info with valid slack user Id",
			expectError:   false,
			userID:        "U013NAYQ1PG",
			mockUser:      buildUserMock(),
			mockSlackUser: buildSlackMock(),
		},
		{
			testName:                "User not found",
			expectError:             true,
			userID:                  "xunda",
			errorMessage:            "User not found",
			getUserInfoContextError: errors.New("User not found"),
		},
	}

	for index, f := range users {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			user, err := getSlackUserInfo(f.ctx, f.mockClient, f.mockLogger, f.userID)

			if f.expectError {
				if err == nil {
					t.Fatal("an error was expected, but not occurred")
				}
				assert.EqualError(t, err, f.errorMessage)
			}

			if !f.expectError {
				if err != nil {
					t.Fatal("an error occurred, but was not expected")
				}
				assert.Equal(t, f.mockUser, user)
			}
		})
	}
}

func buildUserMock() *model.User {

	return &model.User{
		Id:      0,
		SlackId: "U013NAYQ1PG",
		Name:    "Natalia Favareto",
		Email:   "natalia.favareto@resultadosdigitais.com.br",
	}
}

func buildSlackMock() *slack.User {

	return &slack.User{
		ID: "U013NAYQ1PG",
		Profile: slack.UserProfile{
			RealName: "Natalia Favareto",
			Email:    "natalia.favareto@resultadosdigitais.com.br",
		},
	}
}
