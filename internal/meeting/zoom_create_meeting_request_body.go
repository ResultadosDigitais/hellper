package meeting

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const InstantMeetingType = 1
const NoRegistrationRequired = 2

type zoomCreateMeetingInputSettingsPayload struct {
	HostVideo        bool `json:"host_video"`
	ParticipantVideo bool `json:"participant_video"`
	JoinBeforeHost   bool `json:"join_before_host"`
	MuteUponEntry    bool `json:"mute_upon_entry"`
	ApprovalType     int  `json:"approval_type"`
	WaitingRoom      bool `json:"waiting_room"`
}

type zoomCreateMeetingInputPayload struct {
	Topic    string                                `json:"topic"`
	Agenda   string                                `json:"agenda"`
	Type     int                                   `json:"type"`
	Settings zoomCreateMeetingInputSettingsPayload `json:"settings"`
}

func (provider zoomProvider) createMeetingRequestBody(channel string) *bytes.Buffer {
	postData := zoomCreateMeetingInputPayload{
		Topic:  fmt.Sprintf("Incident reported on #%s", channel),
		Agenda: fmt.Sprintf("Meeting for incident resolution reported on #%s", channel),
		Type:   InstantMeetingType,
		Settings: zoomCreateMeetingInputSettingsPayload{
			HostVideo:        false,
			ParticipantVideo: false,
			JoinBeforeHost:   true,
			MuteUponEntry:    true,
			ApprovalType:     NoRegistrationRequired,
			WaitingRoom:      false,
		},
	}

	jsonValue, _ := json.Marshal(postData)

	return bytes.NewBuffer(jsonValue)
}
