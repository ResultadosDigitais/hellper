package bot

import (
	"context"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func NewClientMock() *ClientMock {
	return new(ClientMock)
}

func (mock *ClientMock) PostMessage(
	channel string, options ...slack.MsgOption,
) (string, string, error) {
	args := mock.Called(channel, options)
	return args.String(0), args.String(1), args.Error(2)
}

func (mock *ClientMock) CreateConversationContext(ctx context.Context, channelName string, isPrivate bool) (*slack.Channel, error) {
	var (
		args   = mock.Called(ctx, channelName, isPrivate)
		result = args.Get(0)
	)
	if result == nil {
		return nil, args.Error(1)
	}

	return result.(*slack.Channel), args.Error(1)
}

func (mock *ClientMock) InviteUsersToConversationContext(ctx context.Context, channelID string, users ...string) (*slack.Channel, error) {
	var (
		args   = mock.Called(ctx, channelID, users)
		result = args.Get(0)
	)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*slack.Channel), args.Error(1)
}

func (mock *ClientMock) ListPins(channelID string) ([]slack.Item, *slack.Paging, error) {
	var (
		args   = mock.Called(channelID)
		result = args.Get(0)
		page   = args.Get(1)
		items  []slack.Item
		paging *slack.Paging
	)
	if result != nil {
		items = result.([]slack.Item)
	}
	if page != nil {
		paging = page.(*slack.Paging)
	}
	return items, paging, args.Error(2)
}

func (mock *ClientMock) GetUserInfoContext(ctx context.Context, userID string) (*slack.User, error) {
	var (
		args   = mock.Called(ctx, userID)
		result = args.Get(0)
	)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*slack.User), args.Error(1)
}

func (mock *ClientMock) GetUsersInConversationContext(ctx context.Context, params *slack.GetUsersInConversationParameters) ([]string, string, error) {
	var (
		args   = mock.Called(ctx, params)
		list   = args.Get(0)
		cursor = args.Get(1)
		err    = args.Error(2)
	)
	return list.([]string), cursor.(string), err
}

func (mock *ClientMock) SetTopicOfConversationContext(ctx context.Context, channelID, topic string) (*slack.Channel, error) {
	var (
		args   = mock.Called(ctx, channelID, topic)
		result = args.Get(0)
	)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*slack.Channel), args.Error(1)
}

func (mock *ClientMock) OpenDialog(triggerID string, dialog slack.Dialog) error {
	args := mock.Called(triggerID, dialog)
	return args.Error(0)
}

func (mock *ClientMock) AddPin(channel string, item slack.ItemRef) error {
	args := mock.Called(channel, item)
	return args.Error(0)
}

func (mock *ClientMock) ArchiveChannel(channelID string) error {
	return nil
}

func (mock *ClientMock) PostEphemeralContext(ctx context.Context, channelID string, userID string, options ...slack.MsgOption) (string, error) {
	args := mock.Called(ctx, channelID, userID, options)
	return args.String(0), args.Error(1)
}

func (mock *ClientMock) JoinConversation(string) (*slack.Channel, string, []string, error) {
	return nil, "", nil, nil
}

func (mock *ClientMock) InviteUsersToConversation(string, ...string) (*slack.Channel, error) {
	return nil, nil
}
