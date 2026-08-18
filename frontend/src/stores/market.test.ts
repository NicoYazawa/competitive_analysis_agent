import { describe, it, expect } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useMarketStore } from './market'

describe('market store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initializes with empty state', () => {
    const store = useMarketStore()

    expect(store.products).toEqual([])
    expect(store.competitors).toEqual([])
    expect(store.trends).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.error).toBe(null)
  })

  it('has fetchProducts function', () => {
    const store = useMarketStore()
    expect(typeof store.fetchProducts).toBe('function')
  })

  it('has fetchCompetitors function', () => {
    const store = useMarketStore()
    expect(typeof store.fetchCompetitors).toBe('function')
  })

  it('has fetchTrends function', () => {
    const store = useMarketStore()
    expect(typeof store.fetchTrends).toBe('function')
  })
})
