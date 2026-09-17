import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import HealthIndicator from '../HealthIndicator.vue'
import { useHealthStore } from '../../stores/health'

describe('HealthIndicator.vue', () => {
  let pinia: any

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)
  })

  it('renders healthy state with green checkmark', () => {
    const store = useHealthStore()
    store.status = 'healthy'

    const wrapper = mount(HealthIndicator, {
      global: {
        plugins: [pinia]
      }
    })

    // Check for green checkmark class
    const icon = wrapper.find('i.fa-check-circle')
    expect(icon.exists()).toBe(true)

    // Check for healthy class on container
    expect(wrapper.find('.health-indicator.is-healthy').exists()).toBe(true)
  })

  it('renders unhealthy state with red X', () => {
    const store = useHealthStore()
    store.status = 'unhealthy'

    const wrapper = mount(HealthIndicator, {
      global: {
        plugins: [pinia]
      }
    })

    // Check for red X icon
    const icon = wrapper.find('i.fa-times-circle')
    expect(icon.exists()).toBe(true)

    // Check for unhealthy class on container
    expect(wrapper.find('.health-indicator.is-unhealthy').exists()).toBe(true)
  })

  it('renders checking state with spinner', () => {
    const store = useHealthStore()
    store.status = 'checking'

    const wrapper = mount(HealthIndicator, {
      global: {
        plugins: [pinia]
      }
    })

    // Check for spinner icon
    const icon = wrapper.find('i.fa-spinner')
    expect(icon.exists()).toBe(true)

    // Check for checking class on container
    expect(wrapper.find('.health-indicator.is-checking').exists()).toBe(true)
  })

  it('displays status text in title attribute for accessibility', () => {
    const store = useHealthStore()
    store.status = 'healthy'

    let wrapper = mount(HealthIndicator, {
      global: {
        plugins: [pinia]
      }
    })
    expect(wrapper.find('.health-indicator').attributes('title')).toContain('healthy')

    store.status = 'unhealthy'
    wrapper = mount(HealthIndicator, {
      global: {
        plugins: [pinia]
      }
    })
    expect(wrapper.find('.health-indicator').attributes('title')).toContain('unreachable')

    store.status = 'checking'
    wrapper = mount(HealthIndicator, {
      global: {
        plugins: [pinia]
      }
    })
    expect(wrapper.find('.health-indicator').attributes('title')).toContain('Checking')
  })

  it('displays screen-reader-only text for accessibility', () => {
    const store = useHealthStore()
    store.status = 'healthy'

    const wrapper = mount(HealthIndicator, {
      global: {
        plugins: [pinia]
      }
    })

    // Check for sr-only text
    const srText = wrapper.find('.sr-only')
    expect(srText.exists()).toBe(true)
    expect(srText.text()).toContain('Backend is healthy')
  })

  it('reacts to store status changes', async () => {
    const store = useHealthStore()
    store.status = 'healthy'

    const wrapper = mount(HealthIndicator, {
      global: {
        plugins: [pinia]
      }
    })

    expect(wrapper.find('.health-indicator.is-healthy').exists()).toBe(true)

    // Change to unhealthy
    store.status = 'unhealthy'
    await wrapper.vm.$nextTick()
    expect(wrapper.find('.health-indicator.is-unhealthy').exists()).toBe(true)

    // Change to checking
    store.status = 'checking'
    await wrapper.vm.$nextTick()
    expect(wrapper.find('.health-indicator.is-checking').exists()).toBe(true)
  })

  it('displays error message when unhealthy', async () => {
    const store = useHealthStore()
    store.status = 'unhealthy'
    store.checkError = 'Connection timeout'

    const wrapper = mount(HealthIndicator, {
      global: {
        plugins: [pinia]
      }
    })

    // Check title includes error
    const title = wrapper.find('.health-indicator').attributes('title')
    expect(title).toContain('Connection timeout')
  })

  it('has proper icon visibility and accessibility', async () => {
    const store = useHealthStore()
    store.status = 'healthy'

    const wrapper = mount(HealthIndicator, {
      global: {
        plugins: [pinia]
      }
    })

    // Icon should have aria-hidden for screen readers
    const icon = wrapper.find('i.fa-check-circle')
    expect(icon.exists()).toBe(true)
    expect(icon.attributes('aria-hidden')).toBe('true')

    // SR text should exist
    const srText = wrapper.find('.sr-only')
    expect(srText.exists()).toBe(true)
    expect(srText.text()).toContain('Backend is healthy')
  })
})
