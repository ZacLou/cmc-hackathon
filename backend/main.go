package main

import (
	"log"

	"github.com/ZacLou/cmc-hackathon/app/api"
	"github.com/ZacLou/cmc-hackathon/app/core"
	"github.com/ZacLou/cmc-hackathon/app/services"
)

func main() {
	cfg := core.Load()

	if cfg.CMCAPIKey == "" {
		log.Println("⚠️  CMC_API_KEY not set — running with mock data mode")
	}

	cmc := services.NewCMCClient(cfg.CMCAPIKey)
	handler := api.NewHandler(cmc)
	r := handler.SetupRouter()

	log.Printf("🚀 CMC AI Agent Backend starting on :%s", cfg.Port)
	log.Printf("   Health:   http://localhost:%s/api/health", cfg.Port)
	log.Printf("   Quotes:   http://localhost:%s/api/quotes?symbol=BTC", cfg.Port)
	log.Printf("   Agent:    http://localhost:%s/api/agent/quotes?symbol=ETH", cfg.Port)
	log.Printf("   Dashboard: http://localhost:%s/api/agent/dashboard", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}