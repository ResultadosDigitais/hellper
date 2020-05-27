package googlecalendar

import (
	"context"
	"fmt"
	"testing"

	"hellper/internal/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type googleCalendarFixture struct {
	testName     string
	expectError  bool
	errorMessage string

	calendarService googleCalendar

	ctx        context.Context
	mockLogger log.Logger
}

func (f *googleCalendarFixture) setup(t *testing.T) {
	f.ctx = context.Background()

	loggerMock := log.NewLoggerMock()
	loggerMock.On("Error", f.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]log.Value")).Return()
	f.mockLogger = loggerMock
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
