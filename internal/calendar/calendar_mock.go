package calendar

import (
	"context"
	"hellper/internal/model"

	"github.com/stretchr/testify/mock"
)

type CalendarMock struct {
	mock.Mock
}

func NewCalendarMock() *CalendarMock {
	return new(CalendarMock)
}

func (mock *CalendarMock) CreateCalendarEvent(ctx context.Context, start, end, summary, commander string, emails []string) (*model.Event, error) {
	var (
		args   = mock.Called(ctx, start, end, summary, commander, emails)
		result = args.Get(0)
	)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*model.Event), args.Error(1)
}
