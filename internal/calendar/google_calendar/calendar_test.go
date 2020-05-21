package googlecalendar

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type googleCalendarFixture struct {
	testName     string
	expectError  bool
	errorMessage string

	calendarService googleCalendar
}

func (f *googleCalendarFixture) setup(t *testing.T) {

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
