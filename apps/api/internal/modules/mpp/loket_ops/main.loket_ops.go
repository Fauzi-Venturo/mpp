package loket_ops

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ndollem/mpp/apps/api/internal/middleware"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/loket_ops/handler"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/loket_ops/repository"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/loket_ops/service"
)

type LoketOpsModule struct {
	Handler    *handler.LoketOpsHandler
	Service    *service.LoketOpsService
	Repository *repository.LoketOpsRepository
}

// Initialize wires the operator (loket) actions of slice 5. No Redis: number
// allocation belongs to the queue module; calling is pure Postgres.
func Initialize(db *pgxpool.Pool) *LoketOpsModule {
	repo := repository.NewLoketOpsRepository(db)
	svc := service.NewLoketOpsService(repo)

	return &LoketOpsModule{
		Handler:    handler.NewLoketOpsHandler(svc),
		Service:    svc,
		Repository: repo,
	}
}

// SetupRoutes registers the operator endpoints from docs/04-api/rest-endpoints.md:49-53.
// Staff authenticate with the same JWTAuth() chain as queue and check-in — it accepts
// both a JWT and an X-API-Key credential.
func (m *LoketOpsModule) SetupRoutes(router *gin.RouterGroup) {
	router.POST("/queue/next", middleware.JWTAuth(), m.Handler.CallNext)

	antrian := router.Group("/antrian/:id", middleware.JWTAuth())
	{
		antrian.POST("/recall", m.Handler.Recall)
		antrian.POST("/start", m.Handler.Start)
		antrian.POST("/skip", m.Handler.Skip)
	}
}
