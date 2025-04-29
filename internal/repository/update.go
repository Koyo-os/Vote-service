package repository

import (
	"context"

	"github.com/Koyo-os/Vote-service/internal/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (repoImpl *VoteRepositoryImpl) Delete(_ context.Context, uid uuid.UUID) error {
	res := repoImpl.db.Delete(&entity.Vote{}, uid)
	if err := res.Error; err != nil {
		repoImpl.logger.Error("error delete vote", zap.String("vote_id", uid), zap.Error(err))

		return err
	}

	repoImpl.logger.Info("successfully deleted vote", zap.String("vote_id", uid))

	return nil
}
