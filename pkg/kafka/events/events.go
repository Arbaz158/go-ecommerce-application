package events

import (
	"encoding/json"
	"time"
)

const (
	UserSignedUpEvent = "user.signed_up"
)

type UserSignedUp struct {
	EventType string    `json:"event_type"`
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *UserSignedUp) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *UserSignedUp) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}
