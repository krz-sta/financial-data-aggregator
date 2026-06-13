package asset

import (
	"financial-data-aggregator-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	assetService service.AssetService
}

func NewHandler(assetService service.AssetService) *Handler {
	return &Handler{
		assetService: assetService,
	}
}

func (h *Handler) GetAssets(ctx *gin.Context) {
	assets := h.assetService.GetSupportedAssets()
	ctx.JSON(http.StatusOK, assets)
}
