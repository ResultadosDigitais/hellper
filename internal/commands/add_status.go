package commands

import (
	"context"
	"hellper/internal/app"
	"hellper/internal/log"
)

func AddStatus(ctx context.Context, app *app.App, channelID string, userID string, userName string, message string) error {
	loggerWritter := app.Logger.With(
		log.NewValue("channelID", channelID),
		log.NewValue("userName", userName),
		log.NewValue("userID", userID),
	)

	loggerWritter.Debug(
		ctx,
		log.Trace(),
		log.Action("AddStatus"),
		log.NewValue("message", message),
	)

	formattedMessage := "<@" + userName + ">: " + message

	msgRef, err := postMessage(app, channelID, formattedMessage)
	if err != nil {
		return err
	}

	pinMessage(app, channelID, *msgRef)
	return nil
}
