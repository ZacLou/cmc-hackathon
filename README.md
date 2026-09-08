# CMC AI Market Agent

**CoinMarketCap API Hackathon 2026 · DoraHacks**  
**Track: AI Agents & Automation**  
**Team: Zac Lou (solo)**

A real-time cryptocurrency market intelligence agent powered by CoinMarketCap API.  
Combines live price quotes, Fear & Greed Index, community trending data, and historical OHLCV data into LLM-optimized endpoints. Ships with a React dashboard for human-facing visualization.

## Features

| Feature | Status |
|---|---|
| BTC/ETH/SOL real-time quotes | ✅ |
| Fear & Greed Index | ✅ |
| Community trending tokens | ✅ |
| Historical OHLCV endpoint | ✅ |
| Agent-optimized JSON endpoints | ✅ |
| React dashboard (dark theme) | ✅ |
| MCP server (stdio, 5 tools) | ✅ |
| Mock data mode (no API key needed) | ✅ |

## Architecture

```
AI Agent (Claude/Cursor/WorkBuddy)
          │ MCP stdio (mcp-server/) or REST /api/agent/*
          ▼
    Go Backend (Gin + Resty) :8080
          │
          ▼
    CMC Pro API (Startup tier)
          │
          ▼
    React SPA Dashboard :5173
```

## Tech Stack

- **Backend**: Go 1.21 + Gin + go-resty
- **Frontend**: React 18 + Vite 5
- **MCP Server**: Node.js + @modelcontextprotocol/sdk (stdio)
- **API**: CoinMarketCap Pro API (72+ endpoints)
- **Agent**: REST endpoints + standard MCP tools for LLM tool-calling

## Quick Start

```bash
# 1. Backend (works without API key — mock mode)
cd backend
# Optional: export CMC_API_KEY=your-cmc-api-key
go run main.go

# 2. MCP server (connects any LLM agent to live market data)
cd mcp-server
npm install
npm start

# 3. Frontend (optional)
cd frontend
npm install && npm run dev
```

Open http://localhost:5173 for the dashboard, http://localhost:8080/api/agent/dashboard for raw agent data.

## MCP Integration (Claude Desktop / Cursor / any MCP client)

Add to your MCP client config:

```json
{
  "mcpServers": {
    "cmc-agent": {
      "command": "node",
      "args": ["/path/to/cmc-hackathon/mcp-server/index.js"],
      "env": { "CMC_BACKEND_URL": "http://localhost:8080" }
    }
  }
}
```

**5 tools exposed:**

| Tool | CMC Endpoint | Description |
|---|---|---|
| `get_crypto_quote` | `/v1/cryptocurrency/quotes/latest` | Live price + changes |
| `get_market_dashboard` | quotes + fear-greed | One-shot market overview |
| `get_fear_greed` | `/v1/fear-and-greed/latest` | Fear & Greed Index |
| `get_trending_tokens` | `/v1/trending/latest` | Community trending |
| `get_ohlcv` | `/v2/cryptocurrency/ohlcv/historical` | Historical candles |

## Judging Criteria (Aligned)

| Criterion | Score | Coverage |
|---|---|---|
| Does it work | 30 | Live CMC API calls, tested MCP server, real-time dashboard |
| Usefulness | 25 | 5 MCP tools + REST endpoints feed any LLM agent |
| Interesting API use | 20 | Quotes + fear-greed + OHLCV combined in agent calls |
| Code quality | 15 | Clean Go project, typed models, standard MCP SDK |
| Presentation | 10 | Dark crypto dashboard + MCP demo flow |

## Roadmap (3 weeks)

- Week 1 ✅: Backend API + React dashboard + CMC live integration
- Week 2 ✅: MCP server wrapper, 5 tools, client test passed
- Week 3: Deploy, demo video, X post, submission (by Sep 30)