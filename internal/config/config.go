package config

import (
	"github.com/paked/configure"
)

//Env contains operating system environment variables.
var Env = newEnvironment()

type environment struct {
	OAuthToken                    string
	SlackSigningSecret            string
	ProductChannelID              string
	ProductList                   string
	Language                      string
	MatrixHost                    string
	SupportTeam                   string
	Messages                      messages
	BindAddress                   string
	Database                      string
	DSN                           string
	GoogleCredentials             string
	GoogleDriveToken              string
	GoogleDriveFileID             string
	GoogleCalendarToken           string
	GoogleCalendarID              string
	PostmortemGapDays             int
	ReminderOpenStatusSeconds     int
	ReminderResolvedStatusSeconds int
	ReminderOpenNotifyMsg         string
	ReminderResolvedNotifyMsg     string
	Environment                   string
	FileStorage                   string
	NotifyOnResolve               bool
	NotifyOnClose                 bool
	NotifyOnCancel                bool
	Timezone                      string
	SLAHoursToClose               int
}

type messages struct {
	IncidentClosed         string
	NoListOpenIncidents    string
	NoTimelineItems        string
	AnswerAnIncident       string
	IncidentChannelCreated string
	BotHelp                string
}

func newEnvironment() environment {
	var (
		vars = configure.New(configure.NewEnvironment())
		env  environment
	)

	vars.StringVar(&env.BindAddress, "hellper_bind_address", ":8080", "Hellper local bind address")
	vars.StringVar(&env.MatrixHost, "hellper_matrix_host", "", "Matrix host")
	vars.StringVar(&env.SupportTeam, "hellper_support_team", "", "Support team identifier")
	vars.StringVar(&env.OAuthToken, "hellper_oauth_token", "", "Token to execute oauth actions")
	vars.StringVar(&env.SlackSigningSecret, "hellper_slack_signing_secret", "", "Slack signs the requests confirm that each request comes from Slack by verifying its unique signature")
	vars.StringVar(&env.ProductChannelID, "hellper_product_channel_id", "", "The Product channel id")
	vars.StringVar(&env.ProductList, "hellper_product_list", "Product A;Product B;Product C;Product D", "List of all products splitted by semicolon")
	vars.StringVar(&env.Database, "hellper_database", "postgres", "Hellper database provider")
	vars.StringVar(&env.DSN, "hellper_dsn", "", "Hellper database provider")
	vars.StringVar(&env.GoogleCredentials, "hellper_google_credentials", "", "Google Credentials")
	vars.StringVar(&env.GoogleDriveToken, "hellper_google_drive_token", "", "Google Drive Token")
	vars.StringVar(&env.GoogleDriveFileID, "hellper_google_drive_file_id", "", "Google Drive FileId")
	vars.StringVar(&env.GoogleCalendarToken, "hellper_google_calendar_token", "", "Google Calendar Token")
	vars.StringVar(&env.GoogleCalendarID, "hellper_google_calendar_id", "", "Calendar ID to create a event")
	vars.IntVar(&env.PostmortemGapDays, "hellper_postmortem_gap_days", 2, "Gap in days between resolve and postmortem event")
	vars.IntVar(&env.ReminderOpenStatusSeconds, "hellper_reminder_open_status_seconds", 7200, "Contains the time for the stat reminder to be triggered when status is open, by default the time is 2 hours if there is no variable")
	vars.IntVar(&env.ReminderResolvedStatusSeconds, "hellper_reminder_resolved_status_seconds", 86400, "Contains the time for the stat reminder to be triggered when status is resolved, by default the time is 24 hours if there is no variable")
	vars.StringVar(&env.ReminderOpenNotifyMsg, "hellper_reminder_open_notify_msg", "Incident Status: Open - Update the status of this incident, just pin a message with status on the channel.", "Notify message when status is open")
	vars.StringVar(&env.ReminderResolvedNotifyMsg, "hellper_reminder_resolved_notify_msg", "Incident Status: Resolved - Update the status of this incident, just pin a message with status on the channel.", "Notify message when status is resolved")
	vars.StringVar(&env.Environment, "hellper_environment", "", "Hellper current environment")
	vars.StringVar(&env.FileStorage, "file_storage", "google_drive", "Hellper file storage for postmortem document")
	vars.BoolVar(&env.NotifyOnResolve, "hellper_notify_on_resolve", true, "Notify the Product channel when resolve the incident")
	vars.BoolVar(&env.NotifyOnClose, "hellper_notify_on_close", true, "Notify the Product channel when close the incident")
	vars.BoolVar(&env.NotifyOnCancel, "hellper_notify_on_cancel", true, "Notify the Product channel when cancel the incident")
	vars.StringVar(&env.Timezone, "timezone", "America/Sao_Paulo", "The local time of a region or a country used to create a event.")
	vars.IntVar(&env.SLAHoursToClose, "hellper_sla_hours_to_close", 168, "SLA hours to close")

	vars.Parse()
	return env
}
