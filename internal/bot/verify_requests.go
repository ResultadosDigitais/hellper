package bot

import (
	"bytes"
	"hellper/internal/config"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

func VerifyRequests(r *http.Request, w http.ResponseWriter, f http.Handler) {
	secretVerifier, err := slack.NewSecretsVerifier(r.Header, config.Env.SlackSigningSecret)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	_, err = secretVerifier.Write(body)
	if err != nil {
		log.Fatal(err)
	}

	err = secretVerifier.Ensure()
	if err != nil {
		log.Fatal(err)
	}

	f.ServeHTTP(w, r)
}
