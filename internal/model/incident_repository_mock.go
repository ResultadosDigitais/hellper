package model

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type IncidentRepositoryMock struct {
	mock.Mock
}

func NewIncidentRepositoryMock() *IncidentRepositoryMock {
	return new(IncidentRepositoryMock)
}

func (mock *IncidentRepositoryMock) SetIncident(inc *Incident) error {
	args := mock.Called(inc)
	return args.Error(0)
}

func (mock *IncidentRepositoryMock) GetIncident(ctx context.Context, channelID string) (Incident, error) {
	args := mock.Called(channelID)
	return args.Get(0).(Incident), args.Error(1)
}

func (mock *IncidentRepositoryMock) ListActiveIncidents(ctx context.Context) ([]Incident, error) {
	var (
		args   = mock.Called()
		result = args.Get(0)
	)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]Incident), args.Error(1)
}

func (mock *IncidentRepositoryMock) AddPostMortemUrl(ctx context.Context, channelName string, postMortemUrl string) error {
	args := mock.Called(channelName, postMortemUrl)
	return args.Error(0)
}

func (mock *IncidentRepositoryMock) CancelIncident(ctx context.Context, inc *Incident) error {
	args := mock.Called(ctx, inc.ChannelID, inc.DescriptionCancelled)
	return args.Error(0)
}

func (mock *IncidentRepositoryMock) CloseIncident(ctx context.Context, inc *Incident) error {
	args := mock.Called(inc)
	return args.Error(0)
}

func (mock *IncidentRepositoryMock) InsertIncident(ctx context.Context, inc *Incident) (int64, error) {
	args := mock.Called(inc)
	return args.Get(0).(int64), args.Error(1)
}

func (mock *IncidentRepositoryMock) ResolveIncident(ctx context.Context, inc *Incident) error {
	args := mock.Called(ctx, inc)
	return args.Error(0)
}

func (mock *IncidentRepositoryMock) UpdateIncidentDates(ctx context.Context, inc *Incident) error {
	args := mock.Called(ctx, inc)
	return args.Error(0)
}

func (mock *IncidentRepositoryMock) PauseNotifyIncident(ctx context.Context, inc *Incident) error {
	args := mock.Called(ctx, inc)
	return args.Error(0)
}
