import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ErrorToast from '../ErrorToast.vue'

describe('ErrorToast.vue', () => {
  it('renders error message', async () => {
    const wrapper = mount(ErrorToast, {
      props: {
        message: 'Test error message',
        autoShow: true,
      },
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Test error message')
  })

  it('shows error toast on mount', async () => {
    const wrapper = mount(ErrorToast, {
      props: {
        message: 'Test error',
        autoShow: true,
      },
    })

    await wrapper.vm.$nextTick()
    const toast = wrapper.find('.notification')
    expect(toast.exists()).toBe(true)
  })

  it('can be hidden when autoShow is false', async () => {
    const wrapper = mount(ErrorToast, {
      props: {
        message: 'Test error',
        autoShow: false,
      },
    })

    await wrapper.vm.$nextTick()
    const toast = wrapper.find('.notification')
    expect(toast.exists()).toBe(false)
  })

  it('closes toast when close button is clicked', async () => {
    const wrapper = mount(ErrorToast, {
      props: {
        message: 'Test error',
        autoShow: true,
      },
    })

    await wrapper.vm.$nextTick()
    const closeButton = wrapper.find('.delete')
    await closeButton.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('auto-dismisses after specified duration', async () => {
    const wrapper = mount(ErrorToast, {
      props: {
        message: 'Test error',
        duration: 50,
      },
    })

    await new Promise((resolve) => setTimeout(resolve, 100))

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('does not auto-dismiss when duration is 0', async () => {
    const wrapper = mount(ErrorToast, {
      props: {
        message: 'Test error',
        duration: 0,
      },
    })

    await new Promise((resolve) => setTimeout(resolve, 100))

    expect(wrapper.emitted('close')).toBeFalsy()
  })
})
