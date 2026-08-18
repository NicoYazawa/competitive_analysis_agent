import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  competitorsAPI,
  trendOverviewAPI,
  trendsAPI,
  strategyAPI,
  productsAPI
} from '@/api/index'
import {
  type CompetitorResponse,
  type DemandSignalsResponse,
  type MarketOverviewResponse,
  type TrendingProductsResponse,
  mapCompetitor,
  mapDemandSignals,
  mapMarketOverview,
  mapTrendingProducts
} from '@/api/dto'

export interface Product {
  id: string
  name: string
  category: string
  brand: string
  description: string
  rating?: number
  reviews?: number
  price?: number
  score?: number
  badge?: string
  badgeType?: 'critical' | 'warning' | 'good' | 'info'
}

export interface Competitor {
  id: string
  name: string
  platform: string
  region: string
  category: string
  currentPrice: number
  currency: string
  rating: number
  reviewCount: number
  status: 'critical' | 'warning' | 'info' | 'good'
  statusText: string
}

export interface MarketTrend {
  id: string
  category: string
  keyword: string
  popularityScore: number
  growthRate: number
}

export interface Alert {
  id: string
  competitor: string
  signalType: string
  time: string
  riskLevel: 'critical' | 'warning' | 'info' | 'good'
  riskText: string
}

export interface SupplyAlert {
  id: string
  riskLevel: 'critical' | 'warning' | 'info'
  category: string
  signalType: string
  detail: string
  market: string
  time: string
}

export const useMarketStore = defineStore('market', () => {
  // Dashboard KPIs
  const kpis = ref({
    competitorsCount: 0,
    competitorsDelta: 0,
    priceChanges: 0,
    priceChangesDelta: 0,
    activeAlerts: 0,
    activeAlertsNew: 0,
    markets: 0,
    marketNames: ''
  })

  // Market Overview
  const marketOverview = ref([
    { name: '北美市场', competition: '高', growth: '+18.2%' },
    { name: '欧洲市场', competition: '中', growth: '+9.7%' },
    { name: '东南亚', competition: '低', growth: '+31.4%' }
  ])

  // Trend data for chart
  const trend30Days = ref({
    xData: ['8/1', '8/4', '8/8', '8/12', '8/16', '8/20', '8/24', '8/28'],
    priceIndex: [100, 98, 95, 92, 88, 85, 82, 78],
    scoreAvg: [72, 71, 70, 70, 69, 68, 67, 66]
  })

  // Hot categories
  const hotCategories = ref([
    { name: '智能手表', score: 92 },
    { name: '无线耳机', score: 87 },
    { name: '便携储能', score: 81 },
    { name: '空气净化器', score: 74 },
    { name: '电动滑板车', score: 68 }
  ])

  // Latest alerts
  const latestAlerts = ref<Alert[]>([])

  // Competitors list
  const competitors = ref<Competitor[]>([])

  // Product selection data
  const productScores = ref<Product[]>([])

  // Price comparison data
  const priceComparison = ref([
    { id: '1', name: 'Apple Watch Ultra 2', category: '智能手表', low: 749, high: 849, avg: 799, trend: 'up', suggestedPrice: 769 },
    { id: '2', name: 'Anker 737 Power Bank', category: '便携储能', low: 119, high: 159, avg: 139, trend: 'down', suggestedPrice: 129 },
    { id: '3', name: 'Sony WH-1000XM5', category: '无线耳机', low: 328, high: 378, avg: 348, trend: 'stable', suggestedPrice: 338 },
    { id: '4', name: 'UGREEN 100W Charger', category: '充电器', low: 29, high: 42, avg: 36, trend: 'up', suggestedPrice: 32 },
    { id: '5', name: 'Jackery 1000 Plus', category: '便携储能', low: 899, high: 1099, avg: 999, trend: 'stable', suggestedPrice: 949 }
  ])

  // Pricing KPIs
  const pricingKpis = ref({
    skuCount: 0,
    underpriceAlerts: 0,
    avgProfitMargin: 0,
    profitDelta: 0,
    pendingAdjustments: 0
  })

  // Pricing trend data (60 days)
  const pricingTrend60Days = ref({
    xData: ['6/20', '7/5', '7/20', '8/4', '8/18'],
    anker: [140, 138, 135, 132, 130],
    jackery: [999, 998, 997, 996, 995],
    ugreen: [36, 37, 38, 39, 40]
  })

  // Pricing suggestions
  const pricingSuggestions = ref([
    { product: 'Anker 737 Power Bank', currentPrice: 139.99, suggestedPrice: 129, reason: '当前价格高于竞品均值 8.5%，建议降价以提升竞争力。', severity: 'critical' },
    { product: 'Jackery 1000 Plus', currentPrice: 999, suggestedPrice: 949, reason: '建议配合限时促销激活销量。', severity: 'warning' },
    { product: 'UGREEN 100W', currentPrice: 35.99, suggestedPrice: 32, reason: '价格竞争激烈，建议维持低价策略守住类目BSR。', severity: 'good' }
  ])

  // Profit simulation
  const profitSimulation = ref({
    suggestedPrice: 129.00,
    cost: 87.00,
    platformFee: 19.35,
    shipping: 8.50,
    profit: 14.15,
    profitRate: 11
  })

  // Supply chain KPIs
  const supplyKpis = ref({
    highRisk: 0,
    mediumRisk: 0,
    resolved: 0,
    categories: 0
  })

  // Supply alerts
  const supplyAlerts = ref<SupplyAlert[]>([])

  // Risk radar data
  const riskRadar = ref({
    rawMaterials: [
      { name: '碳酸锂', change: '+23%', level: 'critical', percent: 85 },
      { name: '钴', change: '+8%', level: 'warning', percent: 55 },
      { name: '芯片', change: '-3%', level: 'good', percent: 25 }
    ],
    logistics: [
      { name: '红海航线', status: '高风险', level: 'critical', percent: 85 },
      { name: '洛杉矶港', status: '正常', level: 'good', percent: 20 },
      { name: '中欧班列', status: '中风险', level: 'warning', percent: 50 }
    ],
    policy: [
      { name: '欧盟认证', cost: '+$2.3/件', level: 'warning', percent: 60 },
      { name: '美国关税', status: '稳定', level: 'good', percent: 20 },
      { name: 'RCEP优惠', status: '稳定', level: 'good', percent: 15 }
    ]
  })

  const loading = ref(false)
  const error = ref<string | null>(null)

  // ---------------------------------------------------------------------------
  // Async actions
  // ---------------------------------------------------------------------------

  async function fetchKpis() {
    try {
      const res = await trendOverviewAPI.get()
      const data = res.data as MarketOverviewResponse
      const mapped = mapMarketOverview(data)
      kpis.value = {
        competitorsCount: mapped.competitorsCount,
        competitorsDelta: mapped.competitorsDelta,
        priceChanges: mapped.priceChanges,
        priceChangesDelta: mapped.priceChangesDelta,
        activeAlerts: mapped.activeAlerts,
        activeAlertsNew: mapped.activeAlertsNew,
        markets: mapped.markets,
        marketNames: mapped.marketNames
      }
    } catch (e) {
      // Keep zeroed KPIs on error
    }
  }

  async function fetchCompetitors(platform?: string) {
    try {
      const res = await competitorsAPI.list(platform)
      const data = res.data.competitors as CompetitorResponse[]
      competitors.value = data.map(mapCompetitor)
    } catch (e) {
      competitors.value = []
    }
  }

  async function fetchLatestAlerts() {
    try {
      const res = await trendsAPI.getDemandSignals()
      const data = res.data as DemandSignalsResponse
      latestAlerts.value = mapDemandSignals(data)
    } catch (e) {
      latestAlerts.value = []
    }
  }

  async function fetchHotCategories() {
    try {
      const res = await trendsAPI.getTrending(5)
      const data = res.data as TrendingProductsResponse
      hotCategories.value = mapTrendingProducts(data).map(p => ({
        name: p.name,
        score: p.score || 0
      }))
    } catch (e) {
      // Keep existing data on error
    }
  }

  async function fetchProductScores() {
    try {
      const res = await productsAPI.getRecommendations()
      const data = res.data as Product[]
      productScores.value = data
    } catch (e) {
      productScores.value = []
    }
  }

  async function fetchSupplyAlerts() {
    try {
      const res = await trendsAPI.getDemandSignals()
      const data = res.data as DemandSignalsResponse
      const signals = mapDemandSignals(data)
      supplyAlerts.value = signals.map((s) => ({
        id: s.id,
        riskLevel: s.riskLevel,
        category: s.competitor,
        signalType: s.signalType,
        detail: '',
        market: '',
        time: s.time
      }))
    } catch (e) {
      supplyAlerts.value = []
    }
  }

  async function fetchTrendData() {
    try {
      const res = await trendsAPI.analyze('market trends')
      const data = res.data as { trend: string; opportunities: string[]; demand_signal: string }
      // Use API data if available
      if (data.trend) {
        trend30Days.value.xData = ['8/1', '8/4', '8/8', '8/12', '8/16', '8/20', '8/24', '8/28']
        // Keep chart data structure but could be extended
      }
    } catch (e) {
      // Keep existing data on error
    }
  }

  async function fetchPricingData() {
    try {
      const res = await strategyAPI.getPricing()
      const data = res.data as { sku_count?: number; underprice_alerts?: number; avg_profit_margin?: number }
      if (data) {
        pricingKpis.value = {
          skuCount: data.sku_count || 0,
          underpriceAlerts: data.underprice_alerts || 0,
          avgProfitMargin: data.avg_profit_margin || 0,
          profitDelta: 0,
          pendingAdjustments: 0
        }
      }
    } catch (e) {
      // Keep default/zeroed values on error
    }
  }

  // Fetch all dashboard data
  async function fetchDashboard() {
    loading.value = true
    error.value = null
    try {
      await Promise.all([
        fetchKpis(),
        fetchCompetitors(),
        fetchLatestAlerts(),
        fetchHotCategories(),
        fetchTrendData()
      ])
    } catch (e: any) {
      error.value = e.message || 'Failed to fetch dashboard data'
    } finally {
      loading.value = false
    }
  }

  return {
    // Dashboard
    kpis,
    marketOverview,
    trend30Days,
    hotCategories,
    latestAlerts,
    // Competitor
    competitors,
    // Product
    productScores,
    priceComparison,
    // Pricing
    pricingKpis,
    pricingTrend60Days,
    pricingSuggestions,
    profitSimulation,
    // Supply chain
    supplyKpis,
    supplyAlerts,
    riskRadar,
    // State
    loading,
    error,
    // Actions
    fetchKpis,
    fetchCompetitors,
    fetchLatestAlerts,
    fetchHotCategories,
    fetchProductScores,
    fetchSupplyAlerts,
    fetchTrendData,
    fetchPricingData,
    fetchDashboard
  }
})
