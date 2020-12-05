package invitation

import (
	"context"
	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"
	"sync"
)

// Inviter is responsible for inviting stakeholders to a slack channel
type Inviter struct {
	logger           log.Logger
	client           bot.Client
	teamRepository   model.TeamRepository
	personRepository model.PersonRepository
}

// NewInviter creates a new instance of an inviter
func NewInviter(
	logger log.Logger, client bot.Client,
	teamRepository model.TeamRepository, personRepository model.PersonRepository,
) Inviter {
	return Inviter{
		logger:           logger,
		client:           client,
		teamRepository:   teamRepository,
		personRepository: personRepository,
	}
}

// CreateStrategy returns the strategy by its name. If no strategy found, an error is returned
func (i *Inviter) CreateStrategy(strategyName string) (Strategy, error) {
	return newStrategy(strategyName, i.logger)
}

// InviteStakeholders executes an strategy of invitation and adds the stakeholders to the slack channel
func (i *Inviter) InviteStakeholders(ctx context.Context, incident model.Incident, strategy Strategy) error {
	serviceInstance, err := i.getServiceInstanceFromIncident(ctx, incident)
	if err != nil {
		return err
	}

	stakeholders, err := strategy.GetStakeholders(ctx, serviceInstance, incident, i.teamRepository)
	if err != nil {
		i.logger.Error(
			ctx,
			"Could not load the stakeholders list",
			log.Action("InviteStakeholders"),
			log.NewValue("slackChannel", incident.ChannelName),
			log.NewValue("error", err),
		)

		return err
	}

	go func(ctx context.Context) {
		i.populateStakeholdersSlackID(ctx, stakeholders, incident)
		i.inviteStakeholdersToChannel(ctx, stakeholders, incident)
	}(context.Background())

	return nil
}

func (i *Inviter) getServiceInstanceFromIncident(ctx context.Context, incident model.Incident) (model.ServiceInstance, error) {
	return model.ServiceInstance{ID: incident.Product}, nil
}

func (i *Inviter) populateStakeholdersSlackID(ctx context.Context, stakeholders []*stakeholder, incident model.Incident) {
	var wg sync.WaitGroup
	for _, stakeholder := range stakeholders {
		wg.Add(1)
		go i.ensureStakeholderHasSlackID(ctx, &wg, stakeholder, incident)
	}

	wg.Wait()
}

func (i *Inviter) ensureStakeholderHasSlackID(
	ctx context.Context, wg *sync.WaitGroup, stakeholder *stakeholder, incident model.Incident,
) error {
	defer wg.Done()

	logWriter := i.logger.With(
		log.Action("ensureStakeholderHasSlackID"),
		log.NewValue("stakeholderEmail", stakeholder.email),
	)

	logWriter.Debug(
		ctx,
		"Inviting stakeholder to channel",
	)
	err := populateSlackIDIfEmpty(ctx, stakeholder, i.client, i.personRepository)

	if err != nil {
		logWriter.Error(
			ctx,
			"Could not get stakeholder slack id",
			log.NewValue("error", err),
		)
		return err
	}

	return nil
}

func (i *Inviter) inviteStakeholdersToChannel(
	ctx context.Context, stakeholders []*stakeholder, incident model.Incident,
) error {
	stakeholdersSlackIds := make([]string, 0, len(stakeholders))
	for _, stakeholder := range stakeholders {
		if stakeholder.slackID != "" {
			stakeholdersSlackIds = append(stakeholdersSlackIds, stakeholder.slackID)
		}
	}

	_, err := i.client.InviteUsersToConversationContext(ctx, incident.ChannelID, stakeholdersSlackIds...)

	if err != nil {
		i.logger.Error(
			ctx,
			"Could not invite stakeholders to slack channel",
			log.Action("inviteStakeholdersToChannel"),
			log.NewValue("slackChannel", incident.ChannelName),
			log.NewValue("error", err),
		)
	}

	return err

}
