package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/checkin/dto"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/checkin/service"
	"github.com/ndollem/mpp/apps/api/internal/shared/response"
)

type CheckinHandler struct {
	checkinService *service.CheckinService
}

func NewCheckinHandler(checkinService *service.CheckinService) *CheckinHandler {
	return &CheckinHandler{checkinService: checkinService}
}

// Create handles POST /mpp/v1/checkin
func (h *CheckinHandler) Create(c *gin.Context) {
	var req dto.CheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	result, err := h.checkinService.Create(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTokenNotFound):
			response.Error(c, http.StatusNotFound, "QR token not found", "")
		case errors.Is(err, service.ErrTokenRejected):
			// Reuse, expiry and wrong-day all reject the state change (409), per
			// docs/04-api/api-conventions.md.
			response.Error(c, http.StatusConflict, "QR token rejected", err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to check in", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Check-in successful", result)
}
