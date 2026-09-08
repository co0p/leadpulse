import { AddEmployee, ListEmployees, RemoveEmployee, TeamPulseOverview } from '../wailsjs/go/app/App'
import type { employees, overview } from '../wailsjs/go/models'

export type Employee = employees.Employee
export type EmployeeInput = employees.Input
export type Overview = overview.Overview
export type HealthDimensions = overview.HealthDimensions
export type EmployeeHealth = overview.EmployeeHealth
export type Anomaly = overview.Anomaly
export type EvidenceGap = overview.EvidenceGap

export const employeeApi = {
  list: ListEmployees,
  add: AddEmployee,
  remove: RemoveEmployee,
  teamPulseOverview: TeamPulseOverview,
}
