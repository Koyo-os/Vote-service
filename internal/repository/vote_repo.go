package repository

import (
	"context"

	"github.com/Koyo-os/Vote-service/internal/entity"
	"github.com/Koyo-os/Vote-service/pkg/config"
	"github.com/Koyo-os/Vote-service/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type VoteRepository interface {
	Add(context.Context, *entity.Vote) error
	Delete(context.Context, string) error
}

type VoteRepositoryImpl struct {
	db     *gorm.DB
	logger *logger.Logger
}

func Init(cfg *config.Config, logger *logger.Logger) (VoteRepository, error) {
	logger.Info("starting connect to db with", zap.String("dsn", cfg.DSN))

	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		logger.Error("failed connect to db", zap.String("err", err.Error()))

		return nil, err
	}

	logger.Info("starting migration")

	if err = db.AutoMigrate(&entity.Event{}); err != nil {
		logger.Error("failed do migrate", zap.String("err", err.Error()))

		return nil, err
	}

	logger.Info("migration successfull")

	return &VoteRepositoryImpl{}, nil
}
