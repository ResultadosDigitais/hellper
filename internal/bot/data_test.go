package bot

import (
	"testing"

	"github.com/slack-go/slack"
)

func TestGetLastPin(t *testing.T) {
	itens := []slack.Item{
		slack.Item{
			Timestamp: "1",
		},
		slack.Item{
			Timestamp: "0",
		},
		slack.Item{
			Timestamp: "2",
		},
	}

	items, _ := getLastPin(itens)
	if items.Timestamp != "2" {
		t.Fatal("Not get last value")
	}
}

func TestEmptyLastPin(t *testing.T) {
	itens := []slack.Item{}

	_, err := getLastPin(itens)
	if err == nil {
		t.Fatal("Not return error with emmpty")
	}
}
