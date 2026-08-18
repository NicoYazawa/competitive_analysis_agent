import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import SupplyChain from './SupplyChain.vue'

describe('SupplyChain', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders page title', () => {
    const wrapper = mount(SupplyChain)
    expect(wrapper.find('h1').text()).toBe('供应链预警')
  })

  it('renders KPI cards', () => {
    const wrapper = mount(SupplyChain)
    const kpiCards = wrapper.findAll('.kpi-card')
    expect(kpiCards.length).toBe(4)
  })

  it('renders supply alert table', () => {
    const wrapper = mount(SupplyChain)
    expect(wrapper.text()).toContain('预警列表')
    expect(wrapper.text()).toContain('便携储能')
  })

  it('renders risk radar cards', () => {
    const wrapper = mount(SupplyChain)
    expect(wrapper.text()).toContain('原材料风险')
    expect(wrapper.text()).toContain('物流风险')
    expect(wrapper.text()).toContain('政策风险')
  })

  it('renders filters', () => {
    const wrapper = mount(SupplyChain)
    expect(wrapper.findAll('.filter-select').length).toBe(2)
  })
})
