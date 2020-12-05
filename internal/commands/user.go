package commands

import (
	"context"

	"hellper/internal/app"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

func getSlackUserInfo(
	ctx context.Context,
	app *app.App,
	userID string,
) (*model.User, error) {

	slackUser, err := app.Client.GetUserInfoContext(ctx, userID)
	if err != nil {
		app.Logger.Error(
			ctx,
			"command/user.getSlackUserInfo error",
			log.NewValue("userID", userID),
			log.NewValue("error", err),
		)

		return nil, err
	}

	user := model.User{
		SlackID:     slackUser.ID,
		Name:        slackUser.Profile.RealName,
		DisplayName: slackUser.Profile.DisplayName,
		Email:       slackUser.Profile.Email,
	}

	return &user, err
}

func getUsersIDsInConversation(
	ctx context.Context,
	app *app.App,
	channelID string,
) (*[]string, error) {

	logWriter := app.Logger.With(
		log.NewValue("channelID", channelID),
	)

	logWriter.Debug(
		ctx,
		"command/user.getUsersIDsInConversation",
	)

	params := slack.GetUsersInConversationParameters{
		ChannelID: channelID,
	}

	var members []string

	for {
		list, cursor, err := app.Client.GetUsersInConversationContext(ctx, &params)
		if err != nil {
			logWriter.Error(
				ctx,
				"command/user.getUsersIDsInConversation GetUsersInConversationContext",
				log.NewValue("error", err),
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
	app *app.App,
	channelID string,
) (*[]model.User, error) {
	var users []model.User

	logWriter := app.Logger.With(
		log.NewValue("channelID", channelID),
	)

	logWriter.Debug(
		ctx,
		"command/user.getUsersInConversation",
	)

	usersIDs, err := getUsersIDsInConversation(ctx, app, channelID)
	if err != nil {
		logWriter.Error(
			ctx,
			"command/user.getUsersInConversation",
			log.Action("getUsersDIsInConversation"),
			log.NewValue("error", err),
		)
		return nil, err
	}

	for _, id := range *usersIDs {
		user, err := getSlackUserInfo(ctx, app, id)
		if err != nil {
			logWriter.Error(
				ctx,
				"command/user.getUsersInConversation",
				log.Action("getSlackUserInfo"),
				log.NewValue("user_id", id),
				log.NewValue("error", err),
			)
			return nil, err
		}
		users = append(users, *user)
	}

	return &users, err
}

func getUsersEmailsInConversation(
	ctx context.Context,
	app *app.App,
	channelID string,
) (*[]string, error) {
	var emails []string

	logWriter := app.Logger.With(
		log.NewValue("channelID", channelID),
	)

	logWriter.Debug(
		ctx,
		"command/user.getUsersEmailsInConversation",
	)

	users, err := getUsersInConversation(ctx, app, channelID)
	if err != nil {
		logWriter.Error(
			ctx,
			"command/user.getUsersEmailsInConversation",
			log.Action("getUsersInConversation"),
			log.NewValue("error", err),
		)
		return nil, err
	}

	for _, user := range *users {
		emails = append(emails, user.Email)
	}

	return &emails, err
}
