import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fetchHealth } from '../api/client'

export type HealthStatus = 'healthy' | 'unhealthy' | 'checking'

export const useHealthStore = defineStore('health', () => {
  const status = ref<HealthStatus>('checking')
  const lastCheckedAt = ref<number>(0)
  const checkError = ref<string | null>(null)

  /**
   * Check health status from backend
   * Sets status to 'checking', then either 'healthy' or 'unhealthy'
   */
  async function checkHealth() {
    status.value = 'checking'
    checkError.value = null
    lastCheckedAt.value = Date.now()

    try {
      const response = await fetchHealth()

      if (response.status === 'ok') {
        status.value = 'healthy'
      } else {
        status.value = 'unhealthy'
        checkError.value = `Unexpected status: ${response.status}`
      }
    } catch (error) {
      status.value = 'unhealthy'
      if (error instanceof Error) {
        checkError.value = error.message
      } else {
        checkError.value = 'Unknown error'
      }

      // Log error to console (dev aid only)
      if (!import.meta.env.PROD) {
        console.debug('[health] Backend health check failed:', checkError.value)
      }
    }
  }

  return {
    status,
    lastCheckedAt,
    checkError,
    checkHealth
  }
})
