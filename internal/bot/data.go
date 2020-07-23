package bot

import (
	"github.com/slack-go/slack"
)

type EmptyPin struct {
	s string
}

func (e *EmptyPin) Error() string {
	return e.s
}

func getLastPin(itens []slack.Item) slack.Item {
	var result slack.Item
	for _, item := range itens {
		if result.Message == nil || result.Message.Msg.Timestamp < item.Message.Msg.Timestamp {
			result = item
		}
	}
	return result
}

//LastPin return last message pinned by date of message
func LastPin(client Client, channelID string) (slack.Item, error) {
	itens, _, err := client.ListPins(channelID)
	if err != nil {
		return slack.Item{}, err
	}

	return getLastPin(itens), err
}
