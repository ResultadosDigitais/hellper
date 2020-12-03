package meeting

import (
	"fmt"
)

type matrixProvider struct {
	matrixURL   string
	channel     string
	environment string
}

func getMatrixMeetingProvider(config map[string]string, additionalConfig map[string]string) matrixProvider {
	return matrixProvider{
		matrixURL:   config["matrixHost"],
		channel:     additionalConfig["channel"],
		environment: additionalConfig["environment"],
	}
}

func (provider matrixProvider) CreateURL() (string, error) {
	var (
		url         = provider.matrixURL
		channelName = provider.channel
		environment = provider.environment
	)

	roomID := channelName
	roomName := channelName

	if environment == "staging" {
		roomID = "dc82e346-639c-44ee-a470-63f7545ae8e4"
		roomName = "hellper-staging"
	}

	return fmt.Sprintf("%s/new?roomId=%s&roomName=%s", url, roomID, roomName), nil
}
