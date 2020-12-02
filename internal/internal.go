package internal

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
	"hellper/internal/log"
	"hellper/internal/log/zap"
	"hellper/internal/model"
	"hellper/internal/model/sql"
	"hellper/internal/model/sql/postgres"
)

func New() (log.Logger, bot.Client, model.IncidentRepository, filestorage.Driver, calendar.Calendar) {
	ctx := context.Background()
	logger := NewLogger()
	return logger, NewClient(logger), NewRepository(logger), NewFileStorage(logger), NewCalendar(ctx, logger)
}

func NewLogger() log.Logger {
	return zap.NewDefault()
}

func NewClient(logger log.Logger) bot.Client {
	return slack.NewClient(config.Env.OAuthToken)
}

func NewRepository(logger log.Logger) model.IncidentRepository {
	fmt.Printf("Configured database: %s", config.Env.Database)
	switch config.Env.Database {
	case "postgres":
		db := sql.NewDBWithDSN(config.Env.Database, config.Env.DSN)
		return postgres.NewIncidentRepository(logger, db)
	default:
		panic(fmt.Sprintf(
			"invalid database option: option=%s valid_options=[postgres]",
			config.Env.Database,
		))
	}
}

// NewFileStorage creates a new connection with the file storage for postmortem document
func NewFileStorage(logger log.Logger) filestorage.Driver {
	fileStorage := config.Env.FileStorage
	switch fileStorage {
	case "google_drive":
		return googledrive.NewFileStorage(logger)
	default:
		panic(fmt.Sprintf(
			"invalid file storage option: option=%s valid_options=[google_drive]",
			fileStorage,
		))
	}
}

// NewCalendar creates a new connection with the calendar service
func NewCalendar(ctx context.Context, logger log.Logger) calendar.Calendar {
	var (
		calendarToken = config.Env.GoogleCalendarToken
		calendarID    = config.Env.GoogleCalendarID
	)
	calendar, err := googlecalendar.NewCalendar(ctx, logger, calendarToken, calendarID)
	if err != nil {
		logger.Error(
			ctx,
			"internal.NewCalendar ERROR",
			log.NewValue("error", err),
		)

		return nil
	}

	return calendar
}
