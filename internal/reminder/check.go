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
func CanSendNotify(ctx context.Context, client bot.Client, logger log.Logger, repository model.IncidentRepository, incident model.Incident) bool {
	logWriter := logger.With(
		log.NewValue("channelID", incident.ChannelId),
		log.NewValue("channelName", incident.ChannelName),
	)

	logWriter.Debug(
		ctx,
		log.Trace(),
		log.Action("running"),
	)

	rules := notifyRules{
		snoozedUntil: hasSnoozedUntil(ctx, logWriter, incident),
		lastPin:      hasLastPin(ctx, client, logWriter, incident),
		slaClose:     hasSLAClose(ctx, client, logWriter, incident),
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
