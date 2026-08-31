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

  await expect(query.getByText('Query Type', { exact: true })).toBeVisible();
  await expect(query.getByText('Group', { exact: true })).toBeVisible();
  await expect(query.getByText('Device', { exact: true })).toBeVisible();
  await expect(query.getByText('Sensor', { exact: true })).toBeVisible();
  await expect(query.getByText('Channel', { exact: true })).toBeVisible();
  await expect(query.getByText('Enable Streaming', { exact: true })).toBeVisible();
});

test('loads group options from a mocked PRTG resource', async ({
  panelEditPage,
  readProvisionedDataSource,
}) => {
  const dataSource = await readProvisionedDataSource({
    fileName: 'datasources.yml',
  });

  const editorPage = panelEditPage.ctx.page;
  const context = editorPage.context();

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

  /*
   * Grafana resource endpoint:
   *
   * /api/datasources/uid/<uid>/resources/groups
   *
   * Query string eklenmesi ihtimalini de destekliyoruz.
   */
  const groupsResource =
    /\/api\/datasources\/uid\/[^/]+\/resources\/groups(?:\?.*)?$/;

  /*
   * Datasource'un gerçekten groups resource'una istek
   * gönderdiğini takip ediyoruz.
   *
   * Bu, CI ortamındaki race condition'ı engeller.
   */
  let resolveGroupsRequest!: () => void;

  const groupsRequestHandled = new Promise<void>((resolve) => {
    resolveGroupsRequest = resolve;
  });

  /*
   * Grafana bazı sürümlerde query editor'ı farklı bir page
   * üzerinden açabildiği için mock'u browser context üzerine
   * bağlıyoruz.
   */
  await context.route(groupsResource, async (route) => {
    await route.fulfill({
      status: 200,
      json: groupsResponse,
    });

    resolveGroupsRequest();
  });

  /*
   * Datasource seçildiğinde QueryEditor mount olur ve
   * getGroups() çağrısının yapılması beklenir.
   */
  await panelEditPage.datasource.set(dataSource.name);

  /*
   * Combobox ile işlem yapmadan önce mock endpoint'in
   * gerçekten çağrıldığından emin oluyoruz.
   */
  await groupsRequestHandled;

  const query = panelEditPage.getQueryEditorRow('A');

  const groupField = query.getByTestId('query-editor-group-field');

  await expect(groupField).toBeVisible();

  const groupCombobox = groupField.getByRole('combobox');

  await expect(groupCombobox).toBeVisible();

  /*
   * Select'i aç.
   */
  await groupCombobox.click();

  /*
   * API cevabının React state'e aktarılması async olabilir.
   * expect(...).toBeVisible() Playwright'in auto-wait
   * mekanizmasını kullanır.
   */
  const productionOption = editorPage.getByRole('option', {
    name: 'Production',
    exact: true,
  });

  await expect(productionOption).toBeVisible();

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

  await query.locator('#query-editor-queryType').click();

  await page
    .getByRole('option', {
      name: 'Text',
      exact: true,
    })
    .click();

  await expect(query.locator('#query-editor-property')).toBeVisible();

  await expect(query.locator('#query-editor-filterProperty')).toBeVisible();
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

  const streamingSwitch = components.switch.within(streamingField);

  await streamingSwitch.check();

  await expect(streamingSwitch).toBeChecked();

  await expect(
    query.locator('#query-editor-stream-interval')
  ).toBeVisible();
});
