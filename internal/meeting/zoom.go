package meeting

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type zoomProvider struct {
	jwtToken    string
	userID      string
	channel     string
	environment string
}

func getZoomMeetingProvider(config, additionalConfig map[string]string) zoomProvider {
	return zoomProvider{
		jwtToken:    config["jwtToken"],
		userID:      config["userId"],
		channel:     additionalConfig["channel"],
		environment: additionalConfig["environment"],
	}
}

func (provider zoomProvider) CreateMeeting() (string, error) {
	var (
		apiBaseURL = "https://api.zoom.us/v2"
		userID     = provider.userID
		channel    = provider.channel
	)

	url := fmt.Sprintf("%s/users/%s/meetings", apiBaseURL, userID)
	postData := provider.createMeetingInput(channel)

	req, err := http.NewRequest("POST", url, postData)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", provider.jwtToken))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return provider.getMeetingURLFromPayload(body)
}

func (provider zoomProvider) getMeetingURLFromPayload(payload []byte) (string, error) {
	type ZoomPayload struct {
		JoinURL string `json:"join_url"`
	}

	var data ZoomPayload

	if err := json.Unmarshal(payload, &data); err != nil {
		return "", err
	}

	if data.JoinURL == "" {
		return "", fmt.Errorf("couldn't retrieve a meeting url, please check you Zoom API settings. payload: %s", string(payload))
	}

	return data.JoinURL, nil
}
