package user

import (
	"financial-data-aggregator-backend/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userServvice service.UserService
}

func NewHandler(userServvice service.UserService) *Handler {
	return &Handler{
		userServvice: userServvice,
	}
}

func (h *Handler) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unauthorized"})
		return
	}

	idStr := fmt.Sprint(userID)

	user, err := h.userServvice.GetUserProfile(idStr)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "useer not found"})
	}

	ctx.JSON(http.StatusOK, user)
}
