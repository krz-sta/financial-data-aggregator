package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewHandler(db *gorm.DB, redis *redis.Client) *Handler {
	return &Handler{
		DB:    db,
		Redis: redis,
	}
}

func (h Handler) DBHealth(ctx *gin.Context) {
	db, err := h.DB.DB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = db.Ping()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "DB is down"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}

func (h Handler) RedisHealth(ctx *gin.Context) {
	if err := h.Redis.Ping(ctx.Request.Context()).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Redis is down"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}
