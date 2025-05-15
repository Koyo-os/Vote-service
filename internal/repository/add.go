package repository

import (
	"context"

	"github.com/Koyo-os/Vote-service/internal/entity"
	"github.com/Koyo-os/Vote-service/internal/service"
	"go.uber.org/zap"
)

var _ service.Repostory = &VoteRepositoryImpl{}

func (repoImpl *VoteRepositoryImpl) Add(ctx context.Context, vote *entity.Vote) error {
	res := repoImpl.db.Create(vote)
	if err := res.Error; err != nil {
		repoImpl.logger.Error("error add vote to db", zap.String("err", err.Error()))

		return err
	}

	repoImpl.logger.Debug("vote successfully added", zap.String("vote_id", vote.ID.String()))

	return nil
}
