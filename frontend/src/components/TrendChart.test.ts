import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import TrendChart from './TrendChart.vue'

// Mock ECharts
vi.mock('echarts', () => ({
  init: vi.fn(() => ({
    setOption: vi.fn()
  }))
}))

describe('TrendChart', () => {
  it('renders title correctly', () => {
    const wrapper = mount(TrendChart, {
      props: {
        data: [],
        title: 'Market Trends'
      }
    })

    expect(wrapper.find('h3').text()).toBe('Market Trends')
  })

  it('renders chart container', () => {
    const wrapper = mount(TrendChart, {
      props: {
        data: [],
        title: 'Test Chart'
      }
    })

    expect(wrapper.find('.chart').exists()).toBe(true)
  })

  it('accepts empty data array', () => {
    const wrapper = mount(TrendChart, {
      props: {
        data: [],
        title: 'Empty Chart'
      }
    })

    expect(wrapper.exists()).toBe(true)
  })
})
