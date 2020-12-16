package bot

import (
	"context"

	"github.com/slack-go/slack"
)

type Client interface {
	PostEphemeralContext(context.Context, string, string, ...slack.MsgOption) (string, error)
	PostMessage(string, ...slack.MsgOption) (string, string, error)
	CreateConversationContext(ctx context.Context, channelName string, isPrivate bool) (*slack.Channel, error)
	InviteUsersToConversationContext(ctx context.Context, channelID string, users ...string) (*slack.Channel, error)
	ListPins(string) ([]slack.Item, *slack.Paging, error)
	GetUserInfoContext(context.Context, string) (*slack.User, error)
	SetTopicOfConversation(channelID, topic string) (*slack.Channel, error)
	OpenDialog(string, slack.Dialog) error
	AddPin(string, slack.ItemRef) error
	ArchiveConversationContext(ctx context.Context, channelID string) error
	JoinConversationContext(ctx context.Context, channelID string) (*slack.Channel, string, []string, error)
	GetUsersInConversationContext(context.Context, *slack.GetUsersInConversationParameters) ([]string, string, error)
}
