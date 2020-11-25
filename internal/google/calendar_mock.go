package google

import (
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/calendar/v3"
	googleapi "google.golang.org/api/googleapi"
)

type CalendarServiceMock struct {
	mock.Mock
}

func NewCalendarServiceMock() *CalendarServiceMock {
	return new(CalendarServiceMock)
}

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

type CalendarEventsInsertCallMock struct {
	mock.Mock
}

func NewCalendarEventsInsertCallMock() *CalendarEventsInsertCallMock {
	return new(CalendarEventsInsertCallMock)
}

func (mock *CalendarEventsInsertCallMock) Do(opts ...googleapi.CallOption) (*calendar.Event, error) {
	var (
		args   = mock.Called()
		result = args.Get(0)
	)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*calendar.Event), args.Error(1)
}
