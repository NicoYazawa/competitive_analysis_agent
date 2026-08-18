// DTO types matching backend Go responses
export interface CompetitorResponse {
  id: string
  name: string
  platform: string
  platform_product_id: string
  current_price: number
  currency: string
  rating: number
  review_count: number
  seller_rating: number
  seller_review_count: number
  source_url: string
  created_at: string
  updated_at: string
  price_history?: PricePoint[]
}

export interface PricePoint {
  price: number
  currency: string
  recorded_at: string
}

export interface ListCompetitorsResponse {
  competitors: CompetitorResponse[]
  total: number
  limit: number
  offset: number
}

export interface PriceChangeResponse {
  competitor_id: string
  name: string
  old_price: number
  new_price: number
  change: number
  change_percent: number
  direction: 'increased' | 'decreased' | 'unchanged'
}

export interface DetectPriceChangesResponse {
  threshold: string
  changes: PriceChangeResponse[]
  count: number
}

export interface DemandSignal {
  keyword: string
  search_volume: string
  trend: string
  competition: string
}

export interface DemandSignalsResponse {
  platform?: string
  signals: DemandSignal[]
  updated_at: string
}

export interface MarketOverviewResponse {
  total_products: number
  active_competitors: number
  avg_price_change_percent: number
  market_sentiment: string
  top_category: string
  data_points: number
}

export interface TrendingProduct {
  rank: number
  name: string
  category: string
  price_change_percent: number
  demand_score: number
  recommendation: string
}

export interface TrendingProductsResponse {
  products: TrendingProduct[]
}

export interface CategoryTrend {
  category: string
  trend: string
  demand_level: string
  growth_rate: string
  top_products: string[]
}

export interface CategoryResponse {
  categories: Array<{ id: string; name: string; icon: string }>
}

export interface MarketTrendResponse {
  trend: string
  opportunities: string[]
  demand_signal: string
}

export interface AggregatedResult {
  summary: string
  market_trend?: {
    trend: string
    opportunities: string[]
    demand_signal: string
  }
  competitor_analysis?: {
    analysis: string
    competitors: Array<{
      name: string
      strength: string
      weakness: string
      strategy: string
    }>
  }
  pricing_strategy?: {
    recommended_price_range: string
    rationale: string
    positioning: string
  }
  supply_chain?: {
    status: string
    risk_level: string
    factors: string[]
  }
  task_count: number
  errors?: string[]
}

// ---------------------------------------------------------------------------
// Mappers to store types
// ---------------------------------------------------------------------------

export function mapCompetitor(r: CompetitorResponse) {
  return {
    id: r.id,
    name: r.name,
    platform: r.platform,
    region: r.platform, // backend uses platform as region proxy
    category: '',
    currentPrice: r.current_price,
    currency: r.currency,
    rating: r.rating,
    reviewCount: r.review_count,
    status: 'good' as const,
    statusText: '正常'
  }
}

export function mapDemandSignals(r: DemandSignalsResponse) {
  return r.signals.map((s, i) => ({
    id: String(i),
    competitor: s.keyword,
    signalType: `${s.trend} (${s.search_volume})`,
    time: r.updated_at,
    riskLevel: s.competition === 'High' ? ('warning' as const) : ('info' as const),
    riskText: s.competition
  }))
}

export function mapMarketOverview(r: MarketOverviewResponse) {
  return {
    competitorsCount: r.active_competitors,
    competitorsDelta: 0,
    priceChanges: 0,
    priceChangesDelta: r.avg_price_change_percent,
    activeAlerts: 0,
    activeAlertsNew: 0,
    markets: r.top_category ? 1 : 0,
    marketNames: r.top_category || ''
  }
}

export function mapTrendingProducts(r: TrendingProductsResponse) {
  return r.products.map(p => ({
    id: String(p.rank),
    name: p.name,
    category: p.category,
    brand: '',
    description: p.recommendation,
    rating: 0,
    reviews: 0,
    price: 0,
    score: Math.round(p.demand_score * 10),
    badge: p.recommendation,
    badgeType: p.recommendation === 'Strong Buy' ? 'good' as const : 'info' as const
  }))
}
