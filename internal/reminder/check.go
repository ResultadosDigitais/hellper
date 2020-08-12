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

	return checkRules(rules)
}

func checkRules(rules notifyRules) bool {
	switch rules.status {
	case model.StatusOpen:
		return rulesInOpenStatus(rules)
	case model.StatusResolved:
		return rulesInResolvedStatus(rules)
	default:
		return false
	}
}

func rulesInOpenStatus(rules notifyRules) bool {
	if rules.snoozedUntil {
		return false
	}

	if rules.lastPin {
		return false
	}

	return true
}

func rulesInResolvedStatus(rules notifyRules) bool {
	if rules.snoozedUntil {
		return false
	}

	if rules.slaClose {
		return false
	}

	return true
}
