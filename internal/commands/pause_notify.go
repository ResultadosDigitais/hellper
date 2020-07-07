package commands

import (
	"context"
	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"

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

	dialog := slack.Dialog{
		CallbackID:     "inc-pausenotify",
		Title:          "Pause Notify",
		SubmitLabel:    "Pause",
		NotifyOnCancel: false,
		Elements:       []slack.DialogElement{pauseNotifyTime},
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
	logger.Info(
		ctx,
		"command/pauseNotify.PauseNotifyIncidentByDialog",
		log.NewValue("incidentDetails", incidentDetails),
		log.NewValue("repository", repository),
		log.NewValue("client", client),
	)

	PostInfoAttachment(ctx, client, incidentDetails.Channel.ID, incidentDetails.User.ID, "Ops", "This command is not ready yet")
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
