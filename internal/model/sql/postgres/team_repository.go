package postgres

import (
	"context"
	"errors"
	"fmt"
	"hellper/internal/log"
	"hellper/internal/model"
	"hellper/internal/model/sql"
	"strings"
)

type teamRepository struct {
	logger log.Logger
	db     sql.DB
}

// NewTeamRepository creates a new instance of a repository for accessing team data
func NewTeamRepository(logger log.Logger, db sql.DB) model.TeamRepository {
	return &teamRepository{logger: logger, db: db}
}

func (r *teamRepository) GetOwnersByServiceInstance(ctx context.Context, serviceInstance model.ServiceInstance) ([]*model.User, error) {
	r.logger.Debug(ctx, "Executing GetOwnersByServiceInstance query")
	query := `
	SELECT
		person.email as email,
		person.slack_member_id as slackID
	FROM public.service_instance
	INNER JOIN public.team team on service_instance.owner_team_id = team.id
	INNER JOIN public.team_member on team_member.team_id = team.id
	INNER JOIN public.person on team_member.person_email = person.email
	WHERE service_instance.id = $1
	`

	serviceInstanceID, err := r.getServiceInstanceIDFromName(ctx, serviceInstance.ID)
	if err != nil {
		return []*model.User{}, err
	}

	rows, err := r.db.Query(query, serviceInstanceID)
	if err != nil {
		r.logger.Error(ctx, "Error while executing query GetOwnersByServiceInstance", log.NewValue("error", err))
		return []*model.User{}, err
	}

	defer rows.Close()
	users := make([]*model.User, 0)
	for rows.Next() {
		user := model.User{}
		rows.Scan(&user.Email, &user.SlackID)

		users = append(users, &user)
	}

	r.logger.Debug(ctx, "Query GetOwnersByServiceInstance was executed", log.NewValue("numberRows", len(users)))

	return users, nil
}

func (r *teamRepository) GetUsersOfServiceInstance(ctx context.Context, serviceInstance model.ServiceInstance) ([]*model.User, error) {
	r.logger.Debug(ctx, "Executing GetUsersOfServiceInstance query")
	query := `
	SELECT
		person.email as email,
		person.slack_member_id as slackID
	FROM public.service_instance
	INNER JOIN service_instance_stakeholder sis on sis.service_instance_id = service_instance.id
	INNER JOIN team on sis.team_id  = team.id
	INNER JOIN public.team_member on team_member.team_id = team.id
	INNER JOIN public.person on team_member.person_email = person.email
	WHERE service_instance.id = $1
	`

	serviceInstanceID, err := r.getServiceInstanceIDFromName(ctx, serviceInstance.ID)
	if err != nil {
		return []*model.User{}, err
	}

	rows, err := r.db.Query(query, serviceInstanceID)
	if err != nil {
		r.logger.Error(ctx, "Error while executing query GetUsersOfServiceInstance", log.NewValue("error", err))
		return []*model.User{}, err
	}

	defer rows.Close()
	users := make([]*model.User, 0)
	for rows.Next() {
		user := model.User{}
		rows.Scan(&user.Email, &user.SlackID)

		users = append(users, &user)
	}

	r.logger.Debug(ctx, "Query GetUsersOfServiceInstance was executed", log.NewValue("numberRows", len(users)))

	return users, nil
}

func (r *teamRepository) getServiceInstanceIDFromName(ctx context.Context, name string) (int64, error) {
	names := strings.Split(name, "/")
	if len(names) < 2 {
		err := errors.New("Unable to parse instance name #" + name)
		r.logger.Error(
			ctx,
			"postgres/incident-repository.GetServiceInstanceOwnerTeamName ERROR",
			log.NewValue("instanceName", name),
			log.NewValue("error", err),
		)
		return -1, err
	}
	serviceName := strings.TrimSpace(names[0])
	serviceInstanceName := strings.TrimSpace(names[1])

	query := `
		SELECT
			service_instance.id as id
		FROM public.service_instance
		INNER JOIN service on service_instance.service_id = service.id
		WHERE service.name = $1 AND service_instance.name = $2
	`

	rows, err := r.db.Query(query, serviceName, serviceInstanceName)
	if err != nil {
		r.logger.Error(ctx, "Error while executing query GetOwnersByServiceInstance", log.NewValue("error", err))
		return -1, err
	}

	defer rows.Close()
	var id int64
	if !rows.Next() {
		return -1, fmt.Errorf("Service instance \"%s/%s\" not found", serviceName, serviceInstanceName)
	}

	rows.Scan(&id)

	return id, nil
}
