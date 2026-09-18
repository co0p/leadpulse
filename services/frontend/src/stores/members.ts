import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as membersAPI from '../api/members'

export interface MemberDTO {
  id: string
  firstName: string
  lastName: string
  seniority: string
  status: 'active' | 'deactivated'
  createdAt: string
}

export type CurrentTab = 'active' | 'deactivated' | 'all'

export const useMembersStore = defineStore('members', () => {
  const members = ref<MemberDTO[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const currentTab = ref<CurrentTab>('active')

  /**
   * Computed: filtered members based on current tab
   * - 'active': members with status === 'active'
   * - 'deactivated': members with status === 'deactivated'
   * - 'all': all members
   */
  const filteredMembers = computed(() => {
    if (currentTab.value === 'all') {
      return members.value
    }
    return members.value.filter((m) => m.status === currentTab.value)
  })

  /**
   * Set current tab and filter members
   */
  function setCurrentTab(tab: CurrentTab) {
    currentTab.value = tab
  }

  /**
   * Set loading state
   */
  function setLoading(isLoading: boolean) {
    loading.value = isLoading
  }

  /**
   * Set error message
   */
  function setError(errorMessage: string | null) {
    error.value = errorMessage
  }

  /**
   * Clear error message
   */
  function clearError() {
    error.value = null
  }

  /**
   * Update members list
   */
  function setMembers(newMembers: MemberDTO[]) {
    members.value = newMembers
  }

  /**
   * Add member to the list
   */
  function addMemberToList(member: MemberDTO) {
    members.value.push(member)
  }

  /**
   * Update a member in the list
   */
  function updateMemberInList(id: string, updates: Partial<MemberDTO>) {
    const index = members.value.findIndex((m) => m.id === id)
    if (index !== -1) {
      members.value[index] = { ...members.value[index], ...updates }
    }
  }

  /**
   * Remove member from the list
   */
  function removeMemberFromList(id: string) {
    members.value = members.value.filter((m) => m.id !== id)
  }

  /**
   * Find member by ID
   */
  function findMemberById(id: string): MemberDTO | undefined {
    return members.value.find((m) => m.id === id)
  }

  /**
   * Load members from API (filter by current tab)
   * @throws Error if API call fails
   */
  async function loadMembers(): Promise<void> {
    setLoading(true)
    clearError()
    try {
      const status = currentTab.value === 'all' ? 'all' : currentTab.value
      const fetchedMembers = await membersAPI.fetchMembers(status)
      setMembers(fetchedMembers as MemberDTO[])
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load members'
      setError(errorMessage)
      throw err
    } finally {
      setLoading(false)
    }
  }

  /**
   * Add a new member via API
   * @param data Member data (firstName, lastName, seniority)
   * @throws Error if API call fails
   */
  async function addMember(data: membersAPI.AddMemberRequest): Promise<MemberDTO> {
    setLoading(true)
    clearError()
    try {
      const newMember = await membersAPI.addMember(data)
      addMemberToList(newMember as MemberDTO)
      return newMember as MemberDTO
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to add member'
      setError(errorMessage)
      throw err
    } finally {
      setLoading(false)
    }
  }

  /**
   * Edit an existing member via API
   * @param id Member ID (UUID)
   * @param data Member data to update
   * @throws Error if API call fails
   */
  async function editMember(id: string, data: membersAPI.EditMemberRequest): Promise<MemberDTO> {
    setLoading(true)
    clearError()
    try {
      const updatedMember = await membersAPI.editMember(id, data)
      updateMemberInList(id, updatedMember as MemberDTO)
      return updatedMember as MemberDTO
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to edit member'
      setError(errorMessage)
      throw err
    } finally {
      setLoading(false)
    }
  }

  /**
   * Deactivate a member via API
   * @param id Member ID (UUID)
   * @throws Error if API call fails
   */
  async function deactivateMember(id: string): Promise<void> {
    setLoading(true)
    clearError()
    try {
      await membersAPI.deactivateMember(id)
      updateMemberInList(id, { status: 'deactivated' })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to deactivate member'
      setError(errorMessage)
      throw err
    } finally {
      setLoading(false)
    }
  }

  /**
   * Reactivate a member via API
   * @param id Member ID (UUID)
   * @throws Error if API call fails
   */
  async function reactivateMember(id: string): Promise<void> {
    setLoading(true)
    clearError()
    try {
      const updatedMember = await membersAPI.reactivateMember(id)
      updateMemberInList(id, updatedMember as MemberDTO)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to reactivate member'
      setError(errorMessage)
      throw err
    } finally {
      setLoading(false)
    }
  }

  return {
    // State
    members,
    loading,
    error,
    currentTab,

    // Computed
    filteredMembers,

    // Methods
    setCurrentTab,
    setLoading,
    setError,
    clearError,
    setMembers,
    addMemberToList,
    updateMemberInList,
    removeMemberFromList,
    findMemberById,

    // Actions
    loadMembers,
    addMember,
    editMember,
    deactivateMember,
    reactivateMember,
  }
})
