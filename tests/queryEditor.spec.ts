import { expect, test } from '@grafana/plugin-e2e';

test('renders the PRTG query editor', async ({ panelEditPage, readProvisionedDataSource }) => {
  const dataSource = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(dataSource.name);
  const query = panelEditPage.getQueryEditorRow('A');

  await expect(query.getByText('Query Type', { exact: true })).toBeVisible();
  await expect(query.getByText('Group', { exact: true })).toBeVisible();
  await expect(query.getByText('Device', { exact: true })).toBeVisible();
  await expect(query.getByText('Sensor', { exact: true })).toBeVisible();
  await expect(query.getByText('Channel', { exact: true })).toBeVisible();
  await expect(query.getByText('Enable Streaming', { exact: true })).toBeVisible();
});

test('loads group options from a mocked PRTG resource', async ({ page, panelEditPage, readProvisionedDataSource }) => {
  const dataSource = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.mockResourceResponse('groups', {
    prtgversion: 'test',
    treesize: 2,
    groups: [
      { group: 'Production', objid: 1 },
      { group: 'Development', objid: 2 },
    ],
  });
  await panelEditPage.datasource.set(dataSource.name);

  const groupField = panelEditPage.getQueryEditorRow('A').getByTestId('query-editor-group-field');
  const groupCombobox = groupField.getByRole('combobox');
  await groupCombobox.click();
  await page.getByRole('option', { name: 'Production', exact: true }).click();

  await expect(groupCombobox).toHaveValue('Production');
});

test('shows text query options when Text mode is selected', async ({
  page,
  panelEditPage,
  readProvisionedDataSource,
}) => {
  const dataSource = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(dataSource.name);
  const query = panelEditPage.getQueryEditorRow('A');

  await query.locator('#query-editor-queryType').click();
  await page.getByRole('option', { name: 'Text', exact: true }).click();

  await expect(query.locator('#query-editor-property')).toBeVisible();
  await expect(query.locator('#query-editor-filterProperty')).toBeVisible();
});

test('shows the update interval when streaming is enabled', async ({
  components,
  panelEditPage,
  readProvisionedDataSource,
}) => {
  const dataSource = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(dataSource.name);
  const query = panelEditPage.getQueryEditorRow('A');
  const streamingField = query.getByTestId('query-editor-streaming-field');
  const streamingSwitch = components.switch.within(streamingField);

  await streamingSwitch.check();

  await expect(streamingSwitch).toBeChecked();
  await expect(query.locator('#query-editor-stream-interval')).toBeVisible();
});
