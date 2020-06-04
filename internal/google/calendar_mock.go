package google

import (
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/calendar/v3"
)

type CalendarEventsServiceMock struct {
	mock.Mock
}

func NewCalendarEventsServiceMock() *CalendarEventsServiceMock {
	return new(CalendarEventsServiceMock)
}

func (mock *CalendarEventsServiceMock) Insert(calendarID string, event *calendar.Event) *calendar.EventsInsertCall {
	args := mock.Called(calendarID, event)
	return args.Get(0).(*calendar.EventsInsertCall)
}

type CalendarServiceMock struct {
	mock.Mock
}

func NewCalendarServiceMock() *CalendarServiceMock {
	return new(CalendarServiceMock)
}
