import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import TrendChart from './TrendChart.vue'

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

describe('TrendChart', () => {
  it('renders title and subtitle', () => {
    const wrapper = mount(TrendChart, {
      props: {
        title: '价格与评分趋势',
        subtitle: '近30天 · 全部监控品类'
      }
    })

    expect(wrapper.find('.chart-title').text()).toBe('价格与评分趋势')
    expect(wrapper.find('.chart-subtitle').text()).toBe('近30天 · 全部监控品类')
  })

  it('renders legend items', () => {
    const wrapper = mount(TrendChart, {
      props: {
        title: 'Test Chart',
        legend: [
          { label: '平均价格指数', color: '#2a78d6' },
          { label: '品类评分均值', color: '#eb6834' }
        ]
      }
    })

    const legendItems = wrapper.findAll('.legend-item')
    expect(legendItems).toHaveLength(2)
    expect(legendItems[0].text()).toContain('平均价格指数')
    expect(legendItems[1].text()).toContain('品类评分均值')
  })

  it('renders chart container', () => {
    const wrapper = mount(TrendChart, {
      props: {
        title: 'Test Chart'
      }
    })

    expect(wrapper.find('.chart-container').exists()).toBe(true)
    expect(wrapper.find('.chart-area').exists()).toBe(true)
  })

  it('renders without legend', () => {
    const wrapper = mount(TrendChart, {
      props: {
        title: 'No Legend Chart'
      }
    })

    expect(wrapper.find('.chart-legend').exists()).toBe(false)
  })

  it('accepts series data', () => {
    const wrapper = mount(TrendChart, {
      props: {
        title: 'Series Chart',
        xData: ['8/1', '8/8', '8/15'],
        series: [
          { name: '价格指数', data: [100, 95, 90], color: '#2a78d6' }
        ]
      }
    })

    expect(wrapper.exists()).toBe(true)
  })
})
