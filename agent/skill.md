# CMC Market Intelligence Agent

**CoinMarketCap API Hackathon 2026 · DoraHacks**
**Track: AI Agents & Automation**

A real-time cryptocurrency market agent powered by CoinMarketCap API. Provides LLM-friendly endpoints for market analysis, price quotes, fear & greed index, trending tokens, and historical OHLCV data. Designed for integration into AI coding agents, chatbots, and trading dashboards via REST and MCP.

## Quick Start

```bash
# Start backend
cd backend && CMC_API_KEY=your-key go run main.go

# Start frontend (optional)
cd frontend && npm install && npm run dev
```

## Agent Endpoints

All endpoints return JSON optimized for LLM consumption.

| Endpoint | Description | Parameters |
|---|---|---|
| `GET /api/agent/quotes` | Latest price + change % for a symbol | `symbol` (default BTC) |
| `GET /api/agent/dashboard` | BTC/ETH/SOL quotes + Fear & Greed | none |
| `GET /api/fear-greed` | Fear & Greed Index | none |
| `GET /api/trending` | Community trending tokens | none |
| `GET /api/ohlcv` | Historical OHLCV data | `symbol`, `period` |

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────────┐
│  AI Agent   │────▶│  Go Backend  │────▶│  CMC Pro API    │
│  (Claude /  │     │  (Gin +      │     │  /v1/ quotes    │
│   GPT /     │     │   Resty)     │     │  /v1/ fear-greed│
│   Cursor)   │     │  :8080       │     │  /v1/ trending   │
└─────────────┘     └──────────────┘     └─────────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │  React SPA   │
                    │  Dashboard   │
                    │  :5173       │
                    └──────────────┘
```

## MCP Integration

Copy this into your AI agent's MCP config to give it live market data:

```json
{
  "mcpServers": {
    "cmc-agent": {
      "command": "curl",
      "args": ["-s", "http://localhost:8080/api/agent/dashboard"],
      "description": "Live BTC/ETH/SOL prices + Fear & Greed index"
    }
  }
}
```

## Tech Stack

- **Backend**: Go 1.21 + Gin + go-resty (REST API layer for CMC)
- **Frontend**: React 18 + Vite 5 (real-time dashboard)
- **Data**: CoinMarketCap Pro API (Startup tier, 10K+ credits/month)
- **Agent**: REST endpoints designed for LLM tool-calling

## Judging Criteria Alignment

| Criterion | How we score |
|---|---|
| **Does it work** (30 pts) | Live CMC API calls, real-time dashboard, working agent endpoints |
| **Usefulness** (25 pts) | Any developer can wrap a trading bot, screener, or research assistant around `/api/agent/*` |
| **Interesting use of API** (20 pts) | Combines quotes + fear-greed + trending in one agent-optimized call |
| **Code quality** (15 pts) | Clean Go project structure, typed API models, Vite proxy setup |
| **Presentation** (10 pts) | Dark-themed crypto dashboard with Fear & Greed visual