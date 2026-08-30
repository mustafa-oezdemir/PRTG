<div align="center">

## Overview

This repository contains a Grafana datasource plugin for PRTG, allowing users to visualize and analyze PRTG metrics within Grafana.

The plugin integrates with PRTG and fetches data from PRTG sensors directly into Grafana dashboards. It provides a practical way to monitor, analyze, and visualize PRTG data with Grafana panels.

Plugin metadata read from `maxmarkusprogram-prtg-datasource`:

| Item               | Value                                |
| ------------------ | ------------------------------------ |
| Plugin ID          | `maxmarkusprogram-prtg-datasource` |
| Plugin type        | Grafana datasource                   |
| Package name       | `prtg`                             |
| Version            | `1.0.0`                            |
| Grafana dependency | `>=10.4.0`                         |
| Backend            | Go / Grafana Plugin SDK              |
| Frontend           | TypeScript, React, Grafana UI        |
| License            | Apache 2.0                           |

## Features

- Grafana datasource plugin for PRTG Network Monitor.
- Backend plugin support with health checks and PRTG API access.
- Metrics query flow: group, device, sensor, and channel selection.
- Raw query flow with property and filter property selection.
- Text query flow for displaying PRTG text data.
- Manual API methods such as `getsensordetails.json` and `getstatus.htm`.
- Optional panel labels for group, device, and sensor names.
- Cache time and timezone configuration.
- Grafana metrics, logs, annotations, alerting, and backend support declared in `plugin.json`.

## Requirements

Install these tools before building or running the plugin:

- [Docker](https://www.docker.com/)
- [Go](https://go.dev/dl/)
- Node.js `>=22`
- npm `11.3.0` or compatible
- Grafana `>=10.4.0`

## Quick Start

Clone the repository:

```sh
git clone https://github.com/1DeliDolu/PRTG.git
```

Navigate to the plugin directory:

```sh
cd PRTG/maxmarkusprogram-prtg-datasource
```

Install frontend dependencies:

```sh
npm install
```

Build the frontend plugin:

```sh
npm run build
```

Build the backend plugin:

```sh
mage
```

Copy the built plugin to Grafana's plugin directory:

```sh
cp -r dist /var/lib/grafana/plugins/PRTG
```

Restart Grafana on Linux or WSL:

```sh
sudo systemctl restart grafana-server
```

Restart Grafana on Windows PowerShell:

```powershell
net stop grafana
net start grafana
```

### Windows Manual Copy Alternative

After cloning and building, copy the `Prtg` folder to:

```text
C:\Program Files\GrafanaLabs\grafana\data\plugins
```

Then restart Grafana.

## Configuration

1. Open Grafana and navigate to the Data Sources page.
2. Click **Add data source** and select **PRTG**.
3. Configure the PRTG datasource with the required connection details.
4. Enter the PRTG server path without `https://`.
5. Enter the API key or API token.
6. Set cache time in seconds. The editor accepts values from `10` seconds upward.
7. Select the timezone. If no timezone is configured, the backend defaults to `Europe/Berlin`.
8. Save and test the datasource.

Provisioning example from `maxmarkusprogram-prtg-datasource/provisioning/datasources/datasources.yml`:

```yaml
apiVersion: 1

datasources:
  - name: 'PRTG'
    type: 'maxmarkusprogram-prtg-datasource'
    access: proxy
    isDefault: false
    orgId: 1
    version: 1
    editable: true
    jsonData:
      path: 'your-prtg-server'
      cacheTime: 6000
    secureJsonData:
      apiKey: 'your-api-token'
```

## Usage

1. Create a new dashboard or open an existing Grafana dashboard.
2. Add a new panel and select the PRTG datasource.
3. Configure the query to fetch data from the desired PRTG sensors.
4. Customize the visualization settings to display the data as needed.

Supported query types:

| Query type  | Purpose                                                           |
| ----------- | ----------------------------------------------------------------- |
| `Metrics` | Select group, device, sensor, and channel for time series panels. |
| `Raw`     | Select PRTG properties and filter properties.                     |
| `Text`    | Display text-based PRTG values.                                   |
| `Manual`  | Call supported manual PRTG API endpoints.                         |

## Development

Run the plugin in development mode:

```sh
npm run dev
```

Start Grafana with Docker for local plugin testing:

```sh
npm run server
```

Run frontend tests:

```sh
npm run test:ci
```

Run E2E tests:

```sh
npm run e2e
```

Run linting:

```sh
npm run lint
```

Sign the plugin when distributing outside local development:

```sh
npm run sign
```

## Troubleshooting

- Ensure the PRTG server path and API key are correctly configured.
- Check Grafana server logs for error messages.
- Verify that the plugin is installed in Grafana's plugin directory.
- Restart Grafana and try again.
- For development builds, make sure Docker, Go, Node.js, npm, and Mage are available.

## Resources

- [Grafana Plugin Development Documentation](https://grafana.com/developers/plugin-tools/)
- [PRTG API Documentation](https://www.paessler.com/manuals/prtg/api)
- [Grafana plugin.json documentation](https://grafana.com/developers/plugin-tools/reference/plugin-json)
- [Grafana plugin signing documentation](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin)

Contributions are welcome through issues and pull requests.

This README is the maintained repository documentation. The original README note to save the content in `README.md` is now reflected here.

## Config Editor

Open `http://localhost:3000/` in a browser or go directly to `http://grafana.prtg:3000/connections/datasources`.

![Open Grafana datasource page](image/README/1739793462631.png)

Press **PRTG**.

![Select PRTG datasource](image/README/1739793798353.png)

Enter your PRTG server.

![Enter PRTG server](image/README/1739793866048.png)

Enter your API token and press **Save & test**.

![Save and test datasource](image/README/1739793921893.png)

Press **Build a dashboard**.

![Build dashboard](image/README/1739794001603.png)

Press **Add visualization**.

![Add visualization](image/README/1739794068185.png)

Press **PRTG**.

![Select PRTG panel datasource](image/README/1739794166798.png)

## Query Metrics

Select **Query Type**.

![Select query type](image/README/1739795234405.png)

Select **Group**.

![Select group](image/README/1739795274666.png)

Select **Device**.

![Select device](image/README/1739795311709.png)

Select **Sensor**.

![Select sensor](image/README/1739795351207.png)

Select **Channel**.

![Select channel](image/README/1739795402834.png)

Look at the panel.

![Metrics panel](image/README/1739795452206.png)

## Options

Add the group name in the panel.

![Add group name](image/README/1739795578616.png)

Add device and sensor names.

![Add device and sensor](image/README/1739795687023.png)

Add a new query.

![Add new query](image/README/1739795739188.png)

Select query, group, device, sensor, and channel.

![Select query details](image/README/1739795941422.png)

Another example.

![Another example](image/README/1739796156994.png)

Fill opacity.

![Fill opacity](image/README/1739796291106.png)

Select **Stat**.

![Select stat visualization](image/README/1739796324396.png)

## Query Raw

Select **Query Raw**, group, device, sensor, property, and filter property.

![Query raw configuration](image/README/1739796514348.png)

Examples:

![Query raw examples](image/README/1739796591456.png)

## Query Text

Query text examples:

![Query text example 1](image/README/1739796808632.png)

![Query text example 2](image/README/1739796830021.png)

## Panel

![Panel example 1](image/README/1739797181230.png)

![Panel example 2](image/README/1739797385328.png)

![Panel example 3](image/README/1739797413883.png)

## Video Walkthrough

Copy the repository:

```text
https://github.com/1DeliDolu/PRTG.git
```

<video width="1000" height="500" controls>
  <source src="./video/clone.mp4" type="video/mp4">
</video>

Open VS Code and a Bash terminal, then clone the repository.

<video width="1000" height="500" controls>
  <source src="./video/vsc_clone.mp4" type="video/mp4">
</video>

Open a new WSL terminal, enter the `maxmarkusprogram-prtg-datasource` folder, and run:

```sh
cd maxmarkusprogram-prtg-datasource
```

<video width="1000" height="500" controls>
  <source src="./video/cd.mp4" type="video/mp4">
</video>

Install dependencies:

```sh
npm install
```

<video width="1000" height="500" controls>
  <source src="./video/npminstall.mp4" type="video/mp4">
</video>

Build the frontend:

```sh
npm run build
```

<video width="1000" height="500" controls>
  <source src="./video/build.mp4" type="video/mp4">
</video>

Build the backend:

```sh
mage
```

<video width="1000" height="500" controls>
  <source src="./video/mage.mp4" type="video/mp4">
</video>

Rename or move `dist` to `Prtg`.

<video width="1000" height="500" controls>
  <source src="./video/prtg.mp4" type="video/mp4">
</video>

Close VS Code.

<video width="1000" height="500" controls>
  <source src="./video/closevsc.mp4" type="video/mp4">
</video>

Copy the `Prtg` folder.

<video width="1000" height="500" controls>
  <source src="./video/copy.mp4" type="video/mp4">
</video>

Paste `Prtg` into the Grafana plugin directory:

```text
C:\Program Files\GrafanaLabs\grafana\data\plugins
```

<video width="1000" height="500" controls>
  <source src="./video/paste.mp4" type="video/mp4">
</video>

Stop and start Grafana with PowerShell:

```powershell
net stop grafana
net start grafana
```

<video width="1000" height="500" controls>
  <source src="./video/stop-start-grafana.mp4" type="video/mp4">
</video>

Sign in. On the first login, Grafana's default username and password are usually:

```text
user: admin
password: admin
```

<video width="1000" height="500" controls>
  <source src="./video/anmeldung.mp4" type="video/mp4">
</video>

Open the datasource panel.

<video width="1000" height="500" controls>
  <source src="./video/datasource.mp4" type="video/mp4">
</video>

Create a query.

<video width="1000" height="500" controls>
  <source src="./video/query.mp4" type="video/mp4">
</video>

Open the dashboard.

<video width="1000" height="500" controls>
  <source src="./video/dashboard.mp4" type="video/mp4">
</video>
