package entity

import "github.com/google/uuid"

type Vote struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PollID  uuid.UUID `json:"poll_id" gorm:"type:uuid"`
	FieldID uint      `json:"field_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func NewVote(PollID string, fieldID uint, userId string) *Vote {
	pollid, _ := uuid.Parse(PollID)
	userid, _ := uuid.Parse(userId)

	return &Vote{
		ID:      uuid.New(),
		PollID:  pollid,
		FieldID: fieldID,
		UserID:  userid,
	}
}
