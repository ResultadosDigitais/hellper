package invitation

import (
	"context"
	"hellper/internal/log"
	"hellper/internal/model"
)

type baseStrategy struct {
	logger log.Logger
}

func (s *baseStrategy) getOwnerTeamMembers(
	ctx context.Context, serviceInstance model.ServiceInstance, teamRepository model.TeamRepository,
) ([]*stakeholder, error) {
	ownerUsers, err := teamRepository.GetOwnersByServiceInstance(ctx, serviceInstance)
	if err != nil {
		return []*stakeholder{}, err
	}

	teamOwners := s.convertUsersToStakeholders(ownerUsers)

	return teamOwners, nil
}

func (s *baseStrategy) getUserTeamsMembers(
	ctx context.Context, serviceInstance model.ServiceInstance, teamRepository model.TeamRepository,
) ([]*stakeholder, error) {
	users, err := teamRepository.GetUsersOfServiceInstance(ctx, serviceInstance)
	if err != nil {
		return []*stakeholder{}, err
	}

	stakeholdersUsers := s.convertUsersToStakeholders(users)

	return stakeholdersUsers, nil
}

func (s *baseStrategy) convertUsersToStakeholders(users []*model.User) []*stakeholder {
	stakeholders := make([]*stakeholder, 0, len(users))
	for _, user := range users {
		stakeholder := stakeholder{email: user.Email, slackID: user.SlackID}
		stakeholders = append(stakeholders, &stakeholder)
	}

	return stakeholders
}

func (s *baseStrategy) removeDuplicatedStakeholders(stakeholders []*stakeholder) []*stakeholder {
	stakeholdersAlreadyInserted := make(map[string]bool)
	nonRepeatedStakeholders := make([]*stakeholder, 0, len(stakeholders))

	for _, stakeholder := range stakeholders {
		if !stakeholdersAlreadyInserted[stakeholder.email] {
			nonRepeatedStakeholders = append(nonRepeatedStakeholders, stakeholder)
			stakeholdersAlreadyInserted[stakeholder.email] = true
		}
	}

	return nonRepeatedStakeholders
}
