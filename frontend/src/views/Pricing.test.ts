import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Pricing from './Pricing.vue'

// Mock PriceTable
vi.mock('@/components/PriceTable.vue', () => ({
  default: { template: '<div class="price-table-mock"></div>', props: ['data'] }
}))

describe('Pricing', () => {
  it('renders page title', () => {
    const wrapper = mount(Pricing, {
      global: { stubs: { PriceTable: true } }
    })
    expect(wrapper.find('h1').text()).toBe('Pricing Strategy')
  })
})
