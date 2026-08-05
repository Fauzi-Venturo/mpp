package service

import (
	"context"
	"errors"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/loket_ops/domain"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/loket_ops/dto"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/loket_ops/repository"
)

var (
	// ErrNoWaiting maps to 404 — the queue is empty, which is not an error state.
	ErrNoWaiting = repository.ErrNoWaiting
	// ErrAntrianNotFound maps to 404.
	ErrAntrianNotFound = repository.ErrAntrianNotFound
	// ErrLoketNotEligible maps to 409 (BR-12/BR-15).
	ErrLoketNotEligible = repository.ErrLoketNotEligible
	// ErrIllegalTransition maps to 409 (NFR-DATA-03).
	ErrIllegalTransition = errors.New("illegal antrian state transition")
	// ErrCallLimit maps to 409 — called 3x already, the ticket must be skipped (BR-16).
	ErrCallLimit = errors.New("antrian has reached the maximum number of calls")
)

type LoketOpsService struct {
	repo *repository.LoketOpsRepository
}

func NewLoketOpsService(repo *repository.LoketOpsRepository) *LoketOpsService {
	return &LoketOpsService{repo: repo}
}

// CallNext hands the next waiting ticket to a loket.
func (s *LoketOpsService) CallNext(ctx context.Context, req dto.CallNextRequest) (*domain.Called, error) {
	return s.repo.CallNext(ctx, req.JenisLayananID, req.LoketID)
}

// Recall re-announces a ticket that is still CALLED, up to MaxCallCount times.
func (s *LoketOpsService) Recall(ctx context.Context, antrianID string) (*domain.Called, error) {
	return s.repo.Recall(ctx, antrianID, func(snap repository.Snapshot) error {
		if snap.Status != domain.StatusCalled {
			return ErrIllegalTransition
		}
		if snap.CallCount >= domain.MaxCallCount {
			return ErrCallLimit
		}
		return nil
	})
}

// Start moves a called ticket into service.
func (s *LoketOpsService) Start(ctx context.Context, antrianID, userID string) (*domain.Called, error) {
	return s.repo.Start(ctx, antrianID, userID, func(snap repository.Snapshot) error {
		if snap.Status != domain.StatusCalled {
			return ErrIllegalTransition
		}
		return nil
	})
}

// Skip closes a ticket as a no-show. Only a called ticket can be a no-show —
// nobody has been announced yet while it is still WAITING.
func (s *LoketOpsService) Skip(ctx context.Context, antrianID string) (*domain.Called, error) {
	return s.repo.Skip(ctx, antrianID, func(snap repository.Snapshot) error {
		if snap.Status != domain.StatusCalled {
			return ErrIllegalTransition
		}
		return nil
	})
}
