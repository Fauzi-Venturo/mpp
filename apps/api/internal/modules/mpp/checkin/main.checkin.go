package checkin

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ndollem/mpp/apps/api/internal/middleware"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/checkin/handler"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/checkin/repository"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/checkin/service"
)

type CheckinModule struct {
	Handler    *handler.CheckinHandler
	Service    *service.CheckinService
	Repository *repository.CheckinRepository
}

// Initialize initializes the check-in module with all dependencies
func Initialize(db *pgxpool.Pool) *CheckinModule {
	checkinRepo := repository.NewCheckinRepository(db)
	checkinService := service.NewCheckinService(checkinRepo)
	checkinHandler := handler.NewCheckinHandler(checkinService)

	return &CheckinModule{
		Handler:    checkinHandler,
		Service:    checkinService,
		Repository: checkinRepo,
	}
}

// SetupRoutes registers check-in routes. Kiosks authenticate with X-API-Key, which
// JWTAuth() resolves before it looks for a Bearer token (internal/middleware/auth.go).
func (m *CheckinModule) SetupRoutes(router *gin.RouterGroup) {
	router.POST("/checkin", middleware.JWTAuth(), m.Handler.Create)
}
