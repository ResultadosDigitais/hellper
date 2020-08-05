package main

import (
	"context"
	"hellper/internal/notify"
)

func main() {
	ctx := context.Background()
	notify.Notify(ctx)
}
