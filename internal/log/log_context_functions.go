package log

import (
	"context"
	"time"
)

type ContextFunction func(context.Context) []Value

func ResolveContextFunctions(ctx context.Context, functions ...ContextFunction) []Value {
	values := make([]Value, len(functions))
	for _, function := range functions {
		values = append(values, function(ctx)...)
	}
	return values
}

func NewTIDContextFunction(key string) ContextFunction {
	return func(ctx context.Context) []Value {
		uuid := TIDFromContext(ctx)
		return []Value{
			{Name: key, Value: uuid},
		}
	}
}

func TimestampContextFunction(key, layout string) ContextFunction {
	return func(ctx context.Context) []Value {
		return []Value{
			{Name: key, Value: time.Now().Format(layout)},
		}
	}
}
