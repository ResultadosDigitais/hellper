package invitation

import "hellper/internal/model"

type stakeholder struct {
	slackID string
	email   string
}

func createStakeholdersFromUsers(users []*model.User) []*stakeholder {
	stakeholders := make([]*stakeholder, 0, len(users))
	for _, user := range users {
		stakeholders = append(stakeholders, createStakeholderFromUser(user))
	}

	return stakeholders
}

func createStakeholderFromUser(user *model.User) *stakeholder {
	return &stakeholder{
		slackID: user.SlackID,
		email:   user.Email,
	}
}
