package config

import (
	"strings"

	"github.com/paked/configure"
)

//Env contains operating system environment variables.
var Env = newEnvironment()

type environment struct {
	OAuthToken            string
	VerificationToken     string
	ProductChannelID      string
	ProductList           string
	Language              string
	MatrixHost            string
	SupportTeam           string
	Messages              messages
	BindAddress           string
	Database              string
	DSN                   string
	GoogleCredentials     string
	GoogleDriveToken      string
	GoogleDriveFileId     string
	GoogleCalendarToken   string
	GoogleCalendarID      string
	ReminderStatusSeconds int
	Environment           string
	FileStorage           string
	NotifyOnResolve       bool
	NotifyOnClose         bool
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

	vars.StringVar(&env.Language, "hellper_language", "pt-br", "Hellper languague")
	vars.StringVar(&env.BindAddress, "hellper_bind_address", ":8080", "Hellper local bind address")
	vars.StringVar(&env.MatrixHost, "hellper_matrix_host", "", "Matrix host")
	vars.StringVar(&env.SupportTeam, "hellper_support_team", "", "Support team identifier")
	vars.StringVar(&env.OAuthToken, "hellper_oauth_token", "", "Token to execute oauth actions")
	vars.StringVar(&env.VerificationToken, "hellper_verification_token", "", "Token to verify external requests")
	vars.StringVar(&env.ProductChannelID, "hellper_product_channel_id", "", "The Product channel id")
	vars.StringVar(&env.ProductList, "hellper_product_list", "Product A;Product B;Product C;Product D", "List of all products splitted by semicolon")
	vars.StringVar(&env.Database, "hellper_database", "postgres", "Hellper database provider")
	vars.StringVar(&env.DSN, "hellper_dsn", "", "Hellper database provider")
	vars.StringVar(&env.GoogleCredentials, "hellper_google_credentials", "", "Google Credentials")
	vars.StringVar(&env.GoogleDriveToken, "hellper_google_drive_token", "", "Google Drive Token")
	vars.StringVar(&env.GoogleDriveFileId, "hellper_google_drive_file_id", "", "Google Drive FileId")
	vars.StringVar(&env.GoogleCalendarToken, "hellper_google_calendar_token", "", "Google Calendar Token")
	vars.StringVar(&env.GoogleCalendarID, "hellper_google_calendar_id", "", "Calendar ID to create a event")
	vars.IntVar(&env.ReminderStatusSeconds, "hellper_reminder_status_seconds", 7200, "Contains the time for the stat reminder to be triggered, by default the time is 2 hours if there is no variable")
	vars.StringVar(&env.Environment, "hellper_environment", "", "Hellper current environment")
	vars.StringVar(&env.FileStorage, "file_storage", "google_drive", "Hellper file storage for postmortem document")
	vars.BoolVar(&env.NotifyOnResolve, "hellper_notify_on_resolve", true, "Notify the main channel when resolve the incident")
	vars.BoolVar(&env.NotifyOnClose, "hellper_notify_on_close", true, "Notify the main channel when close the incident")

	vars.Parse()
	env.Messages = newMessages(env.Language)
	return env
}

func newMessages(language string) messages {
	language = strings.ToLower(language)
	if language == "pt-br" {
		return messages{
			IncidentClosed:         "Incidente <#%s> encerrado.",
			NoListOpenIncidents:    "Não há incidentes ativos!",
			NoTimelineItems:        "Não há items na timeline",
			AnswerAnIncident:       "Resposta a incidente: %s",
			IncidentChannelCreated: "Canal do incidente criado: <#%s>",
			BotHelp: `
hellper
Um bot para ajudar no tratamento de incidentes

Comandos disponíves:
 help      Mostra esta mensagem
 ping      Testa a conectividade com o bot
 list      Lista os incidentes ativos
 state     Mostra a situação e a linha do tempo do incidente
`,
		}
	}
	if language == "en" {
		return messages{
			IncidentClosed:         "Incident <#%s> closed.",
			NoListOpenIncidents:    "No active incidents!",
			NoTimelineItems:        "No timeline items",
			AnswerAnIncident:       "Response to incident: %s",
			IncidentChannelCreated: "Channel of incident: <#%s>",
			BotHelp: `
hellper
A bot to help the incident treatment

Available commands:
 help      Show this help
 ping      Test bot connectivity
 list      List all active incidents
 state     Show incident state and timeline
`,
		}
	}
	if language == "es" {
		return messages{
			IncidentClosed:         "Incidente <#%s> cerrado.",
			NoListOpenIncidents:    "No hay incidentes activos!",
			NoTimelineItems:        "No timeline items",
			AnswerAnIncident:       "Respuesta a incidente: %s",
			IncidentChannelCreated: "Canal del incidente: <#%s>",
			BotHelp: `
hellper
A bot to help the incident treatment

Available commands:
 help      Show this help
 ping      Test bot connectivity
 list      List all active incidents
 state     Show incident state and timeline
`,
		}
	}
	return messages{}
}
