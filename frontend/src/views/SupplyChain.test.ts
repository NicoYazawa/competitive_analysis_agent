import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SupplyChain from './SupplyChain.vue'

describe('SupplyChain', () => {
  it('renders page title', () => {
    const wrapper = mount(SupplyChain)
    expect(wrapper.find('h1').text()).toBe('Supply Chain Alerts')
  })

  it('renders alert list container', () => {
    const wrapper = mount(SupplyChain)
    expect(wrapper.find('.alert-list').exists()).toBe(true)
  })
})
