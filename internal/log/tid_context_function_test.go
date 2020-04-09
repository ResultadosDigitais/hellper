package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionFromContext(t *testing.T) {
	var (
		ctx = context.Background()
		tid = TIDFromContext(ctx)
	)
	require.NotZero(t, tid, "invalid tid instance")

	ctx = ContextWithTID(ctx, tid)
	tid2 := TIDFromContext(ctx)
	require.Equal(t, tid2, tid, "invalid tid value")

	tid = ""
	ctx = ContextWithTID(ctx, tid)
	tid = TIDFromContext(ctx)
	require.NotZero(t, tid, "blank tid")
}
