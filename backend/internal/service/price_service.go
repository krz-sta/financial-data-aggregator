package service

import (
	"context"
	"encoding/json"
	"financial-data-aggregator-backend/internal/models"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type PriceService interface {
	StartWorker(ctx context.Context)
}

type priceService struct {
	redisClient  *redis.Client
	assetService AssetService
}

func NewPriceService(redisClient *redis.Client, assetService AssetService) PriceService {
	return &priceService{
		redisClient:  redisClient,
		assetService: assetService,
	}
}

func (s *priceService) StartWorker(ctx context.Context) {

	assets := s.assetService.GetSupportedAssets()

	//fetch on startup
	s.fetchCrypto(ctx, assets)
	s.fetchFiat(ctx, assets)

	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.fetchCrypto(ctx, assets)
				s.fetchFiat(ctx, assets)

			case <-ctx.Done():
				ticker.Stop()
				log.Println("fetch worker stopped")
				return
			}
		}
	}()
}

func (s *priceService) fetchCrypto(ctx context.Context, assets []models.AssetInfo) {
	var cryptoIDS []string
	var idToSymbol = make(map[string]string) // bitcoin -> btc

	for _, a := range assets {
		if a.Type == "crypto" {
			cryptoIDS = append(cryptoIDS, a.ApiID)
			idToSymbol[a.ApiID] = a.Symbol
		}
	}

	if len(cryptoIDS) == 0 {
		fmt.Printf("no cryptos")
		return
	}

	idsParam := strings.Join(cryptoIDS, ",")
	cgUrl := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", idsParam)

	res, err := http.Get(cgUrl)
	if err != nil {
		log.Printf("coinngecko fetch error: %v", err.Error())
		return
	}
	defer res.Body.Close()

	var result map[string]map[string]float64 // {"bitcoin": {"usd": 10000.00}}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Printf("coingecko decode failed: %v", err.Error())
		return
	}

	if len(result) <= 0 {
		log.Printf("no cryptoresults fetched")
		return
	}

	for apiID, rates := range result {
		symbol := idToSymbol[apiID] // btc -> bitcoin
		rate := rates["usd"]

		key := fmt.Sprintf("crypto:%s", symbol) //crypto:btc

		if err := s.redisClient.Set(ctx, key, rate, 10*time.Minute).Err(); err != nil {
			log.Printf("redis save err for %s: %v", key, err.Error())
		}
	}

	log.Println("crypto updated")
}

func (s *priceService) fetchFiat(ctx context.Context, assets []models.AssetInfo) {
	fiats := make(map[string]bool)
	var quotes []string

	for _, a := range assets {
		if a.Type == "fiat" {
			fiats[a.Symbol] = true
			quotes = append(quotes, a.Symbol)
		}
	}

	symbolParam := strings.Join(quotes, ",")
	fUrl := fmt.Sprintf("https://api.frankfurter.dev/v2/rates?base=PLN&quotes=%s", symbolParam)
	res, err := http.Get(fUrl)
	if err != nil {
		log.Printf("frankfurt fetch failed: %v", err.Error())
		return
	}
	defer res.Body.Close()

	var result []models.FrankfurterAsset
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Printf("frankfurt decode failed: %v", err)
		return
	}

	if len(result) == 0 {
		log.Printf("no frankfurt results fetched")
		return
	}

	for _, asset := range result {
		if fiats[asset.Quote] == true {
			key := fmt.Sprintf("fiat:%s", asset.Quote)

			plnRate := 1 / asset.Rate

			if err := s.redisClient.Set(ctx, key, plnRate, 10*time.Minute).Err(); err != nil {
				log.Printf("redis save err for %s: %v", key, err.Error())
			}
		}
	}
	log.Println("fiat updated")
}
