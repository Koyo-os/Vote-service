package service

import (
	"context"
	"errors"

	"github.com/Koyo-os/Vote-service/internal/entity"
	"github.com/Koyo-os/Vote-service/pkg/config"
	"github.com/Koyo-os/Vote-service/pkg/errs"
	"github.com/Koyo-os/Vote-service/pkg/retrier"
	"github.com/google/uuid"
)

type Service struct {
	repo      Repostory
	publisher VoidPublisher
	casher    Casher
	errs      ErrHandler

	ctx context.Context
	cfg *config.Config
}

func Init(repo Repostory, pub VoidPublisher, cash Casher, errs ErrHandler) *Service {
	return &Service{
		repo:      repo,
		publisher: pub,
		casher:    cash,
		errs:      errs,
	}
}

func (s *Service) checkOnDeletePerms(authorID, pollID uuid.UUID) (bool, error) {
	vote, err := s.repo.GetByID(s.ctx, pollID)
	if err != nil {
		return false, err
	}

	if vote.UserID != authorID {
		return false, errors.New("you do not have permision to do this")
	}

	return true, nil
}

func (s *Service) checkOn(userID, pollID string) (bool, error) {
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return false, err
	}

	pollUID, err := uuid.Parse(pollID)
	if err != nil {
		return false, err
	}

	_, err = s.repo.GetByUserAndPollID(s.ctx, userUID, pollUID)
	switch s.errs.HandleError(err) {
	case errs.ErrNotFound:
		return true, nil
	case nil:
		return false, errors.New("vote already exist")
	default:
		return false, err
	}
}

func (s *Service) CreateVote(userId string, vote *entity.Vote) error {
	ok, err := s.checkOn(userId, vote.ID.String())
	if err != nil && !ok {
		return err
	}
	if err := s.repo.Add(s.ctx, vote); err != nil {
		return err
	}

	if err := s.publisher.Publish(s.ctx, vote, vote.ID.String()); err != nil {
		return err
	}

	return s.casher.AddToCash(s.ctx, vote.ID.String(), vote)
}

func (s *Service) DeleteVote(authorID, voteId string) error {
	uid, err := uuid.Parse(voteId)
	if err != nil {
		return err
	}

	var cherr chan error

	authorUID, err := uuid.Parse(authorID)
	if err != nil {
		return err
	}

	if ok, err := s.checkOnDeletePerms(authorUID, uid); !ok || err != nil {
		return err
	}

	if err = s.repo.Delete(s.ctx, uid); err != nil {
		return err
	}

	go func() {
		cherr <- retrier.Do(3, 5, func() error {
			return s.publisher.Publish(s.ctx, struct {
				VoteID string `json:"vote_id"`
			}{
				VoteID: voteId,
			}, "vote.deleted")
		})
	}()

	return retrier.Do(3, 5, func() error {
		return s.casher.RemoveFromCash(s.ctx, voteId)
	})
}
