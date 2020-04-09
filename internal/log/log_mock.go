package log

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type LoggerMock struct {
	mock.Mock
}

func NewLoggerMock() *LoggerMock {
	return new(LoggerMock)
}

func (mock *LoggerMock) Debug(ctx context.Context, msg string, values ...Value) {
	mock.Called(ctx, msg, values)
}

func (mock *LoggerMock) Info(ctx context.Context, msg string, values ...Value) {
	mock.Called(ctx, msg, values)
}

func (mock *LoggerMock) Error(ctx context.Context, msg string, values ...Value) {
	mock.Called(ctx, msg, values)
}
