// MCP client test: spawn cmc-agent MCP server, call get_market_dashboard tool
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";

const transport = new StdioClientTransport({
  command: process.execPath,
  args: ["/Users/mac/Desktop/dev/cmc-hackathon/mcp-server/index.js"],
});

const client = new Client({ name: "test-client", version: "1.0" });
await client.connect(transport);

// List tools
const tools = await client.listTools();
console.log("=== Tools available ===");
for (const t of tools.tools) {
  console.log(`- ${t.name}: ${t.description?.slice(0, 60)}...`);
}

// Call get_crypto_quote
console.log("\n=== get_crypto_quote BTC ===");
const quote = await client.callTool({ name: "get_crypto_quote", arguments: { symbol: "BTC" } });
console.log(quote.content[0].text);

// Call get_market_dashboard
console.log("\n=== get_market_dashboard ===");
const dash = await client.callTool({ name: "get_market_dashboard", arguments: {} });
console.log(dash.content[0].text);

await client.close();
console.log("\n=== MCP TEST PASSED ===");
