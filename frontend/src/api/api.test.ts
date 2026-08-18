import { describe, it, expect } from 'vitest'
import apiClient, { healthAPI, productsAPI, competitorsAPI, trendsAPI, strategyAPI } from './index'

describe('API module', () => {
  it('exports apiClient', () => {
    expect(apiClient).toBeDefined()
    expect(apiClient.get).toBeDefined()
    expect(apiClient.post).toBeDefined()
  })

  it('apiClient has baseURL configured', () => {
    expect(apiClient.defaults.baseURL).toBe('/api/v1')
  })

  it('healthAPI has check function', () => {
    expect(healthAPI.check).toBeDefined()
    expect(typeof healthAPI.check).toBe('function')
  })

  it('productsAPI has list and get functions', () => {
    expect(productsAPI.list).toBeDefined()
    expect(productsAPI.get).toBeDefined()
    expect(typeof productsAPI.list).toBe('function')
    expect(typeof productsAPI.get).toBe('function')
  })

  it('competitorsAPI has list and get functions', () => {
    expect(competitorsAPI.list).toBeDefined()
    expect(competitorsAPI.get).toBeDefined()
  })

  it('trendsAPI has list function', () => {
    expect(trendsAPI.list).toBeDefined()
  })

  it('strategyAPI has getPricing function', () => {
    expect(strategyAPI.getPricing).toBeDefined()
  })
})
