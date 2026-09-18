/**
 * API client for backend communication
 * Handles health checks and other API calls
 * 
 * In production (Docker): uses /api which nginx proxies to backend:8080
 * In development: uses configured VITE_API_URL or /api (dev server proxy)
 */

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'

export interface HealthResponse {
  status: 'ok' | string
  version?: string
}

/**
 * Fetch health status from backend
 * @throws Error if the request fails or times out
 * @returns Promise<HealthResponse>
 */
export async function fetchHealth(): Promise<HealthResponse> {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), 5000) // 5 second timeout

  try {
    const response = await fetch(`${API_BASE_URL}/health`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      },
      signal: controller.signal
    })

    clearTimeout(timeoutId)

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    const data: HealthResponse = await response.json()
    return data
  } catch (error) {
    clearTimeout(timeoutId)
    if (error instanceof Error) {
      throw error
    }
    throw new Error('Unknown error fetching health')
  }
}
