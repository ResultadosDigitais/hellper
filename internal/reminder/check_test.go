package reminder_test

import (
	"context"
	"database/sql"
	"fmt"
	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/model"
	"hellper/internal/reminder"
	"strconv"
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type checkReminderFixture struct {
	testName       string
	expected       bool
	ctx            context.Context
	mockLogger     log.Logger
	mockClient     bot.Client
	mockRepository model.Repository
	mockIncident   model.Incident

	channelID     string
	listPinsError error
	lastPin       []slack.Item
}

func (f *checkReminderFixture) setup(t *testing.T) {
	var (
		loggerMock     = log.NewLoggerMock()
		clientMock     = bot.NewClientMock()
		repositoryMock = model.NewRepositoryMock()
	)

	f.ctx = context.Background()

	loggerMock.On("Info", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()
	clientMock.On("ListPins", mock.AnythingOfType("string")).Return(f.lastPin, nil, nil)

	f.mockLogger = loggerMock
	f.mockClient = clientMock
	f.mockRepository = repositoryMock

}
func TestCanSendNotify(t *testing.T) {

	table := []checkReminderFixture{
		{
			testName: "Notify when status is open",
			expected: true,
			lastPin:  lastPin(config.Env.ReminderOpenStatusSeconds, -30),
			mockIncident: model.Incident{
				EndedAt:      &time.Time{},
				SnoozedUntil: sql.NullTime{},
				Status:       "open",
			},
		},
		{
			testName: "Do not notify when status is open and has pinned message",
			expected: false,
			lastPin:  lastPin(config.Env.ReminderOpenStatusSeconds, 30),
			mockIncident: model.Incident{
				EndedAt:      &time.Time{},
				SnoozedUntil: sql.NullTime{},
				Status:       "open",
			},
		},
		{
			testName: "Notify when status is resolved and SLA > 7 days",
			expected: true,
			mockIncident: model.Incident{
				EndedAt:      &[]time.Time{time.Now().AddDate(0, 0, -8)}[0],
				SnoozedUntil: sql.NullTime{},
				Status:       "resolved",
			},
		},
		{
			testName: "Do not notify when status is resolved and SLA <= 7 days",
			expected: false,
			mockIncident: model.Incident{
				EndedAt:      &[]time.Time{time.Now().AddDate(0, 0, -7)}[0],
				SnoozedUntil: sql.NullTime{},
				Status:       "resolved",
			},
		},
		{
			testName: "Do not notify when status is open and notify is paused",
			expected: false,
			mockIncident: model.Incident{
				EndedAt:      &time.Time{},
				SnoozedUntil: sql.NullTime{Time: time.Now().AddDate(0, 0, 1)},
				Status:       "open",
			},
		},
		{
			testName: "Do not notify when status is resolved and SLA > 7 days and notify is paused",
			expected: false,
			mockIncident: model.Incident{
				EndedAt:      &[]time.Time{time.Now().AddDate(0, 0, -8)}[0],
				SnoozedUntil: sql.NullTime{Time: time.Now().AddDate(0, 0, 3)},
				Status:       "resolved",
			},
		},
		{
			testName: "Do not notify when status is closed",
			expected: false,
			mockIncident: model.Incident{
				EndedAt:      &[]time.Time{time.Now()}[0],
				SnoozedUntil: sql.NullTime{},
				Status:       "closed",
			},
		},
		{
			testName: "Do not notify when status is canceled",
			expected: false,
			mockIncident: model.Incident{
				EndedAt:      &[]time.Time{time.Now()}[0],
				SnoozedUntil: sql.NullTime{},
				Status:       "canceled",
			},
		},
		{
			testName: "Do not notify in any other status",
			expected: false,
			mockIncident: model.Incident{
				EndedAt:      &[]time.Time{time.Now()}[0],
				SnoozedUntil: sql.NullTime{},
				Status:       "xyzxyzxyz",
			},
		},
	}

	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)
			hasNotify := reminder.CanSendNotify(f.ctx, f.mockClient, f.mockLogger, f.mockRepository, f.mockIncident)
			assert.Equal(t, f.expected, hasNotify, "they should be equal")
		})
	}

}

func lastPin(env int, diff int64) []slack.Item {
	return []slack.Item{
		{
			Message: &slack.Message{
				Msg: slack.Msg{
					Timestamp: strconv.Itoa(int(time.Now().Add((time.Second*-time.Duration(env))+time.Second*time.Duration(diff)).Unix())) + ".0",
				},
			},
		},
	}
}
