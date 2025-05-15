package repository

import (
	"github.com/Koyo-os/Vote-service/pkg/logger"
	"gorm.io/gorm"
)

type VoteRepositoryImpl struct {
	db     *gorm.DB
	logger *logger.Logger
}

func Init(db *gorm.DB, logger *logger.Logger) *VoteRepositoryImpl {
	return &VoteRepositoryImpl{
		db:     db,
		logger: logger,
	}
}
