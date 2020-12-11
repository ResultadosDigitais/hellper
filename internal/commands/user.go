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
			log.Trace(),
			log.Action("client.GetUserInfoContext"),
			log.Reason(err.Error()),
			log.NewValue("userID", userID),
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

func getUsersIDsInConversation(
	ctx context.Context,
	client bot.Client,
	logger log.Logger,
	channelID string,
) (*[]string, error) {
	logger.Info(
		ctx,
		log.Trace(),
		log.NewValue("params", channelID),
	)

	params := slack.GetUsersInConversationParameters{
		ChannelID: channelID,
	}

	var members []string

	for {
		list, cursor, err := client.GetUsersInConversationContext(ctx, &params)
		if err != nil {
			logger.Error(
				ctx,
				log.Trace(),
				log.Action("client.GetUsersInConversationContext"),
				log.Reason(err.Error()),
				log.NewValue("channelID", channelID),
			)
			return nil, err
		}
		members = append(members, list...)
		if cursor == "" {
			break
		} else {
			params.Cursor = cursor
		}
	}

	return &members, nil
}

func getUsersInConversation(
	ctx context.Context,
	client bot.Client,
	logger log.Logger,
	channelID string,
) (*[]model.User, error) {
	var users []model.User

	usersIDs, err := getUsersIDsInConversation(ctx, client, logger, channelID)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Action("getUsersIDsInConversation"),
			log.Reason(err.Error()),
			log.NewValue("channel_id", channelID),
		)
		return nil, err
	}

	for _, id := range *usersIDs {
		user, err := getSlackUserInfo(ctx, client, logger, id)
		if err != nil {
			logger.Error(
				ctx,
				log.Trace(),
				log.Action("getSlackUserInfo"),
				log.Reason(err.Error()),
				log.NewValue("user_id", id),
			)
			return nil, err
		}
		users = append(users, *user)
	}

	return &users, err
}

func getUsersEmailsInConversation(
	ctx context.Context,
	client bot.Client,
	logger log.Logger,
	channelID string,
) (*[]string, error) {
	var emails []string

	users, err := getUsersInConversation(ctx, client, logger, channelID)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Action("getUsersInConversation"),
			log.Reason(err.Error()),
			log.NewValue("channel_id", channelID),
		)
		return nil, err
	}

	for _, user := range *users {
		emails = append(emails, user.Email)
	}

	return &emails, err
}
