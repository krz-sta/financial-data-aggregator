package service

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedisContainer(t *testing.T) (*redis.Client, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	endpoint, err := redisC.Endpoint(ctx, "")
	assert.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	cleanup := func() {
		redisC.Terminate(ctx)
	}
	return client, cleanup
}

func TestPriceService_FullCycle(t *testing.T) {
	client, cleanup := setupRedisContainer(t)
	defer cleanup()

	assetSvc := NewAssetService()
	priceSvc := NewPriceService(client, assetSvc)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	priceSvc.StartWorker(ctx)

	rates := priceSvc.GetRates(ctx)

	assert.NotEmpty(t, rates)
	assert.Equal(t, 1.0, rates["PLN"])

	assert.Greater(t, rates["BTC"], 0.0)
	assert.Greater(t, rates["ETH"], 0.0)

	assert.Greater(t, rates["USD"], 0.0)
	assert.Greater(t, rates["EUR"], 0.0)
}
