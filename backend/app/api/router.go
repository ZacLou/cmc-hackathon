package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ZacLou/cmc-hackathon/app/services"
)

type Handler struct {
	CMC *services.CMCClient
}

func NewHandler(cmc *services.CMCClient) *Handler {
	return &Handler{CMC: cmc}
}

func (h *Handler) SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/api/health", h.health)

	api := r.Group("/api")
	{
		api.GET("/quotes", h.getQuotes)
		api.GET("/listings", h.getListings)
		api.GET("/ohlcv", h.getOHLCV)
		api.GET("/fear-greed", h.getFearGreed)
		api.GET("/trending", h.getTrending)
		// Agent 专用端点，返回结构化的行情摘要
		api.GET("/agent/quotes", h.agentQuote)
		api.GET("/agent/dashboard", h.agentDashboard)
	}

	return r
}

func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) getQuotes(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "BTC,ETH")
	resp, err := h.CMC.GetLatestQuotes(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getListings(c *gin.Context) {
	resp, err := h.CMC.GetLatestListings(10, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getOHLCV(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "BTC")
	period := c.DefaultQuery("period", "daily")
	resp, err := h.CMC.GetOHLCVHistorical(symbol, period, 7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getFearGreed(c *gin.Context) {
	resp, err := h.CMC.GetFearAndGreed()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getTrending(c *gin.Context) {
	resp, err := h.CMC.GetTrending()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// agentQuote 返回 LLM 友好的行情摘要
func (h *Handler) agentQuote(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "BTC")
	resp, err := h.CMC.GetLatestQuotes(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data, ok := resp.Data[symbol]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "symbol not found: " + symbol})
		return
	}
	usd := data.Quote["USD"]
	c.JSON(http.StatusOK, gin.H{
		"symbol":     data.Symbol,
		"name":       data.Name,
		"price":      usd.Price,
		"market_cap": usd.MarketCap,
		"volume_24h": usd.Volume24h,
		"change_1h":  usd.PercentChange1h,
		"change_24h": usd.PercentChange24h,
		"change_7d":  usd.PercentChange7d,
	})
}

// agentDashboard 返回一站式仪表盘数据
func (h *Handler) agentDashboard(c *gin.Context) {
	quotes, _ := h.CMC.GetLatestQuotes("BTC,ETH,SOL")
	fg, _ := h.CMC.GetFearAndGreed()

	result := gin.H{
		"quotes": quotes,
	}
	if fg != nil {
		result["fear_greed"] = fg.Data
	}
	c.JSON(http.StatusOK, result)
}