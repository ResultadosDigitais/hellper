package googlecalendar

import (
	"context"
	"hellper/internal/calendar"
	"hellper/internal/config"
	googleauth "hellper/internal/google_auth"
	"hellper/internal/log"

	gCalendar "google.golang.org/api/calendar/v3"
)

type googleCalendar struct {
	logger          log.Logger
	calendarService *gCalendar.Service
}

//NewCalendar initialize the file storage service
func NewCalendar(ctx context.Context, logger log.Logger) calendar.Calendar {
	calendarTokenBytes := []byte(config.Env.GoogleCalendarToken)

	gClient, err := googleauth.Struct.GetGClient(ctx, logger, calendarTokenBytes, gCalendar.CalendarScope)
	if err != nil {
		logger.Error(
			ctx,
			"googlecalendar/calendar.NewCalendar GetGClient ERROR",
			log.NewValue("error", err),
		)

		return nil
	}

	calendarService, err := gCalendar.New(gClient)
	if err != nil {
		logger.Error(
			ctx,
			"googlecalendar/calendar.NewCalendar gCalendar.New ERROR",
			log.NewValue("error", err),
		)

		return nil
	}

	return &googleCalendar{
		logger:          logger,
		calendarService: calendarService,
	}
}

//CreateCalendarEvent creates a event in Google Calendar
func (*googleCalendar) CreateCalendarEvent() error {
	return nil
}
