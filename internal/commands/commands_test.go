package commands

import (
	"context"
	"fmt"
	"testing"
	"time"

	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/log/zap"
	"hellper/internal/model"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testCommand struct {
	ctx            context.Context
	name           string
	command        string
	trigger        TriggerEvent
	mockLogger     log.Logger
	mockClient     bot.Client
	mockRepository model.IncidentRepository
}

func newTestCommand(t *testing.T, name, cmd string, trigger TriggerEvent) testCommand {
	return testCommand{
		name:    name,
		command: cmd,
		trigger: trigger,
	}
}

func (scenario *testCommand) setup(t *testing.T) {
	var (
		slackMock      = bot.NewClientMock()
		repositoryMock = model.NewIncidentRepositoryMock()
		mockChannel    = slack.Channel{}
	)
	scenario.ctx = context.Background()

	mockChannel.ID = "mockChannel"
	mockChannel.Name = "Mock Channel Name"
	slackMock.On("CreateConversationContext", scenario.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("bool")).Return(&mockChannel, nil)
	slackMock.On("PostMessage", mock.AnythingOfType("string"), mock.Anything).Return(
		scenario.trigger.Channel, time.Now().Format(time.RFC3339), nil,
	)
	slackMock.On(
		"InviteUsersToConversationContext", scenario.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string"),
	).Return(&mockChannel, nil)
	slackMock.On("ListPins", mock.AnythingOfType("string")).Return([]slack.Item{}, nil, nil)

	repositoryMock.On("SetIncident", mock.AnythingOfType("*model.Incident")).Return(nil)
	repositoryMock.On("GetIncident", mock.AnythingOfType("string")).Return(model.Incident{}, nil)
	repositoryMock.On("ListActiveIncidents").Return([]model.Incident{}, nil)

	scenario.mockLogger = zap.NewDefault()
	scenario.mockClient = slackMock
	scenario.mockRepository = repositoryMock
}

func TestCommand(test *testing.T) {
	scenarios := []testCommand{
		newTestCommand(
			test,
			"When empty command is submitted",
			"<@mockbot>",
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When ping command is submitted",
			"<@mockbot> ping",
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When list command is submitted",
			"<@mockbot> list",
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When beer command is submitted",
			"<@mockbot> beer",
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When start command is submitted",
			"<@mockbot> start mock-incident",
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When state command is submitted",
			"<@mockbot> state",
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When close command is submitted",
			"<@mockbot> close",
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When add user command is submitted",
			"<@mockbot> adduser <@userid1>",
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When add users command is submitted",
			"<@mockbot> addusers <@userid1> <@userid2> <@userid3>",
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When set description command is submitted",
			"<@mockbot> setdescription my long description of my",
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When set description alias command is submitted",
			"<@mockbot> setdesc my long description with alias command",
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When start command with all arguments is submitted",
			`<@mockbot> start mock_inc-name_0987654321 'Mock Incident, a better vision of madness!' --users=<@user1> <@user2> <@user3>`,
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When start command with a mass of arguments is submitted",
			`<@bot99fgx>                      start     mybug               'meu camelo 2222 _!!!!!!!??????????,@@@@#####$$$$$%&&&&&,,******,,((())))),,, -- meu bode minha cabr a' --users=dadsdasdas lkçdkkaçslklça çldkasçkçldaskçl`,
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
		newTestCommand(
			test,
			"When start command with a mass of arguments is submitted",
			`<@bot99fgx>start bug_name'meu camelo 2222 _!!!!!!!??????????,@@@@#####$$$$$%&&&&&,,******,,((())))),,, -- meu bode minha cabr a' --users=dadsdasdas,lkçdkkaçslklça çldkasçkçldaskçl`,
			TriggerEvent{
				Type:    "mockType",
				User:    "mockUser",
				Channel: "mockChannel",
			},
		),
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("%d-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				executor := NewEventExecutor(scenario.mockLogger, scenario.mockClient, scenario.mockRepository)
				err := executor.ExecuteEventCommand(context.Background(), scenario.command, scenario.trigger)
				require.Nil(t, err, "execute command error")
			},
		)
	}
}
