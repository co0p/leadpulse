import type { Employee, EmployeeInput } from './employeeApi'

export interface EmployeeListState {
  employees: Employee[]
  loading: boolean
  error: string
}

export const emptyEmployeeListState = (): EmployeeListState => ({
  employees: [],
  loading: false,
  error: '',
})

export const addEmployeeToState = (state: EmployeeListState, employee: Employee): EmployeeListState => ({
  ...state,
  employees: [...(state.employees ?? []), employee],
})

export const removeEmployeeFromState = (state: EmployeeListState, id: number): EmployeeListState => ({
  ...state,
  employees: (state.employees ?? []).filter((employee) => employee.id !== id),
})

export const employeeInput = (firstName: string, secondName: string, seniority: string, startDate: string): EmployeeInput => ({
  firstName,
  secondName,
  seniority,
  startDate,
})
