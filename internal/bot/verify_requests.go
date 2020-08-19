package bot

import (
	"bytes"
	"context"
	"hellper/internal/config"
	"hellper/internal/log"
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack"
)

func VerifyRequests(r *http.Request, w http.ResponseWriter, f http.Handler) {
	var logger log.Logger

	secretVerifier, err := slack.NewSecretsVerifier(r.Header, config.Env.SlackSigningSecret)
	if err != nil {
		logger.Error(context.Background(), log.Trace(), log.Reason(err.Error()))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(context.Background(), log.Trace(), log.Reason(err.Error()))
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	_, err = secretVerifier.Write(body)
	if err != nil {
		logger.Error(context.Background(), log.Trace(), log.Reason(err.Error()))
	}

	err = secretVerifier.Ensure()
	if err != nil {
		logger.Error(context.Background(), log.Trace(), log.Reason(err.Error()))
	}

	f.ServeHTTP(w, r)
}
