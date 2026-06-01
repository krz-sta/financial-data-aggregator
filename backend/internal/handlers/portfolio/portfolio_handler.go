package portfolio

import (
	"financial-data-aggregator-backend/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addInput struct {
	Symbol string  `json:"symbol" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

type Handler struct {
	portfolioService service.PortfolioService
}

func NewHandler(portfolioService service.PortfolioService) *Handler {
	return &Handler{portfolioService: portfolioService}
}

func (h *Handler) AddPortfolioItem(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var addInput addInput
	if err := ctx.ShouldBindJSON(&addInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	idStr := fmt.Sprint(userID)

	err := h.portfolioService.AddItem(idStr, addInput.Symbol, addInput.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "item added to portoflio"})
}

func (h *Handler) DeletePortfolioItem(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userStr := fmt.Sprint(userID)
	itemID := ctx.Param("id")

	err := h.portfolioService.DeleteItem(userStr, itemID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "item deleted from portoflio"})
}
