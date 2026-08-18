import { describe, it, expect } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useMarketStore } from './market'

describe('market store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initializes with dashboard KPIs', () => {
    const store = useMarketStore()

    expect(store.kpis.competitorsCount).toBe(247)
    expect(store.kpis.priceChanges).toBe(1842)
    expect(store.kpis.activeAlerts).toBe(38)
    expect(store.kpis.markets).toBe(12)
  })

  it('has market overview data', () => {
    const store = useMarketStore()

    expect(store.marketOverview).toHaveLength(3)
    expect(store.marketOverview[0].name).toBe('北美市场')
    expect(store.marketOverview[0].growth).toBe('+18.2%')
  })

  it('has trend data for chart', () => {
    const store = useMarketStore()

    expect(store.trend30Days.xData).toBeDefined()
    expect(store.trend30Days.priceIndex).toBeDefined()
    expect(store.trend30Days.scoreAvg).toBeDefined()
    expect(store.trend30Days.xData.length).toBe(store.trend30Days.priceIndex.length)
  })

  it('has hot categories', () => {
    const store = useMarketStore()

    expect(store.hotCategories).toHaveLength(5)
    expect(store.hotCategories[0].name).toBe('智能手表')
    expect(store.hotCategories[0].score).toBe(92)
  })

  it('has latest alerts', () => {
    const store = useMarketStore()

    expect(store.latestAlerts).toHaveLength(4)
    expect(store.latestAlerts[0].competitor).toBe('Anker 737 Power Bank')
    expect(store.latestAlerts[0].riskLevel).toBe('critical')
  })

  it('has competitors list', () => {
    const store = useMarketStore()

    expect(store.competitors).toHaveLength(6)
    expect(store.competitors[0].name).toBe('Anker 737 Power Bank')
    expect(store.competitors[0].rating).toBe(4.8)
    expect(store.competitors[0].currentPrice).toBe(139.99)
  })

  it('has product scores', () => {
    const store = useMarketStore()

    expect(store.productScores).toHaveLength(6)
    expect(store.productScores[0].score).toBe(97)
    expect(store.productScores[0].badge).toBe('高利润')
  })

  it('has price comparison data', () => {
    const store = useMarketStore()

    expect(store.priceComparison).toHaveLength(5)
    expect(store.priceComparison[0].suggestedPrice).toBe(769)
  })

  it('has pricing KPIs', () => {
    const store = useMarketStore()

    expect(store.pricingKpis.skuCount).toBe(1284)
    expect(store.pricingKpis.avgProfitMargin).toBe(34)
    expect(store.pricingKpis.pendingAdjustments).toBe(8)
  })

  it('has pricing trend data', () => {
    const store = useMarketStore()

    expect(store.pricingTrend60Days.anker).toBeDefined()
    expect(store.pricingTrend60Days.jackery).toBeDefined()
    expect(store.pricingTrend60Days.ugreen).toBeDefined()
  })

  it('has pricing suggestions', () => {
    const store = useMarketStore()

    expect(store.pricingSuggestions).toHaveLength(3)
    expect(store.pricingSuggestions[0].severity).toBe('critical')
  })

  it('has profit simulation', () => {
    const store = useMarketStore()

    expect(store.profitSimulation.suggestedPrice).toBe(129)
    expect(store.profitSimulation.profit).toBe(14.15)
    expect(store.profitSimulation.profitRate).toBe(11)
  })

  it('has supply chain KPIs', () => {
    const store = useMarketStore()

    expect(store.supplyKpis.highRisk).toBe(8)
    expect(store.supplyKpis.mediumRisk).toBe(21)
    expect(store.supplyKpis.resolved).toBe(47)
  })

  it('has supply alerts', () => {
    const store = useMarketStore()

    expect(store.supplyAlerts).toHaveLength(6)
    expect(store.supplyAlerts[0].riskLevel).toBe('critical')
    expect(store.supplyAlerts[0].signalType).toBe('原材料涨价')
  })

  it('has risk radar data', () => {
    const store = useMarketStore()

    expect(store.riskRadar.rawMaterials).toHaveLength(3)
    expect(store.riskRadar.logistics).toHaveLength(3)
    expect(store.riskRadar.policy).toHaveLength(3)
  })

  it('has loading and error state', () => {
    const store = useMarketStore()

    expect(store.loading).toBe(false)
    expect(store.error).toBe(null)
  })
})
