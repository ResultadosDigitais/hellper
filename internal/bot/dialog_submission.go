package bot

type DialogSubmission struct {
	Type        string     `json:"type"`
	Token       string     `json:"token"`
	ActionTs    string     `json:"action_ts"`
	Team        Team       `json:"team"`
	User        User       `json:"user"`
	Channel     Channel    `json:"channel"`
	Submission  Submission `json:"submission"`
	CallbackID  string     `json:"callback_id"`
	ResponseURL string     `json:"response_url"`
	State       string     `json:"state"`
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

type Submission struct {
	IncidentTitle       string `json:"incident_title"`
	ChannelName         string `json:"channel_name"`
	IncidentRoomURL     string `json:"incident_room_url"`
	SeverityLevel       string `json:"severity_level"`
	Responsibility      string `json:"responsibility"`
	Product             string `json:"product"`
	IncidentCommander   string `json:"incident_commander"`
	IncidentDescription string `json:"incident_description"`
	InitDate            string `json:"init_date"`
	IdentificationDate  string `json:"identification_date"`
	EndDate             string `json:"end_date"`
	TimeZone            string `json:"time_zone"`
	Feature             string `json:"feature"`
	Team                string `json:"owner_team"`
	RootCause           string `json:"root_cause"`
	Impact              string `json:"impact"`
	StatusIO            string `json:"status_io"`
	PostMortemMeeting   string `json:"post_mortem_meeting"`
	PauseNotifyTime     string `json:"pause_notify_time"`
	PauseNotifyReason   string `json:"pause_notify_reason"`
}
