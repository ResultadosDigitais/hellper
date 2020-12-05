package commands

import (
	"context"
	"fmt"
	"hellper/internal/app"
	"hellper/internal/concurrence"
	"strconv"
	"strings"
	"sync"
	"time"

	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

// ResolveIncidentDialog opens a dialog on Slack, so the user can resolve an incident
func ResolveIncidentDialog(app *app.App, triggerID string) error {
	description := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Solution Description",
			Name:        "incident_description",
			Type:        "textarea",
			Placeholder: "Brief description on what was done to solve this incident. eg. The incident was solved in PR #42",
			Optional:    false,
		},
		MaxLength: 500,
	}

	postMortemMeeting := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Should I schedule a Post Mortem meeting?",
			Name:        "post_mortem_meeting",
			Type:        "select",
			Placeholder: "Post Mortem Meeting",
			Optional:    false,
		},
		Value: "false",
		Options: []slack.DialogSelectOption{
			{
				Label: "Yes",
				Value: "true",
			},
			{
				Label: "No",
				Value: "false",
			},
		},
		OptionGroups: []slack.DialogOptionGroup{},
	}

	dialogElements := []slack.DialogElement{
		description,
	}
	if app.Calendar != nil {
		dialogElements = append(dialogElements, postMortemMeeting)
	}

	dialog := slack.Dialog{
		CallbackID:     "inc-resolve",
		Title:          "Resolve an Incident",
		SubmitLabel:    "Resolve",
		NotifyOnCancel: false,
		Elements:       dialogElements,
	}

	return app.Client.OpenDialog(triggerID, dialog)
}

// ResolveIncidentByDialog resolves an incident after receiving data from a Slack dialog
func ResolveIncidentByDialog(
	ctx context.Context,
	app *app.App,
	incidentDetails bot.DialogSubmission,
) error {
	app.Logger.Debug(
		ctx,
		"command/resolve.ResolveIncidentByDialog",
		log.NewValue("incident_resolve_details", incidentDetails),
	)

	var (
		now               = time.Now().UTC()
		channelID         = incidentDetails.Channel.ID
		channelName       = incidentDetails.Channel.Name
		userID            = incidentDetails.User.ID
		userName          = incidentDetails.User.Name
		submissions       = incidentDetails.Submission
		description       = submissions["incident_description"]
		postMortemMeeting = submissions["post_mortem_meeting"]
		notifyOnResolve   = config.Env.NotifyOnResolve
		productChannelID  = config.Env.ProductChannelID

		calendarEvent *model.Event
	)

	incident := model.Incident{
		ChannelID:           channelID,
		EndTimestamp:        &now,
		DescriptionResolved: description,
	}

	logWriter := app.Logger.With(
		log.NewValue("incident", incident),
		log.NewValue("channelID", channelID),
	)

	logWriter.Debug(
		ctx,
		log.Trace(),
		log.Action("running"),
	)

	err := app.IncidentRepository.ResolveIncident(ctx, &incident)
	if err != nil {
		logWriter.Error(
			ctx,
			log.Trace(),
			log.Reason("ResolveIncident"),
			log.NewValue("error", err),
		)
		return err
	}

	inc, err := app.IncidentRepository.GetIncident(ctx, channelID)
	if err != nil {
		logWriter.Error(
			ctx,
			log.Trace(),
			log.Reason("GetIncident"),
			log.NewValue("error", err),
		)
		return err
	}

	hasPostMortemMeeting := false
	if app.Calendar != nil {
		hasPostMortemMeeting, err = strconv.ParseBool(postMortemMeeting)
		if err != nil {
			logWriter.Error(
				ctx,
				log.Trace(),
				log.Reason("strconv.ParseBool"),
				log.NewValue("error", err),
			)
			return err
		}
	}

	if hasPostMortemMeeting {
		calendarEvent, err = getCalendarEvent(ctx, app, incident.EndTimestamp, channelName, channelID)
		if err != nil {
			logWriter.Error(
				ctx,
				log.Trace(),
				log.Reason("getCalendarEvent"),
				log.NewValue("error", err),
			)
			return err
		}
	}

	channelAttachment := createResolveChannelAttachment(inc, userName, calendarEvent)
	privateAttachment := createResolvePrivateAttachment(incident, calendarEvent)
	message := "The Incident <#" + incident.ChannelID + "> has been resolved by <@" + userName + ">"

	var waitgroup sync.WaitGroup
	defer waitgroup.Wait()

	concurrence.WithWaitGroup(&waitgroup, func() {
		postAndPinMessage(
			app,
			channelID,
			message,
			channelAttachment,
		)
	})
	if notifyOnResolve {
		concurrence.WithWaitGroup(&waitgroup, func() {
			postAndPinMessage(
				app,
				productChannelID,
				message,
				channelAttachment,
			)
		})
	}
	postMessage(app, userID, "", privateAttachment)

	return nil
}

func setMeetingDate(ctx context.Context, logger log.Logger, d *time.Time, postMortemGapDays int, timezone string) (string, string, error) {
	previewPostMortemDate := d.AddDate(0, 0, postMortemGapDays)

	switch previewPostMortemDate.Weekday() {
	case time.Saturday:
		postMortemGapDays += 2
	case time.Sunday:
		postMortemGapDays++
	default:
		break
	}

	utc, err := time.LoadLocation(timezone)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Reason("time.LoadLocation"),
			log.NewValue("timezone", timezone),
			log.NewValue("error", err),
		)
		return "", "", err
	}

	setMeetingHour := time.Date(d.Year(), d.Month(), d.Day(), 15, 0, 0, 0, utc)

	startMeeting := setMeetingHour.AddDate(0, 0, postMortemGapDays)
	endMeeting := startMeeting.Add(time.Hour).Format(time.RFC3339)

	return startMeeting.Format(time.RFC3339), endMeeting, err
}

func getCalendarEvent(
	ctx context.Context,
	app *app.App,
	t *time.Time,
	channelName string,
	channelID string,
) (*model.Event, error) {
	logWriter := app.Logger.With(
		log.NewValue("channelID", channelID),
	)

	if app.Calendar == nil {
		logWriter.Error(
			ctx,
			log.Trace(),
			log.Reason("noCalendarConfigured"),
		)
		return nil, fmt.Errorf("Calendar is not available")
	}

	startMeeting, endMeeting, err := setMeetingDate(ctx, app.Logger, t, config.Env.PostmortemGapDays, config.Env.Timezone)
	if err != nil {
		logWriter.Error(
			ctx,
			log.Trace(),
			log.Reason("setMeetingDate"),
			log.NewValue("error", err),
		)
		return nil, err
	}

	summary := "[Post Mortem] " + channelName
	emails, _ := getUsersEmailsInConversation(ctx, app, channelID)

	inc, err := app.IncidentRepository.GetIncident(ctx, channelID)
	if err != nil {
		logWriter.Error(
			ctx,
			log.Trace(),
			log.Reason("GetIncident"),
			log.NewValue("error", err),
		)
		return nil, err
	}

	calendarEvent, err := app.Calendar.CreateCalendarEvent(ctx, startMeeting, endMeeting, summary, inc.CommanderEmail, *emails)
	if err != nil {
		logWriter.Error(
			ctx,
			log.Trace(),
			log.Reason("CreateCalendarEvent"),
			log.NewValue("error", err),
		)
		return nil, err
	}
	return calendarEvent, err
}

func createResolveChannelAttachment(inc model.Incident, userName string, event *model.Event) slack.Attachment {
	var (
		endDateText       = inc.EndTimestamp.Format(time.RFC1123)
		postMortemMessage string
		messageText       strings.Builder
	)

	messageText.WriteString("The Incident <#" + inc.ChannelID + "> has been resolved by <@" + userName + ">\n\n")
	messageText.WriteString("*End date:* <#" + endDateText + ">\n")
	messageText.WriteString("*Description:* `" + inc.DescriptionResolved + "`\n")
	if event == nil {
		messageText.WriteString("*Post Mortem:* A Post Mortem Meeting was not scheduled, but be sure to fill up the Post Mortem document.\n")
		postMortemMessage = "A Post Mortem Meeting was not scheduled, but be sure to fill up the Post Mortem document."
	} else {
		messageText.WriteString("*Post Mortem Meeting Link:* `" + event.EventURL + "`\n\n")
		postMortemMessage = "I have scheduled a Post Mortem Meeting for you!\nIt will be on `" + event.Start.Format(time.RFC1123) + "`.\nHere is the link: `" + event.EventURL + "`\n"
	}

	return slack.Attachment{
		Pretext:  "",
		Fallback: messageText.String(),
		Text:     "",
		Color:    "#1164A3",
		Fields: []slack.AttachmentField{
			{
				Title: "Incident ID",
				Value: strconv.FormatInt(inc.ID, 10),
			},
			{
				Title: "Incident Channel",
				Value: "<#" + inc.ChannelID + ">",
			},
			{
				Title: "Incident Title",
				Value: inc.Title,
			},
			{
				Title: "End date",
				Value: endDateText,
			},
			{
				Title: "Description",
				Value: inc.DescriptionResolved,
			},
			{
				Title: "Post Mortem",
				Value: postMortemMessage,
			},
		},
	}
}

func createResolvePrivateAttachment(inc model.Incident, event *model.Event) slack.Attachment {
	var (
		postMortemMessage string
		privateText       strings.Builder
	)

	privateText.WriteString("The Incident <#" + inc.ChannelID + "> has been resolved by you\n\n")
	if event == nil {
		privateText.WriteString("*Post Mortem:* A Post Mortem Meeting was not scheduled, but be sure to fill up the Post Mortem document.\n")
		postMortemMessage = "A Post Mortem Meeting was not scheduled, but be sure to fill up the Post Mortem document."
	} else {
		privateText.WriteString("*Post Mortem Meeting Link:* `" + event.EventURL + "`\n\n")
		postMortemMessage = "I have scheduled a Post Mortem Meeting for you!\nIt will be on `" + event.Start.Format(time.RFC1123) + "`.\nHere is the link: `" + event.EventURL + "`\n"
	}

	return slack.Attachment{
		Pretext:  "The Incident <#" + inc.ChannelID + "> has been resolved by you",
		Fallback: privateText.String(),
		Text:     "",
		Color:    "#1164A3",
		Fields: []slack.AttachmentField{
			{
				Title: "Post Mortem",
				Value: postMortemMessage,
			},
		},
	}
}
