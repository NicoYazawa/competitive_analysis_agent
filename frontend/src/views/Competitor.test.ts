import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import Competitor from './Competitor.vue'

describe('Competitor', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders page title', () => {
    const wrapper = mount(Competitor)
    expect(wrapper.find('h1').text()).toBe('竞品情报')
  })

  it('renders filter bar', () => {
    const wrapper = mount(Competitor)
    expect(wrapper.find('.filter-bar').exists()).toBe(true)
    expect(wrapper.find('.filter-search').exists()).toBe(true)
    expect(wrapper.findAll('.filter-select').length).toBeGreaterThan(0)
  })

  it('renders competitor table', () => {
    const wrapper = mount(Competitor)
    expect(wrapper.find('.data-table').exists()).toBe(true)
    expect(wrapper.text()).toContain('Anker 737 Power Bank')
  })

  it('renders pagination', () => {
    const wrapper = mount(Competitor)
    expect(wrapper.find('.pagination').exists()).toBe(true)
  })

  it('renders filter chips', () => {
    const wrapper = mount(Competitor)
    const chips = wrapper.findAll('.filter-chip')
    expect(chips.length).toBe(3)
  })
})
