import { expect, test } from '@playwright/test'
import { addEmployee } from './fixtures'

test('product owner can add employees and view their details', async ({ page }) => {
  await page.goto('/')

  await expect(page.locator('#employee-navigation')).toBeVisible()
  await expect(page.getByRole('main', { name: 'LeadPulse employee list' })).toBeVisible()
  await expect(page.locator('#team-pulse-overview')).toBeVisible()

  await addEmployee(page, 'Grace', 'Hopper', 'Principal', '2018-03-01')
  await addEmployee(page, 'Alan', 'Turing', 'Senior', '2019-06-15')

  await expect(page.getByRole('button', { name: 'Grace Hopper' })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Alan Turing' })).toBeVisible()

  await addEmployee(page, 'Ada', 'Lovelace', 'Senior', '2020-01-15')
  await page.getByRole('button', { name: 'Ada Lovelace' }).click()
  await expect(page.getByRole('region', { name: 'Employee details' })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Ada Lovelace' })).toBeVisible()
  await expect(page.getByRole('region', { name: 'Employee details' })).toContainText('Senior')
  await expect(page.getByRole('region', { name: 'Employee details' })).toContainText('2020-01-15')
})
