import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Competitor from './Competitor.vue'

// Mock store
vi.mock('@/stores/market', () => ({
  useMarketStore: vi.fn(() => ({
    fetchCompetitors: vi.fn(),
    competitors: [],
    loading: false,
    error: null
  }))
}))

vi.mock('@/components/AlertBadge.vue', () => ({
  default: { template: '<div class="alert-badge-mock"></div>', props: ['count', 'type'] }
}))

describe('Competitor', () => {
  it('renders page title', () => {
    const wrapper = mount(Competitor, {
      global: { stubs: { AlertBadge: true } }
    })
    expect(wrapper.find('h1').text()).toBe('Competitor Intelligence')
  })

  it('renders competitor list container', () => {
    const wrapper = mount(Competitor, {
      global: { stubs: { AlertBadge: true } }
    })
    expect(wrapper.find('.competitor-list').exists()).toBe(true)
  })
})
