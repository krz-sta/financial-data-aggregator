package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB: db,
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
