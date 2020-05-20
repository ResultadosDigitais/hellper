package googleauth

import (
	"context"
	"encoding/json"
	"hellper/internal/config"
	"hellper/internal/log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleAuthStruct struct{}

// Interface interfaces the public methods from package
type Interface interface {
	GetGClient(context.Context, log.Logger, []byte, string) (*http.Client, error)
}

var (
	//Struct creates the interface for the usage of googleauth package
	Struct Interface = &googleAuthStruct{}
)

// GetGClient generates a google Client, given a token and a scope
func (gs *googleAuthStruct) GetGClient(ctx context.Context, logger log.Logger, token []byte, scope string) (*http.Client, error) {
	googleCredentialBytes := []byte(config.Env.GoogleCredentials)

	gConfig, err := google.ConfigFromJSON(googleCredentialBytes, scope)
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
