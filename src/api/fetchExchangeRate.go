package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

// MonobankResponse represents the structure of the response from Monobank API
type MonobankResponse struct {
	CurrencyCodeA int     `json:"currencyCodeA"`
	CurrencyCodeB int     `json:"currencyCodeB"`
	RateSell      float64 `json:"rateSell"`
}

var (
	cachedRate    float64
	cacheTime     time.Time
	cacheDuration = 5 * time.Minute
	cacheMutex    sync.Mutex
)

// FetchExchangeRate fetches the exchange rate from Monobank API
func FetchExchangeRate() (float64, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Check if the cache is still valid
	if time.Since(cacheTime) < cacheDuration {
		return cachedRate, nil
	}

	// Fetch new rate from Monobank API
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.monobank.ua/bank/currency")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("failed to get a valid response from Monobank API")
	}

	var rates []MonobankResponse
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return 0, err
	}

	for _, rate := range rates {
		if rate.CurrencyCodeA == 840 && rate.CurrencyCodeB == 980 {
			// Update the cache
			cachedRate = rate.RateSell
			cacheTime = time.Now()
			return cachedRate, nil
		}
	}

	return 0, errors.New("exchange rate not found")
}
