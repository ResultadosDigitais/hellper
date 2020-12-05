package app

import (
	"context"
	"fmt"
	"hellper/internal/bot"
	"hellper/internal/bot/slack"
	"hellper/internal/calendar"
	googlecalendar "hellper/internal/calendar/google_calendar"
	"hellper/internal/config"
	filestorage "hellper/internal/file_storage"
	googledrive "hellper/internal/file_storage/google_drive"
	"hellper/internal/invitation"
	"hellper/internal/log"
	"hellper/internal/log/zap"
	"hellper/internal/model"
	"hellper/internal/model/sql"
	"hellper/internal/model/sql/postgres"
)

type App struct {
	Logger             log.Logger
	Client             bot.Client
	IncidentRepository model.IncidentRepository
	ServiceRepository  model.ServiceRepository
	TeamRepository     model.TeamRepository
	PersonRepository   model.PersonRepository
	FileStorage        filestorage.Driver
	Calendar           calendar.Calendar
	Inviter            invitation.Inviter
}

func NewApp() App {
	ctx := context.Background()
	logger := NewLogger()
	defer logger.Info(ctx, "Application configured")

	db := NewDatabase(ctx, logger)
	teamRepository := NewTeamRepository(ctx, logger, db)
	personRepository := NewPersonRepository(ctx, logger, db)
	client := NewClient(ctx, logger)

	return App{
		Logger:             logger,
		Client:             client,
		IncidentRepository: NewIncidentRepository(ctx, logger, db),
		ServiceRepository:  NewServiceRepository(ctx, logger, db),
		TeamRepository:     teamRepository,
		PersonRepository:   personRepository,
		Inviter:            NewInviter(ctx, logger, client, teamRepository, personRepository),
		FileStorage:        NewFileStorage(ctx, logger),
		Calendar:           NewCalendar(ctx, logger),
	}
}

func NewLogger() log.Logger {
	configuredLogger := config.Env.Logger
	fmt.Printf("internal.NewLogger initializing logger: %s\n", configuredLogger)
	switch configuredLogger {
	case LoggerZap:
		return zap.NewDefault()
	default:
		panic(fmt.Sprintf(
			"internal.NewLogger invalid logger option: option=%s valid_options=[%s]",
			configuredLogger, LoggerZap,
		))
	}
}

func NewClient(ctx context.Context, logger log.Logger) bot.Client {
	configuredClient := config.Env.Client
	logger.Debug(ctx, fmt.Sprintf(
		"internal.NewClient initializing client connection: %s", configuredClient,
	))
	switch configuredClient {
	case ClientSlack:
		return slack.NewClient(config.Env.OAuthToken)
	default:
		panic(fmt.Sprintf(
			"internal.NewClient invalid client option: option=%s valid_options=[%s]\n",
			configuredClient, ClientSlack,
		))
	}
}

// NewDatabase creates a mew database connection
func NewDatabase(ctx context.Context, logger log.Logger) sql.DB {
	configuredDatabase := config.Env.Database
	logger.Info(ctx, fmt.Sprintf(
		"internal.NewDatabase initializing database connection: %s", configuredDatabase,
	))
	switch configuredDatabase {
	case DatabasePostgres:
		return sql.NewDBWithDSN(config.Env.Database, config.Env.DSN)
	default:
		panic(fmt.Sprintf(
			"internal.NewDatabase invalid database option: option=%s valid_options=[%s]\n",
			configuredDatabase, DatabasePostgres,
		))
	}
}

// NewIncidentRepository creates a new connection with the database for incidents
func NewIncidentRepository(ctx context.Context, logger log.Logger, db sql.DB) model.IncidentRepository {
	return postgres.NewIncidentRepository(logger, db)
}

// NewServiceRepository creates a new connection with the database for services
func NewServiceRepository(ctx context.Context, logger log.Logger, db sql.DB) model.ServiceRepository {
	return postgres.NewServiceRepository(logger, db)
}

// NewPersonRepository creates a new connection with the database for persons
func NewPersonRepository(ctx context.Context, logger log.Logger, db sql.DB) model.PersonRepository {
	return postgres.NewPersonRepository(logger, db)
}

// NewTeamRepository creates a new connection with the database for teams
func NewTeamRepository(ctx context.Context, logger log.Logger, db sql.DB) model.TeamRepository {
	return postgres.NewTeamRepository(logger, db)
}

func NewInviter(
	ctx context.Context,
	logger log.Logger,
	client bot.Client,
	teamRepository model.TeamRepository,
	personRepository model.PersonRepository,
) invitation.Inviter {
	logger.Info(ctx, "internal.NewInviter initializing inviter")
	return invitation.NewInviter(logger, client, teamRepository, personRepository)
}

// NewFileStorage creates a new connection with the file storage for postmortem document
func NewFileStorage(ctx context.Context, logger log.Logger) filestorage.Driver {
	configuredFileStorage := config.Env.FileStorage
	logger.Debug(
		ctx, fmt.Sprintf("internal.NewFileStorage initializing file storage connection: %s", configuredFileStorage))
	switch configuredFileStorage {
	case FileStorageGoogleDrive:
		return googledrive.NewFileStorage(logger)
	case FileStorageNone:
		return nil
	default:
		logger.Error(ctx, fmt.Sprintf(
			"internal.NewFileStorage invalid file storage option: option=%s valid_options=[%s,%s]",
			configuredFileStorage, FileStorageGoogleDrive, FileStorageNone,
		))
		return nil
	}
}

// NewCalendar creates a new connection with the calendar service
func NewCalendar(ctx context.Context, logger log.Logger) calendar.Calendar {
	configuredCalendar := config.Env.Calendar
	logger.Debug(
		ctx,
		fmt.Sprintf("internal.NewCalendar initializing calendar connection: %s", configuredCalendar),
	)
	switch configuredCalendar {
	case CalendarGoogle:
		var (
			calendarToken = config.Env.GoogleCalendarToken
			calendarID    = config.Env.GoogleCalendarID
		)
		if !googlecalendar.ValidateParameters(calendarID, calendarToken) {
			logger.Error(
				ctx,
				"internal.NewCalendar[Google] parameters not configured: calendarID/calendarToken",
			)
			return nil
		}
		calendar, err := googlecalendar.NewCalendar(ctx, logger, calendarToken, calendarID)
		if err != nil {
			logger.Error(
				ctx,
				"internal.NewCalendar[Google] ERROR",
				log.NewValue("error", err),
			)
			return nil
		}
		return calendar
	case CalendarNone:
		return nil
	default:
		logger.Error(
			ctx, fmt.Sprintf(
				"internal.NewCalendar invalid calendar option: option=%s valid_options=[%s,%s]",
				configuredCalendar, CalendarGoogle, CalendarNone,
			),
		)
		return nil
	}

}
