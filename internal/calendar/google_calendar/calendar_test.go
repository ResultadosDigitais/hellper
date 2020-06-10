package googlecalendar

import (
	"context"
	"fmt"
	"testing"

	"hellper/internal/google"
	"hellper/internal/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	gCalendar "google.golang.org/api/calendar/v3"
)

type googleCalendarFixture struct {
	testName     string
	expectError  bool
	errorMessage string

	calendarService *googleCalendar

	mockLogger                log.Logger
	mockCalendarService       google.CalendarService
	mockCalendarEventsService google.CalendarEventsService

	ctx                 context.Context
	startDateTime       string
	endDateTime         string
	emails              []string
	commander           string
	summary             string
	mockEvent           *gCalendar.Event
	mockEventInsertCall *gCalendar.EventsInsertCall
	calendarID          string
}

func (f *googleCalendarFixture) setup(t *testing.T) {
	f.ctx = context.Background()

	loggerMock := log.NewLoggerMock()
	loggerMock.On("Error", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()
	f.mockLogger = loggerMock

	f.mockCalendarService = google.NewCalendarServiceMock()

	calendarEventsServiceMock := google.NewCalendarEventsServiceMock()
	calendarEventsServiceMock.On("Insert", f.calendarID, f.mockEvent).Return(new(gCalendar.EventsInsertCall))
	f.mockCalendarEventsService = calendarEventsServiceMock

	f.calendarService = calendarServiceMock(f.mockLogger, f.mockCalendarService, f.mockCalendarEventsService, f.calendarID)
}

func TestInsertEvent(t *testing.T) {
	f := googleCalendarFixture{
		testName:   "The InsertCall is created without problem",
		mockEvent:  newEventMock(),
		calendarID: "lucas.feijo@resultaosdigitais.com.br",
	}
	t.Run("Create event struct", func(t *testing.T) {
		f.setup(t)
		insertCall, _ := f.calendarService.insertEvent(f.ctx, f.mockLogger, f.mockEvent)
		assert.IsType(t, new(gCalendar.EventsInsertCall), insertCall)
	})
}
func TestEvent(t *testing.T) {
	f := googleCalendarFixture{
		startDateTime: `2020-05-27T16:00:00-07:00`,
		endDateTime:   `2020-05-27T17:00:00-07:00`,
		emails:        []string{},
		commander:     `guilherme.fonseca@resultadosdigitais.com.br`,
		summary:       `Test postmortem event`,
		mockEvent:     newEventMock(),
	}

	t.Run("Create event struct", func(t *testing.T) {
		event := event(f.startDateTime, f.endDateTime, f.summary, f.emails, f.commander)
		ok := assert.Equal(t, f.mockEvent, event)
		if !ok {
			t.Fatal("fail")
		}
	})
}

func TestCreateCalendarEvent(t *testing.T) {
	table := []googleCalendarFixture{}

	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			err := f.calendarService.CreateCalendarEvent()

			if f.expectError {
				if err == nil {
					t.Fatal("an error was expected, but not occurred")
				}

				assert.EqualError(t, err, f.errorMessage)
			}

			if !f.expectError && err != nil {
				t.Fatal("an error occurred, but was not expected")
			}
		})
	}
}

func newEventMock() *gCalendar.Event {
	return &gCalendar.Event{
		Attendees: []*gCalendar.EventAttendee{
			{
				Email:     `guilherme.fonseca@resultadosdigitais.com.br`,
				Organizer: true,
			},
		},
		Start: &gCalendar.EventDateTime{
			DateTime: `2020-05-27T16:00:00-07:00`,
		},
		End: &gCalendar.EventDateTime{
			DateTime: `2020-05-27T17:00:00-07:00`,
		},
		Summary: `Test postmortem event`,
	}
}

func calendarServiceMock(
	logger log.Logger,
	service google.CalendarService,
	eventsService google.CalendarEventsService,
	calendarID string,
) *googleCalendar {

	return &googleCalendar{
		logger:          logger,
		calendarService: service,
		eventsService:   eventsService,
		calendarID:      calendarID,
	}
}
