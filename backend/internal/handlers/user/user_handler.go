package user

import (
	"fmt"
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

func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user doesn't exists"})
		return
	}

	idStr := fmt.Sprint(userID)
	c.JSON(http.StatusOK, gin.H{"data": "hello user, " + idStr})
}
