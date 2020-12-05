package model

import "context"

// TeamRepository wraps all data access related to teams
type TeamRepository interface {
	GetOwnersByServiceInstance(ctx context.Context, serviceInstance ServiceInstance) ([]*User, error)
	GetUsersOfServiceInstance(ctx context.Context, ServiceInstance ServiceInstance) ([]*User, error)
}
