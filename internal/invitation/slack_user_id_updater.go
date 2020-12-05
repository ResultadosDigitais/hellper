package invitation

import (
	"context"
	"hellper/internal/bot"
	"hellper/internal/model"
)

func populateSlackIDIfEmpty(
	ctx context.Context, stakeholder *stakeholder, client bot.Client, personRepository model.PersonRepository,
) error {
	if stakeholder.slackID != "" {
		return nil
	}

	user, err := client.GetUserByEmailContext(ctx, stakeholder.email)
	if err != nil {
		return err
	}

	stakeholder.slackID = user.ID

	stakeholderUser := model.User{Email: stakeholder.email, SlackID: stakeholder.slackID}
	return personRepository.UpdatePersonSlackID(&stakeholderUser)
}
