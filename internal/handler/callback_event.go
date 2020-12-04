package handler

import (
	"context"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/log"

	"github.com/slack-go/slack/slackevents"
)

func replyCallbackEvent(
	ctx context.Context, app *app.App, event slackevents.EventsAPIEvent,
) error {
	var (
		innerEvent = event.InnerEvent

		err     error
		cmdLine string
		trigger commands.TriggerEvent
	)

	switch callbackEvent := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		app.Logger.Debug(
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
		app.Logger.Info(
			ctx,
			"handler/event.message",
			log.NewValue("callbackEvent", callbackEvent),
		)
		return nil
	case *slackevents.AppUninstalledEvent:
		app.Logger.Debug(
			ctx,
			"handler/event.appunistalled",
			log.NewValue("callbackEvent", callbackEvent),
		)
		return nil
	case *slackevents.LinkSharedEvent:
		app.Logger.Debug(
			ctx,
			"handler/event.linkshared",
			log.NewValue("callbackEvent", callbackEvent),
		)
		return nil
	default:
		app.Logger.Debug(
			ctx,
			"handler/event.unkown_event",
			log.NewValue("callbackEvent", callbackEvent),
		)
		return nil
	}

	executor := commands.NewEventExecutor(app)
	err = executor.ExecuteEventCommand(ctx, cmdLine, trigger)
	if err != nil {
		return err
	}
	return err
}
