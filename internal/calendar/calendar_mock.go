package calendar

import "github.com/stretchr/testify/mock"

type CalendarMock struct {
	mock.Mock
}

func NewCalendarMock() *CalendarMock {
	return new(CalendarMock)
}

func (*CalendarMock) CreateCalendarEvent() {

}
