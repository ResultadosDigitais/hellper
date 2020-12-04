package log

import (
	"context"
	"runtime"
)

const (
	//STDOUT any message to stdout
	STDOUT Out = "stdout"

	//ERROR is the error level logger
	ERROR Level = "error"
	// WARN is the warning level logger
	WARN Level = "warn"
	//INFO is the info level logger
	INFO Level = "info"
	//DEBUG is the debug level logger
	DEBUG Level = "debug"
)

//Out is the type for logger writer config
type Out string

func (o Out) String() string {
	return string(o)
}

// Set is a utility method for flag system usage
func (o *Out) Set(value string) error {
	switch value {
	case "stdout", "STDOUT", "":
		*o = STDOUT
	default:
		*o = Out(value)
	}
	return nil
}

//Level is the threshold of the logger
type Level string

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	return string(l)
}

// Set is a utility method for flag system usage
func (l *Level) Set(value string) error {
	switch value {
	case "info", "INFO":
		*l = INFO
	case "error", "ERROR":
		*l = ERROR
	default:
		*l = DEBUG
	}
	return nil
}

type Value struct {
	Name  string
	Value interface{}
}

func NewValue(name string, value interface{}) Value {
	return Value{Name: name, Value: value}
}

type Logger interface {
	Debug(context.Context, string, ...Value)
	Info(context.Context, string, ...Value)
	Warn(context.Context, string, ...Value)
	Error(context.Context, string, ...Value)
	With(...Value) Logger
}

func Trace() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}

func Action(action string) Value {
	return NewValue("action", action)
}

func Reason(reason string) Value {
	return NewValue("reason", reason)
}
