package price

import (
	"financial-data-aggregator-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	priceService service.PriceService
}

func NewHandler(priceService service.PriceService) *Handler {
	return &Handler{priceService: priceService}
}

func (h *Handler) GetRates(ctx *gin.Context) {
	rates := h.priceService.GetRates(ctx.Request.Context())
	ctx.JSON(http.StatusOK, rates)
}

func (h *Handler) GetHistory(ctx *gin.Context) {
	symbol := ctx.Param("symbol")

	history, err := h.priceService.GetHistory(ctx.Request.Context(), symbol)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, history)
}
