import { expect, test } from '@playwright/test'
import { addEmployee } from './fixtures'

test('team pulse overview renders on startup with empty state', async ({ page }) => {
  await page.goto('/')

  await expect(page.locator('#team-pulse-overview')).toBeVisible()
  await expect(page.locator('#team-pulse-metrics')).toContainText('Headcount')
  await expect(page.locator('#team-pulse-metrics')).toContainText('Rolling Team Morale')
  await expect(page.locator('#team-pulse-metrics')).toContainText('High Capacity Risk')
  await expect(page.locator('#team-pulse-metrics')).toContainText('Team Multiplier')

  await expect(page.locator('#team-pulse-health')).toContainText('Execution')
  await expect(page.locator('#team-pulse-health')).toContainText('Impact')
  await expect(page.locator('#team-pulse-health')).toContainText('Multiplier')
  await expect(page.locator('#team-pulse-health')).toContainText('Growth')
  await expect(page.locator('#team-pulse-health')).toContainText('Culture')

  await expect(page.locator('#team-pulse-health')).toContainText('No anomalies detected.')
})

test('team pulse overview reflects added employees', async ({ page }) => {
  await page.goto('/')
  await addEmployee(page, 'Grace', 'Hopper', 'Principal', '2018-03-01')

  await expect(page.locator('#team-pulse-metrics')).toContainText('1')
})
