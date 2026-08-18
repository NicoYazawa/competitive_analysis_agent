import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ScoreCard from './ScoreCard.vue'

describe('ScoreCard', () => {
  it('renders title and value correctly', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        title: 'Total Products',
        value: 42,
        icon: '📦'
      }
    })

    expect(wrapper.find('.title').text()).toBe('Total Products')
    expect(wrapper.find('.value').text()).toBe('42')
    expect(wrapper.find('.icon').text()).toBe('📦')
  })

  it('renders zero value', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        title: 'Price Changes',
        value: 0,
        icon: '📊'
      }
    })

    expect(wrapper.find('.value').text()).toBe('0')
  })

  it('renders large numbers', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        title: 'Reviews',
        value: 999999,
        icon: '⭐'
      }
    })

    expect(wrapper.find('.value').text()).toBe('999999')
  })
})
