package serving

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ndollem/mpp/apps/api/internal/middleware"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/serving/handler"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/serving/repository"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/serving/service"
)

type ServingModule struct {
	Handler    *handler.ServingHandler
	Service    *service.ServingService
	Repository *repository.ServingRepository
}

// Initialize initializes the serving module: closing a service and showing what is
// currently being served.
func Initialize(db *pgxpool.Pool) *ServingModule {
	servingRepo := repository.NewServingRepository(db)
	servingService := service.NewServingService(servingRepo)
	servingHandler := handler.NewServingHandler(servingService)

	return &ServingModule{
		Handler:    servingHandler,
		Service:    servingService,
		Repository: servingRepo,
	}
}

// SetupRoutes registers the closing half of the serving stage. Both the loket app
// and the TV authenticate through the same JWTAuth() chain (X-API-Key for devices).
func (m *ServingModule) SetupRoutes(router *gin.RouterGroup) {
	router.POST("/antrian/:id/done", middleware.JWTAuth(), m.Handler.Done)
	router.GET("/display", middleware.JWTAuth(), m.Handler.Display)
}
