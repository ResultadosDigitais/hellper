package googlecalendar

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	gCalendar "google.golang.org/api/calendar/v3"
)

type googleCalendarFixture struct {
	testName     string
	expectError  bool
	errorMessage string

	calendarService googleCalendar

	startDateTime string
	endDateTime   string
	emails        []string
	commander     string
	summary       string
	mockEvent     *gCalendar.Event
}

func (f *googleCalendarFixture) setup(t *testing.T) {

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
