package repository

import (
	"context"

	"github.com/Koyo-os/Vote-service/internal/entity"
	"github.com/google/uuid"
)

func (repoImpl *VoteRepositoryImpl) GetByUserAndPollID(ctx context.Context, userID, pollID uuid.UUID) ([]entity.Vote, error) {
	var votes []entity.Vote

	res := repoImpl.db.Where(&entity.Vote{
		PollID: pollID,
		UserID: userID,
	}).Find(&votes)

	return votes, res.Error
}
