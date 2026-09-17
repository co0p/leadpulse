import { test, expect } from '@playwright/test';

test.describe('App Shell', () => {
  test('app shell loads with sidebar, top bar, footer', async ({ page }) => {
    // Navigate to home
    await page.goto('/');

    // Verify sidebar is visible
    const sidebar = page.locator('[data-testid="sidebar"]');
    await expect(sidebar).toBeVisible();

    // Verify topbar is visible
    const topbar = page.locator('[data-testid="topbar"]');
    await expect(topbar).toBeVisible();

    // Verify footer is visible
    const footer = page.locator('[data-testid="footer"]');
    await expect(footer).toBeVisible();

    // Verify main heading is present
    const heading = page.locator('h2:has-text("Welcome to LeadPulse")');
    await expect(heading).toBeVisible();
  });

  test('health indicator shows green when backend responds', async ({ page }) => {
    // Navigate to home
    await page.goto('/');

    // Wait for page to load and health check to complete
    await page.waitForTimeout(2000);

    // Verify health indicator is present and shows healthy status
    const healthIndicator = page.locator('[data-testid="health-indicator"]');
    await expect(healthIndicator).toBeVisible();

    // Should contain success tag
    const successTag = healthIndicator.locator('.tag.is-success');
    await expect(successTag).toBeVisible();
  });

  test('no console errors on page load', async ({ page }) => {
    // Collect console messages
    const consoleLogs: string[] = [];
    const consoleErrors: string[] = [];

    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
      consoleLogs.push(`[${msg.type()}] ${msg.text()}`);
    });

    // Navigate to home
    await page.goto('/');

    // Wait for page to fully load
    await page.waitForTimeout(2000);

    // Print console logs for debugging
    consoleLogs.forEach((log) => console.log(log));

    // Assert no errors in console
    expect(consoleErrors).toHaveLength(0);
  });
});
