package log

import (
	"context"

	"github.com/google/uuid"
)

type key int

const tidKey = 0

func ContextWithTID(ctx context.Context, tid string) context.Context {
	if tid == "" {
		tid = uuid.New().String()
	}

	return context.WithValue(ctx, tidKey, tid)
}

func TIDFromContext(ctx context.Context) (tid string) {
	switch v := ctx.Value(tidKey).(type) {
	case string:
		tid = v
	default:
		tid = ""
	}
	if tid == "" {
		tid = uuid.New().String()
	}
	return
}
