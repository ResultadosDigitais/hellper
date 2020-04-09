package commands

// TriggerEvent represents some slack event request parameters
type TriggerEvent struct {
	Type    string
	User    string
	Channel string
}
