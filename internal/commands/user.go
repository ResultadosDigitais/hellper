package commands

import (
	"context"

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
		SlackId: slackUser.ID,
		Name:    slackUser.Profile.RealName,
		Email:   slackUser.Profile.Email,
	}

	return &user, err
}
