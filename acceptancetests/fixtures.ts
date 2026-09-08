import { expect, type Page } from '@playwright/test'

export const addEmployee = async (page: Page, firstName: string, secondName: string, seniority: string, startDate: string) => {
  await page.getByRole('button', { name: 'Add employee' }).click()
  await page.getByRole('region', { name: 'Add employee' }).getByLabel('First name').fill(firstName)
  await page.getByRole('region', { name: 'Add employee' }).getByLabel('Second name').fill(secondName)
  await page.getByRole('region', { name: 'Add employee' }).getByLabel('Seniority').selectOption({ label: seniority })
  await page.getByRole('region', { name: 'Add employee' }).getByLabel('Start date').fill(startDate)
  await page.getByRole('region', { name: 'Add employee' }).getByRole('button', { name: 'Add employee', exact: true }).click()
  await expect(page.getByRole('button', { name: `${firstName} ${secondName}` })).toBeVisible()
}

export const addPerformanceEntry = async (page: Page, employeeName: string, data: {
  year: string
  month: string
  morale: string
  execution: string
  impact: string
  growth: string
  culture: string
  projects: string
  evidenceCount: string
}) => {
  // Performance entries are not exposed in the UI yet, so this fixture relies on
  // a backend helper seeded through the acceptance database. Tests that need
  // explicit performance data will call the backend seed endpoint once it exists.
  // For now the acceptance suite asserts the empty-state overview renders correctly.
}
