package postgres

import (
	"context"
	"errors"
	"hellper/internal/log"
	"hellper/internal/model"
	"hellper/internal/model/sql"
	"strings"

	_ "github.com/lib/pq"
)

type serviceRepository struct {
	logger log.Logger
	db     sql.DB
}

// NewServiceRepository creates a new instance of a repository to handle services information
func NewServiceRepository(logger log.Logger, db sql.DB) model.ServiceRepository {
	return &serviceRepository{
		logger,
		db,
	}
}

// ListServiceInstances returns all service instances registered in the database
func (r *serviceRepository) ListServiceInstances(ctx context.Context) ([]*model.ServiceInstance, error) {
	query := `
	SELECT
		service_instance.id as id,
		(service.name || ' / ' || service_instance.name) as name
	FROM public.service
	INNER JOIN public.service_instance on service_instance.service_id = service.id
	`

	rows, err := r.db.Query(query)

	if err != nil {
		r.logger.Error(
			ctx,
			"postgres/service-repository.ListServiceInstances Query ERROR",
			log.NewValue("Error", err),
		)

		return nil, err
	}

	defer rows.Close()

	serviceInstances := make([]*model.ServiceInstance, 0)
	for rows.Next() {
		instance := model.ServiceInstance{}
		rows.Scan(&instance.ID, &instance.Name)

		serviceInstances = append(serviceInstances, &instance)
	}

	return serviceInstances, nil
}

// GetServiceInstanceOwner returns the owner team name of a service instance registered in the database
func (r *serviceRepository) GetServiceInstanceOwnerTeamName(
	ctx context.Context, instanceName string,
) (string, error) {
	names := strings.Split(instanceName, "/")
	if len(names) < 2 {
		err := errors.New("Unable to parse instance name #" + instanceName)
		r.logger.Error(
			ctx,
			"postgres/incident-repository.GetServiceInstanceOwnerTeamName ERROR",
			log.NewValue("instanceName", instanceName),
			log.NewValue("error", err),
		)
		return "", err
	}

	serviceName := strings.TrimSpace(names[0])
	serviceInstanceName := strings.TrimSpace(names[1])

	query := `
	SELECT
    team.name as team_name
	FROM public.service
	INNER JOIN public.service_instance on service.id = service_instance.service_id
	INNER JOIN public.team on service_instance.owner_team_id = team.id
  WHERE service.name = $1 AND service_instance.name = $2
	`

	rows, err := r.db.Query(query, serviceName, serviceInstanceName)
	if err != nil {
		r.logger.Error(
			ctx,
			"postgres/service-repository.GetServiceInstanceOwnerTeamName Query ERROR",
			log.NewValue("instanceName", instanceName),
			log.NewValue("serviceName", serviceName),
			log.NewValue("serviceInstanceName", serviceInstanceName),
			log.NewValue("Error", err),
		)
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		err = errors.New("Owner team of service instance" + instanceName + " not found")
		r.logger.Error(
			ctx,
			"postgres/incident-repository.GetServiceInstanceOwnerTeamName ERROR",
			log.NewValue("instanceName", instanceName),
			log.NewValue("serviceName", serviceName),
			log.NewValue("serviceInstanceName", serviceInstanceName),
			log.NewValue("error", err),
		)
		return "", err
	}

	var ownerTeamName string
	rows.Scan(&ownerTeamName)
	r.logger.Info(
		ctx,
		"postgres/incident-repository.GetServiceInstanceOwnerTeamName SUCCESS",
		log.NewValue("instanceName", instanceName),
		log.NewValue("serviceName", serviceName),
		log.NewValue("serviceInstanceName", serviceInstanceName),
		log.NewValue("ownerTeamName", ownerTeamName),
	)

	return ownerTeamName, nil
}
