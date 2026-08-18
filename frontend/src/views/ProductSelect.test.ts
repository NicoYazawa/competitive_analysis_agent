import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import ProductSelect from './ProductSelect.vue'

// Mock store
vi.mock('@/stores/market', () => ({
  useMarketStore: vi.fn(() => ({
    fetchProducts: vi.fn(),
    products: [],
    loading: false,
    error: null
  }))
}))

describe('ProductSelect', () => {
  it('renders page title', () => {
    const wrapper = mount(ProductSelect)
    expect(wrapper.find('h1').text()).toBe('Product Selection')
  })

  it('renders search input', () => {
    const wrapper = mount(ProductSelect)
    expect(wrapper.find('.search-input').exists()).toBe(true)
  })

  it('renders product list container', () => {
    const wrapper = mount(ProductSelect)
    expect(wrapper.find('.product-list').exists()).toBe(true)
  })
})
