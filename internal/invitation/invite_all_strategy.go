package invitation

import (
	"context"
	"fmt"
	"hellper/internal/log"
	"hellper/internal/model"
)

type inviteAllStrategy struct {
	baseStrategy
}

func newInviteAllStrategy(logger log.Logger) Strategy {
	return &inviteAllStrategy{
		baseStrategy{logger: logger},
	}
}

func (s *inviteAllStrategy) GetStakeholders(
	ctx context.Context, serviceInstance model.ServiceInstance, incident model.Incident, teamRepository model.TeamRepository,
) ([]*stakeholder, error) {
	allStakeholders := make([]*stakeholder, 0)

	commander := stakeholder{slackID: incident.CommanderID, email: incident.CommanderEmail}
	allStakeholders = append(allStakeholders, &commander)

	ownerTeamMembers, err := s.getOwnerTeamMembers(ctx, serviceInstance, teamRepository)
	if err != nil {
		return allStakeholders, err
	}
	allStakeholders = append(allStakeholders, ownerTeamMembers...)

	userTeamStakeholders, err := s.getUserTeamsMembers(ctx, serviceInstance, teamRepository)
	if err != nil {
		return allStakeholders, err
	}
	allStakeholders = append(allStakeholders, userTeamStakeholders...)
	nonDuplicatedStakeholders := s.removeDuplicatedStakeholders(allStakeholders)

	s.logger.Debug(
		ctx,
		fmt.Sprintf("%d stakeholders found for incident, %d are unique", len(allStakeholders), len(nonDuplicatedStakeholders)),
		log.Action("GetStakeholders"),
	)

	return nonDuplicatedStakeholders, nil
}
