package service

import (
	"context"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/serving/domain"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/serving/repository"
)

// nextUpSize is how many upcoming numbers a TV shows under the current call.
const nextUpSize = 5

var (
	// ErrAntrianNotFound maps to 404.
	ErrAntrianNotFound = repository.ErrAntrianNotFound
	// ErrNotServing maps to 409 — an illegal state transition.
	ErrNotServing = repository.ErrNotServing
)

type ServingService struct {
	servingRepo *repository.ServingRepository
}

func NewServingService(servingRepo *repository.ServingRepository) *ServingService {
	return &ServingService{servingRepo: servingRepo}
}

// Finish closes a service: SERVING -> DONE, session closed, duration recorded (BR-19).
func (s *ServingService) Finish(ctx context.Context, antrianID string) (*domain.Completion, error) {
	return s.servingRepo.Finish(ctx, antrianID)
}

// Display builds one TV snapshot. tts_text is generated here rather than on the TV
// so every screen pronounces the number identically (BR-18 keeps them in one audio
// queue; identical text keeps them saying the same thing).
func (s *ServingService) Display(ctx context.Context, instansiID string) (*domain.Display, error) {
	current, err := s.servingRepo.CurrentCall(ctx, instansiID)
	if err != nil {
		return nil, err
	}
	if current != nil {
		current.TTSText = domain.BuildTTSText(current.Nomor, current.Loket)
	}

	nextUp, err := s.servingRepo.NextUp(ctx, instansiID, nextUpSize)
	if err != nil {
		return nil, err
	}

	return &domain.Display{InstansiID: instansiID, Current: current, NextUp: nextUp}, nil
}
