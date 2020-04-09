package zap

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type mockZapLogger struct {
	mock.Mock
}

func newMockZapLogger() *mockZapLogger {
	return new(mockZapLogger)
}

func (mock *mockZapLogger) Check(level zapcore.Level, msg string) zapWriter {
	args := mock.Called(level, msg)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(zapWriter)
}

func (mock *mockZapLogger) Sync() error {
	args := mock.Called()
	return args.Error(0)
}

type mockZapWriter struct {
	mock.Mock
}

func newMockZapWriter() *mockZapWriter {
	return new(mockZapWriter)
}

func (mock *mockZapWriter) Write(fields ...zap.Field) {
	mock.Called(fields)
}
