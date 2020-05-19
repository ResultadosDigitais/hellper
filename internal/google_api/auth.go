package googleapi

import (
	"context"
	"encoding/json"
	"hellper/internal/config"
	"hellper/internal/log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GetGClient generates a google Client, given a token and a scope
func GetGClient(ctx context.Context, logger log.Logger, token []byte, scope string) (*http.Client, error) {
	driveCredentialBytes := []byte(config.Env.GoogleDriveCredentials)

	gConfig, err := google.ConfigFromJSON(driveCredentialBytes, scope)
	if err != nil {
		logger.Error(
			ctx,
			"googleApi/auth.GetGClient ConfigFromJSON error",
			log.NewValue("error", err),
		)

		return nil, err
	}

	googleToken := &oauth2.Token{}
	err = json.Unmarshal(token, googleToken)
	if err != nil {
		logger.Error(
			ctx,
			"googleApi/auth.getGClient Unmarshal error",
			log.NewValue("error", err),
		)

		return nil, err
	}

	gClient := gConfig.Client(ctx, googleToken)

	return gClient, nil
}
