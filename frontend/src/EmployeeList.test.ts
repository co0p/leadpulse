import { describe, expect, it } from 'vitest'
import { employeeListViewMarkers } from './employeeListView'

describe('employee list view', () => {
  it('contains the employee management controls', () => {
    expect(employeeListViewMarkers.id).toBe('employee-list')
    expect(employeeListViewMarkers.addLabel).toBe('Add employee')
    expect(employeeListViewMarkers.removeLabel).toBe('Remove')
  })
})
