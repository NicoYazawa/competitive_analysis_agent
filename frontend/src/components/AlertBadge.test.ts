import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AlertBadge from './AlertBadge.vue'

describe('AlertBadge', () => {
  it('renders with warning type', () => {
    const wrapper = mount(AlertBadge, {
      props: {
        count: 5,
        type: 'warning' as const
      }
    })

    expect(wrapper.find('.alert-badge').exists()).toBe(true)
    expect(wrapper.find('.alert-badge.warning').exists()).toBe(true)
    expect(wrapper.text()).toBe('5')
  })

  it('renders with danger type', () => {
    const wrapper = mount(AlertBadge, {
      props: {
        count: 10,
        type: 'danger' as const
      }
    })

    expect(wrapper.find('.alert-badge.danger').exists()).toBe(true)
  })

  it('renders with info type', () => {
    const wrapper = mount(AlertBadge, {
      props: {
        count: 3,
        type: 'info' as const
      }
    })

    expect(wrapper.find('.alert-badge.info').exists()).toBe(true)
  })

  it('renders with success type', () => {
    const wrapper = mount(AlertBadge, {
      props: {
        count: 0,
        type: 'success' as const
      }
    })

    expect(wrapper.find('.alert-badge.success').exists()).toBe(true)
  })

  it('renders zero count', () => {
    const wrapper = mount(AlertBadge, {
      props: {
        count: 0,
        type: 'info' as const
      }
    })

    expect(wrapper.text()).toBe('0')
  })
})
