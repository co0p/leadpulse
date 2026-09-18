import { test, expect } from '@playwright/test';

test.describe('Members Screen - Happy Path', () => {
  test('add member → view → edit → deactivate → reactivate', async ({ page }) => {
    // Navigate to members page
    await page.goto('/members');
    await page.waitForLoadState('networkidle');

    // Verify members page is loaded
    const pageTitle = page.locator('h1:has-text("Members")');
    await expect(pageTitle).toBeVisible();

    // Verify tabs are visible
    const activeTabs = page.locator('text=Active');
    await expect(activeTabs).toBeVisible();
    const deactivatedTab = page.locator('text=Deactivated');
    await expect(deactivatedTab).toBeVisible();

    // ── Step 1: Add Member ──
    const addMemberButton = page.locator('button:has-text("Add Member")');
    await addMemberButton.click();
    await page.waitForURL('**/members/add');

    // Fill form
    await page.locator('input[placeholder="Enter first name"]').fill('John');
    await page.locator('input[placeholder="Enter last name"]').fill('Doe');
    await page.locator('select').selectOption('Senior');

    // Submit form
    const submitButton = page.locator('button:has-text("Create Member")');
    await submitButton.click();
    await page.waitForURL('**/members');
    await page.waitForTimeout(500); // Wait for member to appear in list

    // ── Step 2: View Member in List ──
    // Verify member appears in active members list
    const memberRow = page.locator('text=John Doe');
    await expect(memberRow).toBeVisible();
    const seniorTag = page.locator('text=Senior');
    await expect(seniorTag).toBeVisible();
    const activeTag = page.locator('.tag:has-text("Active")');
    await expect(activeTag).toBeVisible();

    // ── Step 3: Edit Member ──
    // Click edit button
    const editButton = page.locator('button:has-text("Edit")').first();
    await editButton.click();
    await page.waitForURL('**/members/*/edit');

    // Verify pre-fill
    const firstNameInput = page.locator('input[placeholder="Enter first name"]');
    await expect(firstNameInput).toHaveValue('John');
    const lastNameInput = page.locator('input[placeholder="Enter last name"]');
    await expect(lastNameInput).toHaveValue('Doe');
    const senioritySelect = page.locator('select');
    await expect(senioritySelect).toHaveValue('Senior');

    // Change seniority
    await senioritySelect.selectOption('Lead');

    // Submit form
    const updateButton = page.locator('button:has-text("Update Member")');
    await updateButton.click();
    await page.waitForURL('**/members');
    await page.waitForTimeout(500);

    // ── Step 4: Verify Edit Applied ──
    // Navigate back to members list and verify the change
    const leadTag = page.locator('.tag:has-text("Lead")');
    await expect(leadTag).toBeVisible();

    // ── Step 5: Deactivate Member ──
    // Find deactivate button (second button for John Doe row)
    const deactivateButton = page.locator('button:has-text("Deactivate")').first();
    await deactivateButton.click();

    // Confirm in dialog
    const confirmDeactivateButton = page.locator('button.is-warning:has-text("Deactivate")');
    await confirmDeactivateButton.click();
    await page.waitForTimeout(500);

    // ── Step 6: Verify Member is Deactivated ──
    // Switch to deactivated tab
    const deactivatedTabLink = page.locator('text=Deactivated').last();
    await deactivatedTabLink.click();
    await page.waitForTimeout(300);

    // Verify member appears in deactivated list
    const deactivatedMember = page.locator('text=John Doe');
    await expect(deactivatedMember).toBeVisible();
    const deactivatedTag = page.locator('.tag:has-text("Deactivated")');
    await expect(deactivatedTag).toBeVisible();

    // ── Step 7: Reactivate Member ──
    // Find reactivate button
    const reactivateButton = page.locator('button:has-text("Reactivate")').first();
    await reactivateButton.click();

    // Confirm in dialog
    const confirmReactivateButton = page.locator('button.is-success:has-text("Reactivate")');
    await confirmReactivateButton.click();
    await page.waitForTimeout(500);

    // ── Step 8: Verify Member is Reactivated ──
    // Switch to active tab
    const activeTabLink = page.locator('text=Active').first();
    await activeTabLink.click();
    await page.waitForTimeout(300);

    // Verify member appears in active list again
    const reactivatedMember = page.locator('text=John Doe');
    await expect(reactivatedMember).toBeVisible();
    const reactivatedActiveTag = page.locator('.tag:has-text("Active")');
    await expect(reactivatedActiveTag).toBeVisible();
  });
});
