import { expect, test } from '@playwright/test'
import { addEmployee } from './fixtures'

test('team pulse overview renders on startup with empty state', async ({ page }) => {
  await page.goto('/')

  await expect(page.locator('#team-pulse-overview')).toBeVisible()
  await expect(page.locator('#team-pulse-metrics')).toContainText('Headcount')
  await expect(page.locator('#team-pulse-metrics')).toContainText('Rolling Team Morale')
  await expect(page.locator('#team-pulse-metrics')).toContainText('High Capacity Risk')
  await expect(page.locator('#team-pulse-metrics')).toContainText('Team Multiplier')

  await expect(page.locator('#team-pulse-health')).toContainText('execution')
  await expect(page.locator('#team-pulse-health')).toContainText('impact')
  await expect(page.locator('#team-pulse-health')).toContainText('multiplier')
  await expect(page.locator('#team-pulse-health')).toContainText('growth')
  await expect(page.locator('#team-pulse-health')).toContainText('culture')

  await expect(page.locator('#team-pulse-health')).toContainText('No anomalies detected.')
})

test('team pulse overview reflects added employees', async ({ page }) => {
  await page.goto('/')
  await addEmployee(page, 'Grace', 'Hopper', 'Principal', '2018-03-01')

  const headcountBox = page.locator('#team-pulse-metrics').locator('.box').filter({ hasText: 'Headcount' }).locator('.title')
  await expect(headcountBox).toBeVisible()
  const value = await headcountBox.textContent()
  expect(value).toMatch(/^[0-9]+$/)
  expect(Number.parseInt(value ?? '0', 10)).toBeGreaterThanOrEqual(1)
})
