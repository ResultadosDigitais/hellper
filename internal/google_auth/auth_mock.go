package googleauth

import (
	"context"
	"hellper/internal/log"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type AuthMock struct {
	mock.Mock
}

func NewAuthMock() *AuthMock {
	return new(AuthMock)
}

func (mock *AuthMock) GetGClient(ctx context.Context, logger log.Logger, token []byte, scope string) (*http.Client, error) {
	var (
		args   = mock.Called(ctx, logger, token, scope)
		result = args.Get(0)
	)

	if result == nil {
		return nil, args.Error(1)
	}

	return result.(*http.Client), args.Error(1)
}
