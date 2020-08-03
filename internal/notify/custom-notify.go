package notify

import (
	"context"
	"hellper/internal/log"
)

func customNotify(ctx context.Context) {
	logger.Info(ctx, log.Trace(), log.Action("running"))

	err := send(arg.toFlag, arg.msgFlag)
	if err != nil {
		logger.Error(ctx, log.Trace(), log.NewValue("error", err))
	}

	logger.Info(ctx, log.Trace(), log.Action("done"))

}
