export type ViewState =
  | { kind: 'empty' }
  | { kind: 'add' }
  | { kind: 'details'; employeeId: number }

export const emptyView = (): ViewState => ({ kind: 'empty' })
export const addView = (): ViewState => ({ kind: 'add' })
export const detailsView = (employeeId: number): ViewState => ({ kind: 'details', employeeId })
