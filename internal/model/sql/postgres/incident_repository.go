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

type incidentRepository struct {
	logger log.Logger
	db     sql.DB
}

// NewIncidentRepository creates a new repository to handle the Incident entity
func NewIncidentRepository(logger log.Logger, db sql.DB) model.IncidentRepository {
	return &incidentRepository{
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
		log.NewValue("startTime", inc.StartTimestamp),
		log.NewValue("identificationTime", inc.IdentificationTimestamp),
		log.NewValue("endTime", inc.EndTimestamp),
		log.NewValue("rootCause", inc.RootCause),
		log.NewValue("meetingURL", inc.MeetingURL),
		log.NewValue("postMortemURL", inc.PostMortemURL),
		log.NewValue("team", inc.Team),
		log.NewValue("product", inc.Product),
		log.NewValue("severityLevel", inc.SeverityLevel),
		log.NewValue("channelName", inc.ChannelName),
		log.NewValue("channelID", inc.ChannelID),
		log.NewValue("commanderID", inc.CommanderID),
		log.NewValue("commanderEmail", inc.CommanderEmail),
	}
}

// InsertIncident inserts a new incident on a database
func (r *incidentRepository) InsertIncident(ctx context.Context, inc *model.Incident) (int64, error) {
	r.logger.Debug(
		ctx,
		"postgres/incident-repository.InsertIncident INFO",
		incidentLogValues(inc)...,
	)

	insertCommand := `INSERT INTO incident
		( title
		, description_started
		, description_cancelled
		, description_resolved
		, start_ts
		, end_ts
		, identification_ts
		, root_cause
		, meeting_url
		, post_mortem_url
		, status
		, product
		, severity_level
		, channel_name
		, channel_id
		, commander_id
		, commander_email)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	RETURNING id`

	id := int64(0)

	idResult := r.db.QueryRow(
		insertCommand,
		inc.Title,
		inc.DescriptionStarted,
		inc.DescriptionCancelled,
		inc.DescriptionResolved,
		inc.StartTimestamp,
		inc.EndTimestamp,
		inc.IdentificationTimestamp,
		inc.RootCause,
		inc.MeetingURL,
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
		r.logger.Debug(
			ctx,
			"postgres/incident-repository.InsertIncident SUCCESS",
			incidentLogValues(inc)...,
		)
		return id, nil
	default:
		r.logger.Error(
			ctx,
			"postgres/incident-repository.InsertIncident ERROR",
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return 0, err
	}
}

// UpdateIncident updates the incident on a database
func (r *incidentRepository) UpdateIncident(ctx context.Context, inc *model.Incident) error {
	r.logger.Info(
		ctx,
		"postgres/incident-repository.UpdateIncident INFO",
		incidentLogValues(inc)...,
	)

	updateCommand := `UPDATE incident SET
		title               = $1,
		description_started = $2,
		start_ts            = $3,
		root_cause          = $4,
		meeting_url         = $5,
		post_mortem_url     = $6,
		product             = $7,
		severity_level      = $8,
		commander_id        = $9,
		commander_email     = $10
	WHERE id = $11
	RETURNING id`

	_, err := r.db.Exec(
		updateCommand,
		inc.Title,
		inc.DescriptionStarted,
		inc.StartTimestamp,
		inc.RootCause,
		inc.MeetingURL,
		inc.PostMortemURL,
		inc.Product,
		inc.SeverityLevel,
		inc.CommanderID,
		inc.CommanderEmail,
		inc.ID,
	)

	if err != nil {
		r.logger.Error(
			ctx,
			"postgres/incident-repository.UpdateIncident ERROR",
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)

		return err
	}

	r.logger.Info(
		ctx,
		"postgres/incident-repository.UpdateIncident SUCCESS",
		incidentLogValues(inc)...,
	)

	return nil
}

// AddPostMortemUrl adds a PostMortemUrl into an incident registerd on the repository
func (r *incidentRepository) AddPostMortemURL(ctx context.Context, channelName string, postMortemURL string) error {
	logWriter := r.logger.With(
		log.NewValue("channelName", channelName),
		log.NewValue("postMortemURL", postMortemURL),
	)
	logWriter.Debug(
		ctx,
		"postgres/incident-repository.AddPostMortemUrl INFO",
	)

	updateCommand := `UPDATE incident SET post_mortem_url = $1 WHERE channel_name = $2`

	_, err := r.db.Exec(
		updateCommand,
		postMortemURL,
		channelName)

	if err != nil {
		logWriter.Error(
			ctx,
			"postgres/incident-repository.AddPostMortemUrl ERROR",
			log.NewValue("error", err),
		)
	}

	logWriter.Debug(
		ctx,
		"postgres/incident-repository.AddPostMortemUrl SUCCESS",
	)

	return err
}

// GetIncident retrieves an incident entity from the repository given a channelID
func (r *incidentRepository) GetIncident(ctx context.Context, channelID string) (inc model.Incident, err error) {
	logWriter := r.logger.With(
		log.NewValue("channelID", channelID),
	)

	logWriter.Debug(
		ctx,
		"postgres/incident-repository.GetIncident INFO",
	)

	rows, err := r.db.Query(
		GetIncidentByChannelID(),
		channelID,
	)
	if err != nil {
		logWriter.Error(
			ctx,
			"postgres/incident-repository.GetIncident Query ERROR",
			log.NewValue("error", err),
		)

		return model.Incident{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		err = errors.New("Incident " + channelID + "not found")
		logWriter.Error(
			ctx,
			"postgres/incident-repository.GetIncident ERROR",
			log.NewValue("error", err),
		)

		return model.Incident{}, err
	}

	rows.Scan(
		&inc.ID,
		&inc.Title,
		&inc.DescriptionStarted,
		&inc.DescriptionCancelled,
		&inc.DescriptionResolved,
		&inc.StartTimestamp,
		&inc.EndTimestamp,
		&inc.IdentificationTimestamp,
		&inc.RootCause,
		&inc.MeetingURL,
		&inc.PostMortemURL,
		&inc.Status,
		&inc.Product,
		&inc.SeverityLevel,
		&inc.ChannelName,
		&inc.ChannelID,
		&inc.CommanderID,
		&inc.CommanderEmail,
	)

	logWriter.Debug(
		ctx,
		"postgres/incident-repository.GetIncident SUCCESS",
	)
	return inc, nil
}

// GetIncidentByChannelID retrieves an Incident given a channelID
func GetIncidentByChannelID() string {
	return `SELECT
		id
		, title
		, CASE WHEN description_started IS NULL THEN '' ELSE description_started END description_started
		, CASE WHEN description_cancelled IS NULL THEN '' ELSE description_cancelled END description_cancelled
		, CASE WHEN description_resolved IS NULL THEN '' ELSE description_resolved END description_resolved
		, start_ts
		, end_ts
		, identification_ts
		, root_cause
		, meeting_url
		, post_mortem_url
		, status
		, product
		, CASE WHEN severity_level IS NULL THEN 0 ELSE severity_level END AS severity_level
		, CASE WHEN channel_name IS NULL THEN '' ELSE channel_name END AS channel_name
		, CASE WHEN channel_id IS NULL THEN '' ELSE channel_id END AS channel_id
		, CASE WHEN commander_id IS NULL THEN '' ELSE commander_id END commander_id
		, CASE WHEN commander_email IS NULL THEN '' ELSE commander_email END commander_email
	FROM incident
	WHERE channel_id = $1
	LIMIT 1`
}

func (r *incidentRepository) CancelIncident(ctx context.Context, inc *model.Incident) error {
	logWriter := r.logger.With(
		log.NewValue("channelID", inc.ChannelID),
		log.NewValue("descriptionCancel", inc.DescriptionCancelled),
	)

	logWriter.Debug(
		ctx,
		"postgres/incident-repository.CancelIncident DEBUG",
	)

	result, err := r.db.Exec(
		`UPDATE incident SET status = $1, description_cancelled = $2 WHERE channel_id = $3`,
		model.StatusCancel,
		inc.DescriptionCancelled,
		inc.ChannelID,
	)

	if err != nil {
		logWriter.Error(
			ctx,
			"postgres/incident-repository.CancelIncident ERROR",
			log.NewValue("error", err),
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		logWriter.Error(
			ctx,
			"postgres/incident-repository.CancelIncident ERROR",
			log.NewValue("error", err),
		)

		return err
	}

	if rowsAffected == 0 {
		err = errors.New("rows not affected")
		logWriter.Error(
			ctx,
			"postgres/incident-repository.CancelIncident ERROR",
			log.NewValue("error", err),
		)

		return err
	}

	logWriter.Info(
		ctx,
		"postgres/incident-repository.CancelIncident SUCCESS",
	)
	return nil
}

func (r *incidentRepository) CloseIncident(ctx context.Context, inc *model.Incident) error {
	//TODO: implement team
	r.logger.Info(
		ctx,
		"postgres/incident-repository.CloseIncident INFO",
		incidentLogValues(inc)...,
	)

	result, err := r.db.Exec(
		`UPDATE incident SET
			root_cause = $1,
			team = $2,
			severity_level = $3,
			status = $4
		WHERE channel_id = $5`,
		inc.RootCause,
		inc.Team,
		inc.SeverityLevel,
		model.StatusClosed,
		inc.ChannelID,
	)

	if err != nil {
		r.logger.Error(
			ctx,
			"postgres/incident-repository.CloseIncident ERROR",
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
			"postgres/incident-repository.CloseIncident ERROR",
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
			"postgres/incident-repository.CloseIncident ERROR",
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}
	r.logger.Info(
		ctx,
		"postgres/incident-repository.CloseIncident SUCCESS",
		incidentLogValues(inc)...,
	)

	return nil
}

func (r *incidentRepository) ResolveIncident(ctx context.Context, inc *model.Incident) error {
	//TODO: implement team
	r.logger.Info(
		ctx,
		"postgres/incident-repository.ResolveIncident INFO",
		incidentLogValues(inc)...,
	)

	result, err := r.db.Exec(
		`UPDATE incident SET
			description_resolved = $1,
			start_ts = $2,
			end_ts = $3,
			status = $4
		WHERE channel_id = $5`,
		inc.DescriptionResolved,
		inc.StartTimestamp,
		inc.EndTimestamp,
		model.StatusResolved,
		inc.ChannelID,
	)

	if err != nil {
		r.logger.Error(
			ctx,
			"postgres/incident-repository.ResolveIncident ERROR",
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
			"postgres/incident-repository.ResolveIncident ERROR",
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
			"postgres/incident-repository.ResolveIncident ERROR",
			append(
				incidentLogValues(inc),
				log.NewValue("error", err),
			)...,
		)
		return err
	}

	r.logger.Info(
		ctx,
		"postgres/incident-repository.ResolveIncident SUCCESS",
		incidentLogValues(inc)...,
	)

	return nil
}

func (r *incidentRepository) ListActiveIncidents(ctx context.Context) ([]model.Incident, error) {
	r.logger.Info(
		ctx,
		"postgres/incident-repository.ListActiveIncidents",
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
			"postgres/incident-repository.ListActiveIncidents Query ERROR",
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
			&inc.StartTimestamp,
			&inc.EndTimestamp,
			&inc.IdentificationTimestamp,
			&inc.RootCause,
			&inc.MeetingURL,
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
				"postgres/incident-repository.ListActiveIncidents Scan ERROR",
				log.NewValue("error", err),
			)

			return nil, err
		}
		logIncidents = append(logIncidents, log.NewValue(fmt.Sprintf("Incident %d", i), incidentLogValues(&inc)))
		incidents = append(incidents, inc)
	}

	r.logger.Info(
		ctx,
		"postgres/incident-repository.ListActiveIncidents SUCCESS",
		logIncidents...,
	)

	return incidents, nil
}

// GetIncidentStatusFilterQuery returns a query to filter incident by status
func GetIncidentStatusFilterQuery() string {
	return `SELECT
		  id
		, title
		, CASE WHEN description_started IS NULL THEN '' ELSE description_started END description_started
		, CASE WHEN description_cancelled IS NULL THEN '' ELSE description_cancelled END description_cancelled
		, CASE WHEN description_resolved IS NULL THEN '' ELSE description_resolved END description_resolved
		, start_ts
		, end_ts
		, identification_ts
		, root_cause
		, meeting_url
		, post_mortem_url
		, status
		, product
		, CASE WHEN severity_level IS NULL THEN 0 ELSE severity_level END AS severity_level
		, CASE WHEN channel_name IS NULL THEN '' ELSE channel_name END AS channel_name
		, CASE WHEN channel_id IS NULL THEN '' ELSE channel_id END AS channel_id
		, CASE WHEN commander_id IS NULL THEN '' ELSE commander_id END commander_id
		, CASE WHEN commander_email IS NULL THEN '' ELSE commander_email END commander_email
	FROM incident
	WHERE status IN ($1, $2)
	LIMIT 100`
}
