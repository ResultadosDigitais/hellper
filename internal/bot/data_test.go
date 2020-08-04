package bot_test

import (
	"errors"
	"fmt"
	"hellper/internal/bot"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

type dataFixture struct {
	testName     string
	expectError  bool
	errorMessage string

	channelID string
	item      slack.Item
	itemList  []slack.Item

	mockClient    bot.Client
	listPinsError error
}

func (f *dataFixture) setup(t *testing.T) {
	clientMock := bot.NewClientMock()
	clientMock.On("ListPins", f.channelID).Return(f.itemList, new(slack.Paging), f.listPinsError)
	f.mockClient = clientMock
}

func TestEmptyLastPin(t *testing.T) {
	table := []dataFixture{
		{
			testName:    "teste 1",
			expectError: false,
			channelID:   "XYZ123",
			item:        buildItem("1595847600"),      //2020-07-27T12:00:00+00:00
			itemList:    buildItemLists("1595847600"), //2020-07-27T12:00:00+00:00
		},
		{
			testName:      "teste 2",
			expectError:   true,
			channelID:     "",
			errorMessage:  "Channel does not exist",
			listPinsError: errors.New("Channel does not exist"),
		},
		{
			testName:    "teste 3",
			expectError: false,
			channelID:   "XYZ123",
			item:        slack.Item{},
			itemList:    []slack.Item{},
		},
	}

	for index, f := range table {
		t.Run(fmt.Sprintf("%v-%v", index, f.testName), func(t *testing.T) {
			f.setup(t)

			lastItem, err := bot.LastPin(f.mockClient, f.channelID)

			if f.expectError {
				if err == nil {
					t.Fatal("an error was expected, but not occurred")
				}

				assert.EqualError(t, err, f.errorMessage)
			} else {
				if err != nil {
					t.Fatal(
						"an error occurred, but was not expected\n",
						"error: ",
						err,
					)
				}

				assert.Equal(t, f.item, lastItem)
			}
		})
	}
}

func buildItem(ts string) slack.Item {
	return slack.Item{
		Type:      "message",
		Channel:   "XYZ123",
		Timestamp: ts,
		Message: &slack.Message{
			Msg: slack.Msg{
				Type:      "Message",
				Channel:   "XYZ123",
				User:      "XUNDA",
				Text:      "status 1",
				Timestamp: ts,
			},
		},
	}
}

func buildItemLists(lastTS string) []slack.Item {
	return []slack.Item{
		buildItem(lastTS),
		buildItem("1595836800"), //2020-07-27T08:00:00+00:00
		buildItem("1595844000"), //2020-07-27T10:00:00+00:00
		buildItem("1595845800"), //2020-07-27T10:30:00+00:00
	}
}
