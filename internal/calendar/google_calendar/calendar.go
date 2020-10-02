package googlecalendar

import (
	"context"
	"hellper/internal/calendar"
	"hellper/internal/google"
	googleauth "hellper/internal/google_auth"
	"hellper/internal/log"
	"hellper/internal/model"
	"time"

	"github.com/google/uuid"

	gCalendar "google.golang.org/api/calendar/v3"
)

type googleCalendar struct {
	logger          log.Logger
	calendarService google.CalendarService
	eventsService   google.CalendarEventsService
	calendarID      string
}

//NewCalendar initialize the calendar service
func NewCalendar(
	ctx context.Context,
	logger log.Logger,
	calendarToken string,
	calendarID string,
) (calendar.Calendar, error) {
	calendarTokenBytes := []byte(calendarToken)

	gClient, err := googleauth.Struct.GetGClient(ctx, logger, calendarTokenBytes, gCalendar.CalendarScope)
	if err != nil {
		logger.Error(ctx, log.Trace(), log.Action("GetGClient"), log.Reason(err.Error()))
		return nil, err
	}

	calendarService, err := google.NewCalendarService(gClient)
	if err != nil {
		logger.Error(ctx, log.Trace(), log.Action("NewCalendarService"), log.Reason(err.Error()))
		return nil, err
	}

	s := calendarService.(*gCalendar.Service)
	eventsService := google.NewCalendarEventsService(s)

	calendar := googleCalendar{
		logger:          logger,
		calendarService: calendarService,
		eventsService:   eventsService,
		calendarID:      calendarID,
	}

	return &calendar, nil
}

func eventAttendee(email string, commander bool) *gCalendar.EventAttendee {
	return &gCalendar.EventAttendee{
		Email:     email,
		Organizer: commander,
	}
}

// eventDateTime receives a date-time value (formatted
// according to RFC3339) with time zone offset.
func eventDateTime(datetime string) *gCalendar.EventDateTime {
	return &gCalendar.EventDateTime{
		DateTime: datetime,
	}
}

func event(start, end, summary, commander string, emails []string) *gCalendar.Event {
	var attendees []*gCalendar.EventAttendee
	attendees = append(attendees, eventAttendee(commander, true))
	for _, email := range emails {
		if email != "" && email != commander {
			attendees = append(attendees, eventAttendee(email, false))
		}
	}

	eventStart := eventDateTime(start)
	eventEnd := eventDateTime(end)
	conferenceData := conferenceData()

	return &gCalendar.Event{
		Attendees:       attendees,
		Start:           eventStart,
		End:             eventEnd,
		Summary:         summary,
		ConferenceData:  conferenceData,
		GuestsCanModify: true,
	}
}

func (gc *googleCalendar) insertEvent(ctx context.Context, event *gCalendar.Event) (*gCalendar.Event, error) {
	call := gc.eventsService.Insert(gc.calendarID, event)
	gcEvent, err := gc.handleInsertEvent(ctx, call)
	if err != nil {
		gc.logger.Error(ctx, log.Trace(), log.Action("gc.handleInsertEvent"), log.Reason(err.Error()))
		return nil, err
	}

	return gcEvent, nil
}

func (gc *googleCalendar) handleInsertEvent(ctx context.Context, insertCall google.CalendarEventsInsertCall) (*gCalendar.Event, error) {
	insertCall = insertCall.Context(ctx)
	insertCall = insertCall.ConferenceDataVersion(1)

	gcEvent, err := insertCall.Do()
	if err != nil {
		gc.logger.Error(ctx, log.Trace(), log.Action("insertCall.Do"), log.Reason(err.Error()))
		return nil, err
	}

	return gcEvent, nil
}

func googleMeetKey() *gCalendar.ConferenceSolutionKey {
	return &gCalendar.ConferenceSolutionKey{
		Type: "hangoutsMeet",
	}
}

func requestID() string {
	return uuid.New().String()
}

func createConferenceRequest() *gCalendar.CreateConferenceRequest {
	return &gCalendar.CreateConferenceRequest{
		ConferenceSolutionKey: googleMeetKey(),
		RequestId:             requestID(),
	}
}

func conferenceData() *gCalendar.ConferenceData {
	return &gCalendar.ConferenceData{
		CreateRequest: createConferenceRequest(),
	}
}

//CreateCalendarEvent creates a event in Google Calendar
func (gc *googleCalendar) CreateCalendarEvent(ctx context.Context, start, end, summary, commander string, emails []string) (*model.Event, error) {
	e := event(start, end, summary, commander, emails)
	googleEvent, err := gc.insertEvent(ctx, e)
	if err != nil {
		gc.logger.Error(ctx, log.Trace(), log.Action("gc.insertEvent"), log.Reason(err.Error()))
		return nil, err
	}

	eventStart, err := time.Parse(time.RFC3339, googleEvent.Start.DateTime)
	if err != nil {
		gc.logger.Error(ctx, log.Trace(), log.Action("time.Parse Start"), log.Reason(err.Error()))
		return nil, err
	}

	eventEnd, err := time.Parse(time.RFC3339, googleEvent.End.DateTime)
	if err != nil {
		gc.logger.Error(ctx, log.Trace(), log.Action("time.Parse End"), log.Reason(err.Error()))
		return nil, err
	}

	modelEvent := &model.Event{
		EventURL: googleEvent.HtmlLink,
		Start:    &eventStart,
		End:      &eventEnd,
		Summary:  googleEvent.Summary,
	}

	return modelEvent, nil
}
