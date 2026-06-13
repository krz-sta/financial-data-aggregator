package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssetService_GetSupportedAssets(t *testing.T) {
	service := NewAssetService()
	assets := service.GetSupportedAssets()

	assert.Len(t, assets, 8)

	assert.Equal(t, "BTC", assets[0].Symbol)
	assert.Equal(t, "crypto", assets[0].Type)

	assert.Equal(t, "USD", assets[5].Symbol)
	assert.Equal(t, "fiat", assets[5].Type)
}
