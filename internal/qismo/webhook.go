package qismo

type WebhookNewSessionRequest struct {
	IsNewSession bool `json:"is_new_session"`
	Payload      struct {
		Room struct {
			ID              string `json:"id"`
			IDStr           string `json:"id_str"`
			IsPublicChannel bool   `json:"is_public_channel"`
			Name            string `json:"name"`
			Options         string `json:"options"`
			Participants    []struct {
				Email string `json:"email"`
			} `json:"participants"`
			RoomAvatar string `json:"room_avatar"`
			TopicID    string `json:"topic_id"`
			TopicIDStr string `json:"topic_id_str"`
			Type       string `json:"type"`
		} `json:"room"`
	} `json:"payload"`
	WebhookType string `json:"webhook_type"`
}
