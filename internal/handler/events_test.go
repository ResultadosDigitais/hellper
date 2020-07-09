package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/log/zap"
	"hellper/internal/model"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testHandler struct {
	name           string
	body           string
	mockClient     *bot.ClientMock
	mockRepository *model.RepositoryMock
	handler        *handlerEvents
	responseStatus int
	request        *http.Request
	response       *httptest.ResponseRecorder
}

func newTestHandler(name, body string, responseStatus int) testHandler {
	return testHandler{
		name:           name,
		body:           body,
		responseStatus: responseStatus,
	}
}

func (scenario *testHandler) setup(*testing.T) {
	msgsCache = map[string]struct{}{}
	config.Env.VerificationToken = "7WV2asfPzOnZyh9JnBwBiUKu"
	slackMock := bot.NewClientMock()
	slackMock.On("PostMessage",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return("mockChannel", time.Now().Format(time.RFC3339), nil)
	slackMock.On("CreateConversationContext", mock.AnythingOfType("context.Context"), mock.AnythingOfType("string"), mock.AnythingOfType("bool")).Return(new(slack.Channel), nil).Once()
	slackMock.On(
		"InviteUsersToConversationContext",
		mock.AnythingOfType("context.Context"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(
		new(slack.Channel),
		nil,
	)

	repositoryMock := model.NewRepositoryMock()
	repositoryMock.On(
		"SetIncident", mock.AnythingOfType("*model.Incident"),
	).Return(nil).Once()
	repositoryMock.On("GetIncident", mock.Anything).Return(model.Incident{}, nil)

	r := httptest.NewRequest(
		"POST", "/events", strings.NewReader(scenario.body),
	)
	r.Header.Set("content-type", "application/json")

	scenario.mockClient = slackMock
	scenario.mockRepository = repositoryMock
	scenario.request = r
	scenario.response = httptest.NewRecorder()
}

func TestHandler(test *testing.T) {
	scenarios := []testHandler{
		newTestHandler(
			"When body is empty",
			``,
			400,
		),
		newTestHandler(
			"When body has an url verification event",
			`{
				"token":"7WV2asfPzOnZyh9JnBwBiUKu",
				"challenge":"8BrllE22VHrB9gY1dOHHSKYpMc8ryS4qjmQyh6ofsOQKdKPBjt7D",
				"type":"url_verification"
			}`,
			200,
		),
		newTestHandler(
			"When body has a callback with a message event",
			`{
				"token":"7WV2asfPzOnZyh9JnBwBiUKu",
				"team_id":"TEK53T5SP",
				"api_app_id":"AEWA14UE6",
				"event":{
					"client_msg_id":"d5651eb7-2609-4a2b-903f-92835c694de8",
					"type":"message",
					"text":"<@UEVHT00G0> ping",
					"user":"UEV85SUTS",
					"ts":"1545096428.000500",
					"channel":"DEVHT026L",
					"event_ts":"1545096428.000500",
					"channel_type":"im"
				},
				"type":"event_callback",
				"event_id":"EvEWB1TQTC",
				"event_time":1545096428,
				"authed_users":["UEVHT00G0"]
			}`,
			202,
		),
		newTestHandler(
			"When body has a callback with a app mention event",
			`{
				"token":"7WV2asfPzOnZyh9JnBwBiUKu",
				"team_id":"TEK53T5SP",
				"api_app_id":"AEWA14UE6",
				"event":{
					"client_msg_id":"78be304a-269d-42e1-b398-37914c5ab031",
					"type":"app_mention",
					"text":"<@UEVHT00G0> ping",
					"user":"UEV85SUTS",
					"ts":"1545096726.001100",
					"channel":"CEVJU8C1E",
					"event_ts":"1545096726.001100"
				},
				"type":"event_callback",
				"event_id":"EvEVJVFXUY",
				"event_time":1545096726,
				"authed_users":["UEVHT00G0"]
			}`,
			202,
		),
		newTestHandler(
			"When body has a callback event with a help command",
			`{
				"token":"7WV2asfPzOnZyh9JnBwBiUKu",
				"team_id":"TEK53T5SP",
				"api_app_id":"AEWA14UE6",
				"event":{
					"client_msg_id":"78be304a-269d-42e1-b398-37914c5ab031",
					"type":"app_mention",
					"text":"<@UEVHT00G0> help",
					"user":"UEV85SUTS",
					"ts":"1545096726.001100",
					"channel":"CEVJU8C1E",
					"event_ts":"1545096726.001100"
				},
				"type":"event_callback",
				"event_id":"EvEVJVFXUY",
				"event_time":1545096726,
				"authed_users":["UEVHT00G0"]
			}`,
			202,
		),
		newTestHandler(
			"When body has a callback event with a list command",
			`{
				"token":"7WV2asfPzOnZyh9JnBwBiUKu",
				"team_id":"TEK53T5SP",
				"api_app_id":"AEWA14UE6",
				"event":{
					"client_msg_id":"78be304a-269d-42e1-b398-37914c5ab031",
					"type":"app_mention",
					"text":"<@UEVHT00G0> ping",
					"user":"UEV85SUTS",
					"ts":"1545096726.001100",
					"channel":"CEVJU8C1E",
					"event_ts":"1545096726.001100"
				},
				"type":"event_callback",
				"event_id":"EvEVJVFXUY",
				"event_time":1545096726,
				"authed_users":["UEVHT00G0"]
			}`,
			202,
		),
		newTestHandler(
			"When body has a callback event with a beer command",
			`{
				"token":"7WV2asfPzOnZyh9JnBwBiUKu",
				"team_id":"TEK53T5SP",
				"api_app_id":"AEWA14UE6",
				"event":{
					"client_msg_id":"78be304a-269d-42e1-b398-37914c5ab031",
					"type":"app_mention",
					"text":"<@UEVHT00G0> beer",
					"user":"UEV85SUTS",
					"ts":"1545096726.001100",
					"channel":"CEVJU8C1E",
					"event_ts":"1545096726.001100"
				},
				"type":"event_callback",
				"event_id":"EvEVJVFXUY",
				"event_time":1545096726,
				"authed_users":["UEVHT00G0"]
			}`,
			202,
		),
		newTestHandler(
			"When body has a callback with a bot message",
			`{
				"token":"7WV2asfPzOnZyh9JnBwBiUKu",
				"team_id":"TEK53T5SP",
				"api_app_id":"AEWA14UE6",
				"event":{
					"type":"message",
					"subtype":"bot_message",
					"text":"pong",
					"ts":"1545120735.003300",
					"username":"rjansen-bot",
					"bot_id":"BEVHT00E4",
					"channel":"DEVHT026L",
					"event_ts":"1545120735.003300",
					"channel_type":"im"
				},
				"type":"event_callback",
				"event_id":"EvEWMZCD45",
				"event_time":1545120735,
				"authed_users":["UEVHT00G0"]
			}`,
			202,
		),
		newTestHandler(
			"When body has a callback with a message event and an unkwon command",
			`{
				"token":"7WV2asfPzOnZyh9JnBwBiUKu",
				"team_id":"TEK53T5SP",
				"api_app_id":"AEWA14UE6",
				"event":{
					"client_msg_id":"d5651eb7-2609-4a2b-903f-92835c694de8",
					"type":"message",
					"text":"<@UEVHT00G0> pingxxxxxx",
					"user":"UEV85SUTS",
					"ts":"1545096428.000500",
					"channel":"DEVHT026L",
					"event_ts":"1545096428.000500",
					"channel_type":"im"
				},
				"type":"event_callback",
				"event_id":"EvEWB1TQTC",
				"event_time":1545096428,
				"authed_users":["UEVHT00G0"]
			}`,
			202,
		),
		newTestHandler(
			"When body has a callback event with a start command",
			`{
				"token":"7WV2asfPzOnZyh9JnBwBiUKu",
				"team_id":"TEK53T5SP",
				"api_app_id":"AEWA14UE6",
				"event":{
					"client_msg_id":"78be304a-269d-42e1-b398-37914c5ab031",
					"type":"app_mention",
					"text":"<@UEVHT00G0> start mock-incident 'mock description' --users=<@user1> <@user2> <@user3>",
					"user":"UEV85SUTS",
					"ts":"1545096726.001100",
					"channel":"CEVJU8C1E",
					"event_ts":"1545096726.001100"
				},
				"type":"event_callback",
				"event_id":"EvEVJVFXUY",
				"event_time":1545096726,
				"authed_users":["UEVHT00G0"]
			}`,
			202,
		),
		newTestHandler(
			"When body has a callback event with a report command",
			`{
				"token":"7WV2asfPzOnZyh9JnBwBiUKu",
				"team_id":"TEK53T5SP",
				"api_app_id":"AEWA14UE6",
				"event":{
					"client_msg_id":"78be304a-269d-42e1-b398-37914c5ab031",
					"type":"app_mention",
					"text":"<@UEVHT00G0> report",
					"user":"UEV85SUTS",
					"ts":"1545096726.001100",
					"channel":"CEVJU8C1E",
					"event_ts":"1545096726.001100"
				},
				"type":"event_callback",
				"event_id":"EvEVJVFXUY",
				"event_time":1545096726,
				"authed_users":["UEVHT00G0"]
			}`,
			202,
		),
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("%d-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				h := newHandlerEvents(zap.NewDefault(), scenario.mockClient, scenario.mockRepository)
				h.ServeHTTP(scenario.response, scenario.request)
				result := scenario.response.Result()
				require.Equal(t, scenario.responseStatus, result.StatusCode, "invalid statuscode value")
			},
		)
	}
}
