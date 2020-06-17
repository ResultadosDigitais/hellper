package commands

import (
	"context"
	"hellper/internal/concurrence"
	"strings"
	"sync"
	"time"

	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

var patternStringDate = "2013-04-01 22:43"

// ResolveIncidentDialog opens a dialog on Slack, so the user can resolve an incident
func ResolveIncidentDialog(client bot.Client, triggerID string) (*slack.ViewResponse, error) {

	description := slack.PlainTextInputBlockElement{
		Type:         "plain_text_input",
		ActionID:     "incident_description",
		Placeholder:  slack.NewTextBlockObject("plain_text", "Description", false, false),
		InitialValue: "Description eg. The incident was resolved after #PR fix",
		Multiline:    true,
		MaxLength:    500,
	}

	statusIO := slack.PlainTextInputBlockElement{
		Type:     "plain_text_input",
		ActionID: "status_io",
		Placeholder: &slack.TextBlockObject{
			Type: "plain_text",
			Text: "Status.io link",
		},
		InitialValue: "status.io/xxxx",
	}

	// postMortemMeeting := &slack.RadioButtonsBlockElement{
	// 	Type:     "radio_buttons",
	// 	ActionID: "post_mortem_meeting",
	// 	Options: []*slack.OptionBlockObject{
	// 		{
	// 			Text: &slack.TextBlockObject{
	// 				Type: "plain_text",
	// 				Text: "Yes",
	// 			},
	// 			Value: "true",
	// 		},
	// 		{
	// 			Text: &slack.TextBlockObject{
	// 				Type: "plain_text",
	// 				Text: "No",
	// 			},
	// 			Value: "false",
	// 		},
	// 	},
	// 	InitialOption: &slack.OptionBlockObject{
	// 		Text: &slack.TextBlockObject{
	// 			Type: "plain_text",
	// 			Text: "No",
	// 		},
	// 		Value: "false",
	// 	},
	// }

	summaryMeeting := &slack.PlainTextInputBlockElement{
		Type:     "plain_text_input",
		ActionID: "summary_meeting",
		Placeholder: &slack.TextBlockObject{
			Type: "plain_text",
			Text: "Summary",
		},
		InitialValue: "Post Mortem - inc-xxxx",
	}

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewDividerBlock(),
			slack.NewInputBlock(description.ActionID, description.Placeholder, description),
			slack.NewInputBlock(statusIO.ActionID, statusIO.Placeholder, statusIO),
			// slack.NewActionBlock(postMortemMeeting.ActionID, postMortemMeeting),
			slack.NewInputBlock(summaryMeeting.ActionID, summaryMeeting.Placeholder, summaryMeeting),
		},
	}

	mod := slack.ModalViewRequest{
		Type:          slack.VTModal,
		CallbackID:    "inc-resolve",
		Title:         slack.NewTextBlockObject("plain_text", "Resolve an Incidente", false, false),
		Submit:        slack.NewTextBlockObject("plain_text", "Resolve", false, false),
		NotifyOnClose: false,
		Blocks:        blocks,
	}

	return client.OpenView(triggerID, mod)
}

// ResolveIncidentByDialog resolves an incident after receiving data from a Slack dialog
func ResolveIncidentByDialog(ctx context.Context, client bot.Client, logger log.Logger, repository model.Repository, incidentDetails bot.DialogSubmission) error {
	logger.Info(
		ctx,
		"command/resolve.ResolveIncidentByDialog",
		log.NewValue("incident_resolve_details", incidentDetails),
	)

	var (
		now              = time.Now().UTC()
		channelID        = incidentDetails.Channel.ID
		userID           = incidentDetails.User.ID
		userName         = incidentDetails.User.Name
		submissions      = incidentDetails.Submission
		description      = submissions.IncidentDescription
		statusPageURL    = submissions.StatusIO
		notifyOnResolve  = config.Env.NotifyOnResolve
		productChannelID = config.Env.ProductChannelID
	)

	incident := model.Incident{
		ChannelId:           channelID,
		EndTimestamp:        &now,
		DescriptionResolved: description,
		StatusPageUrl:       statusPageURL,
	}

	err := repository.ResolveIncident(ctx, &incident)
	if err != nil {
		return err
	}

	channelAttachment := createResolveChannelAttachment(incident, userName)
	privateAttachment := createResolvePrivateAttachment(incident)

	var waitgroup sync.WaitGroup
	defer waitgroup.Wait()

	concurrence.WithWaitGroup(&waitgroup, func() {
		postAndPinMessage(client, channelID, "", channelAttachment)
	})
	if notifyOnResolve {
		concurrence.WithWaitGroup(&waitgroup, func() {
			postAndPinMessage(client, productChannelID, "", channelAttachment)
		})
	}
	postMessage(client, userID, "", privateAttachment)

	return nil
}

func createResolveChannelAttachment(inc model.Incident, userName string) slack.Attachment {
	endDateText := inc.EndTimestamp.Format(time.RFC3339)

	var messageText strings.Builder
	messageText.WriteString("The Incident <#" + inc.ChannelId + "> has been resolved by <@" + userName + ">\n\n")
	messageText.WriteString("*End date:* <#" + endDateText + ">\n")
	messageText.WriteString("*Status.io link:* `" + inc.StatusPageUrl + "`\n")
	messageText.WriteString("*Description:* `" + inc.DescriptionResolved + "`\n\n")

	return slack.Attachment{
		Pretext:  "The Incident <#" + inc.ChannelId + "> has been resolved by <@" + userName + ">",
		Fallback: messageText.String(),
		Text:     "",
		Color:    "#1164A3",
		Fields: []slack.AttachmentField{
			{
				Title: "Incident",
				Value: "<#" + inc.ChannelId + ">",
			},
			{
				Title: "End date",
				Value: endDateText,
			},
			{
				Title: "Status.io link",
				Value: inc.StatusPageUrl,
			},
			{
				Title: "Description",
				Value: inc.DescriptionResolved,
			},
		},
	}
}

func createResolvePrivateAttachment(inc model.Incident) slack.Attachment {
	var privateText strings.Builder
	privateText.WriteString("The Incident <#" + inc.ChannelId + "> has been resolved by you\n\n")
	privateText.WriteString("*Status.io:* Be sure to update the incident status on" + inc.StatusPageUrl + "\n")
	privateText.WriteString("*Post Mortem:* Don't forget to bookmark Post Mortem for the incident <#" + inc.ChannelId + ">\n\n")

	return slack.Attachment{
		Pretext:  "The Incident <#" + inc.ChannelId + "> has been resolved by you",
		Fallback: privateText.String(),
		Text:     "",
		Color:    "#FE4D4D",
		Fields: []slack.AttachmentField{
			{
				Title: "Status.io",
				Value: "Be sure to update the incident status on " + inc.StatusPageUrl,
			},
			{
				Title: "Post Mortem",
				Value: "Don't forget to bookmark Post Mortem for the incident <#" + inc.ChannelId + ">",
			},
		},
	}
}
