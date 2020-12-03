package postgres

import (
	"context"
	"hellper/internal/log"
	"hellper/internal/model"
	"hellper/internal/model/sql"

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
		(service.name || '/' || service_instance.name) as name
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
