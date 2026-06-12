package service

import "financial-data-aggregator-backend/internal/models"

var supportedAssets = []models.AssetInfo{
	// krypto
	{Symbol: "BTC", Name: "Bitcoin", Type: "crypto", ApiID: "bitcoin"},
	{Symbol: "ETH", Name: "Ethereum", Type: "crypto", ApiID: "ethereum"},
	{Symbol: "SOL", Name: "Solana", Type: "crypto", ApiID: "solana"},
	{Symbol: "ADA", Name: "Cardano", Type: "crypto", ApiID: "cardano"},
	{Symbol: "DOGE", Name: "Dogecoin", Type: "crypto", ApiID: "dogecoin"},

	// fiat (normalne)
	{Symbol: "USD", Name: "Dolar Amerykański", Type: "fiat", ApiID: ""},
	{Symbol: "EUR", Name: "Euro", Type: "fiat", ApiID: ""},
	{Symbol: "GBP", Name: "Funt Szterling", Type: "fiat", ApiID: ""},
}

type AssetService interface {
	GetSupportedAssets() []models.AssetInfo
}

type assetService struct{}

func NewAssetService() AssetService {
	return &assetService{}
}

func (s *assetService) GetSupportedAssets() []models.AssetInfo {
	return supportedAssets
}
