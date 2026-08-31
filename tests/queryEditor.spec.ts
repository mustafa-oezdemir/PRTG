import { expect, test } from '@grafana/plugin-e2e';

test.beforeEach(async ({ panelEditPage, readProvisionedDataSource }) => {
  const dataSource = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(dataSource.name);
});

test('renders the PRTG query editor', async ({ panelEditPage }) => {
  const query = panelEditPage.getQueryEditorRow('A');

  await expect(query.locator('#query-editor-queryType')).toBeVisible();
  await expect(query.locator('#query-editor-group')).toBeVisible();
  await expect(query.locator('#query-editor-device')).toBeVisible();
  await expect(query.locator('#query-editor-sensor')).toBeVisible();
  await expect(query.getByText('Channel', { exact: true })).toBeVisible();
  await expect(query.locator('#query-editor-is-stream')).toBeVisible();
});

test('shows text query options when Text mode is selected', async ({ page, panelEditPage }) => {
  const query = panelEditPage.getQueryEditorRow('A');

  await query.locator('#query-editor-queryType').click();
  await page.getByRole('option', { name: 'Text', exact: true }).click();

  await expect(query.locator('#query-editor-property')).toBeVisible();
  await expect(query.locator('#query-editor-filterProperty')).toBeVisible();
});

test('shows the update interval when streaming is enabled', async ({ panelEditPage }) => {
  const query = panelEditPage.getQueryEditorRow('A');
  const streaming = query.locator('#query-editor-is-stream');

  await query.getByText('Enable Streaming', { exact: true }).click();

  await expect(streaming).toBeChecked();
  await expect(query.locator('#query-editor-stream-interval')).toBeVisible();
});
