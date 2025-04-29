package service

import (
	"context"

	"github.com/Koyo-os/Vote-service/internal/entity"
	"github.com/google/uuid"
)

type (
	Repostory interface {
		GetByUserAndPollID(context.Context, uuid.UUID, uuid.UUID) ([]entity.Vote, error)
		Add(context.Context, *entity.Vote) error
		Delete(context.Context, uuid.UUID) error
		GetByID(context.Context, uuid.UUID) (*entity.Vote, error)
	}

	VoidPublisher interface {
		Publish(context.Context, any, string) error
	}

	ErrHandler interface {
		HandleError(error) error
	}
	Casher interface {
		RemoveFromCash(context.Context, string) error
		AddToCash(context.Context, string, any) error
	}
)
