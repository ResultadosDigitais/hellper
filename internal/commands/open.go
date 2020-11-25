package commands

import (
	"context"
	"fmt"
	"hellper/internal/concurrence"
	"strconv"
	"strings"
	"sync"
	"time"

	"hellper/internal/bot"
	"hellper/internal/config"
	filestorage "hellper/internal/file_storage"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

// OpenStartIncidentDialog opens a dialog on Slack, so the user can start an incident
func OpenStartIncidentDialog(client bot.Client, triggerID string) error {
	productList := []slack.DialogSelectOption{}

	for _, product := range strings.Split(config.Env.ProductList, ";") {
		productList = append(productList, slack.DialogSelectOption{
			Label: product,
			Value: product,
		})
	}

	incidentTitle := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Incident Title",
			Name:        "incident_title",
			Type:        "text",
			Placeholder: "My Incident Title",
		},
		MaxLength: 100,
	}

	channelName := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Channel name",
			Name:        "channel_name",
			Type:        "text",
			Placeholder: "inc-my-incident",
		},
		MaxLength: 30,
	}

	meeting := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "War Room URL",
			Name:        "war_room_url",
			Type:        "text",
			Placeholder: "War Room URL eg. Matrix/Meeting/Zoom",
			Optional:    true,
		},
	}

	severityLevel := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Severity level",
			Name:        "severity_level",
			Type:        "select",
			Placeholder: "Set the severity level",
			Optional:    false,
		},
		Options: []slack.DialogSelectOption{
			{
				Label: "SEV0 - All hands on deck",
				Value: "0",
			},
			{
				Label: "SEV1 - Critical impact to many users",
				Value: "1",
			},
			{
				Label: "SEV2 - Minor issue that impacts ability to use product",
				Value: "2",
			},
			{
				Label: "SEV3 - Minor issue not impacting ability to use product",
				Value: "3",
			},
		},
		OptionGroups: []slack.DialogOptionGroup{},
	}

	product := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Product",
			Name:        "product",
			Type:        "select",
			Placeholder: "Set the product",
			Optional:    false,
		},
		Options:      productList,
		OptionGroups: []slack.DialogOptionGroup{},
	}

	commander := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Incident commander",
			Name:        "incident_commander",
			Type:        "select",
			Placeholder: "Set the Incident commander",
			Optional:    false,
		},
		DataSource:   "users",
		OptionGroups: []slack.DialogOptionGroup{},
	}

	description := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Incident description",
			Name:        "incident_description",
			Type:        "textarea",
			Placeholder: "Incident description eg. We're having a delay email campaign delivery",
			Optional:    false,
		},
		MaxLength: 500,
	}

	dialog := slack.Dialog{
		CallbackID:     "inc-open",
		Title:          "Start an Incident",
		SubmitLabel:    "Start",
		NotifyOnCancel: false,
		Elements: []slack.DialogElement{
			incidentTitle,
			channelName,
			meeting,
			severityLevel,
			product,
			commander,
			description,
		},
	}

	return client.OpenDialog(triggerID, dialog)
}

// StartIncidentByDialog starts an incident after receiving data from a Slack dialog
func StartIncidentByDialog(
	ctx context.Context,
	client bot.Client,
	logger log.Logger,
	repository model.Repository,
	fileStorage filestorage.Driver,
	incidentDetails bot.DialogSubmission,
) error {
	logger.Info(
		ctx,
		"command/open.StartIncidentByDialog",
		log.NewValue("incident_open_details", incidentDetails),
	)

	var (
		now              = time.Now().UTC()
		incidentAuthor   = incidentDetails.User.ID
		submission       = incidentDetails.Submission
		incidentTitle    = submission.IncidentTitle
		channelName      = submission.ChannelName
		warRoomURL       = submission.WarRoomURL
		severityLevel    = submission.SeverityLevel
		product          = submission.Product
		commander        = submission.IncidentCommander
		description      = submission.IncidentDescription
		environment      = config.Env.Environment
		matrixURL        = config.Env.MatrixHost
		supportTeam      = config.Env.SupportTeam
		productChannelID = config.Env.ProductChannelID
		stagingRoom      = "dc82e346-639c-44ee-a470-63f7545ae8e4"
	)

	user, err := getSlackUserInfo(ctx, client, logger, commander)
	if err != nil {
		return fmt.Errorf("commands.StartIncidentByDialog.get_slack_user_info: incident=%v commanderId=%v error=%v", channelName, commander, err)
	}

	channel, err := client.CreateConversationContext(ctx, channelName, false)
	if err != nil {
		return fmt.Errorf("commands.StartIncidentByDialog.create_conversation_context: incident=%v error=%v", channelName, err)
	}

	severityLevelInt64, err := getStringInt64(severityLevel)
	if err != nil {
		return err
	}

	incident := model.Incident{
		ChannelName:             channelName,
		ChannelId:               channel.ID,
		Title:                   incidentTitle,
		Product:                 product,
		DescriptionStarted:      description,
		Status:                  model.StatusOpen,
		IdentificationTimestamp: &now,
		SeverityLevel:           severityLevelInt64,
		IncidentAuthor:          incidentAuthor,
		CommanderId:             user.SlackID,
		CommanderEmail:          user.Email,
	}

	incidentID, err := repository.InsertIncident(ctx, &incident)
	if err != nil {
		return err
	}

	if warRoomURL == "" {
		if environment == "production" {
			warRoomURL = matrixURL + "/new?roomId=" + channelName + "&roomName=" + channelName
		} else if environment == "staging" {
			warRoomURL = matrixURL + "/new?roomId=" + stagingRoom + "&roomName=hellper-staging"
		}
	}

	attachment := createOpenAttachment(incident, incidentID, warRoomURL, supportTeam)
	message := "An Incident has been opened by <@" + incident.IncidentAuthor + ">"

	var waitgroup sync.WaitGroup
	defer waitgroup.Wait()

	concurrence.WithWaitGroup(&waitgroup, func() {
		postAndPinMessage(client, channel.ID, message, attachment)
	})
	concurrence.WithWaitGroup(&waitgroup, func() {
		postAndPinMessage(client, productChannelID, message, attachment)
	})

	//We need run that without wait because the modal need close in only 3s
	go createPostMortemAndUpdateTopic(ctx, logger, client, fileStorage, incident, incidentID, repository, channel, warRoomURL)

	// startReminderStatusJob(ctx, logger, client, repository, incident)

	_, err = client.InviteUsersToConversationContext(ctx, channel.ID, commander)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Reason("InviteUsersToConversationContext"),
			log.NewValue("channel.ID", channel.ID),
			log.NewValue("commander", commander),
			log.NewValue("error", err),
		)
		return err
	}

	return nil
}

func createPostMortemAndUpdateTopic(ctx context.Context, logger log.Logger, client bot.Client, fileStorage filestorage.Driver, incident model.Incident, incidentID int64, repository model.Repository, channel *slack.Channel, warRoomURL string) {
	postMortemURL, err := createPostMortem(ctx, logger, client, fileStorage, incidentID, incident.Title, repository, channel.Name)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Reason("createPostMortem"),
			log.NewValue("channel.Name", channel.Name),
			log.NewValue("error", err),
		)
		return
	}

	var topic strings.Builder
	topic.WriteString("*WarRoom:* " + warRoomURL + "\n\n")
	topic.WriteString("*PostMortem:* " + postMortemURL + "\n\n")
	topic.WriteString("*Commander:* <@" + incident.CommanderId + ">\n\n")

	_, err = client.SetTopicOfConversation(channel.ID, topic.String())
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Reason("SetTopicOfConversation"),
			log.NewValue("channel.ID", channel.ID),
			log.NewValue("topic.String", topic.String()),
			log.NewValue("error", err),
		)
	}
}

func createOpenAttachment(incident model.Incident, incidentID int64, warRoomURL string, supportTeam string) slack.Attachment {
	var messageText strings.Builder
	messageText.WriteString("An Incident has been opened by <@" + incident.IncidentAuthor + ">\n\n")
	messageText.WriteString("*Title:* " + incident.Title + "\n")
	messageText.WriteString("*Severity:* " + getSeverityLevelText(incident.SeverityLevel) + "\n\n")
	messageText.WriteString("*Product:* " + incident.Product + "\n")
	messageText.WriteString("*Channel:* <#" + incident.ChannelId + ">\n")
	messageText.WriteString("*Commander:* <@" + incident.CommanderId + ">\n\n")
	messageText.WriteString("*Description:* `" + incident.DescriptionStarted + "`\n\n")
	messageText.WriteString("*War Room:* " + warRoomURL + "\n")
	messageText.WriteString("*cc:* <@" + supportTeam + ">\n")

	return slack.Attachment{
		Pretext:  "*cc:* <!subteam^" + supportTeam + ">",
		Fallback: messageText.String(),
		Text:     "",
		Color:    "#FE4D4D",
		Fields: []slack.AttachmentField{
			{
				Title: "Incident ID",
				Value: strconv.FormatInt(incidentID, 10),
			},
			{
				Title: "Incident Channel",
				Value: "<#" + incident.ChannelId + ">",
			},
			{
				Title: "Incident Title",
				Value: incident.Title,
			},
			{
				Title: "Severity",
				Value: getSeverityLevelText(incident.SeverityLevel),
			},
			{
				Title: "Product",
				Value: incident.Product,
			},
			{
				Title: "Commander",
				Value: "<@" + incident.CommanderId + ">",
			},
			{
				Title: "Description",
				Value: "```" + incident.DescriptionStarted + "```",
			},
			{
				Title: "War Room",
				Value: warRoomURL,
			},
		},
	}
}
