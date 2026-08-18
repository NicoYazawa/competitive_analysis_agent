import { describe, it, expect } from 'vitest'
import router from './index'

describe('router', () => {
  it('has routes defined', () => {
    expect(router.getRoutes()).toBeDefined()
    expect(router.getRoutes().length).toBeGreaterThan(0)
  })

  it('has Dashboard route', () => {
    const routes = router.getRoutes()
    const dashboardRoute = routes.find(r => r.name === 'Dashboard')
    expect(dashboardRoute).toBeDefined()
    expect(dashboardRoute?.path).toBe('/')
  })

  it('has ProductSelect route', () => {
    const routes = router.getRoutes()
    const productRoute = routes.find(r => r.name === 'ProductSelect')
    expect(productRoute).toBeDefined()
    expect(productRoute?.path).toBe('/products')
  })

  it('has Pricing route', () => {
    const routes = router.getRoutes()
    const pricingRoute = routes.find(r => r.name === 'Pricing')
    expect(pricingRoute).toBeDefined()
    expect(pricingRoute?.path).toBe('/pricing')
  })

  it('has Competitor route', () => {
    const routes = router.getRoutes()
    const competitorRoute = routes.find(r => r.name === 'Competitor')
    expect(competitorRoute).toBeDefined()
    expect(competitorRoute?.path).toBe('/competitors')
  })

  it('has SupplyChain route', () => {
    const routes = router.getRoutes()
    const supplyChainRoute = routes.find(r => r.name === 'SupplyChain')
    expect(supplyChainRoute).toBeDefined()
    expect(supplyChainRoute?.path).toBe('/supply-chain')
  })
})
