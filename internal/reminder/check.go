package reminder

import (
	"context"
	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"
)

type notifyRules struct {
	snoozedUntil bool
	lastPin      bool
	slaClose     bool
	status       string
}

// CanSendNotify checks the notification rules
func CanSendNotify(ctx context.Context, client bot.Client, logger log.Logger, repository model.Repository, incident model.Incident) bool {
	logger.Info(
		ctx,
		log.Trace(),
		log.Action("running"),
		log.NewValue("channelID", incident.ChannelId),
		log.NewValue("channelName", incident.ChannelName),
	)

	rules := notifyRules{
		snoozedUntil: hasSnoozedUntil(ctx, logger, incident),
		lastPin:      hasLastPin(ctx, client, logger, incident),
		slaClose:     hasSLAClose(ctx, client, logger, incident),
		status:       incident.Status,
	}

	return rules.checkRules()
}

func (rules notifyRules) checkRules() bool {
	switch rules.status {
	case model.StatusOpen:
		return rules.rulesInOpenStatus()
	case model.StatusResolved:
		return rules.rulesInResolvedStatus()
	default:
		return false
	}
}

func (rules notifyRules) rulesInOpenStatus() bool {
	if rules.snoozedUntil {
		return false
	}

	if rules.lastPin {
		return false
	}

	return true
}

func (rules notifyRules) rulesInResolvedStatus() bool {
	if rules.snoozedUntil {
		return false
	}

	if rules.slaClose {
		return false
	}

	return true
}
