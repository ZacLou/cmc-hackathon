import { useState, useEffect, useCallback } from 'react'
import Dashboard from './components/Dashboard'
import AgentPanel from './components/AgentPanel'

const API_BASE = '/api'

export default function App() {
  const [marketData, setMarketData] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const fetchDashboard = useCallback(async () => {
    try {
      setLoading(true)
      const res = await fetch(`${API_BASE}/agent/dashboard`)
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      const data = await res.json()
      setMarketData(data)
      setError(null)
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchDashboard()
    const id = setInterval(fetchDashboard, 30000)
    return () => clearInterval(id)
  }, [fetchDashboard])

  return (
    <div className="app">
      <header className="app-header">
        <h1>🤖 CMC AI Market Agent</h1>
        <span className="subtitle">CoinMarketCap API Hackathon · DoraHacks 2026</span>
      </header>
      <main className="app-main">
        {loading && !marketData && <div className="loading">Loading market data...</div>}
        {error && !marketData && <div className="error">Error: {error}</div>}
        {marketData && (
          <>
            <Dashboard data={marketData} />
            <AgentPanel />
          </>
        )}
      </main>
    </div>
  )
}