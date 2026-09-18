import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import MemberTabs from '../MemberTabs.vue'

describe('MemberTabs.vue', () => {
  it('renders all tabs', () => {
    const wrapper = mount(MemberTabs, {
      props: {
        currentTab: 'active',
      },
    })

    expect(wrapper.text()).toContain('Active')
    expect(wrapper.text()).toContain('Deactivated')
    expect(wrapper.text()).toContain('All')
  })

  it('marks active tab as current', () => {
    const wrapper = mount(MemberTabs, {
      props: {
        currentTab: 'active',
      },
    })

    const tabs = wrapper.findAll('li')
    expect(tabs[0].classes()).toContain('is-active')
    expect(tabs[1].classes()).not.toContain('is-active')
    expect(tabs[2].classes()).not.toContain('is-active')
  })

  it('marks deactivated tab as current', () => {
    const wrapper = mount(MemberTabs, {
      props: {
        currentTab: 'deactivated',
      },
    })

    const tabs = wrapper.findAll('li')
    expect(tabs[0].classes()).not.toContain('is-active')
    expect(tabs[1].classes()).toContain('is-active')
    expect(tabs[2].classes()).not.toContain('is-active')
  })

  it('marks all tab as current', () => {
    const wrapper = mount(MemberTabs, {
      props: {
        currentTab: 'all',
      },
    })

    const tabs = wrapper.findAll('li')
    expect(tabs[0].classes()).not.toContain('is-active')
    expect(tabs[1].classes()).not.toContain('is-active')
    expect(tabs[2].classes()).toContain('is-active')
  })

  it('emits selectTab event when tab is clicked', async () => {
    const wrapper = mount(MemberTabs, {
      props: {
        currentTab: 'active',
      },
    })

    const tabs = wrapper.findAll('li')
    await tabs[1].find('a').trigger('click')

    expect(wrapper.emitted('selectTab')).toBeTruthy()
    const emitted = wrapper.emitted('selectTab')?.[0]
    expect(emitted).toEqual(['deactivated'])
  })

  it('emits correct tab value for all tab', async () => {
    const wrapper = mount(MemberTabs, {
      props: {
        currentTab: 'active',
      },
    })

    const tabs = wrapper.findAll('li')
    await tabs[2].find('a').trigger('click')

    const emitted = wrapper.emitted('selectTab')?.[0]
    expect(emitted).toEqual(['all'])
  })
})
