import { expect, test } from '@grafana/plugin-e2e';

test('renders the PRTG configuration editor', async ({ createDataSourceConfigPage }) => {
  const configPage = await createDataSourceConfigPage({ type: 'maxmarkusprogram-prtg-datasource' });
  const page = configPage.ctx.page;

  await expect(page.locator('#config-editor-path')).toBeVisible();
  await expect(page.locator('#config-editor-api-key')).toBeVisible();
  await expect(page.locator('#config-editor-cache-time')).toHaveValue('6000');
  await expect(page.locator('#config-editor-timezone')).toBeVisible();
});

test('loads the provisioned PRTG configuration', async ({ gotoDataSourceConfigPage, readProvisionedDataSource }) => {
  const dataSource = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  const configPage = await gotoDataSourceConfigPage(dataSource.uid);
  const page = configPage.ctx.page;

  await expect(page.locator('#config-editor-path')).toBeEditable();
  await expect(page.locator('#config-editor-cache-time')).toHaveValue(/^\d+$/);
  await expect(page.locator('#config-editor-timezone')).toBeVisible();
});

test('saves valid settings after a successful health check', async ({ createDataSourceConfigPage }) => {
  const configPage = await createDataSourceConfigPage({ type: 'maxmarkusprogram-prtg-datasource' });
  const page = configPage.ctx.page;

  await page.locator('#config-editor-path').fill('prtg.example.com');
  await page.locator('#config-editor-api-key').fill('test-token');
  await configPage.mockHealthCheckResponse({ status: 'OK', message: 'PRTG connection is healthy' });

  await expect(configPage.saveAndTest()).toBeOK();
  await expect(configPage).toHaveAlert('success');
});
