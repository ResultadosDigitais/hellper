package bot

import (
	"bytes"
	"context"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/log/zap"
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack"
)

func VerifyRequests(r *http.Request, w http.ResponseWriter, f http.Handler) {
	var logger log.Logger = zap.NewDefault()

	secretVerifier, err := slack.NewSecretsVerifier(r.Header, config.Env.SlackSigningSecret)
	if err != nil {
		logger.Error(context.Background(), log.Trace(), log.Action("NewSecretsVerifier"), log.Reason(err.Error()))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(context.Background(), log.Trace(), log.Action("ReadAll"), log.Reason(err.Error()))
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	_, err = secretVerifier.Write(body)
	if err != nil {
		logger.Error(context.Background(), log.Trace(), log.Action("Write"), log.Reason(err.Error()))
		return
	}

	err = secretVerifier.Ensure()
	if err != nil {
		logger.Error(context.Background(), log.Trace(), log.Action("Ensure"), log.Reason(err.Error()))
		return
	}

	f.ServeHTTP(w, r)
}
