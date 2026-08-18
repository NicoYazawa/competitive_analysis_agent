import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import PriceTable from './PriceTable.vue'

describe('PriceTable', () => {
  it('renders table headers', () => {
    const wrapper = mount(PriceTable, {
      props: {
        data: []
      }
    })

    expect(wrapper.find('th:nth-child(1)').text()).toBe('Competitor')
    expect(wrapper.find('th:nth-child(2)').text()).toBe('Price')
    expect(wrapper.find('th:nth-child(3)').text()).toBe('Rating')
    expect(wrapper.find('th:nth-child(4)').text()).toBe('Reviews')
  })

  it('renders empty table', () => {
    const wrapper = mount(PriceTable, {
      props: {
        data: []
      }
    })

    expect(wrapper.find('tbody tr').exists()).toBe(false)
  })

  it('renders table with data', () => {
    const data = [
      { id: '1', name: 'Amazon', price: 99.99, rating: 4.5, reviews: 1000 },
      { id: '2', name: 'eBay', price: 89.99, rating: 4.2, reviews: 500 }
    ]

    const wrapper = mount(PriceTable, {
      props: { data }
    })

    const rows = wrapper.findAll('tbody tr')
    expect(rows).toHaveLength(2)

    expect(rows[0].find('td:nth-child(1)').text()).toBe('Amazon')
    expect(rows[0].find('td:nth-child(2)').text()).toBe('$99.99')
    expect(rows[0].find('td:nth-child(3)').text()).toBe('4.5 ★')
    expect(rows[0].find('td:nth-child(4)').text()).toBe('1000')
  })
})
