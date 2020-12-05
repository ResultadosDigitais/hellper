package model

import "context"

// ServiceRepository wraps services data from the database
type ServiceRepository interface {
	ListServiceInstances(ctx context.Context) ([]*ServiceInstance, error)
	GetServiceInstanceOwnerTeamName(ctx context.Context, instanceName string) (string, error)
}
