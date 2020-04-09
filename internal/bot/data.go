package bot

import (
	"errors"

	"github.com/slack-go/slack"
)

type EmptyPin struct {
	s string
}

func (e *EmptyPin) Error() string {
	return e.s
}

func getLastPin(itens []slack.Item) (result slack.Item, err error) {
	if len(itens) == 0 {
		return slack.Item{}, &EmptyPin{s: "itens is empty"}
	}
	for _, item := range itens {
		if result.Message == nil || result.Message.Msg.Timestamp < item.Message.Msg.Timestamp {
			result = item
		}
	}
	return
}

//LastPin return last message pinned by date of message
func LastPin(client Client, channelID string) (result slack.Item, err error) {
	itens, _, err := client.ListPins(channelID)
	if len(itens) == 0 {
		return slack.Item{}, errors.New("Lista está vazia")
	}

	if err != nil {
		return slack.Item{}, err
	}
	return getLastPin(itens)
}
