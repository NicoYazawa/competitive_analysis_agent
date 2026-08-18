import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import Pricing from './Pricing.vue'

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

describe('Pricing', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders page title', () => {
    const wrapper = mount(Pricing)
    expect(wrapper.find('h1').text()).toBe('定价策略')
  })

  it('renders KPI cards', () => {
    const wrapper = mount(Pricing)
    const kpiCards = wrapper.findAll('.kpi-card')
    expect(kpiCards.length).toBe(4)
  })

  it('renders price trend chart', () => {
    const wrapper = mount(Pricing)
    expect(wrapper.find('.chart-container').exists()).toBe(true)
  })

  it('renders pricing suggestions', () => {
    const wrapper = mount(Pricing)
    expect(wrapper.text()).toContain('智能定价建议')
    expect(wrapper.text()).toContain('Anker 737')
  })

  it('renders profit simulation', () => {
    const wrapper = mount(Pricing)
    expect(wrapper.text()).toContain('利润模拟')
    expect(wrapper.text()).toContain('预估利润')
  })
})
