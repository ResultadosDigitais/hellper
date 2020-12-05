package invitation

import (
	"context"
	"fmt"
	"hellper/internal/log"
	"hellper/internal/model"
)

// Strategy is an abstraction of who the bot should invite to the incident slack channel
type Strategy interface {
	GetStakeholders(
		ctx context.Context, serviceInstance model.ServiceInstance,
		incident model.Incident, teamRepository model.TeamRepository,
	) ([]*stakeholder, error)
}

// NewStrategy creates a strategy instance based on its name
func newStrategy(strategyName string, logger log.Logger) (Strategy, error) {
	switch strategyName {
	case "invite_all":
		return newInviteAllStrategy(logger), nil
	default:
		return nil, fmt.Errorf("Invalid invitation strategy: %s", strategyName)
	}
}
