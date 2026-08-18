import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Dashboard from './Dashboard.vue'

// Mock child components
vi.mock('@/components/ScoreCard.vue', () => ({
  default: { template: '<div class="score-card-mock">{{ title }}:{{ value }}</div>', props: ['title', 'value', 'icon'] }
}))

vi.mock('@/components/TrendChart.vue', () => ({
  default: { template: '<div class="trend-chart-mock">{{ title }}</div>', props: ['data', 'title'] }
}))

describe('Dashboard', () => {
  it('renders dashboard title', () => {
    const wrapper = mount(Dashboard, {
      global: { stubs: { ScoreCard: true, TrendChart: true } }
    })
    expect(wrapper.find('h1').text()).toBe('Market Overview')
  })

  it('renders metrics section', () => {
    const wrapper = mount(Dashboard, {
      global: { stubs: { ScoreCard: true, TrendChart: true } }
    })
    expect(wrapper.find('.metrics').exists()).toBe(true)
  })

  it('renders chart container', () => {
    const wrapper = mount(Dashboard, {
      global: { stubs: { ScoreCard: true, TrendChart: true } }
    })
    expect(wrapper.find('.chart-container').exists()).toBe(true)
  })
})
