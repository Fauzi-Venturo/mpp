package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ndollem/mpp/apps/api/internal/middleware"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/loket_ops/domain"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/loket_ops/dto"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/loket_ops/service"
	"github.com/ndollem/mpp/apps/api/internal/shared/response"
)

type LoketOpsHandler struct {
	svc *service.LoketOpsService
}

func NewLoketOpsHandler(svc *service.LoketOpsService) *LoketOpsHandler {
	return &LoketOpsHandler{svc: svc}
}

// CallNext handles POST /mpp/v1/queue/next
func (h *LoketOpsHandler) CallNext(c *gin.Context) {
	var req dto.CallNextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	called, err := h.svc.CallNext(c.Request.Context(), req)
	if err != nil {
		h.fail(c, err, "Failed to call the next antrian")
		return
	}

	response.Success(c, http.StatusOK, "Antrian called", called)
}

// Recall handles POST /mpp/v1/antrian/:id/recall
func (h *LoketOpsHandler) Recall(c *gin.Context) {
	h.transition(c, h.svc.Recall, "Antrian recalled", "Failed to recall the antrian")
}

// Skip handles POST /mpp/v1/antrian/:id/skip
func (h *LoketOpsHandler) Skip(c *gin.Context) {
	h.transition(c, h.svc.Skip, "Antrian skipped", "Failed to skip the antrian")
}

// Start handles POST /mpp/v1/antrian/:id/start
func (h *LoketOpsHandler) Start(c *gin.Context) {
	// The operator is whoever the credential belongs to; an API-key device simply
	// records no user on the serving session.
	userID, _ := middleware.GetUserID(c)

	called, err := h.svc.Start(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		h.fail(c, err, "Failed to start serving")
		return
	}

	response.Success(c, http.StatusOK, "Serving started", called)
}

type transitionFunc func(ctx context.Context, antrianID string) (*domain.Called, error)

func (h *LoketOpsHandler) transition(c *gin.Context, fn transitionFunc, okMsg, failMsg string) {
	called, err := fn(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.fail(c, err, failMsg)
		return
	}

	response.Success(c, http.StatusOK, okMsg, called)
}

// fail maps domain errors to the status codes docs/04-api/api-conventions.md:62
// mandates: 409 for an illegal transition or an ineligible loket, 404 when there
// is simply nothing to act on.
func (h *LoketOpsHandler) fail(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrNoWaiting):
		response.Error(c, http.StatusNotFound, "No waiting antrian for this service", "")
	case errors.Is(err, service.ErrAntrianNotFound):
		response.Error(c, http.StatusNotFound, "Antrian not found", "")
	case errors.Is(err, service.ErrLoketNotEligible):
		response.Error(c, http.StatusConflict, "Loket is closed or does not serve this service", "")
	case errors.Is(err, service.ErrCallLimit):
		response.Error(c, http.StatusConflict, "Antrian has already been called 3 times", "")
	case errors.Is(err, service.ErrIllegalTransition):
		response.Error(c, http.StatusConflict, "Illegal antrian state transition", "")
	default:
		response.Error(c, http.StatusInternalServerError, fallback, err.Error())
	}
}
