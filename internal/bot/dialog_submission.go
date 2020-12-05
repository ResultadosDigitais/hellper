package bot

type DialogSubmission struct {
	Type        string            `json:"type"`
	Token       string            `json:"token"`
	ActionTs    string            `json:"action_ts"`
	Team        Team              `json:"team"`
	User        User              `json:"user"`
	Channel     Channel           `json:"channel"`
	Submission  map[string]string `json:"submission"`
	CallbackID  string            `json:"callback_id"`
	ResponseURL string            `json:"response_url"`
	State       string            `json:"state"`
}

type Team struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
