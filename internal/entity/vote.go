package entity

import "github.com/google/uuid"

type Vote struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PollID  uuid.UUID `json:"poll_id" gorm:"type:uuid"`
	FieldID uint      `json:"field_id"`
}

func NewVote(PollID string, fieldID uint) *Vote {
	pollid, _ := uuid.Parse(PollID)
	return &Vote{
		ID:      uuid.New(),
		PollID:  pollid,
		FieldID: fieldID,
	}
}
