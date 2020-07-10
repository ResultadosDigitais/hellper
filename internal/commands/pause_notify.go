package commands

import (
	"context"
	"database/sql"
	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

// PauseNotifyIncidentDialog opens a dialog on Slack, so the user can pause notify
func PauseNotifyIncidentDialog(ctx context.Context, logger log.Logger, client bot.Client, repository model.Repository, channelID string, userID string, triggerID string) error {

	inc, err := repository.GetIncident(ctx, channelID)
	if err != nil {
		logger.Error(
			ctx,
			"command/pauseNotify.PauseNotifyIncidentDialog GetIncident ERROR",
			log.NewValue("channelID", channelID),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	if inc.Status == model.StatusClosed || inc.Status == model.StatusCancel {
		PostInfoAttachment(ctx, client, channelID, userID, "Ops! That's not possible", "The incident status is: "+inc.Status)
		return nil
	}

	pauseNotifyTime := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "How long time would you like to pause?",
			Name:        "pause_notify_time",
			Type:        "select",
			Placeholder: "Select an option",
			Optional:    false,
		},
		Value:        "1",
		Options:      optionsPauseNotify(inc.Status),
		OptionGroups: []slack.DialogOptionGroup{},
	}

	reason := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Reason",
			Name:        "pause_notify_reason",
			Type:        "textarea",
			Placeholder: "Reason",
			Optional:    false,
		},
		MaxLength: 500,
	}

	dialog := slack.Dialog{
		CallbackID:     "inc-pausenotify",
		Title:          "Pause Notify",
		SubmitLabel:    "Pause",
		NotifyOnCancel: false,
		Elements:       []slack.DialogElement{pauseNotifyTime, reason},
	}

	return client.OpenDialog(triggerID, dialog)
}

// PauseNotifyIncidentByDialog Pause a notify from a Slack dialog
func PauseNotifyIncidentByDialog(
	ctx context.Context,
	client bot.Client,
	logger log.Logger,
	repository model.Repository,
	incidentDetails bot.DialogSubmission,
) error {

	var (
		channelID             = incidentDetails.Channel.ID
		userID                = incidentDetails.User.ID
		submissions           = incidentDetails.Submission
		pauseNotifyTimeText   = submissions.PauseNotifyTime
		pauseNotifyReasonText = submissions.PauseNotifyReason

		pauseNotifyTime sql.NullTime
	)

	days, err := strconv.Atoi(pauseNotifyTimeText)
	if err != nil {
		logger.Error(
			ctx,
			"command/pauseNotify.PauseNotifyIncidentByDialog strconv.Atoi ERROR",
			log.NewValue("channelID", channelID),
			log.NewValue("pauseNotifyTimeText", pauseNotifyTimeText),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	pauseNotifyTime.Time = time.Now().AddDate(0, 0, days)

	logger.Info(
		ctx,
		"command/pauseNotify.PauseNotifyIncidentByDialog",
		log.NewValue("pauseNotifyTimeText", pauseNotifyTimeText),
		log.NewValue("pauseNotifyReasonText", pauseNotifyReasonText),
		log.NewValue("pauseNotifyTime", pauseNotifyTime),
	)

	incident := model.Incident{
		ChannelId: channelID,
		SnoozedAt: pauseNotifyTime,
	}

	err = repository.PauseNotifyIncident(ctx, &incident)
	if err != nil {
		logger.Error(
			ctx,
			"command/pauseNotify.PauseNotifyIncidentByDialog PauseNotifyIncident ERROR",
			log.NewValue("incident", incident),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	PostInfoAttachment(ctx, client, channelID, userID, "Ops", "This command is not ready yet")
	return nil
}

func optionsPauseNotify(status string) (option []slack.DialogSelectOption) {
	switch status {
	case model.StatusOpen:
		option = []slack.DialogSelectOption{
			{Label: "1 day", Value: "1"},
		}
	case model.StatusResolved:
		option = []slack.DialogSelectOption{
			{Label: "1 day", Value: "1"},
			{Label: "2 days", Value: "2"},
			{Label: "3 days", Value: "3"},
		}
	default:
		option = []slack.DialogSelectOption{
			{Label: "1 day", Value: "1"},
		}
	}

	return
}
