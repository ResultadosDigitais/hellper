package googleauth

import (
	"context"
	"encoding/json"
	"hellper/internal/config"
	googleapi "hellper/internal/google/google_api"
	"hellper/internal/log"
	"net/http"

	"golang.org/x/oauth2"
)

type googleAuthStruct struct{}

type GoogleAuthInterface interface {
	GetGClient(context.Context, log.Logger, []byte, string) (*http.Client, error)
}

var (
	GoogleAuthStruct GoogleAuthInterface = &googleAuthStruct{}
)

// GetGClient generates a google Client, given a token and a scope
func (gs *googleAuthStruct) GetGClient(ctx context.Context, logger log.Logger, token []byte, scope string) (*http.Client, error) {
	driveCredentialBytes := []byte(config.Env.GoogleDriveCredentials)

	gConfig, err := googleapi.GoogleStruct.ConfigFromJSON(driveCredentialBytes, scope)
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
