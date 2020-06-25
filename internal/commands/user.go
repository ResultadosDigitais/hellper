package commands

import (
	"context"

	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
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

func getUsersInConversation(
	ctx context.Context,
	client bot.Client,
	logger log.Logger,
	channelID string,
) (*[]model.User, error) {
	logger.Info(
		ctx,
		"command/user.getUsersInConversation",
		log.NewValue("params", channelID),
	)

	params := slack.GetUsersInConversationParameters{
		ChannelID: channelID,
	}

	var (
		membersID []string
		users     []model.User
	)

	for {
		list, cursor, err := client.GetUsersInConversationContext(ctx, &params)
		if err != nil {
			logger.Error(
				ctx,
				"command/user.getUsersInConversation",
				log.NewValue("channelID", channelID),
				log.NewValue("error", err),
			)
			return nil, err
		}
		membersID = append(membersID, list...)
		if cursor == "" {
			break
		} else {
			params.Cursor = cursor
		}
	}

	return &users, nil
}
