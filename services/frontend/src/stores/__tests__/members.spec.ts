import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useMembersStore } from '../members'
import * as membersAPI from '../../api/members'

// Mock the API module
vi.mock('../../api/members')

describe('useMembersStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initializes with empty members and active tab', () => {
    const store = useMembersStore()
    expect(store.members).toEqual([])
    expect(store.currentTab).toBe('active')
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
  })

  it('filters members by active tab', () => {
    const store = useMembersStore()
    const activeMember = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }
    const deactivatedMember = {
      id: '2',
      firstName: 'Bob',
      lastName: 'Jones',
      seniority: 'Mid',
      status: 'deactivated' as const,
      createdAt: '2026-09-17T10:00:00Z',
    }

    store.setMembers([activeMember, deactivatedMember])
    store.setCurrentTab('active')
    expect(store.filteredMembers).toHaveLength(1)
    expect(store.filteredMembers[0].id).toBe('1')
  })

  it('filters members by deactivated tab', () => {
    const store = useMembersStore()
    const activeMember = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }
    const deactivatedMember = {
      id: '2',
      firstName: 'Bob',
      lastName: 'Jones',
      seniority: 'Mid',
      status: 'deactivated' as const,
      createdAt: '2026-09-17T10:00:00Z',
    }

    store.setMembers([activeMember, deactivatedMember])
    store.setCurrentTab('deactivated')
    expect(store.filteredMembers).toHaveLength(1)
    expect(store.filteredMembers[0].id).toBe('2')
  })

  it('returns all members for all tab', () => {
    const store = useMembersStore()
    const activeMember = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }
    const deactivatedMember = {
      id: '2',
      firstName: 'Bob',
      lastName: 'Jones',
      seniority: 'Mid',
      status: 'deactivated' as const,
      createdAt: '2026-09-17T10:00:00Z',
    }

    store.setMembers([activeMember, deactivatedMember])
    store.setCurrentTab('all')
    expect(store.filteredMembers).toHaveLength(2)
  })

  it('sets and clears error message', () => {
    const store = useMembersStore()
    expect(store.error).toBeNull()

    store.setError('Test error')
    expect(store.error).toBe('Test error')

    store.clearError()
    expect(store.error).toBeNull()
  })

  it('manages loading state', () => {
    const store = useMembersStore()
    expect(store.loading).toBe(false)

    store.setLoading(true)
    expect(store.loading).toBe(true)

    store.setLoading(false)
    expect(store.loading).toBe(false)
  })

  it('adds member to list', () => {
    const store = useMembersStore()
    const member = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }

    store.addMemberToList(member)
    expect(store.members).toHaveLength(1)
    expect(store.members[0].id).toBe('1')
  })

  it('updates member in list', () => {
    const store = useMembersStore()
    const member = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }

    store.setMembers([member])
    store.updateMemberInList('1', { seniority: 'Lead' })
    expect(store.members[0].seniority).toBe('Lead')
  })

  it('removes member from list', () => {
    const store = useMembersStore()
    const member = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }

    store.setMembers([member])
    store.removeMemberFromList('1')
    expect(store.members).toHaveLength(0)
  })

  it('finds member by ID', () => {
    const store = useMembersStore()
    const member = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }

    store.setMembers([member])
    const found = store.findMemberById('1')
    expect(found).toBeDefined()
    expect(found?.firstName).toBe('Alice')
  })

  it('returns undefined when finding non-existent member', () => {
    const store = useMembersStore()
    const found = store.findMemberById('non-existent')
    expect(found).toBeUndefined()
  })
})

describe('useMembersStore - Async Actions', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('loads members from API on current tab', async () => {
    const mockMembers = [
      {
        id: '1',
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Senior',
        status: 'active' as const,
        createdAt: '2026-09-18T10:00:00Z',
      },
    ]

    vi.mocked(membersAPI.fetchMembers).mockResolvedValue(mockMembers)

    const store = useMembersStore()
    store.setCurrentTab('active')
    await store.loadMembers()

    expect(membersAPI.fetchMembers).toHaveBeenCalledWith('active')
    expect(store.members).toHaveLength(1)
    expect(store.members[0].id).toBe('1')
    expect(store.loading).toBe(false)
  })

  it('clears error when loading members successfully', async () => {
    vi.mocked(membersAPI.fetchMembers).mockResolvedValue([])

    const store = useMembersStore()
    store.setError('Previous error')
    await store.loadMembers()

    expect(store.error).toBeNull()
  })

  it('sets error when loading members fails', async () => {
    vi.mocked(membersAPI.fetchMembers).mockRejectedValue(new Error('Network error'))

    const store = useMembersStore()
    try {
      await store.loadMembers()
    } catch (e) {
      // Expected to throw
    }

    expect(store.error).toBe('Network error')
    expect(store.loading).toBe(false)
  })

  it('adds member via API', async () => {
    const newMember = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }

    vi.mocked(membersAPI.addMember).mockResolvedValue(newMember)

    const store = useMembersStore()
    const result = await store.addMember({
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
    })

    expect(membersAPI.addMember).toHaveBeenCalled()
    expect(store.members).toHaveLength(1)
    expect(store.members[0].id).toBe('1')
    expect(result.id).toBe('1')
    expect(store.loading).toBe(false)
  })

  it('sets error when adding member fails', async () => {
    vi.mocked(membersAPI.addMember).mockRejectedValue(new Error('Add failed'))

    const store = useMembersStore()
    try {
      await store.addMember({
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Senior',
      })
    } catch (e) {
      // Expected to throw
    }

    expect(store.error).toBe('Add failed')
    expect(store.members).toHaveLength(0)
  })

  it('edits member via API', async () => {
    const initialMember = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }
    const updatedMember = {
      ...initialMember,
      seniority: 'Lead',
    }

    vi.mocked(membersAPI.editMember).mockResolvedValue(updatedMember)

    const store = useMembersStore()
    store.addMemberToList(initialMember)
    const result = await store.editMember('1', {
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Lead',
    })

    expect(membersAPI.editMember).toHaveBeenCalledWith('1', {
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Lead',
    })
    expect(store.members[0].seniority).toBe('Lead')
    expect(result.seniority).toBe('Lead')
    expect(store.loading).toBe(false)
  })

  it('sets error when editing member fails', async () => {
    const initialMember = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }

    vi.mocked(membersAPI.editMember).mockRejectedValue(new Error('Edit failed'))

    const store = useMembersStore()
    store.addMemberToList(initialMember)
    try {
      await store.editMember('1', {
        firstName: 'Alice',
        lastName: 'Smith',
        seniority: 'Lead',
      })
    } catch (e) {
      // Expected to throw
    }

    expect(store.error).toBe('Edit failed')
    expect(store.members[0].seniority).toBe('Senior') // Unchanged
  })

  it('deactivates member via API', async () => {
    const member = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }

    vi.mocked(membersAPI.deactivateMember).mockResolvedValue(undefined)

    const store = useMembersStore()
    store.addMemberToList(member)
    await store.deactivateMember('1')

    expect(membersAPI.deactivateMember).toHaveBeenCalledWith('1')
    expect(store.members[0].status).toBe('deactivated')
    expect(store.loading).toBe(false)
  })

  it('sets error when deactivating member fails', async () => {
    const member = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'active' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }

    vi.mocked(membersAPI.deactivateMember).mockRejectedValue(
      new Error('Deactivate failed')
    )

    const store = useMembersStore()
    store.addMemberToList(member)
    try {
      await store.deactivateMember('1')
    } catch (e) {
      // Expected to throw
    }

    expect(store.error).toBe('Deactivate failed')
    expect(store.members[0].status).toBe('active') // Unchanged
  })

  it('reactivates member via API', async () => {
    const member = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'deactivated' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }
    const reactivatedMember = {
      ...member,
      status: 'active' as const,
    }

    vi.mocked(membersAPI.reactivateMember).mockResolvedValue(reactivatedMember)

    const store = useMembersStore()
    store.addMemberToList(member)
    await store.reactivateMember('1')

    expect(membersAPI.reactivateMember).toHaveBeenCalledWith('1')
    expect(store.members[0].status).toBe('active')
    expect(store.loading).toBe(false)
  })

  it('sets error when reactivating member fails', async () => {
    const member = {
      id: '1',
      firstName: 'Alice',
      lastName: 'Smith',
      seniority: 'Senior',
      status: 'deactivated' as const,
      createdAt: '2026-09-18T10:00:00Z',
    }

    vi.mocked(membersAPI.reactivateMember).mockRejectedValue(
      new Error('Reactivate failed')
    )

    const store = useMembersStore()
    store.addMemberToList(member)
    try {
      await store.reactivateMember('1')
    } catch (e) {
      // Expected to throw
    }

    expect(store.error).toBe('Reactivate failed')
    expect(store.members[0].status).toBe('deactivated') // Unchanged
  })
})
