package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/dto"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/service"
	"github.com/ndollem/mpp/apps/api/internal/shared/response"
)

type QueueHandler struct {
	queueService *service.QueueService
}

func NewQueueHandler(queueService *service.QueueService) *QueueHandler {
	return &QueueHandler{queueService: queueService}
}

// List handles GET /mpp/v1/queue?layanan_id=...
func (h *QueueHandler) List(c *gin.Context) {
	var query dto.QueueQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	stream, err := h.queueService.ListWaiting(c.Request.Context(), query.LayananID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to load the queue", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Queue retrieved successfully",
		dto.NewQueueStream(query.LayananID, stream))
}
