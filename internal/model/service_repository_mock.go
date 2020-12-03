package model

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ServiceRepositoryMock struct {
	mock.Mock
}

func NewServiceRepositoryMock() *ServiceRepositoryMock {
	return new(ServiceRepositoryMock)
}

func (r *ServiceRepositoryMock) ListServiceInstances(ctx context.Context) ([]*ServiceInstance, error) {
	args := r.Mock.Called(ctx)
	result := args.Get(0)

	if result == nil {
		return nil, args.Error(1)
	}

	return result.([]*ServiceInstance), args.Error(1)
}
