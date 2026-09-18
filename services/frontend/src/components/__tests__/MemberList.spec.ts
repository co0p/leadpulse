import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import MemberList from '../MemberList.vue'
import type { MemberDTO } from '../../stores/members'

const router = createRouter({
  history: createMemoryHistory(),
  routes: [
    {
      path: '/members/:id/edit',
      name: 'EditMember',
      component: { template: '<div></div>' },
    },
  ],
})

describe('MemberList.vue', () => {
  it('renders empty state when no members', () => {
    const wrapper = mount(MemberList, {
      props: {
        members: [],
      },
      global: {
        plugins: [router],
      },
    })

    expect(wrapper.text()).toContain('No members found')
  })

  it('renders member table with members', async () => {
    const members: MemberDTO[] = [
      {
        id: '1',
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Senior',
        status: 'active',
        createdAt: '2026-09-18T10:00:00Z',
      },
    ]

    const wrapper = mount(MemberList, {
      props: {
        members,
      },
      global: {
        plugins: [router],
        stubs: {
          RouterLink: { template: '<a href="#"><slot /></a>' },
        },
      },
    })

    expect(wrapper.find('table').exists()).toBe(true)
    expect(wrapper.text()).toContain('Alice Smith')
    expect(wrapper.text()).toContain('Senior')
    expect(wrapper.text()).toContain('Active')
  })

  it('shows deactivate button for active member', async () => {
    const members: MemberDTO[] = [
      {
        id: '1',
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Senior',
        status: 'active',
        createdAt: '2026-09-18T10:00:00Z',
      },
    ]

    const wrapper = mount(MemberList, {
      props: {
        members,
      },
      global: {
        plugins: [router],
        stubs: {
          RouterLink: { template: '<a href="#"><slot /></a>' },
        },
      },
    })

    expect(wrapper.text()).toContain('Deactivate')
    expect(wrapper.text()).not.toContain('Reactivate')
  })

  it('shows reactivate button for deactivated member', async () => {
    const members: MemberDTO[] = [
      {
        id: '1',
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Senior',
        status: 'deactivated',
        createdAt: '2026-09-18T10:00:00Z',
      },
    ]

    const wrapper = mount(MemberList, {
      props: {
        members,
      },
      global: {
        plugins: [router],
        stubs: {
          RouterLink: { template: '<a href="#"><slot /></a>' },
        },
      },
    })

    const reactivateButtons = wrapper.findAll('button').filter((b) => b.text().includes('Reactivate'))
    expect(reactivateButtons.length).toBeGreaterThan(0)
    
    const deactivateButtons = wrapper.findAll('button').filter((b) => b.text() === 'Deactivate')
    expect(deactivateButtons.length).toBe(0)
  })

  it('emits deactivate event when deactivate is confirmed', async () => {
    const members: MemberDTO[] = [
      {
        id: '1',
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Senior',
        status: 'active',
        createdAt: '2026-09-18T10:00:00Z',
      },
    ]

    const wrapper = mount(MemberList, {
      props: {
        members,
      },
      global: {
        plugins: [router],
        stubs: {
          RouterLink: { template: '<a href="#"><slot /></a>' },
          ConfirmDialog: { template: '<div @confirm="$emit(\'confirm\')" @cancel="$emit(\'cancel\')"></div>' },
        },
      },
    })

    const deactivateButton = wrapper.findAll('button').filter((b) => b.text().includes('Deactivate') && !b.classes().includes('is-light'))[0]
    await deactivateButton.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.showConfirmDialog).toBe(true)
    expect(wrapper.vm.actionType).toBe('deactivate')
    
    // Call handleConfirm directly
    await wrapper.vm.handleConfirm()

    expect(wrapper.emitted('deactivate')).toBeTruthy()
    const emitted = wrapper.emitted('deactivate')?.[0]
    expect(emitted).toEqual(['1'])
  })

  it('emits reactivate event when reactivate is confirmed', async () => {
    const members: MemberDTO[] = [
      {
        id: '1',
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Senior',
        status: 'deactivated',
        createdAt: '2026-09-18T10:00:00Z',
      },
    ]

    const wrapper = mount(MemberList, {
      props: {
        members,
      },
      global: {
        plugins: [router],
        stubs: {
          RouterLink: { template: '<a href="#"><slot /></a>' },
          ConfirmDialog: { template: '<div @confirm="$emit(\'confirm\')" @cancel="$emit(\'cancel\')"></div>' },
        },
      },
    })

    const reactivateButton = wrapper.findAll('button').filter((b) => b.text().includes('Reactivate') && !b.classes().includes('is-light'))[0]
    await reactivateButton.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.showConfirmDialog).toBe(true)
    expect(wrapper.vm.actionType).toBe('reactivate')
    
    // Call handleConfirm directly
    await wrapper.vm.handleConfirm()

    expect(wrapper.emitted('reactivate')).toBeTruthy()
    const emitted = wrapper.emitted('reactivate')?.[0]
    expect(emitted).toEqual(['1'])
  })

  it('closes dialog when cancel is clicked', async () => {
    const members: MemberDTO[] = [
      {
        id: '1',
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Senior',
        status: 'active',
        createdAt: '2026-09-18T10:00:00Z',
      },
    ]

    const wrapper = mount(MemberList, {
      props: {
        members,
      },
      global: {
        plugins: [router],
        stubs: {
          RouterLink: { template: '<a href="#"><slot /></a>' },
          ConfirmDialog: { template: '<div @confirm="$emit(\'confirm\')" @cancel="$emit(\'cancel\')"></div>' },
        },
      },
    })

    const deactivateButton = wrapper.findAll('button').filter((b) => b.text().includes('Deactivate') && !b.classes().includes('is-light'))[0]
    await deactivateButton.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.showConfirmDialog).toBe(true)

    // Call cancelAction directly
    wrapper.vm.cancelAction()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.showConfirmDialog).toBe(false)
    expect(wrapper.emitted('deactivate')).toBeFalsy()
  })

  it('renders correct status tag colors', async () => {
    const members: MemberDTO[] = [
      {
        id: '1',
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Senior',
        status: 'active',
        createdAt: '2026-09-18T10:00:00Z',
      },
      {
        id: '2',
        firstName: 'Bob',
        lastName: 'Jones',
        seniority: 'Mid',
        status: 'deactivated',
        createdAt: '2026-09-17T10:00:00Z',
      },
    ]

    const wrapper = mount(MemberList, {
      props: {
        members,
      },
      global: {
        plugins: [router],
        stubs: {
          RouterLink: { template: '<a href="#"><slot /></a>' },
        },
      },
    })

    const tags = wrapper.findAll('.tag')
    expect(tags[0].classes()).toContain('is-success')
    expect(tags[1].classes()).toContain('is-warning')
  })
})
