import { test, expect } from '@playwright/test'

test.describe('Smoke Tests', () => {
  test('dashboard loads without console errors', async ({ page }) => {
    const errors: string[] = []
    page.on('console', msg => {
      if (msg.type() === 'error') {
        errors.push(msg.text())
      }
    })

    await page.goto('/')
    await page.waitForLoadState('networkidle')

    // Check page title/heading
    await expect(page.locator('h1')).toBeVisible()

    // Check KPI cards are rendered
    await expect(page.locator('.kpi-card').first()).toBeVisible()

    // No console errors (ignore API errors when backend returns 503, and network failures)
    const realErrors = errors.filter(e =>
      !e.includes('/api/') &&
      !e.includes('fetch') &&
      !e.includes('503') &&
      !e.includes('Service Unavailable')
    )
    expect(realErrors).toHaveLength(0)
  })

  test('navigation between pages', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')

    // Navigate to Competitor page
    await page.click('text=竞品情报')
    await expect(page.locator('h1')).toContainText('竞品情报')

    // Navigate to Pricing page
    await page.click('text=定价策略')
    await expect(page.locator('h1')).toContainText('定价策略')

    // Navigate to Supply Chain page
    await page.click('text=供应链预警')
    await expect(page.locator('h1')).toContainText('供应链预警')

    // Navigate to Product Selection page
    await page.click('text=选品分析')
    await expect(page.locator('h1')).toContainText('选品分析')

    // Navigate back to Dashboard
    await page.click('text=Dashboard')
    await expect(page.locator('h1')).toContainText('监控总览')
  })

  test('competitor list page renders', async ({ page }) => {
    await page.goto('/competitors')
    await page.waitForLoadState('networkidle')

    await expect(page.locator('.competitor')).toBeVisible()
    await expect(page.locator('table.data-table')).toBeVisible()
  })

  test('pricing page shows data', async ({ page }) => {
    await page.goto('/pricing')
    await page.waitForLoadState('networkidle')

    await expect(page.locator('.pricing')).toBeVisible()
    await expect(page.locator('.kpi-card').first()).toBeVisible()
  })
})
