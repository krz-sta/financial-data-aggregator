package service

import (
	"context"
	"encoding/json"
	"errors"
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
	GetRates(ctx context.Context) map[string]float64
	GetHistory(ctx context.Context, symbol string) ([]models.HistoryPoint, error)
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

	//pobranie history (1 dziń interwał)
	go s.fetchCryptoHistory(ctx, assets)

	//pobrani aktualbyh cen
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

func (s *priceService) fetchCryptoHistory(ctx context.Context, assets []models.AssetInfo) {
	for _, a := range assets {
		if a.Type == "crypto" {
			_, err := s.GetHistory(ctx, a.Symbol)
			if err != nil {
				log.Printf("failed to fetch history for %s: %v", a.Symbol, err)
			}
			time.Sleep(2 * time.Second)
		}
	}
	log.Println("crypto history fetched")
}

func (s *priceService) GetRates(ctx context.Context) map[string]float64 {
	rates := make(map[string]float64)
	assets := s.assetService.GetSupportedAssets()

	rates["PLN"] = 1

	for _, a := range assets {
		key := fmt.Sprintf("%s:%s", a.Type, a.Symbol)

		rate, err := s.redisClient.Get(ctx, key).Float64()
		if err != nil {
			log.Printf("cant get rate for %s: %v", key, err.Error())
		}

		rates[a.Symbol] = rate
	}

	return rates
}

func (s *priceService) GetHistory(ctx context.Context, symbol string) ([]models.HistoryPoint, error) {
	days := 30
	key := fmt.Sprintf("history:%s:%d", symbol, days)

	cached, err := s.redisClient.Get(ctx, key).Result()
	if err == nil {
		var history []models.HistoryPoint
		if json.Unmarshal([]byte(cached), &history) == nil {
			return history, nil
		}
	}

	var apiID string
	for _, a := range s.assetService.GetSupportedAssets() {
		if a.Symbol == symbol {
			if a.Type != "crypto" {
				return []models.HistoryPoint{}, nil //historia tylko dla fiatow
			}
			apiID = a.ApiID
			break
		}
	}

	//pobranie danych historycznych
	if apiID == "" {
		return nil, errors.New("unsupported asset")
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart?vs_currency=pln&days=%d", apiID, days)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var result struct {
		Prices [][]float64 `json:"prices"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var history []models.HistoryPoint
	for _, p := range result.Prices {
		if len(p) == 2 {
			history = append(history, models.HistoryPoint{
				Timestamp: int64(p[0]),
				Price:     p[1],
			})
		}
	}

	data, _ := json.Marshal(history)
	s.redisClient.Set(ctx, key, data, 24*time.Hour)

	return history, nil
}

func (s *priceService) fetchCrypto(ctx context.Context, assets []models.AssetInfo) {
	var cryptoIDS []string
	var idToSymbol = make(map[string]string)

	for _, a := range assets {
		if a.Type == "crypto" {
			cryptoIDS = append(cryptoIDS, a.ApiID)
			idToSymbol[a.ApiID] = a.Symbol
		}
	}

	if len(cryptoIDS) == 0 {
		return
	}

	idsParam := strings.Join(cryptoIDS, ",")
	cgUrl := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=pln", idsParam)

	res, err := http.Get(cgUrl)
	if err != nil {
		log.Printf("coinngecko fetch error: %v", err.Error())
		return
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var result map[string]map[string]float64
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Printf("coingecko decode failed: %v", err.Error())
		return
	}

	nowTimestamp := time.Now().UnixMilli()

	for apiID, rates := range result {
		symbol := idToSymbol[apiID]
		rate := rates["pln"]
		key := fmt.Sprintf("crypto:%s", symbol)

		//zapis aktualnej ceny
		if err := s.redisClient.Set(ctx, key, rate, 10*time.Minute).Err(); err != nil {
			log.Printf("redis save err for %s: %v", key, err.Error())
		}

		//akrualna do historii
		historyKey := fmt.Sprintf("history:%s:30", symbol)
		cachedHistory, err := s.redisClient.Get(ctx, historyKey).Result()

		if err == nil {
			var history []models.HistoryPoint
			if err := json.Unmarshal([]byte(cachedHistory), &history); err == nil {

				history = append(history, models.HistoryPoint{
					Timestamp: nowTimestamp,
					Price:     rate,
				})

				data, _ := json.Marshal(history)
				s.redisClient.Set(ctx, historyKey, data, 24*time.Hour)
			}
		}
	}

	log.Println("crypto updated & history appended")
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

	defer func() {
		_ = res.Body.Close()
	}()

	var result []models.FrankfurterAsset
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Printf("frankfurt decode failed: %v", err)
		return
	}

	if len(result) == 0 {
		return
	}

	for _, asset := range result {
		if fiats[asset.Quote] {
			key := fmt.Sprintf("fiat:%s", asset.Quote)
			plnRate := 1 / asset.Rate

			if err := s.redisClient.Set(ctx, key, plnRate, 10*time.Minute).Err(); err != nil {
				log.Printf("redis save err for %s: %v", key, err.Error())
			}
		}
	}
	log.Println("fiat updated")
}
