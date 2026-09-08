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
| MCP integration ready | ✅ |
| Mock data mode (no API key needed) | ✅ |

## Architecture

```
AI Agent (Claude/GPT/Cursor)
          │  REST /api/agent/*
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
- **API**: CoinMarketCap Pro API (72+ endpoints)
- **Agent**: REST endpoints designed for LLM tool-calling

## Quick Start

```bash
# Backend
cd backend
export CMC_API_KEY=your-cmc-api-key
go run main.go

# Frontend (optional)
cd frontend
npm install && npm run dev
```

Open http://localhost:5173 for the dashboard, http://localhost:8080/api/agent/dashboard for raw agent data.

## Judging Criteria (Aligned)

| Criterion | Score | Coverage |
|---|---|---|
| Does it work | 30 | Live CMC API calls, real-time dashboard, working endpoints |
| Usefulness | 25 | `/api/agent/*` endpoints feed any LLM tool-calling pipeline |
| Interesting API use | 20 | Combines quotes + fear-greed + trending in a single agent call |
| Code quality | 15 | Clean Go project, typed API models, Vite proxy |
| Presentation | 10 | Dark-themed crypto dashboard, Fear & Greed visual |

## Roadmap (3 weeks)

- Week 1 (now): Backend API + React dashboard + CMC integration
- Week 2: MCP server wrapper, agent tool-calling demo
- Week 3: Polish, deploy, demo video, submission