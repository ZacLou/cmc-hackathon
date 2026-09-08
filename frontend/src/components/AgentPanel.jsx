export default function AgentPanel() {
  return (
    <section className="agent-panel">
      <h2>🧠 AI Agent Analysis</h2>
      <p className="agent-intro">
        Ask the AI agent about market conditions, token performance, and trends.
        The agent calls CoinMarketCap API in real-time via <code>/api/agent/*</code> endpoints.
      </p>
      <div className="agent-actions">
        <a href="/api/agent/quotes?symbol=BTC" target="_blank" rel="noreferrer" className="agent-btn">
          🔍 BTC Quote
        </a>
        <a href="/api/agent/dashboard" target="_blank" rel="noreferrer" className="agent-btn">
          📊 Dashboard JSON
        </a>
        <a href="https://github.com/ZacLou/cmc-hackathon" target="_blank" rel="noreferrer" className="agent-btn">
          📁 GitHub Repo
        </a>
      </div>
    </section>
  )
}