package meeting

import (
	"hellper/internal/config"
)

const MeetingProviderZoom = "zoom"
const MeetingProviderMatrix = "matrix"
const MeetingProviderNone = "none"

// Provider is an interface that defines how a meeting will be created
type Provider interface {
	CreateMeeting() (string, error)
}

// CreateMeeting creates a meeting and returns its url based on Hellper configs
func CreateMeeting(options map[string]string) (string, error) {
	provider := getMeetingProvider(options)

	if provider == nil {
		return "", nil
	}

	return provider.CreateMeeting()
}

func getMeetingProvider(additionalConfig map[string]string) Provider {
	var (
		providerName   = config.Env.MeetingConfig.ProviderName
		providerConfig = config.Env.MeetingConfig.ProviderConfig
	)

	if providerName == MeetingProviderZoom {
		return getZoomMeetingProvider(providerConfig, additionalConfig)
	} else if providerName == MeetingProviderMatrix {
		return getMatrixMeetingProvider(providerConfig, additionalConfig)
	}

	return nil
}
