import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import PriceTable from './PriceTable.vue'

describe('PriceTable', () => {
  const basicColumns = [
    { key: 'name', label: '产品' },
    { key: 'category', label: '品类' },
    { key: 'price', label: '价格', format: 'currency' as const }
  ]

  it('renders table headers', () => {
    const wrapper = mount(PriceTable, {
      props: {
        columns: basicColumns,
        data: []
      }
    })

    expect(wrapper.findAll('th').at(0)?.text()).toBe('产品')
    expect(wrapper.findAll('th').at(1)?.text()).toBe('品类')
    expect(wrapper.findAll('th').at(2)?.text()).toBe('价格')
  })

  it('renders empty table body', () => {
    const wrapper = mount(PriceTable, {
      props: {
        columns: basicColumns,
        data: []
      }
    })

    expect(wrapper.find('tbody tr').exists()).toBe(false)
  })

  it('renders table with data rows', () => {
    const data = [
      { id: '1', name: 'Apple Watch', category: '智能手表', price: 799 },
      { id: '2', name: 'Anker 737', category: '便携储能', price: 139.99 }
    ]

    const wrapper = mount(PriceTable, {
      props: {
        columns: basicColumns,
        data
      }
    })

    const rows = wrapper.findAll('tbody tr')
    expect(rows).toHaveLength(2)
    expect(rows[0].find('td:nth-child(1)').text()).toBe('Apple Watch')
    expect(rows[1].find('td:nth-child(1)').text()).toBe('Anker 737')
  })

  it('formats currency columns with $', () => {
    const data = [
      { name: 'Test', price: 129.99 }
    ]

    const wrapper = mount(PriceTable, {
      props: {
        columns: [{ key: 'name', label: 'Name' }, { key: 'price', label: 'Price', format: 'currency' as const }],
        data
      }
    })

    expect(wrapper.find('.num').text()).toBe('$129.99')
  })

  it('formats numbers with commas', () => {
    const data = [
      { name: 'Popular Item', reviews: 12847 }
    ]

    const columns = [
      { key: 'name', label: 'Name' },
      { key: 'reviews', label: 'Reviews', format: 'number' as const }
    ]

    const wrapper = mount(PriceTable, {
      props: { columns, data }
    })

    expect(wrapper.find('.num').text()).toBe('12,847')
  })

  it('renders card title when provided', () => {
    const wrapper = mount(PriceTable, {
      props: {
        columns: basicColumns,
        data: [],
        title: '价格对比表'
      }
    })

    expect(wrapper.find('.card-title').text()).toBe('价格对比表')
  })

  it('renders pagination when enabled', () => {
    const data = Array.from({ length: 15 }, (_, i) => ({
      id: String(i + 1),
      name: `Product ${i + 1}`,
      price: 100 + i
    }))

    const wrapper = mount(PriceTable, {
      props: {
        columns: basicColumns,
        data,
        pagination: true,
        total: 15,
        pageSize: 10
      }
    })

    expect(wrapper.find('.pagination').exists()).toBe(true)
    expect(wrapper.find('.pagination-info').text()).toContain('1-10')
    expect(wrapper.find('.pagination-info').text()).toContain('15')
  })

  it('pagination shows correct page info', () => {
    const wrapper = mount(PriceTable, {
      props: {
        columns: basicColumns,
        data: [],
        pagination: true,
        total: 50,
        pageSize: 10
      }
    })

    expect(wrapper.find('.pagination-info').text()).toContain('50')
  })
})
