import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ConfirmDialog from '../ConfirmDialog.vue'

describe('ConfirmDialog.vue', () => {
  it('renders dialog when isOpen is true', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        isOpen: true,
        title: 'Confirm Action',
        message: 'Are you sure?',
      },
    })

    expect(wrapper.find('.modal').classes()).toContain('is-active')
    expect(wrapper.text()).toContain('Confirm Action')
    expect(wrapper.text()).toContain('Are you sure?')
  })

  it('does not render dialog when isOpen is false', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        isOpen: false,
        title: 'Confirm Action',
        message: 'Are you sure?',
      },
    })

    await wrapper.vm.$nextTick()
    const modal = wrapper.find('.modal')
    expect(modal.exists()).toBe(false)
  })

  it('emits confirm event when confirm button is clicked', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        isOpen: true,
        title: 'Confirm Action',
        message: 'Are you sure?',
      },
    })

    const confirmButton = wrapper.findAll('button').find((b) => b.text().includes('Confirm'))
    await confirmButton?.trigger('click')

    expect(wrapper.emitted('confirm')).toBeTruthy()
  })

  it('emits cancel event when cancel button is clicked', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        isOpen: true,
        title: 'Confirm Action',
        message: 'Are you sure?',
      },
    })

    const cancelButton = wrapper.findAll('button').find((b) => b.text() === 'Cancel')
    await cancelButton?.trigger('click')

    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('emits cancel event when close button is clicked', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        isOpen: true,
        title: 'Confirm Action',
        message: 'Are you sure?',
      },
    })

    const closeButton = wrapper.find('.modal-card-head .delete')
    await closeButton.trigger('click')

    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('emits cancel event when background is clicked', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        isOpen: true,
        title: 'Confirm Action',
        message: 'Are you sure?',
      },
    })

    const background = wrapper.find('.modal-background')
    await background.trigger('click')

    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('disables buttons when isLoading is true', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        isOpen: true,
        title: 'Confirm Action',
        message: 'Are you sure?',
        isLoading: true,
      },
    })

    const buttons = wrapper.findAll('button').filter((b) => b.classes().includes('button'))
    buttons.forEach((button) => {
      expect(button.attributes('disabled')).toBeDefined()
    })
  })

  it('uses custom button labels', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        isOpen: true,
        title: 'Confirm Action',
        message: 'Are you sure?',
        confirmButtonLabel: 'Delete',
        cancelButtonLabel: 'Keep',
      },
    })

    expect(wrapper.text()).toContain('Delete')
    expect(wrapper.text()).toContain('Keep')
  })
})
