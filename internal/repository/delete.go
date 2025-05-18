package repository

import (
	"context"

	"github.com/Koyo-os/Vote-service/internal/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (repo *VoteRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	res := repo.db.Delete(&entity.Vote{}, id)

	if err := res.Error; err != nil {
		repo.logger.Error("error delete vote",
			zap.String("vote_id", id.String()),
			zap.Error(err))

		return err
	}

	return nil
}

func (repo *VoteRepositoryImpl) DeleteByPollID(
	ctx context.Context,
	pollID uuid.UUID,
) ([]entity.Vote, error) {
	var votes []entity.Vote

	res := repo.db.Where(&entity.Vote{
		PollID: pollID,
	}).Find(votes)
	if err := res.Error; err != nil {
		repo.logger.Error("error get votes",
			zap.String("poll_id", pollID.String()),
			zap.Error(err))

		return nil, err
	}

	res = repo.db.Where(&entity.Vote{
		PollID: pollID,
	}).Delete(&entity.Vote{})

	if err := res.Error; err != nil {
		repo.logger.Error("error delete votes",
			zap.String("poll_id", pollID.String()),
			zap.Error(err))
		return nil, err
	}

	return votes, nil
}
