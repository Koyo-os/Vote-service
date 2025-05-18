package repository

import (
	"context"

	"github.com/Koyo-os/Vote-service/internal/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (repoImpl *VoteRepositoryImpl) GetByUserAndPollID(
	ctx context.Context,
	userID, pollID uuid.UUID,
) ([]entity.Vote, error) {
	var votes []entity.Vote

	res := repoImpl.db.Where(&entity.Vote{
		PollID: pollID,
		UserID: userID,
	}).Find(&votes)

	return votes, res.Error
}

func (repoImpl *VoteRepositoryImpl) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.Vote, error) {
	var vote entity.Vote

	res := repoImpl.db.Where(&entity.Event{
		ID: id.String(),
	}).Find(&vote)

	if err := res.Error; err != nil {
		repoImpl.logger.Error(
			"error get vote from db",
			zap.String("vote_id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	return &vote, nil
}
