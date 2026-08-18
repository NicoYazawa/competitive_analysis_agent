import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AlertBadge from './AlertBadge.vue'

describe('AlertBadge', () => {
  it('renders with critical type', () => {
    const wrapper = mount(AlertBadge, {
      props: {
        type: 'critical'
      },
      slots: {
        default: '严重'
      }
    })

    expect(wrapper.find('.alert-badge.critical').exists()).toBe(true)
    expect(wrapper.text()).toBe('严重')
  })

  it('renders with warning type', () => {
    const wrapper = mount(AlertBadge, {
      props: {
        type: 'warning'
      },
      slots: {
        default: '警告'
      }
    })

    expect(wrapper.find('.alert-badge.warning').exists()).toBe(true)
    expect(wrapper.text()).toBe('警告')
  })

  it('renders with good type', () => {
    const wrapper = mount(AlertBadge, {
      props: {
        type: 'good'
      },
      slots: {
        default: '正常'
      }
    })

    expect(wrapper.find('.alert-badge.good').exists()).toBe(true)
    expect(wrapper.text()).toBe('正常')
  })

  it('renders with info type', () => {
    const wrapper = mount(AlertBadge, {
      props: {
        type: 'info'
      },
      slots: {
        default: '情报'
      }
    })

    expect(wrapper.find('.alert-badge.info').exists()).toBe(true)
    expect(wrapper.text()).toBe('情报')
  })

  it('renders dot when dot prop is true', () => {
    const wrapper = mount(AlertBadge, {
      props: {
        type: 'critical',
        dot: true
      },
      slots: {
        default: '严重'
      }
    })

    expect(wrapper.find('.alert-dot').exists()).toBe(true)
  })

  it('does not render dot when dot prop is false', () => {
    const wrapper = mount(AlertBadge, {
      props: {
        type: 'warning',
        dot: false
      },
      slots: {
        default: '警告'
      }
    })

    expect(wrapper.find('.alert-dot').exists()).toBe(false)
  })

  it('applies default type when type is not specified', () => {
    const wrapper = mount(AlertBadge, {
      slots: {
        default: 'Test'
      }
    })

    expect(wrapper.find('.alert-badge').exists()).toBe(true)
  })
})
