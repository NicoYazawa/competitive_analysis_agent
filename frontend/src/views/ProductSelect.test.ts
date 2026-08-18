import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import ProductSelect from './ProductSelect.vue'

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

describe('ProductSelect', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders page title', () => {
    const wrapper = mount(ProductSelect)
    expect(wrapper.find('h1').text()).toBe('选品分析')
  })

  it('renders product score cards', () => {
    const wrapper = mount(ProductSelect)
    const scoreCards = wrapper.findAll('.score-card')
    expect(scoreCards.length).toBeGreaterThan(0)
  })

  it('renders price comparison table', () => {
    const wrapper = mount(ProductSelect)
    expect(wrapper.text()).toContain('价格对比表')
    expect(wrapper.text()).toContain('Apple Watch Ultra 2')
  })

  it('renders section head with filter', () => {
    const wrapper = mount(ProductSelect)
    expect(wrapper.find('.section-head').exists()).toBe(true)
    expect(wrapper.find('.filter-select').exists()).toBe(true)
  })

  it('renders score card with rank badges', () => {
    const wrapper = mount(ProductSelect)
    expect(wrapper.find('.score-rank').exists()).toBe(true)
    expect(wrapper.text()).toContain('智能手表')
  })
})
