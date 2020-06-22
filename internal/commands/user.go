package commands

import (
	"context"

	"github.com/slack-go/slack"

	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"
)

func getSlackUserInfo(
	ctx context.Context,
	client bot.Client,
	logger log.Logger,
	userID string,
) (*model.User, error) {

	slackUser, err := client.GetUserInfoContext(ctx, userID)
	if err != nil {
		logger.Error(
			ctx,
			"command/user.getSlackUserInfo error",
			log.NewValue("userID", userID),
			log.NewValue("error", err),
		)

		return nil, err
	}

	user := model.User{
		SlackID: slackUser.ID,
		Name:    slackUser.Profile.RealName,
		Email:   slackUser.Profile.Email,
	}

	return &user, err
}

func getUsersInConversationParameters(
	ctx context.Context,
	client bot.Client,
	logger log.Logger,
) *slack.GetUsersInConversationParameters {

	return nil
}
