package handler

import (
	"context"

	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack/slackevents"
)

func replyCallbackEvent(
	ctx context.Context, logger log.Logger, client bot.Client, repository model.IncidentRepository, event slackevents.EventsAPIEvent,
) error {
	var (
		innerEvent = event.InnerEvent

		err     error
		cmdLine string
		trigger commands.TriggerEvent
	)

	switch callbackEvent := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		logger.Info(
			ctx,
			"handler/event.appmention",
			log.NewValue("callbackEvent", callbackEvent),
		)

		cmdLine = callbackEvent.Text
		trigger = commands.TriggerEvent{
			Type:    callbackEvent.Type,
			Channel: callbackEvent.Channel,
			User:    callbackEvent.User,
		}
	case *slackevents.MessageEvent:
		logger.Info(
			ctx,
			"handler/event.message",
			log.NewValue("callbackEvent", callbackEvent),
		)
		return nil
	case *slackevents.AppUninstalledEvent:
		logger.Info(
			ctx,
			"handler/event.appunistalled",
			log.NewValue("callbackEvent", callbackEvent),
		)
		return nil
	case *slackevents.LinkSharedEvent:
		logger.Info(
			ctx,
			"handler/event.linkshared",
			log.NewValue("callbackEvent", callbackEvent),
		)
		return nil
	default:
		logger.Info(
			ctx,
			"handler/event.unkown_event",
			log.NewValue("callbackEvent", callbackEvent),
		)
		return nil
	}

	executor := commands.NewEventExecutor(logger, client, repository)
	err = executor.ExecuteEventCommand(ctx, cmdLine, trigger)
	if err != nil {
		return err
	}
	return err
}
