package commands

import (
	"context"
	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

// UpdateDatesDialog opens a dialog on Slack, so the user can update the dates of an incident
func UpdateDatesDialog(ctx context.Context, logger log.Logger, client bot.Client, repository model.Repository, channelID string, userID string, triggerID string) error {
	var (
		dateLayout          = "02/01/2006 15:04:05"
		initValue           = ""
		identificationValue = ""
		endValue            = ""
	)

	inc, err := repository.GetIncident(ctx, channelID)
	if err != nil {
		logger.Error(
			ctx,
			"command/dates.UpdateDatesDialog GetIncident ERROR",
			log.NewValue("channelID", channelID),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	logger.Info(
		ctx,
		"command/close.UpdateDatesDialog INFO",
		log.NewValue("channelID", channelID),
		log.NewValue("startDate", inc.StartTimestamp),
		log.NewValue("identificationDate", inc.IdentificationTimestamp),
		log.NewValue("endDate", inc.EndTimestamp),
	)

	if inc.StartTimestamp != nil {
		initValue = inc.StartTimestamp.Format(dateLayout)
	}

	if inc.IdentificationTimestamp != nil {
		identificationValue = inc.IdentificationTimestamp.Format(dateLayout)
	}

	if inc.EndTimestamp != nil {
		endValue = inc.EndTimestamp.Format(dateLayout)
	}

	timeZone := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Time Zone",
			Name:        "time_zone",
			Type:        "select",
			Placeholder: "Choose your time zone",
			Optional:    false,
		},
		Value: "0",
		Options: []slack.DialogSelectOption{
			{
				Label: "UTC",
				Value: "0",
			},
			{
				Label: "UTC-2h",
				Value: "-2",
			},
			{
				Label: "UTC-3h",
				Value: "-3",
			},
		},
		OptionGroups: []slack.DialogOptionGroup{},
	}
	startDate := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Start date (" + dateLayout + ")",
			Name:        "init_date",
			Type:        "text",
			Placeholder: dateLayout,
			Optional:    false,
		},
		Value: initValue,
	}
	identificationDate := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Identification date (" + dateLayout + ")",
			Name:        "identification_date",
			Type:        "text",
			Placeholder: dateLayout,
			Optional:    false,
		},
		Value: identificationValue,
	}
	endDate := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "End date (" + dateLayout + ")",
			Name:        "end_date",
			Type:        "text",
			Placeholder: dateLayout,
			Optional:    false,
		},
		Value: endValue,
	}

	dialog := slack.Dialog{
		CallbackID:     "inc-dates",
		Title:          "Update Incident's dates",
		SubmitLabel:    "Update",
		NotifyOnCancel: false,
		Elements: []slack.DialogElement{
			timeZone,
			startDate,
			identificationDate,
			endDate,
		},
	}

	return client.OpenDialog(triggerID, dialog)
}

// UpdateDatesByDialog updates the dates of an incident after receiving data from a Slack dialog
func UpdateDatesByDialog(ctx context.Context, client bot.Client, logger log.Logger, repository model.Repository, incidentDetails bot.DialogSubmission) error {
	logger.Info(
		ctx,
		"command/close.CloseIncidentByDialog INFO",
		log.NewValue("incident_close_details", incidentDetails),
	)

	var (
		dateLayout             = "02/01/2006 15:04:05"
		channelID              = incidentDetails.Channel.ID
		userID                 = incidentDetails.User.ID
		userName               = incidentDetails.User.Name
		submissions            = incidentDetails.Submission
		timeZoneString         = submissions.TimeZone
		initDateText           = submissions.InitDate
		identificationDateText = submissions.IdentificationDate
		endDateText            = submissions.EndDate

		initDate           time.Time
		identificationDate time.Time
		endDate            time.Time
	)

	location, err := parseTimeZone(timeZoneString)
	if err != nil {
		logger.Error(
			ctx,
			"command/dates.UpdateDatesByDialog parseTimeZone ERROR",
			log.NewValue("channelID", channelID),
			log.NewValue("timeZoneString", timeZoneString),
			log.NewValue("initDateText", initDateText),
			log.NewValue("identificationDateText", identificationDateText),
			log.NewValue("endDateText", endDateText),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	initDate, err = time.ParseInLocation(dateLayout, initDateText, location)
	if err != nil {
		logger.Error(
			ctx,
			"command/dates.UpdateDatesByDialog ParseIn ERROR",
			log.NewValue("channelID", channelID),
			log.NewValue("timeZoneString", timeZoneString),
			log.NewValue("initDateText", initDateText),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	identificationDate, err = time.ParseInLocation(dateLayout, identificationDateText, location)
	if err != nil {
		logger.Error(
			ctx,
			"command/dates.UpdateDatesByDialog ParseIn ERROR",
			log.NewValue("channelID", channelID),
			log.NewValue("timeZoneString", timeZoneString),
			log.NewValue("identificationDateText", identificationDateText),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	endDate, err = time.ParseInLocation(dateLayout, endDateText, location)
	if err != nil {
		logger.Error(
			ctx,
			"command/dates.UpdateDatesByDialog ParseIn ERROR",
			log.NewValue("channelID", channelID),
			log.NewValue("timeZoneString", timeZoneString),
			log.NewValue("endDateText", endDateText),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	incident := model.Incident{
		ChannelId:               channelID,
		StartTimestamp:          &initDate,
		IdentificationTimestamp: &identificationDate,
		EndTimestamp:            &endDate,
	}

	err = repository.UpdateIncidentDates(ctx, &incident)
	if err != nil {
		logger.Error(
			ctx,
			"command/dates.UpdateDatesByDialog UpdateIncidentDates ERROR",
			log.NewValue("incident", incident),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	successAttach := createDatesSuccessAttachment(incident, userName)
	postMessage(client, incident.ChannelId, "", successAttach)

	return nil
}

func parseTimeZone(timeZoneString string) (*time.Location, error) {
	if timeZoneString == "0" {
		return time.UTC, nil
	}

	timeZoneInt, err := strconv.Atoi(timeZoneString)
	if err != nil {
		return nil, err
	}

	loc := time.FixedZone("Custom/Location", timeZoneInt*60*60)
	return loc, nil
}

func createDatesSuccessAttachment(inc model.Incident, userName string) slack.Attachment {
	var (
		dateLayout  = time.RFC1123
		messageText strings.Builder
	)

	messageText.WriteString("The dates of Incident <#" + inc.ChannelId + "> has been updated by <@" + userName + ">\n\n")

	return slack.Attachment{
		Pretext:  "The dates of Incident <#" + inc.ChannelId + "> has been updated by <@" + userName + ">",
		Fallback: messageText.String(),
		Text:     "",
		Color:    "#6fff47",
		Fields: []slack.AttachmentField{
			{
				Title: "Start Date:",
				Value: inc.StartTimestamp.Format(dateLayout),
			},
			{
				Title: "Identification Date:",
				Value: inc.IdentificationTimestamp.Format(dateLayout),
			},
			{
				Title: "End Date:",
				Value: inc.EndTimestamp.Format(dateLayout),
			},
		},
	}
}
