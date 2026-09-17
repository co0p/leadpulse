import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useHealthStore } from '../health'

describe('useHealthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initializes with checking status', () => {
    const store = useHealthStore()
    expect(store.status).toBe('checking')
    expect(store.lastCheckedAt).toBe(0)
    expect(store.checkError).toBeNull()
  })

  it('marks store as healthy after successful health check', async () => {
    const store = useHealthStore()

    // Mock fetch
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ status: 'ok' })
      })
    )

    await store.checkHealth()

    expect(store.status).toBe('healthy')
    expect(store.lastCheckedAt).toBeTruthy()
    expect(store.checkError).toBeNull()
  })

  it('marks store as unhealthy after failed health check', async () => {
    const store = useHealthStore()

    // Mock fetch to reject
    global.fetch = vi.fn(() =>
      Promise.reject(new Error('Network error'))
    )

    await store.checkHealth()

    expect(store.status).toBe('unhealthy')
    expect(store.checkError).toBeTruthy()
  })

  it('marks store as unhealthy for non-ok response', async () => {
    const store = useHealthStore()

    // Mock fetch with non-200 status
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: false,
        status: 503,
        json: () => Promise.resolve({ status: 'unavailable' })
      })
    )

    await store.checkHealth()

    expect(store.status).toBe('unhealthy')
    expect(store.checkError).toBeTruthy()
  })

  it('updates lastCheckedAt timestamp on successful check', async () => {
    const store = useHealthStore()
    const beforeCheck = Date.now()

    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ status: 'ok' })
      })
    )

    await store.checkHealth()

    const afterCheck = Date.now()
    expect(store.lastCheckedAt).toBeGreaterThanOrEqual(beforeCheck)
    expect(store.lastCheckedAt).toBeLessThanOrEqual(afterCheck)
  })

  it('stores error message on failed check', async () => {
    const store = useHealthStore()
    const errorMessage = 'Connection refused'

    global.fetch = vi.fn(() =>
      Promise.reject(new Error(errorMessage))
    )

    await store.checkHealth()

    expect(store.checkError).toContain(errorMessage)
  })

  it('handles timeout gracefully', async () => {
    const store = useHealthStore()

    // Mock fetch to reject with timeout
    global.fetch = vi.fn(() =>
      Promise.reject(new Error('Timeout'))
    )

    await store.checkHealth()

    expect(store.status).toBe('unhealthy')
    expect(store.checkError).toContain('Timeout')
  })
})
