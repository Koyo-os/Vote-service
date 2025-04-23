package repository

import (
	"github.com/Koyo-os/Vote-service/internal/entity"
	"go.uber.org/zap"
)

func (repoImpl *VoteRepositoryImpl) Add(vote *entity.Vote) error {
	res := repoImpl.db.Create(vote)
	if err := res.Error; err != nil {
		repoImpl.logger.Error("error add vote to db", zap.String("err", err.Error()))

		return err
	}

	repoImpl.logger.Debug("vote successfully added", zap.String("vote_id", vote.ID.String()))

	return nil
}
