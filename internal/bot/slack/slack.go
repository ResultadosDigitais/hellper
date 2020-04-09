package slack

import (
	"hellper/internal/bot"

	"github.com/slack-go/slack"
)

func NewClient(token string) bot.Client {
	return slack.New(token)
}
