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
