import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ScoreCard from './ScoreCard.vue'

describe('ScoreCard', () => {
  it('renders product name and score', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        name: 'Apple Watch Ultra 2',
        score: 97,
        rank: 1,
        rating: 4.9,
        reviews: 6291,
        price: 799
      }
    })

    expect(wrapper.find('.score-name').text()).toBe('Apple Watch Ultra 2')
    expect(wrapper.find('.score-num').text()).toBe('97')
    expect(wrapper.find('.score-rank').text()).toBe('1')
  })

  it('renders badge with correct type', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        name: 'Test Product',
        score: 85,
        badge: '高利润',
        badgeType: 'good'
      }
    })

    expect(wrapper.find('.alert-badge').exists()).toBe(true)
    expect(wrapper.find('.alert-badge.good').exists()).toBe(true)
  })

  it('renders rank badge with gold class for rank 1', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        name: 'Gold Product',
        score: 100,
        rank: 1
      }
    })

    expect(wrapper.find('.score-rank.gold').exists()).toBe(true)
  })

  it('renders rank badge with silver class for rank 2', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        name: 'Silver Product',
        score: 90,
        rank: 2
      }
    })

    expect(wrapper.find('.score-rank.silver').exists()).toBe(true)
  })

  it('renders rank badge with bronze class for rank 3', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        name: 'Bronze Product',
        score: 85,
        rank: 3
      }
    })

    expect(wrapper.find('.score-rank.bronze').exists()).toBe(true)
  })

  it('renders description text', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        name: 'Test Product',
        score: 80,
        description: '综合得分 = 品类热度×0.4 + 利润空间×0.3 + 竞争度×0.3'
      }
    })

    expect(wrapper.find('.score-description').text()).toContain('综合得分')
  })

  it('formats large review numbers', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        name: 'Popular Product',
        score: 95,
        reviews: 12847
      }
    })

    expect(wrapper.text()).toContain('12.8k')
  })

  it('has hoverable class when hoverable prop is true', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        name: 'Hoverable Product',
        score: 88,
        hoverable: true
      }
    })

    expect(wrapper.find('.score-card.hoverable').exists()).toBe(true)
  })

  it('does not show badge when badge prop is not provided', () => {
    const wrapper = mount(ScoreCard, {
      props: {
        name: 'No Badge Product',
        score: 75
      }
    })

    expect(wrapper.find('.alert-badge').exists()).toBe(false)
  })
})
