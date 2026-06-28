import { test as setup } from '@playwright/test';

const authFile = 'playwright/.auth/admin.json';

setup('authenticate', async ({ page }) => {
  setup.setTimeout(120000);

  // Create the auth directory if it doesn't exist using Playwright's built-in capabilities
  const path = require('path');
  const fs = require('fs');
  const adminUser = process.env.GRAFANA_USERNAME || process.env.GF_SECURITY_ADMIN_USER || 'admin';
  const adminPassword = process.env.GRAFANA_PASSWORD || process.env.GF_SECURITY_ADMIN_PASSWORD || 'admin';

  const authDir = path.dirname(authFile);
  if (!fs.existsSync(authDir)) {
    fs.mkdirSync(authDir, { recursive: true });
  }

  try {
    console.log('🔐 Starting authentication process...');

    // First check if Grafana is accessible
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    const currentUrl = page.url();
    console.log(`📍 Current URL: ${currentUrl}`);

    // Check if we're already authenticated (anonymous access enabled)
    const isOnLoginPage =
      currentUrl.includes('/login') ||
      (await page.locator('input[type="password"]').isVisible());

    if (!isOnLoginPage) {
      console.log('✅ Already authenticated or anonymous access enabled');
      await page.context().storageState({ path: authFile });
      return;
    }

    console.log('🔑 Login required, attempting authentication...');

    // Navigate to login page if not already there
    if (!currentUrl.includes('/login')) {
      await page.goto('/login');
      await page.waitForLoadState('networkidle');
    }

    // Try multiple selectors for username field
    const usernameSelectors = [
      'input[name="user"]',
      'input[placeholder*="email"]',
      'input[placeholder*="username"]',
      'input[aria-label*="username"]',
      'input[aria-label*="email"]',
      '[data-testid="username"]',
      '[data-testid="email"]',
    ];

    const passwordSelectors = [
      'input[name="password"]',
      'input[type="password"]',
      'input[aria-label*="password"]',
      '[data-testid="password"]',
    ];

    let usernameField, passwordField;

    // Find username field
    for (const selector of usernameSelectors) {
      try {
        const field = page.locator(selector).first();
        if (await field.isVisible({ timeout: 2000 })) {
          usernameField = field;
          console.log(`📝 Found username field: ${selector}`);
          break;
        }
      } catch (e) {
        // Continue to next selector
      }
    }

    // Find password field
    for (const selector of passwordSelectors) {
      try {
        const field = page.locator(selector).first();
        if (await field.isVisible({ timeout: 2000 })) {
          passwordField = field;
          console.log(`🔒 Found password field: ${selector}`);
          break;
        }
      } catch (e) {
        // Continue to next selector
      }
    }

    if (!usernameField || !passwordField) {
      throw new Error('❌ Could not find login form fields');
    }

    // Fill login form
    await usernameField.fill(adminUser);
    await passwordField.fill(adminPassword);

    console.log('✏️ Filled login credentials');

    // Find and click submit button
    const submitSelectors = [
      'button[type="submit"]',
      'input[type="submit"]',
      'button:has-text("Log in")',
      'button:has-text("Sign in")',
      'button:has-text("Login")',
      '[data-testid="login-button"]',
    ];

    let submitButton;
    for (const selector of submitSelectors) {
      try {
        const button = page.locator(selector).first();
        if (await button.isVisible({ timeout: 2000 })) {
          submitButton = button;
          console.log(`🎯 Found submit button: ${selector}`);
          break;
        }
      } catch (e) {
        // Continue to next selector
      }
    }

    if (!submitButton) {
      // Try pressing Enter as fallback
      console.log('⏎ Submit button not found, trying Enter key');
      await passwordField.press('Enter');
    } else {
      await submitButton.click();
    }

    console.log('🚀 Submitted login form');

    // Wait for URL to leave login page first; this is the fastest reliable success check.
    console.log('🔍 Waiting for redirect away from login page...');
    try {
      await page.waitForURL((url) => {
        const href = url.toString();
        return !href.includes('/login') && !href.includes('/auth');
      }, { timeout: 20000 });
      console.log('✅ URL indicates successful login');
    } catch (_e) {
      const loginError = await page
        .locator('.alert-error, [role="alert"], [aria-live="assertive"], .css-1w5x5d2')
        .first()
        .textContent()
        .catch(() => null);

      if (loginError) {
        throw new Error(`❌ Authentication failed: ${loginError.trim()}`);
      }
    }

    // Final check
    const finalUrl = page.url();
    console.log(`📍 Final URL: ${finalUrl}`);

    if (finalUrl.includes('/login') || finalUrl.includes('/auth')) {
      throw new Error('❌ Authentication failed - still on login/auth page');
    }

    // Save authentication state
    await page.context().storageState({ path: authFile });
    console.log('💾 Saved authentication state');
  } catch (error) {
    console.error(
      '❌ Authentication error:',
      error instanceof Error ? error.message : String(error)
    );

    // Create screenshots directory if it doesn't exist
    const screenshotDir = 'tests/screenshots';
    const fs = require('fs');
    if (!fs.existsSync(screenshotDir)) {
      fs.mkdirSync(screenshotDir, { recursive: true });
    }

    // Take screenshot for debugging if page/context is still alive.
    const timestamp = Date.now();
    if (!page.isClosed()) {
      try {
        await page.screenshot({
          path: `${screenshotDir}/auth-error-${timestamp}.png`,
          fullPage: true,
        });
        console.log(`📸 Screenshot saved: auth-error-${timestamp}.png`);
      } catch (screenshotError) {
        console.log(
          `⚠️ Could not capture screenshot: ${screenshotError instanceof Error ? screenshotError.message : String(screenshotError)}`
        );
      }
    }

    // Check if we can proceed anyway (maybe anonymous access works)
    const currentUrl = page.url();

    if (!currentUrl.includes('/login') && !currentUrl.includes('/auth')) {
      console.log('⚠️ Error occurred but not on login page, proceeding...');

      // Create a basic auth state
      const basicAuthState = {
        cookies: [],
        origins: [
          {
            origin: 'http://localhost:3000',
            localStorage: [],
          },
        ],
      };

      fs.writeFileSync(authFile, JSON.stringify(basicAuthState, null, 2));
      console.log('💾 Created fallback auth state');
      return;
    }

    throw error;
  }
});
