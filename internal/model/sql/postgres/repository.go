package postgres

import (
	"context"
	"errors"
	"fmt"

	"hellper/internal/log"
	"hellper/internal/model"
	"hellper/internal/model/sql"

	_ "github.com/lib/pq"
)

type repository struct {
	logger log.Logger
	db     sql.DB
}

func NewRepository(logger log.Logger, db sql.DB) model.Repository {
	return &repository{
		logger: logger,
		db:     db,
	}
}

func incidentLogValues(inc *model.Incident) []log.Value {
	return []log.Value{
		log.NewValue("title", inc.Title),
		log.NewValue("descriptionStarted", inc.DescriptionStarted),
		log.NewValue("descriptionCancelled", inc.DescriptionCancelled),
		log.NewValue("descriptionResolved", inc.DescriptionResolved),
		log.NewValue("startedAt", inc.StartedAt),
		log.NewValue("identifiedAt", inc.IdentifiedAt),
		log.NewValue("endedAt", inc.EndedAt),
		log.NewValue("snoozedTime", inc.SnoozedUntil),
		log.NewValue("responsibility", inc.Responsibility),
		log.NewValue("functionality", inc.Functionality),
		log.NewValue("rootCause", inc.RootCause),
		log.NewValue("customerImpact", inc.CustomerImpact),
		log.NewValue("statusPageURL", inc.StatusPageURL),
		log.NewValue("postMortemURL", inc.PostMortemURL),
		log.NewValue("team", inc.Team),
		log.NewValue("product", inc.Product),
		log.NewValue("severityLevel", inc.SeverityLevel),
		log.NewValue("severityLevel", inc.SeverityLevel),
		log.NewValue("channelName", inc.ChannelName),
		log.NewValue("channelID", inc.ChannelID),
		log.NewValue("commanderID", inc.CommanderID),
		log.NewValue("commanderEmail", inc.CommanderEmail),
	}
}

func (r *repository) InsertIncident(ctx context.Context, inc *model.Incident) (int64, error) {
	r.logger.Info(
		ctx,
		log.Trace(),
		incidentLogValues(inc)...,
	)

	insertCommand := `INSERT INTO incidents
		( title
		, description_started
		, description_cancelled
		, description_resolved
		, started_at
		, ended_at
		, identified_at
		, responsibility
		, functionality
		, root_cause
		, customer_impact
		, status_page_url
		, post_mortem_url
		, status
		, product
		, severity_level
		, channel_name
		, channel_id
		, commander_id
		, commander_email)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	RETURNING id`

	id := int64(0)

	idResult := r.db.QueryRow(
		insertCommand,
		inc.Title,
		inc.DescriptionStarted,
		inc.DescriptionCancelled,
		inc.DescriptionResolved,
		inc.StartedAt,
		inc.EndedAt,
		inc.IdentifiedAt,
		inc.Responsibility,
		inc.Functionality,
		inc.RootCause,
		inc.CustomerImpact,
		inc.StatusPageURL,
		inc.PostMortemURL,
		inc.Status,
		inc.Product,
		inc.SeverityLevel,
		inc.ChannelName,
		inc.ChannelID,
		inc.CommanderID,
		inc.CommanderEmail)

	switch err := idResult.Scan(&id); err {
	case nil:
		r.logger.Info(
			ctx,
			log.Trace(),
			incidentLogValues(inc)...,
		)
		return id, nil
	default:
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return 0, err
	}
}

func (r *repository) AddPostMortemUrl(ctx context.Context, channelName string, postMortemURL string) error {
	r.logger.Info(
		ctx,
		log.Trace(),
		log.NewValue("channelName", channelName),
		log.NewValue("postMortemURL", postMortemURL),
	)

	updateCommand := `UPDATE incidents SET post_mortem_url = $1 WHERE channel_name = $2`

	_, err := r.db.Exec(
		updateCommand,
		postMortemURL,
		channelName)

	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			log.Action("r.db.Exec"),
			log.Reason(err.Error()),
			log.NewValue("channelName", channelName),
			log.NewValue("postMortemURL", postMortemURL),
		)
	} else {
		r.logger.Info(
			ctx,
			log.Trace(),
			log.NewValue("channelName", channelName),
			log.NewValue("postMortemURL", postMortemURL),
		)
	}

	return err
}

func (r *repository) GetIncident(ctx context.Context, channelID string) (inc model.Incident, err error) {
	r.logger.Info(
		ctx,
		log.Trace(),
		log.NewValue("channelID", channelID),
	)

	rows, err := r.db.Query(
		GetIncidentByChannelID(),
		channelID,
	)
	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			log.Action("r.db.Query"),
			log.Reason(err.Error()),
			log.NewValue("channelID", channelID),
		)

		return model.Incident{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		err = errors.New("Incident " + channelID + "not found")
		r.logger.Error(
			ctx,
			log.Trace(),
			log.Action("rows.Next"),
			log.Reason(err.Error()),
			log.NewValue("channelID", channelID),
		)

		return model.Incident{}, err
	}

	rows.Scan(
		&inc.ID,
		&inc.Title,
		&inc.DescriptionStarted,
		&inc.DescriptionCancelled,
		&inc.DescriptionResolved,
		&inc.StartedAt,
		&inc.EndedAt,
		&inc.IdentifiedAt,
		&inc.SnoozedUntil,
		&inc.Responsibility,
		&inc.Functionality,
		&inc.RootCause,
		&inc.CustomerImpact,
		&inc.StatusPageURL,
		&inc.PostMortemURL,
		&inc.Status,
		&inc.Product,
		&inc.SeverityLevel,
		&inc.ChannelName,
		&inc.ChannelID,
		&inc.CommanderID,
		&inc.CommanderEmail,
	)

	r.logger.Info(
		ctx,
		log.Trace(),
		log.NewValue("channelID", channelID),
	)
	return inc, nil
}

func GetIncidentByChannelID() string {
	return `SELECT
		id
		, title
		, CASE WHEN description_started IS NULL THEN '' ELSE description_started END description_started
		, CASE WHEN description_cancelled IS NULL THEN '' ELSE description_cancelled END description_cancelled
		, CASE WHEN description_resolved IS NULL THEN '' ELSE description_resolved END description_resolved
		, started_at
		, ended_at
		, identified_at
    , snoozed_until
    , responsibility
		, functionality
		, root_cause
		, customer_impact
		, status_page_url
		, post_mortem_url
		, status
		, product
		, CASE WHEN severity_level IS NULL THEN 0 ELSE severity_level END AS severity_level
		, CASE WHEN channel_name IS NULL THEN '' ELSE channel_name END AS channel_name
		, CASE WHEN channel_id IS NULL THEN '' ELSE channel_id END AS channel_id
		, CASE WHEN commander_id IS NULL THEN '' ELSE commander_id END commander_id
		, CASE WHEN commander_email IS NULL THEN '' ELSE commander_email END commander_email
	FROM incidents
	WHERE channel_id = $1
	LIMIT 1`
}

func (r *repository) UpdateIncidentDates(ctx context.Context, inc *model.Incident) error {
	r.logger.Info(
		ctx,
		log.Trace(),
		incidentLogValues(inc)...,
	)

	result, err := r.db.Exec(
		`UPDATE incidents SET
			started_at = $1,
			identified_at = $2,
			ended_at = $3
		WHERE channel_id = $4`,
		inc.StartedAt,
		inc.IdentifiedAt,
		inc.EndedAt,
		inc.ChannelID,
	)
	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	if rowsAffected == 0 {
		err = errors.New("rows not affected")
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	r.logger.Info(
		ctx,
		log.Trace(),
		incidentLogValues(inc)...,
	)

	return nil
}

func (r *repository) CancelIncident(ctx context.Context, inc *model.Incident) error {
	r.logger.Info(
		ctx,
		log.Trace(),
		log.NewValue("channelID", inc.ChannelID),
		log.NewValue("descriptionCancel", inc.DescriptionCancelled),
	)
	result, err := r.db.Exec(
		`UPDATE incidents SET status = $1, description_cancelled = $2 WHERE channel_id = $3`,
		model.StatusCancel,
		inc.DescriptionCancelled,
		inc.ChannelID,
	)

	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			log.Action("r.db.Exec"),
			log.Reason(err.Error()),
			log.NewValue("channelID", inc.ChannelID),
			log.NewValue("description", inc.DescriptionCancelled),
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			log.Action("result.RowsAffected"),
			log.Reason(err.Error()),
			log.NewValue("channelID", inc.ChannelID),
			log.NewValue("error", err),
		)

		return err
	}

	if rowsAffected == 0 {
		err = errors.New("rows not affected")
		r.logger.Error(
			ctx,
			log.Trace(),
			log.Action("rowsAffected"),
			log.Reason(err.Error()),
			log.NewValue("channelID", inc.ChannelID),
			log.NewValue("description", inc.DescriptionCancelled),
		)

		return err
	}

	r.logger.Info(
		ctx,
		log.Trace(),
		log.NewValue("channelID", inc.ChannelID),
		log.NewValue("descriptionCancel", inc.DescriptionCancelled),
	)
	return nil
}

func (r *repository) CloseIncident(ctx context.Context, inc *model.Incident) error {
	//TODO: implement team
	r.logger.Info(
		ctx,
		log.Trace(),
		incidentLogValues(inc)...,
	)

	result, err := r.db.Exec(
		`UPDATE incidents SET
			root_cause = $1,
			functionality = $2,
			team = $3,
			customer_impact = $4,
			severity_level = $5,
			status = $6,
			responsibility = $7
		WHERE channel_id = $8`,
		inc.RootCause,
		inc.Functionality,
		inc.Team,
		inc.CustomerImpact.Int64,
		inc.SeverityLevel,
		model.StatusClosed,
		inc.Responsibility,
		inc.ChannelID,
	)

	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	if rowsAffected == 0 {
		err = errors.New("rows not affected")
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}
	r.logger.Info(
		ctx,
		log.Trace(),
		incidentLogValues(inc)...,
	)

	return nil
}

func (r *repository) ResolveIncident(ctx context.Context, inc *model.Incident) error {
	//TODO: implement team
	r.logger.Info(
		ctx,
		log.Trace(),
		incidentLogValues(inc)...,
	)

	result, err := r.db.Exec(
		`UPDATE incidents SET
			status_page_url = $1,
			description_resolved = $2,
			started_at = $3,
			ended_at = $4,
			status = $5
		WHERE channel_id = $6`,
		inc.StatusPageURL,
		inc.DescriptionResolved,
		inc.StartedAt,
		inc.EndedAt,
		model.StatusResolved,
		inc.ChannelID,
	)

	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	if rowsAffected == 0 {
		err = errors.New("rows not affected")
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	r.logger.Info(
		ctx,
		log.Trace(),
		incidentLogValues(inc)...,
	)

	return nil
}

func (r *repository) ListActiveIncidents(ctx context.Context) ([]model.Incident, error) {
	r.logger.Info(
		ctx,
		log.Trace(),
	)
	var (
		incidents    []model.Incident
		logIncidents []log.Value
	)

	rows, err := r.db.Query(
		GetIncidentStatusFilterQuery(),
		model.StatusOpen,
		model.StatusResolved,
	)
	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			log.NewValue("error", err),
		)

		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		i++
		var inc model.Incident
		err := rows.Scan(
			&inc.ID,
			&inc.Title,
			&inc.DescriptionStarted,
			&inc.DescriptionCancelled,
			&inc.DescriptionResolved,
			&inc.StartedAt,
			&inc.EndedAt,
			&inc.IdentifiedAt,
			&inc.SnoozedUntil,
			&inc.Responsibility,
			&inc.Functionality,
			&inc.RootCause,
			&inc.CustomerImpact,
			&inc.StatusPageURL,
			&inc.PostMortemURL,
			&inc.Status,
			&inc.Product,
			&inc.SeverityLevel,
			&inc.ChannelName,
			&inc.ChannelID,
			&inc.CommanderID,
			&inc.CommanderEmail,
		)
		if err != nil {
			r.logger.Error(
				ctx,
				log.Trace(),
				log.NewValue("error", err),
			)

			return nil, err
		}
		logIncidents = append(logIncidents, log.NewValue(fmt.Sprintf("Incident %d", i), incidentLogValues(&inc)))
		incidents = append(incidents, inc)
	}

	r.logger.Info(
		ctx,
		log.Trace(),
		logIncidents...,
	)

	return incidents, nil
}

func GetIncidentStatusFilterQuery() string {
	return `SELECT
		  id
		, title
		, CASE WHEN description_started IS NULL THEN '' ELSE description_started END description_started
		, CASE WHEN description_cancelled IS NULL THEN '' ELSE description_cancelled END description_cancelled
		, CASE WHEN description_resolved IS NULL THEN '' ELSE description_resolved END description_resolved
		, started_at
		, ended_at
		, identified_at
    , snoozed_until
		, responsibility
		, functionality
		, root_cause
		, customer_impact
		, status_page_url
		, post_mortem_url
		, status
		, product
		, CASE WHEN severity_level IS NULL THEN 0 ELSE severity_level END AS severity_level
		, CASE WHEN channel_name IS NULL THEN '' ELSE channel_name END AS channel_name
		, CASE WHEN channel_id IS NULL THEN '' ELSE channel_id END AS channel_id
		, CASE WHEN commander_id IS NULL THEN '' ELSE commander_id END commander_id
		, CASE WHEN commander_email IS NULL THEN '' ELSE commander_email END commander_email
	FROM incidents
	WHERE status IN ($1, $2)
	LIMIT 100`
}

func (r *repository) PauseNotifyIncident(ctx context.Context, inc *model.Incident) error {
	r.logger.Info(
		ctx,
		log.Trace(),
		incidentLogValues(inc)...,
	)

	result, err := r.db.Exec(
		`UPDATE incidents SET
			snoozed_until = $1
		WHERE channel_id = $2`,
		inc.SnoozedUntil.Time,
		inc.ChannelID,
	)
	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	if rowsAffected == 0 {
		err = errors.New("rows not affected")
		r.logger.Error(
			ctx,
			log.Trace(),
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	r.logger.Info(
		ctx,
		log.Trace(),
		incidentLogValues(inc)...,
	)

	return nil
}
