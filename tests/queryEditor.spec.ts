import { expect, test } from '@grafana/plugin-e2e';

test('renders the PRTG query editor', async ({
  panelEditPage,
  readProvisionedDataSource,
}) => {
  const dataSource = await readProvisionedDataSource({
    fileName: 'datasources.yml',
  });

  await panelEditPage.datasource.set(dataSource.name);

  const query = panelEditPage.getQueryEditorRow('A');

  await expect(
    query.getByText('Query Type', { exact: true })
  ).toBeVisible();

  await expect(
    query.getByText('Group', { exact: true })
  ).toBeVisible();

  await expect(
    query.getByText('Device', { exact: true })
  ).toBeVisible();

  await expect(
    query.getByText('Sensor', { exact: true })
  ).toBeVisible();

  await expect(
    query.getByText('Channel', { exact: true })
  ).toBeVisible();

  await expect(
    query.getByText('Enable Streaming', { exact: true })
  ).toBeVisible();
});

test('loads group options from a mocked PRTG resource', async ({
  page,
  panelEditPage,
  readProvisionedDataSource,
}) => {
  const dataSource = await readProvisionedDataSource({
    fileName: 'datasources.yml',
  });

  const groupsResponse = {
    prtgversion: 'test',
    treesize: 2,
    groups: [
      {
        group: 'Production',
        objid: 1,
      },
      {
        group: 'Development',
        objid: 2,
      },
    ],
  };

  /**
   * IMPORTANT:
   *
   * usePrtgSelectionLists() calls datasource.getGroups()
   * immediately when the QueryEditor mounts.
   *
   * Therefore the mock must be registered BEFORE selecting
   * the datasource.
   *
   * Do not use context.route() here. @grafana/plugin-e2e
   * already provides the correct resource mocking helper.
   */
  await panelEditPage.mockResourceResponse(
    'groups',
    groupsResponse
  );

  /**
   * Selecting the datasource mounts QueryEditor.
   *
   * QueryEditor -> usePrtgSelectionLists()
   *             -> datasource.getGroups()
   *             -> mocked /resources/groups
   */
  await panelEditPage.datasource.set(dataSource.name);

  const query = panelEditPage.getQueryEditorRow('A');

  const groupField = query.getByTestId(
    'query-editor-group-field'
  );

  await expect(groupField).toBeVisible();

  const groupCombobox = groupField.getByRole('combobox');

  await expect(groupCombobox).toBeVisible();
  await expect(groupCombobox).toBeEnabled();

  /**
   * Open the dropdown.
   */
  await groupCombobox.click();

  /**
   * Combobox options are rendered in a portal, therefore
   * search from the page instead of inside groupField.
   */
  const productionOption = page.getByRole('option', {
    name: 'Production',
    exact: true,
  });

  await expect(productionOption).toBeVisible({
    timeout: 10000,
  });

  await productionOption.click();

  await expect(groupCombobox).toHaveValue('Production');
});

test('shows text query options when Text mode is selected', async ({
  page,
  panelEditPage,
  readProvisionedDataSource,
}) => {
  const dataSource = await readProvisionedDataSource({
    fileName: 'datasources.yml',
  });

  await panelEditPage.datasource.set(dataSource.name);

  const query = panelEditPage.getQueryEditorRow('A');

  const queryTypeCombobox = query.locator(
    '#query-editor-queryType'
  );

  await expect(queryTypeCombobox).toBeVisible();

  await queryTypeCombobox.click();

  const textOption = page.getByRole('option', {
    name: 'Text',
    exact: true,
  });

  await expect(textOption).toBeVisible();

  await textOption.click();

  await expect(
    query.locator('#query-editor-property')
  ).toBeVisible();

  await expect(
    query.locator('#query-editor-filterProperty')
  ).toBeVisible();
});

test('shows the update interval when streaming is enabled', async ({
  components,
  panelEditPage,
  readProvisionedDataSource,
}) => {
  const dataSource = await readProvisionedDataSource({
    fileName: 'datasources.yml',
  });

  await panelEditPage.datasource.set(dataSource.name);

  const query = panelEditPage.getQueryEditorRow('A');

  const streamingField = query.getByTestId(
    'query-editor-streaming-field'
  );

  await expect(streamingField).toBeVisible();

  const streamingSwitch =
    components.switch.within(streamingField);

  await streamingSwitch.check();

  await expect(streamingSwitch).toBeChecked();

  await expect(
    query.locator('#query-editor-stream-interval')
  ).toBeVisible();
});
