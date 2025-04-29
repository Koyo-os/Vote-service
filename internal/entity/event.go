package entity

import "github.com/google/uuid"

type Event struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

func NewEvent(Type string, payload []byte) *Event {
	return &Event{
		ID:      uuid.New().String(),
		Type:    Type,
		Payload: payload,
	}
}
