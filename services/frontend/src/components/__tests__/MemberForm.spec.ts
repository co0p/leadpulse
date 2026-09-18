import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import MemberForm from '../MemberForm.vue'

describe('MemberForm.vue', () => {
  it('renders form fields with labels', () => {
    const wrapper = mount(MemberForm)
    expect(wrapper.text()).toContain('First Name')
    expect(wrapper.text()).toContain('Last Name')
    expect(wrapper.text()).toContain('Seniority')
    expect(wrapper.find('input[type="text"]').exists()).toBe(true)
    expect(wrapper.find('select').exists()).toBe(true)
  })

  it('disables submit when firstName is empty', async () => {
    const wrapper = mount(MemberForm)
    const submitButton = wrapper.find('button[type="submit"]')
    
    // All fields empty initially
    expect(submitButton.attributes('disabled')).toBeDefined()
  })

  it('disables submit when lastName is empty', async () => {
    const wrapper = mount(MemberForm)
    const firstNameInput = wrapper.find('input[placeholder="Enter first name"]')
    
    await firstNameInput.setValue('Alice')
    await wrapper.vm.$nextTick()
    
    const submitButton = wrapper.find('button[type="submit"]')
    // lastName and seniority still empty
    expect(submitButton.attributes('disabled')).toBeDefined()
  })

  it('disables submit when seniority is empty', async () => {
    const wrapper = mount(MemberForm)
    const firstNameInput = wrapper.find('input[placeholder="Enter first name"]')
    const lastNameInput = wrapper.find('input[placeholder="Enter last name"]')
    
    await firstNameInput.setValue('Alice')
    await lastNameInput.setValue('Smith')
    await wrapper.vm.$nextTick()
    
    const submitButton = wrapper.find('button[type="submit"]')
    // seniority still empty
    expect(submitButton.attributes('disabled')).toBeDefined()
  })

  it('enables submit when all fields are valid', async () => {
    const wrapper = mount(MemberForm)
    const firstNameInput = wrapper.find('input[placeholder="Enter first name"]')
    const lastNameInput = wrapper.find('input[placeholder="Enter last name"]')
    const senioritySelect = wrapper.find('select')
    
    await firstNameInput.setValue('Alice')
    await lastNameInput.setValue('Smith')
    await senioritySelect.setValue('Senior')
    await wrapper.vm.$nextTick()
    
    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.attributes('disabled')).toBeUndefined()
  })

  it('shows error message when firstName is blank and field loses focus', async () => {
    const wrapper = mount(MemberForm)
    const firstNameInput = wrapper.find('input[placeholder="Enter first name"]')
    
    await firstNameInput.trigger('blur')
    await wrapper.vm.$nextTick()
    
    const errorText = wrapper.text()
    expect(errorText).toContain('First name is required')
  })

  it('shows error message when lastName is blank and field loses focus', async () => {
    const wrapper = mount(MemberForm)
    const lastNameInput = wrapper.find('input[placeholder="Enter last name"]')
    
    await lastNameInput.trigger('blur')
    await wrapper.vm.$nextTick()
    
    const errorText = wrapper.text()
    expect(errorText).toContain('Last name is required')
  })

  it('shows error message when seniority is not selected and field loses focus', async () => {
    const wrapper = mount(MemberForm)
    const senioritySelect = wrapper.find('select')
    
    await senioritySelect.trigger('blur')
    await wrapper.vm.$nextTick()
    
    const errorText = wrapper.text()
    expect(errorText).toContain('Seniority level is required')
  })

  it('emits submit event with form data when form is valid', async () => {
    const wrapper = mount(MemberForm)
    const firstNameInput = wrapper.find('input[placeholder="Enter first name"]')
    const lastNameInput = wrapper.find('input[placeholder="Enter last name"]')
    const senioritySelect = wrapper.find('select')
    
    await firstNameInput.setValue('Alice')
    await lastNameInput.setValue('Smith')
    await senioritySelect.setValue('Senior')
    await wrapper.vm.$nextTick()
    
    const form = wrapper.find('form')
    await form.trigger('submit')
    
    expect(wrapper.emitted('submit')).toBeTruthy()
    const submitEvent = wrapper.emitted('submit')?.[0]
    expect(submitEvent).toEqual([
      {
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Senior',
      },
    ])
  })

  it('pre-fills form fields from initialData prop', async () => {
    const wrapper = mount(MemberForm, {
      props: {
        initialData: {
          firstName: 'Bob',
          lastName: 'Jones',
          seniority: 'Mid',
        },
      },
    })
    
    await wrapper.vm.$nextTick()
    
    const firstNameInput = wrapper.find('input[placeholder="Enter first name"]') as any
    const lastNameInput = wrapper.find('input[placeholder="Enter last name"]') as any
    const senioritySelect = wrapper.find('select') as any
    
    expect(firstNameInput.element.value).toBe('Bob')
    expect(lastNameInput.element.value).toBe('Jones')
    expect(senioritySelect.element.value).toBe('Mid')
  })

  it('enables submit button when pre-filled with valid data', async () => {
    const wrapper = mount(MemberForm, {
      props: {
        initialData: {
          firstName: 'Bob',
          lastName: 'Jones',
          seniority: 'Mid',
        },
      },
    })
    
    await wrapper.vm.$nextTick()
    
    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.attributes('disabled')).toBeUndefined()
  })

  it('emits cancel event when cancel button is clicked', async () => {
    const wrapper = mount(MemberForm)
    const cancelButton = wrapper.findAll('button').find(b => b.text() === 'Cancel')
    
    await cancelButton?.trigger('click')
    
    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('disables submit button when isSubmitting is true', async () => {
    const wrapper = mount(MemberForm, {
      props: {
        isSubmitting: true,
        initialData: {
          firstName: 'Alice',
          lastName: 'Smith',
          seniority: 'Senior',
        },
      },
    })
    
    await wrapper.vm.$nextTick()
    
    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.attributes('disabled')).toBeDefined()
  })

  it('trims whitespace from firstName and lastName on submit', async () => {
    const wrapper = mount(MemberForm)
    const firstNameInput = wrapper.find('input[placeholder="Enter first name"]')
    const lastNameInput = wrapper.find('input[placeholder="Enter last name"]')
    const senioritySelect = wrapper.find('select')
    
    await firstNameInput.setValue('  Alice  ')
    await lastNameInput.setValue('  Smith  ')
    await senioritySelect.setValue('Senior')
    await wrapper.vm.$nextTick()
    
    const form = wrapper.find('form')
    await form.trigger('submit')
    
    const submitEvent = wrapper.emitted('submit')?.[0]
    expect(submitEvent).toEqual([
      {
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Senior',
      },
    ])
  })
})
