package model

// User represents a model of an user on Hellper
type User struct {
	// ID is the user unique identifier on Hellper
	ID      int64
	SlackID string
	Name    string
	// DisplayName represents the Slack handle for an user. This is a virtual field.
	DisplayName string
	Email       string
}
