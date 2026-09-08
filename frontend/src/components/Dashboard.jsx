import { useMemo } from 'react'

function formatUSD(n) {
  if (!n) return '—'
  if (n >= 1e12) return `$${(n / 1e12).toFixed(2)}T`
  if (n >= 1e9) return `$${(n / 1e9).toFixed(2)}B`
  if (n >= 1e6) return `$${(n / 1e6).toFixed(2)}M`
  return `$${n.toLocaleString()}`
}

function formatChange(v) {
  const color = v >= 0 ? '#e74c3c' : '#27ae60'
  return <span style={{ color, fontWeight: 700 }}>{v >= 0 ? '+' : ''}{v.toFixed(2)}%</span>
}

function FearGreed({ data }) {
  if (!data || !data.fear_greed) return null
  const fg = data.fear_greed
  const color = fg.value <= 25 ? '#e74c3c' : fg.value >= 75 ? '#27ae60' : '#f39c12'
  return (
    <div className="card fear-greed">
      <h3>😱 Fear & Greed</h3>
      <div className="fg-value" style={{ color }}>{fg.value}</div>
      <div className="fg-label">{fg.classification}</div>
    </div>
  )
}

function QuoteCard({ symbol, data }) {
  const d = data?.quotes?.data?.[symbol]
  if (!d) return null
  const usd = d.quote?.USD
  if (!usd) return null
  return (
    <div className="card quote-card">
      <div className="coin-header">
        <span className="coin-symbol">{d.symbol}</span>
        <span className="coin-name">{d.name}</span>
      </div>
      <div className="coin-price">{usd.price < 1 ? `$${usd.price.toFixed(6)}` : `$${usd.price.toLocaleString()}`}</div>
      <div className="coin-changes">
        <span>1h {formatChange(usd.percent_change_1h)}</span>
        <span>24h {formatChange(usd.percent_change_24h)}</span>
        <span>7d {formatChange(usd.percent_change_7d)}</span>
      </div>
      <div className="coin-meta">
        <span>MCap {formatUSD(usd.market_cap)}</span>
        <span>Vol {formatUSD(usd.volume_24h)}</span>
      </div>
    </div>
  )
}

export default function Dashboard({ data }) {
  const symbols = useMemo(() => ['BTC', 'ETH', 'SOL'], [])

  return (
    <section className="dashboard">
      <h2>📊 Market Dashboard</h2>
      <div className="grid">
        {symbols.map(s => (
          <QuoteCard key={s} symbol={s} data={data} />
        ))}
        <FearGreed data={data} />
      </div>
    </section>
  )
}