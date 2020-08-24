package notify

import (
	"context"
	"errors"
	"flag"
	"hellper/internal"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

var (
	logger     = internal.NewLogger()
	client     = internal.NewClient(logger)
	repository = internal.NewRepository(logger)

	arg opt
)

type opt struct {
	typeFlag   string
	toFlag     string
	msgFlag    string
	statusFlag string
}

func init() {
	var typeFlag, toFlag, msgFlag, statusFlag string
	flag.StringVar(&typeFlag, "type", "", "[channels|report|custom]")
	flag.StringVar(&toFlag, "to", "", "[channel id|user id]")
	flag.StringVar(&msgFlag, "msg", "", "A text message")
	flag.StringVar(&statusFlag, "status", "", "[all|open|resolved]")
	flag.Parse()

	arg = opt{typeFlag, toFlag, msgFlag, statusFlag}
}

// Notify is a CLI responsible for sending messages
func Notify(ctx context.Context) {
	logger.Info(ctx, log.Trace(), log.Action("running"))

	switch arg.typeFlag {
	case "channels":
		channelsNotify(ctx)
	case "report":
		reportNotify(ctx)
	case "custom":
		customNotify(ctx)
	default:
		logger.Error(ctx, log.Trace(), log.NewValue("error", errors.New("Must have a type")))
		return
	}

}

func send(to, msg string) error {

	if to == "" {
		return errors.New("Must have a destination")
	}

	if msg == "" {
		return errors.New("Must have a message")
	}

	_, _, err := client.PostMessage(to, slack.MsgOptionText(msg, false))
	if err != nil {
		return err
	}

	return nil

}

func statusNotify(incident model.Incident) string {
	switch incident.Status {
	case model.StatusOpen:
		return config.Env.ReminderOpenNotifyMsg
	case model.StatusResolved:
		return config.Env.ReminderResolvedNotifyMsg
	}
	return ""
}
