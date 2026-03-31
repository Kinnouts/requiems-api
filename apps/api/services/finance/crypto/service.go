package crypto

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"requiems-api/platform/httpx"
)

const (
	cacheTTL      = 5 * time.Minute
	coinGeckoURL  = "https://api.coingecko.com/api/v3"
	httpTimeout   = 10 * time.Second
)

type coinInfo struct {
	id   string
	name string
}

// coinMap maps uppercase ticker symbols to CoinGecko IDs and display names.
var coinMap = map[string]coinInfo{
	"BTC":   {id: "bitcoin", name: "Bitcoin"},
	"ETH":   {id: "ethereum", name: "Ethereum"},
	"BNB":   {id: "binancecoin", name: "BNB"},
	"XRP":   {id: "ripple", name: "XRP"},
	"ADA":   {id: "cardano", name: "Cardano"},
	"SOL":   {id: "solana", name: "Solana"},
	"DOGE":  {id: "dogecoin", name: "Dogecoin"},
	"DOT":   {id: "polkadot", name: "Polkadot"},
	"MATIC": {id: "matic-network", name: "Polygon"},
	"AVAX":  {id: "avalanche-2", name: "Avalanche"},
	"LINK":  {id: "chainlink", name: "Chainlink"},
	"LTC":   {id: "litecoin", name: "Litecoin"},
	"UNI":   {id: "uniswap", name: "Uniswap"},
	"ATOM":  {id: "cosmos", name: "Cosmos"},
	"TRX":   {id: "tron", name: "TRON"},
	"XLM":   {id: "stellar", name: "Stellar"},
	"ALGO":  {id: "algorand", name: "Algorand"},
	"NEAR":  {id: "near", name: "NEAR Protocol"},
	"FTM":   {id: "fantom", name: "Fantom"},
	"SHIB":  {id: "shiba-inu", name: "Shiba Inu"},
}

// Getter is the interface used by the HTTP transport layer.
type Getter interface {
	GetPrice(ctx context.Context, symbol string) (CryptoPrice, error)
}

// Service fetches cryptocurrency prices from CoinGecko and caches them in Redis.
type Service struct {
	rdb        *redis.Client
	httpClient *http.Client
	baseURL    string
}

// NewService creates a Service backed by the CoinGecko API.
func NewService(rdb *redis.Client) *Service {
	return &Service{
		rdb:        rdb,
		httpClient: &http.Client{Timeout: httpTimeout},
		baseURL:    coinGeckoURL,
	}
}

// newServiceWithClient is used in tests to inject a custom HTTP client and base URL.
func newServiceWithClient(rdb *redis.Client, client *http.Client, baseURL string) *Service {
	return &Service{rdb: rdb, httpClient: client, baseURL: baseURL}
}

// GetPrice returns current price data for the given ticker symbol.
func (s *Service) GetPrice(ctx context.Context, symbol string) (CryptoPrice, error) {
	coin, ok := coinMap[symbol]
	if !ok {
		return CryptoPrice{}, &httpx.AppError{
			Status:  http.StatusUnprocessableEntity,
			Code:    "unknown_symbol",
			Message: fmt.Sprintf("unsupported symbol: %s", symbol),
		}
	}

	if s.rdb != nil {
		if p, ok := s.fromCache(ctx, symbol); ok {
			return p, nil
		}
	}

	price, err := s.fetchPrice(ctx, coin.id, symbol, coin.name)
	if err != nil {
		return CryptoPrice{}, err
	}

	if s.rdb != nil {
		s.toCache(ctx, symbol, price)
	}

	return price, nil
}

func (s *Service) fromCache(ctx context.Context, symbol string) (CryptoPrice, bool) {
	val, err := s.rdb.Get(ctx, cacheKey(symbol)).Result()
	if err != nil {
		return CryptoPrice{}, false
	}

	var p CryptoPrice
	if err := json.Unmarshal([]byte(val), &p); err != nil {
		return CryptoPrice{}, false
	}

	return p, true
}

func (s *Service) toCache(ctx context.Context, symbol string, p CryptoPrice) {
	b, err := json.Marshal(p)
	if err != nil {
		return
	}
	s.rdb.Set(ctx, cacheKey(symbol), string(b), cacheTTL)
}

// coinGeckoResponse is the JSON shape returned by CoinGecko /simple/price.
type coinGeckoResponse map[string]struct {
	USD          float64 `json:"usd"`
	USD24hChange float64 `json:"usd_24h_change"`
	USDMarketCap float64 `json:"usd_market_cap"`
	USD24hVol    float64 `json:"usd_24h_vol"`
}

func (s *Service) fetchPrice(ctx context.Context, coinID, symbol, name string) (CryptoPrice, error) {
	url := fmt.Sprintf(
		"%s/simple/price?ids=%s&vs_currencies=usd&include_market_cap=true&include_24hr_vol=true&include_24hr_change=true",
		s.baseURL, coinID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return CryptoPrice{}, fmt.Errorf("crypto: build request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return CryptoPrice{}, &httpx.AppError{
			Status:  http.StatusServiceUnavailable,
			Code:    "upstream_error",
			Message: "crypto price service unavailable",
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CryptoPrice{}, &httpx.AppError{
			Status:  http.StatusServiceUnavailable,
			Code:    "upstream_error",
			Message: "crypto price service unavailable",
		}
	}

	var body coinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return CryptoPrice{}, &httpx.AppError{
			Status:  http.StatusServiceUnavailable,
			Code:    "upstream_error",
			Message: "crypto price service unavailable",
		}
	}

	data, ok := body[coinID]
	if !ok {
		return CryptoPrice{}, &httpx.AppError{
			Status:  http.StatusServiceUnavailable,
			Code:    "upstream_error",
			Message: "crypto price service unavailable",
		}
	}

	return CryptoPrice{
		Symbol:    strings.ToUpper(symbol),
		Name:      name,
		PriceUSD:  data.USD,
		Change24h: data.USD24hChange,
		MarketCap: data.USDMarketCap,
		Volume24h: data.USD24hVol,
	}, nil
}

func cacheKey(symbol string) string {
	return fmt.Sprintf("crypto:%s", strings.ToUpper(symbol))
}
