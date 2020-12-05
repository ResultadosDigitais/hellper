package commands

import (
	"context"
	"strconv"
	"strings"

	"hellper/internal/app"
	"hellper/internal/log"

	"github.com/slack-go/slack"
)

func createPostMortem(
	ctx context.Context,
	app *app.App,
	incidentID int64,
	incidentName string,
	channelName string,
) (string, error) {

	postMortemName := strconv.FormatInt(incidentID, 10) + " - PostMortem - " + incidentName
	postMortemURL, err := app.FileStorage.CreatePostMortemDocument(ctx, postMortemName)
	if err != nil {
		app.Logger.Error(
			ctx,
			"command/open.create_post_mortem_document ERROR",
			log.NewValue("incident_id", incidentID),
			log.NewValue("incident_name", incidentName),
			log.NewValue("channel_name", channelName),
			log.NewValue("error", err),
		)
		return "", err
	}
	addPostMortemURLToDB(ctx, app, channelName, postMortemURL)

	var messageText strings.Builder
	messageText.WriteString("*Post Mortem URL:* " + postMortemURL + "\n")

	attachment := slack.Attachment{
		Pretext:  "",
		Fallback: messageText.String(),
		Text:     "",
		Color:    "#FE4D4D",
		Fields: []slack.AttachmentField{
			{
				Title: "Post Mortem",
				Value: postMortemURL,
			},
		},
	}

	postAndPinMessage(app, channelName, "Post Mortem document created", attachment)
	return postMortemURL, nil
}

func addPostMortemURLToDB(ctx context.Context, app *app.App, channelName string, postMortemURL string) {
	err := app.IncidentRepository.AddPostMortemURL(ctx, channelName, postMortemURL)
	if err != nil {
		app.Logger.Debug(
			ctx,
			"Post Mortem could not be inserted to DB",
			log.NewValue("channelName", channelName),
			log.NewValue("postMortemURL", postMortemURL),
		)
	}
}
