package queue

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"

	"github.com/ndollem/mpp/apps/api/internal/middleware"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/handler"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/repository"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/service"
)

type QueueModule struct {
	Handler    *handler.QueueHandler
	Service    *service.QueueService
	Repository *repository.AntrianRepository
}

// Initialize initializes the queue module. Redis is a hard dependency: it holds
// the per-service daily counters that make number allocation atomic.
func Initialize(db *pgxpool.Pool, rdb *goredis.Client) *QueueModule {
	antrianRepo := repository.NewAntrianRepository(db)
	queueService := service.NewQueueService(antrianRepo, rdb)
	queueHandler := handler.NewQueueHandler(queueService)

	return &QueueModule{
		Handler:    queueHandler,
		Service:    queueService,
		Repository: antrianRepo,
	}
}

// SetupRoutes registers queue routes. Staff and devices both reach this with the
// same JWTAuth() chain used by check-in.
func (m *QueueModule) SetupRoutes(router *gin.RouterGroup) {
	router.GET("/queue", middleware.JWTAuth(), m.Handler.List)
}
