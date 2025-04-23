package repository

import (
	"context"
	"errors"

	"github.com/Koyo-os/Vote-service/internal/entity"
	"go.uber.org/zap"
)

func (repoImpl *VoteRepositoryImpl) Delete(_ context.Context, uid string) error {
	if len(uid) == 0 {
		return errors.New("uid is empty")
	}

	res := repoImpl.db.Delete(&entity.Vote{}, uid)
	if err := res.Error; err != nil {
		repoImpl.logger.Error("error delete vote", zap.String("vote_id", uid), zap.Error(err))

		return err
	}

	repoImpl.logger.Info("successfully deleted vote", zap.String("vote_id", uid))

	return nil
}
