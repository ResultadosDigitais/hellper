package commands

import (
	"errors"
	"fmt"
	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

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
		SlackID: "U013NAYQ1PG",
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

func TestGetStringInt64(test *testing.T) {
	scenarios := []struct {
		name         string
		input        string
		output       int64
		expectError  bool
		errorMessage string
	}{
		{
			name:   "Convert string to int",
			input:  "0",
			output: 0,
		},
		{
			name:         "Returns error when try to convert",
			input:        "1.2",
			expectError:  true,
			errorMessage: "strconv.ParseInt: parsing \"1.2\": invalid syntax",
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				result, err := getStringInt64(scenario.input)
				if scenario.expectError {
					if err == nil {
						t.Fatal("an error was expected, but not occurred")
					}
					assert.EqualError(t, err, scenario.errorMessage)
				} else {
					if err != nil {
						t.Fatal("an error occurred, but was not expected")
					}
					assert.Equal(t, scenario.output, result, "The return should be the same")
				}

			},
		)
	}
}

func TestGetSeverityLevelText(test *testing.T) {
	scenarios := []struct {
		name   string
		input  int64
		output string
	}{
		{
			name:   "Test case 0",
			input:  0,
			output: "SEV0 - All hands on deck",
		},
		{
			name:   "Test case 1",
			input:  1,
			output: "SEV1 - Critical impact to many users",
		},
		{
			name:   "Test case 2",
			input:  2,
			output: "SEV2 - Minor issue that impacts ability to use product",
		},
		{
			name:   "Test case 3",
			input:  3,
			output: "SEV3 - Minor issue not impacting ability to use product",
		},
		{
			name:   "Test case default",
			input:  123456,
			output: "",
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				result := getSeverityLevelText(scenario.input)
				require.Equal(t, scenario.output, result, "The return should be the same")
			},
		)
	}
}

func TestConvertTimestamp(test *testing.T) {
	scenarios := []struct {
		name         string
		input        string
		output       time.Time
		expectError  bool
		errorMessage string
	}{
		{
			name:         "When the string is empty",
			input:        "",
			expectError:  true,
			errorMessage: "Empty Timestamp",
		},
		{
			name:   "Convert string to timestamp",
			input:  "1512085950.000216",
			output: time.Unix(1512085950, 216),
		},
		{
			name:         "An error is returned when seconds is invalid",
			input:        "1512085950.abcde",
			expectError:  true,
			errorMessage: "strconv.ParseInt: parsing \"abcde\": invalid syntax",
		},
		{
			name:         "An error is returned when minutes is invalid",
			input:        "abcde.000216",
			expectError:  true,
			errorMessage: "strconv.ParseInt: parsing \"abcde\": invalid syntax",
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				result, err := convertTimestamp(scenario.input)
				if scenario.expectError {
					if err == nil {
						t.Fatal("an error was expected, but not occurred")
					}
					assert.EqualError(t, err, scenario.errorMessage)
				} else {
					if err != nil {
						t.Fatal("an error occurred, but was not expected")
					}
					assert.Equal(t, scenario.output, result, "The return should be the same")
				}

			},
		)
	}
}
