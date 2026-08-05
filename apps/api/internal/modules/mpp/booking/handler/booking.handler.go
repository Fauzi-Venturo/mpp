package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/dto"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/service"
	"github.com/ndollem/mpp/apps/api/internal/shared/response"
)

// CompanySlugHeader identifies the MPP building (tenant) on public endpoints.
const CompanySlugHeader = "X-Company-Slug"

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

// Create handles POST /mpp/v1/booking
func (h *BookingHandler) Create(c *gin.Context) {
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	booking, err := h.bookingService.Create(c.Request.Context(), c.GetHeader(CompanySlugHeader), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTenantRequired):
			response.Error(c, http.StatusBadRequest, "X-Company-Slug header is required", "")
		case errors.Is(err, service.ErrServiceNotFound):
			response.Error(c, http.StatusNotFound, "Service not found for this tenant", "")
		case errors.Is(err, service.ErrQuotaFull):
			response.Error(c, http.StatusConflict, "Booking quota for this date is full", "")
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to create booking", err.Error())
		}
		return
	}

	response.Success(c, http.StatusCreated, "Booking created successfully", booking)
}

// GetByID handles GET /mpp/v1/booking/:id — the ticket screen reads the QR token here.
func (h *BookingHandler) GetByID(c *gin.Context) {
	booking, err := h.bookingService.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookingNotFound):
			response.Error(c, http.StatusNotFound, "Booking not found", "")
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to load booking", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Booking retrieved successfully", booking)
}
