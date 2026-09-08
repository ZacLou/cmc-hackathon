package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const baseURL = "https://pro-api.coinmarketcap.com"

type CMCClient struct {
	APIKey   string
	mockMode bool
	http     *resty.Client
}

// QuoteResponse — /v1/cryptocurrency/quotes/latest
type QuoteResponse struct {
	Status Status                       `json:"status"`
	Data   map[string]CryptocurrencyData `json:"data"`
}

type Status struct {
	Timestamp    time.Time `json:"timestamp"`
	ErrorCode    int       `json:"error_code"`
	ErrorMessage string    `json:"error_message"`
	CreditCount  int       `json:"credit_count"`
}

type CryptocurrencyData struct {
	ID     int               `json:"id"`
	Name   string            `json:"name"`
	Symbol string            `json:"symbol"`
	Slug   string            `json:"slug"`
	Quote  map[string]Quote  `json:"quote"`
}

type Quote struct {
	Price             float64 `json:"price"`
	Volume24h         float64 `json:"volume_24h"`
	MarketCap         float64 `json:"market_cap"`
	PercentChange1h   float64 `json:"percent_change_1h"`
	PercentChange24h  float64 `json:"percent_change_24h"`
	PercentChange7d   float64 `json:"percent_change_7d"`
	PercentChange30d  float64 `json:"percent_change_30d"`
	LastUpdated       string  `json:"last_updated"`
}

// ListingResponse — /v1/cryptocurrency/listings/latest
type ListingResponse struct {
	Status Status              `json:"status"`
	Data   []ListingItem       `json:"data"`
}

type ListingItem struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Slug   string `json:"slug"`
	CMCrank int   `json:"cmc_rank"`
	Quote  map[string]ListingQuote `json:"quote"`
}

type ListingQuote struct {
	Price      float64 `json:"price"`
	Volume24h  float64 `json:"volume_24h"`
	MarketCap  float64 `json:"market_cap"`
}

// OHLCVResponse — /v2/cryptocurrency/ohlcv/historical
type OHLCVResponse struct {
	Status Status      `json:"status"`
	Data   OHLCVData    `json:"data"`
}

type OHLCVData struct {
	ID     int              `json:"id"`
	Name   string           `json:"name"`
	Symbol string           `json:"symbol"`
	Quotes []OHLCVQuote     `json:"quotes"`
}

type OHLCVQuote struct {
	Timestamp time.Time `json:"timestamp"`
	Quote     map[string]OHLCV `json:"quote"`
}

type OHLCV struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// FearGreedResponse — /v1/fear-and-greed/latest
type FearGreedResponse struct {
	Status Status     `json:"status"`
	Data   FearGreed  `json:"data"`
}

type FearGreed struct {
	Value          int    `json:"value"`
	Classification string `json:"value_classification"`
	Timestamp      string `json:"timestamp"`
}

// TrendingResponse — /v1/trending/latest
type TrendingResponse struct {
	Status Status     `json:"status"`
	Data   []Trending `json:"data"`
}

type Trending struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Slug   string `json:"slug"`
}

func NewCMCClient(apiKey string) *CMCClient {
	if apiKey == "" {
		apiKey = "mock"
	}
	return &CMCClient{
		APIKey:   apiKey,
		mockMode: apiKey == "mock",
		http:     resty.New().SetTimeout(10 * time.Second),
	}
}

// mockQuote 返回模拟报价数据
func mockQuote(symbol, name string, price, mcap, vol float64, chg1h, chg24h, chg7d float64) CryptocurrencyData {
	return CryptocurrencyData{
		Symbol: symbol,
		Name:   name,
		Quote: map[string]Quote{
			"USD": {
				Price:             price,
				MarketCap:         mcap,
				Volume24h:         vol,
				PercentChange1h:   chg1h,
				PercentChange24h:  chg24h,
				PercentChange7d:   chg7d,
			},
		},
	}
}

func (c *CMCClient) headers() map[string]string {
	return map[string]string{
		"X-CMC_PRO_API_KEY": c.APIKey,
		"Accept":            "application/json",
	}
}

// GetLatestQuotes — 获取指定币种的最新报价
func (c *CMCClient) GetLatestQuotes(symbols string) (*QuoteResponse, error) {
	if c.mockMode {
		symList := strings.Split(symbols, ",")
		data := make(map[string]CryptocurrencyData, len(symList))
		for _, s := range symList {
			s = strings.TrimSpace(strings.ToUpper(s))
			switch s {
			case "BTC":
				data["BTC"] = mockQuote("BTC", "Bitcoin", 42000.0, 820_000_000_000, 25_000_000_000, -0.3, 2.1, 8.5)
			case "ETH":
				data["ETH"] = mockQuote("ETH", "Ethereum", 2250.0, 270_000_000_000, 12_000_000_000, -0.1, 1.5, 5.2)
			case "SOL":
				data["SOL"] = mockQuote("SOL", "Solana", 98.5, 42_000_000_000, 3_000_000_000, 0.5, 3.2, 12.0)
			case "XRP":
				data["XRP"] = mockQuote("XRP", "XRP", 0.58, 32_000_000_000, 1_500_000_000, -0.2, 0.8, 3.1)
			case "DOGE":
				data["DOGE"] = mockQuote("DOGE", "Dogecoin", 0.074, 10_800_000_000, 450_000_000, 0.1, -1.2, 4.5)
			case "ADA":
				data["ADA"] = mockQuote("ADA", "Cardano", 0.35, 12_500_000_000, 380_000_000, -0.4, 1.8, 6.7)
			default:
				data[s] = mockQuote(s, s, 1.0, 100_000_000, 5_000_000, 0, 0, 0)
			}
		}
		return &QuoteResponse{Data: data}, nil
	}
	var resp QuoteResponse
	_, err := c.http.R().
		SetHeaders(c.headers()).
		SetQueryParam("symbol", symbols).
		SetResult(&resp).
		Get(baseURL + "/v1/cryptocurrency/quotes/latest")
	if err != nil {
		return nil, fmt.Errorf("cmc quotes: %w", err)
	}
	if resp.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("cmc quotes error %d: %s", resp.Status.ErrorCode, resp.Status.ErrorMessage)
	}
	return &resp, nil
}

// GetLatestListings — 获取市值排行榜
func (c *CMCClient) GetLatestListings(limit, start int) (*ListingResponse, error) {
	if c.mockMode {
		items := []ListingItem{
			{ID: 1, Name: "Bitcoin", Symbol: "BTC", CMCrank: 1, Quote: map[string]ListingQuote{"USD": {Price: 42000, Volume24h: 25e9, MarketCap: 820e9}}},
			{ID: 1027, Name: "Ethereum", Symbol: "ETH", CMCrank: 2, Quote: map[string]ListingQuote{"USD": {Price: 2250, Volume24h: 12e9, MarketCap: 270e9}}},
			{ID: 5426, Name: "Solana", Symbol: "SOL", CMCrank: 5, Quote: map[string]ListingQuote{"USD": {Price: 98.5, Volume24h: 3e9, MarketCap: 42e9}}},
		}
		return &ListingResponse{Data: items}, nil
	}
	var resp ListingResponse
	_, err := c.http.R().
		SetHeaders(c.headers()).
		SetQueryParams(map[string]string{
			"start": fmt.Sprintf("%d", start),
			"limit": fmt.Sprintf("%d", limit),
			"sort":  "market_cap",
		}).
		SetResult(&resp).
		Get(baseURL + "/v1/cryptocurrency/listings/latest")
	if err != nil {
		return nil, fmt.Errorf("cmc listings: %w", err)
	}
	if resp.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("cmc listings error %d: %s", resp.Status.ErrorCode, resp.Status.ErrorMessage)
	}
	return &resp, nil
}

// GetOHLCVHistorical — 获取历史K线数据
func (c *CMCClient) GetOHLCVHistorical(symbol, timePeriod string, count int) (*OHLCVResponse, error) {
	if c.mockMode {
		quotes := make([]OHLCVQuote, count)
		now := time.Now()
		for i := 0; i < count; i++ {
			t := now.Add(-time.Duration(count-1-i) * 24 * time.Hour)
			quotes[i] = OHLCVQuote{
				Timestamp: t,
				Quote: map[string]OHLCV{
					"USD": {Open: 41800, High: 42500, Low: 41500, Close: 42000, Volume: 25e9},
				},
			}
		}
		return &OHLCVResponse{Data: OHLCVData{Symbol: symbol, Quotes: quotes}}, nil
	}
	var resp OHLCVResponse
	_, err := c.http.R().
		SetHeaders(c.headers()).
		SetQueryParams(map[string]string{
			"symbol":      symbol,
			"time_period": timePeriod,
			"count":       fmt.Sprintf("%d", count),
		}).
		SetResult(&resp).
		Get(baseURL + "/v2/cryptocurrency/ohlcv/historical")
	if err != nil {
		return nil, fmt.Errorf("cmc ohlcv: %w", err)
	}
	if resp.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("cmc ohlcv error %d: %s", resp.Status.ErrorCode, resp.Status.ErrorMessage)
	}
	return &resp, nil
}

// GetFearAndGreed — 获取恐惧贪婪指数（Startup 计划不含此端点，出错时 fallback mock）
func (c *CMCClient) GetFearAndGreed() (*FearGreedResponse, error) {
	if c.mockMode {
		return &FearGreedResponse{
			Data: FearGreed{Value: 45, Classification: "Neutral", Timestamp: time.Now().Format(time.RFC3339)},
		}, nil
	}
	resp, err := c.tryFearGreed()
	if err != nil {
		// Startup 计划不含此端点，fallback mock
		return &FearGreedResponse{
			Data: FearGreed{Value: 45, Classification: "Neutral (mock)", Timestamp: time.Now().Format(time.RFC3339)},
		}, nil
	}
	return resp, nil
}

func (c *CMCClient) tryFearGreed() (*FearGreedResponse, error) {
	var resp FearGreedResponse
	r, err := c.http.R().
		SetHeaders(c.headers()).
		SetResult(&resp).
		Get(baseURL + "/v1/fear-and-greed/latest")
	if err != nil {
		return nil, err
	}
	if r.StatusCode() == 404 || resp.Data.Value == 0 {
		return nil, fmt.Errorf("endpoint unavailable (status %d)", r.StatusCode())
	}
	if resp.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("cmc error %d", resp.Status.ErrorCode)
	}
	return &resp, nil
}

// GetTrending — 获取社区热搜榜（Startup 计划不含此端点，出错时 fallback mock）
func (c *CMCClient) GetTrending() (*TrendingResponse, error) {
	if c.mockMode {
		return &TrendingResponse{
			Data: []Trending{
				{ID: 1, Name: "Bitcoin", Symbol: "BTC", Slug: "bitcoin"},
				{ID: 1027, Name: "Ethereum", Symbol: "ETH", Slug: "ethereum"},
				{ID: 5426, Name: "Solana", Symbol: "SOL", Slug: "solana"},
				{ID: 52, Name: "XRP", Symbol: "XRP", Slug: "xrp"},
			},
		}, nil
	}
	resp, err := c.tryTrending()
	if err != nil {
		return &TrendingResponse{
			Data: []Trending{
				{ID: 1, Name: "Bitcoin", Symbol: "BTC", Slug: "bitcoin"},
				{ID: 1027, Name: "Ethereum", Symbol: "ETH", Slug: "ethereum"},
				{ID: 5426, Name: "Solana", Symbol: "SOL", Slug: "solana"},
				{ID: 52, Name: "XRP", Symbol: "XRP", Slug: "xrp"},
			},
		}, nil
	}
	return resp, nil
}

func (c *CMCClient) tryTrending() (*TrendingResponse, error) {
	var resp TrendingResponse
	r, err := c.http.R().
		SetHeaders(c.headers()).
		SetResult(&resp).
		Get(baseURL + "/v1/trending/latest")
	if err != nil {
		return nil, err
	}
	if r.StatusCode() == 404 || len(resp.Data) == 0 {
		return nil, fmt.Errorf("endpoint unavailable (status %d)", r.StatusCode())
	}
	if resp.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("cmc error %d", resp.Status.ErrorCode)
	}
	return &resp, nil
}

// GetQuoteJSON 返回原始 JSON，供 agent/MCP 消费
func (c *CMCClient) GetQuoteJSON(symbol string) (json.RawMessage, error) {
	resp, err := c.GetLatestQuotes(symbol)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal quote: %w", err)
	}
	return b, nil
}