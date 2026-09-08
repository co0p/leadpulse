import { describe, expect, it } from 'vitest'
import { uiShellMarkers } from './uiShellMarkers'

describe('Bulma employee navigation shell', () => {
  it('defines the left navigation and accessible add action', () => {
    expect(uiShellMarkers.navigationId).toBe('employee-navigation')
    expect(uiShellMarkers.addActionLabel).toBe('Add employee')
    expect(uiShellMarkers.bulmaLayoutClass).toBe('columns is-gapless is-fullheight')
  })
})
