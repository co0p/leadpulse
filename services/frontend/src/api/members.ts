/**
 * Members API client
 * Handles CRUD operations for team members
 * 
 * In production (Docker): uses /api which nginx proxies to backend:8080
 * In development: uses configured VITE_API_URL or /api (dev server proxy)
 */

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'

export interface AddMemberRequest {
  firstName: string
  lastName: string
  seniority: string
}

export interface EditMemberRequest {
  firstName: string
  lastName: string
  seniority: string
}

export interface MemberResponse {
  id: string
  firstName: string
  lastName: string
  seniority: string
  status: string
  createdAt: string
}

export interface GetMembersResponse {
  members: MemberResponse[]
}

/**
 * Fetch members with optional status filter
 * @param status 'active', 'deactivated', or 'all'
 * @returns Promise<MemberResponse[]>
 * @throws Error if the request fails
 */
export async function fetchMembers(status: string = 'active'): Promise<MemberResponse[]> {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), 5000) // 5 second timeout

  try {
    const query = status === 'all' ? '?status=all' : `?status=${status}`
    const response = await fetch(`${API_BASE_URL}/members${query}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal: controller.signal,
    })

    clearTimeout(timeoutId)

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    const data: GetMembersResponse = await response.json()
    return data.members
  } catch (error) {
    clearTimeout(timeoutId)
    if (error instanceof Error) {
      throw error
    }
    throw new Error('Unknown error fetching members')
  }
}

/**
 * Add a new member
 * @param data Member data (firstName, lastName, seniority)
 * @returns Promise<MemberResponse>
 * @throws Error if the request fails
 */
export async function addMember(data: AddMemberRequest): Promise<MemberResponse> {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), 5000) // 5 second timeout

  try {
    const response = await fetch(`${API_BASE_URL}/members`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
      signal: controller.signal,
    })

    clearTimeout(timeoutId)

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to add member')
    }

    const member: MemberResponse = await response.json()
    return member
  } catch (error) {
    clearTimeout(timeoutId)
    if (error instanceof Error) {
      throw error
    }
    throw new Error('Unknown error adding member')
  }
}

/**
 * Edit an existing member
 * @param id Member ID (UUID)
 * @param data Member data to update
 * @returns Promise<MemberResponse>
 * @throws Error if the request fails
 */
export async function editMember(id: string, data: EditMemberRequest): Promise<MemberResponse> {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), 5000) // 5 second timeout

  try {
    const response = await fetch(`${API_BASE_URL}/members/${id}`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
      signal: controller.signal,
    })

    clearTimeout(timeoutId)

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to edit member')
    }

    const member: MemberResponse = await response.json()
    return member
  } catch (error) {
    clearTimeout(timeoutId)
    if (error instanceof Error) {
      throw error
    }
    throw new Error('Unknown error editing member')
  }
}

/**
 * Deactivate a member
 * @param id Member ID (UUID)
 * @returns Promise<void>
 * @throws Error if the request fails
 */
export async function deactivateMember(id: string): Promise<void> {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), 5000) // 5 second timeout

  try {
    const response = await fetch(`${API_BASE_URL}/members/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      signal: controller.signal,
    })

    clearTimeout(timeoutId)

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to deactivate member')
    }
  } catch (error) {
    clearTimeout(timeoutId)
    if (error instanceof Error) {
      throw error
    }
    throw new Error('Unknown error deactivating member')
  }
}

/**
 * Reactivate a member
 * @param id Member ID (UUID)
 * @returns Promise<MemberResponse>
 * @throws Error if the request fails
 */
export async function reactivateMember(id: string): Promise<MemberResponse> {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), 5000) // 5 second timeout

  try {
    const response = await fetch(`${API_BASE_URL}/members/${id}/reactivate`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      signal: controller.signal,
    })

    clearTimeout(timeoutId)

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to reactivate member')
    }

    const member: MemberResponse = await response.json()
    return member
  } catch (error) {
    clearTimeout(timeoutId)
    if (error instanceof Error) {
      throw error
    }
    throw new Error('Unknown error reactivating member')
  }
}
