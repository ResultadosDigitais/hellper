package model

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func NewRepositoryMock() *RepositoryMock {
	return new(RepositoryMock)
}

func (mock *RepositoryMock) SetIncident(inc *Incident) error {
	args := mock.Called(inc)
	return args.Error(0)
}

func (mock *RepositoryMock) GetIncident(ctx context.Context, channelID string) (Incident, error) {
	args := mock.Called(channelID)
	return args.Get(0).(Incident), args.Error(1)
}

func (mock *RepositoryMock) ListActiveIncidents(ctx context.Context) ([]Incident, error) {
	var (
		args   = mock.Called()
		result = args.Get(0)
	)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]Incident), args.Error(1)
}

func (mock *RepositoryMock) AddPostMortemUrl(ctx context.Context, channelName string, postMortemUrl string) error {
	args := mock.Called(channelName, postMortemUrl)
	return args.Error(1)
}

func (mock *RepositoryMock) CancelIncident(ctx context.Context, channelID string, description string) error {
	args := mock.Called(channelID, description)
	return args.Error(1)
}

func (mock *RepositoryMock) CloseIncident(ctx context.Context, inc *Incident) error {
	args := mock.Called(inc)
	return args.Error(1)
}

func (mock *RepositoryMock) InsertIncident(ctx context.Context, inc *Incident) (int64, error) {
	args := mock.Called(inc)
	return args.Get(0).(int64), args.Error(0)
}

func (mock *RepositoryMock) ResolveIncident(ctx context.Context, inc *Incident) error {
	args := mock.Called(ctx, inc)
	return args.Error(0)
}

func (mock *RepositoryMock) UpdateIncidentDates(ctx context.Context, inc *Incident) error {
	args := mock.Called(ctx, inc)
	return args.Error(0)
}

func (mock *RepositoryMock) PauseNotifyIncident(ctx context.Context, inc *Incident) error {
	args := mock.Called(ctx, inc)
	return args.Error(0)
}
