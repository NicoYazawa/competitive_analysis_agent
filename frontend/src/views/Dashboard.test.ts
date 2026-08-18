import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import Dashboard from './Dashboard.vue'

// Mock ECharts
vi.mock('echarts', () => ({
  init: vi.fn(() => ({
    setOption: vi.fn(),
    resize: vi.fn(),
    dispose: vi.fn()
  })),
  graphic: {
    LinearGradient: vi.fn()
  }
}))

describe('Dashboard', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders dashboard title', () => {
    const wrapper = mount(Dashboard)
    expect(wrapper.find('h1').text()).toBe('监控总览')
  })

  it('renders KPI cards', () => {
    const wrapper = mount(Dashboard)
    const kpiCards = wrapper.findAll('.kpi-card')
    expect(kpiCards.length).toBeGreaterThan(0)
  })

  it('renders trend chart', () => {
    const wrapper = mount(Dashboard)
    expect(wrapper.find('.chart-container').exists()).toBe(true)
  })

  it('renders market overview section', () => {
    const wrapper = mount(Dashboard)
    expect(wrapper.text()).toContain('市场概览')
    expect(wrapper.text()).toContain('北美市场')
  })

  it('renders latest alerts section', () => {
    const wrapper = mount(Dashboard)
    expect(wrapper.text()).toContain('最新预警')
    expect(wrapper.text()).toContain('Anker 737 Power Bank')
  })

  it('renders hot categories', () => {
    const wrapper = mount(Dashboard)
    expect(wrapper.text()).toContain('TOP 5 热门品类')
    expect(wrapper.text()).toContain('智能手表')
  })
})
