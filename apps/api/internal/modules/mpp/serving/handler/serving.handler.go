package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/serving/dto"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/serving/service"
	"github.com/ndollem/mpp/apps/api/internal/shared/response"
)

type ServingHandler struct {
	servingService *service.ServingService
}

func NewServingHandler(servingService *service.ServingService) *ServingHandler {
	return &ServingHandler{servingService: servingService}
}

// Done handles POST /mpp/v1/antrian/:id/done
func (h *ServingHandler) Done(c *gin.Context) {
	completion, err := h.servingService.Finish(c.Request.Context(), c.Param("id"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAntrianNotFound):
			response.Error(c, http.StatusNotFound, "Antrian not found", "")
		case errors.Is(err, service.ErrNotServing):
			// DONE is only reachable from SERVING — an illegal transition is a 409
			// (docs/04-api/api-conventions.md).
			response.Error(c, http.StatusConflict, "Antrian is not being served", "")
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to finish the antrian", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Antrian finished successfully", completion)
}

// Display handles GET /mpp/v1/display?instansi_id=...
func (h *ServingHandler) Display(c *gin.Context) {
	var query dto.DisplayQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	snapshot, err := h.servingService.Display(c.Request.Context(), query.InstansiID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to build the display snapshot", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Display snapshot retrieved successfully", snapshot)
}
