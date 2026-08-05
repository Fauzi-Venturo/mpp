package booking

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/handler"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/repository"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/booking/service"
)

type BookingModule struct {
	Handler    *handler.BookingHandler
	Service    *service.BookingService
	Repository *repository.BookingRepository
}

// Initialize initializes the booking module with all dependencies
func Initialize(db *pgxpool.Pool) *BookingModule {
	bookingRepo := repository.NewBookingRepository(db)
	bookingService := service.NewBookingService(bookingRepo)
	bookingHandler := handler.NewBookingHandler(bookingService)

	return &BookingModule{
		Handler:    bookingHandler,
		Service:    bookingService,
		Repository: bookingRepo,
	}
}

// SetupRoutes registers booking routes. Registration is public (see
// docs/04-api/rest-endpoints.md): citizens book without an account, and the
// tenant comes from the X-Company-Slug header instead of a JWT claim.
func (m *BookingModule) SetupRoutes(router *gin.RouterGroup) {
	router.POST("/booking", m.Handler.Create)
	router.GET("/booking/:id", m.Handler.GetByID)
}
