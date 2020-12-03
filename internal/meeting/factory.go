package meeting

import (
	"hellper/internal/config"
)

// Provider is an interface that defines how a meetubg will be created
type Provider interface {
	CreateURL() (string, error)
}

// CreateMeeting creates a meeting and return its url based on Hellper configs
func CreateMeeting(options map[string]string) (string, error) {
	provider := getMeetingProvider(options)
	return provider.CreateURL()
}

func getMeetingProvider(additionalConfig map[string]string) Provider {
	var (
		providerName   = config.Env.MeetingConfig.ProviderName
		providerConfig = config.Env.MeetingConfig.ProviderConfig
	)

	if providerName == "zoom" {
		return getZoomMeetingProvider(providerConfig, additionalConfig)
	}

	return getMatrixMeetingProvider(providerConfig, additionalConfig)
}
