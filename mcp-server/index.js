#!/usr/bin/env node
/**
 * CMC Agent MCP Server
 * CoinMarketCap API Hackathon 2026 · DoraHacks · AI Agents & Automation track
 *
 * Wraps the Go backend (localhost:8080) into standard MCP tools,
 * so any LLM agent (Claude Desktop, Cursor, WorkBuddy...) can query
 * live crypto market data through tool-calling.
 */

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";

const BACKEND = process.env.CMC_BACKEND_URL || "http://localhost:8080";

async function callAPI(path) {
  const res = await fetch(`${BACKEND}${path}`);
  if (!res.ok) {
    throw new Error(`backend ${res.status}: ${await res.text()}`);
  }
  return res.json();
}

function fmtUSD(n) {
  if (n === null || n === undefined) return "—";
  if (n >= 1e12) return `$${(n / 1e12).toFixed(2)}T`;
  if (n >= 1e9) return `$${(n / 1e9).toFixed(2)}B`;
  if (n >= 1e6) return `$${(n / 1e6).toFixed(2)}M`;
  return `$${n.toLocaleString()}`;
}

const server = new McpServer({
  name: "cmc-agent",
  version: "1.0.0",
});

// ---- Tool 1: get_crypto_quote ----
server.tool(
  "get_crypto_quote",
  "Get the latest price quote for a cryptocurrency (BTC, ETH, SOL, ...). Returns price, market cap, 24h volume, and 1h/24h/7d percent changes.",
  { symbol: z.string().default("BTC").describe("Crypto symbol, e.g. BTC, ETH, SOL (comma-separated for multiple)") },
  async ({ symbol }) => {
    const data = await callAPI(`/api/agent/quotes?symbol=${encodeURIComponent(symbol)}`);
    const lines = [
      `${data.name} (${data.symbol})`,
      `Price: ${fmtUSD(data.price)}`,
      `Market Cap: ${fmtUSD(data.market_cap)}`,
      `24h Volume: ${fmtUSD(data.volume_24h)}`,
      `1h: ${data.change_1h?.toFixed(2)}% | 24h: ${data.change_24h?.toFixed(2)}% | 7d: ${data.change_7d?.toFixed(2)}%`,
    ];
    return { content: [{ type: "text", text: lines.join("\n") }] };
  }
);

// ---- Tool 2: get_market_dashboard ----
server.tool(
  "get_market_dashboard",
  "Get a one-shot market overview: BTC/ETH/SOL quotes plus the Fear & Greed Index. Best for 'how is the market doing' questions.",
  {},
  async () => {
    const data = await callAPI(`/api/agent/dashboard`);
    const quotes = data.quotes?.data || {};
    const lines = ["=== Market Dashboard ==="];
    for (const [sym, d] of Object.entries(quotes)) {
      const usd = d.quote?.USD;
      if (!usd) continue;
      lines.push(`${d.name} (${sym}): ${fmtUSD(usd.price)}  24h ${usd.percent_change_24h?.toFixed(2)}%`);
    }
    if (data.fear_greed) {
      lines.push(`Fear & Greed: ${data.fear_greed.value} (${data.fear_greed.value_classification})`);
    }
    return { content: [{ type: "text", text: lines.join("\n") }] };
  }
);

// ---- Tool 3: get_fear_greed ----
server.tool(
  "get_fear_greed",
  "Get the current crypto market Fear & Greed Index (0=Extreme Fear, 100=Extreme Greed).",
  {},
  async () => {
    const data = await callAPI(`/api/fear-greed`);
    const fg = data.data || {};
    return {
      content: [{
        type: "text",
        text: `Fear & Greed Index: ${fg.value} (${fg.value_classification})`,
      }],
    };
  }
);

// ---- Tool 4: get_trending_tokens ----
server.tool(
  "get_trending_tokens",
  "Get the list of currently trending tokens on CoinMarketCap community (most visited coins).",
  {},
  async () => {
    const data = await callAPI(`/api/trending`);
    const items = data.data || [];
    const lines = items.length
      ? items.map((t, i) => `${i + 1}. ${t.name} (${t.symbol})`)
      : ["No trending data"];
    return { content: [{ type: "text", text: `Trending on CMC:\n${lines.join("\n")}` }] };
  }
);

// ---- Tool 5: get_ohlcv ----
server.tool(
  "get_ohlcv",
  "Get historical OHLCV (candlestick) data for a cryptocurrency. Useful for trend and volatility analysis.",
  {
    symbol: z.string().default("BTC").describe("Crypto symbol, e.g. BTC"),
    period: z.enum(["daily", "hourly", "weekly"]).default("daily").describe("Candle period"),
  },
  async ({ symbol, period }) => {
    const data = await callAPI(`/api/ohlcv?symbol=${encodeURIComponent(symbol)}&period=${period}`);
    const quotes = data.data?.quotes || [];
    const lines = [`${symbol} ${period} OHLCV (last ${quotes.length} candles):`];
    for (const q of quotes.slice(-7)) {
      const usd = q.quote?.USD;
      if (!usd) continue;
      const ts = new Date(q.timestamp).toISOString().slice(0, 10);
      lines.push(`${ts}  O ${usd.open?.toFixed(0)}  H ${usd.high?.toFixed(0)}  L ${usd.low?.toFixed(0)}  C ${usd.close?.toFixed(0)}`);
    }
    return { content: [{ type: "text", text: lines.join("\n") }] };
  }
);

// ---- Start ----
const transport = new StdioServerTransport();
await server.connect(transport);
console.error(`[cmc-agent-mcp] connected, backend=${BACKEND}`);
