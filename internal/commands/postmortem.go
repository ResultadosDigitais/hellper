package commands

import (
	"context"
	"strconv"
	"strings"

	"hellper/internal/bot"
	filestorage "hellper/internal/file_storage"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

func createPostMortem(
	ctx context.Context,
	logger log.Logger,
	client bot.Client,
	fileStorage filestorage.Driver,
	incidentID int64,
	incidentName string,
	repository model.Repository,
	channelName string,
) (string, error) {

	postMortemName := strconv.FormatInt(incidentID, 10) + " - PostMortem - " + incidentName
	postMortemURL, err := fileStorage.CreatePostMortemDocument(ctx, postMortemName)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Action("fileStorage.CreatePostMortemDocument"),
			log.Reason(err.Error()),
			log.NewValue("incident_id", incidentID),
			log.NewValue("incident_name", incidentName),
			log.NewValue("channel_name", channelName),
		)
		return "", err
	}
	addPostMortemURLToDB(ctx, logger, repository, channelName, postMortemURL)

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

	postAndPinMessage(client, channelName, "Post Mortem document created", attachment)
	return postMortemURL, nil
}

func addPostMortemURLToDB(ctx context.Context, logger log.Logger, repository model.Repository, channelName string, postMortemURL string) {
	err := repository.AddPostMortemUrl(ctx, channelName, postMortemURL)
	if err != nil {
		logger.Info(
			ctx,
			log.Trace(),
			log.NewValue("channelName", channelName),
			log.NewValue("postMortemURL", postMortemURL),
		)
	}
}
