package log

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockLogger(t *testing.T) {
	logger := NewLoggerMock()
	assert.Implements(t, (*Logger)(nil), logger, "logger bad instance")

	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything).Once()
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Once()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Once()

	logger.Debug(context.Background(), "debug", NewValue("key", "value"))
	logger.Debug(context.Background(), "info", NewValue("key", "value"))
	logger.Error(context.Background(), "error",
		NewValue("key", "value"), NewValue("error", errors.New("err_mock")),
	)

	logger.AssertExpectations(t)
}
