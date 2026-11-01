# PRTG Datasource for Grafana

Explore PRTG groups, devices, sensors, and channels in Grafana dashboards. The Go backend handles PRTG requests, so the API token stays out of the browser.

[![CI](https://github.com/mustafa-oezdemir/PRTG/actions/workflows/ci.yml/badge.svg)](https://github.com/mustafa-oezdemir/PRTG/actions/workflows/ci.yml)
[![Grafana OSS](https://img.shields.io/badge/Grafana%20OSS-13.3.0%20tested-F46800?logo=grafana&logoColor=white)](https://github.com/mustafa-oezdemir/PRTG/actions/workflows/ci.yml)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue)](./LICENSE)
[![Sponsor via PayPal](https://img.shields.io/badge/Sponsor-via%20PayPal-EC4899?logo=paypal&logoColor=white)](https://paypal.me/pehlion)

**Plugin ID:** `maxmarkusprogram-prtg-datasource` · **Version:** `1.0.0` · **Grafana:** OSS 12.0.0 or newer

[Quick start](#quick-start) · [Configuration](#configuration) · [Query modes](#query-modes) · [Development](#development) · [Support](#support-the-project)

## Features

- **Metrics, raw, text, and manual queries** for PRTG objects and API endpoints.
- **Multiple channels and flexible series labels** for useful dashboard panels.
- **Backend health checks, caching, and timezone controls** in the datasource settings.
- **Grafana integrations** for metrics, logs, annotations, alerting, and streaming.

## Requirements

### Run

- Grafana OSS 12.0.0 or newer. CI tests the current supported matrix, including 13.3.0.
- A reachable PRTG Network Monitor server and a token with access to the objects you want to query.

### Build and develop

| Tool    | Version                                |
| ------- | -------------------------------------- |
| Node.js | 26.13.0 LTS ([`.nvmrc`](./.nvmrc))     |
| npm     | 11.3.0 (`package.json#packageManager`) |
| Go      | 1.27.3 ([`go.mod`](./go.mod))          |
| Mage    | 1.17.2                                 |
| Docker  | Compose v2                             |

<details>
<summary>Core dependency versions</summary>

| Dependency                |    Version |
| ------------------------- | ---------: |
| React / React DOM         |   `18.2.0` |
| TypeScript                |    `5.5.4` |
| Grafana Data              |   `12.4.5` |
| Grafana Runtime           |   `12.4.5` |
| Grafana Schema            |   `12.4.5` |
| Grafana UI                |   `12.4.5` |
| Emotion CSS               |  `11.10.6` |
| Grafana Plugin SDK for Go | `v0.292.0` |
| Playwright                |   `1.61.1` |
| Jest                      |   `29.7.0` |

Versions in this table reflect the current lockfiles. Declared compatible ranges remain in `package.json` and `go.mod`.

</details>

## Quick start

### 1. Clone and install

```bash
git clone https://github.com/mustafa-oezdemir/PRTG.git
cd PRTG
nvm use
npm ci
go install github.com/magefile/mage@v1.17.2
```

### 2. Build the frontend and backend

```bash
npm run build
mage -v build:linux
```

On Windows, run `mage -v build:windows`. Use `mage -v buildAll` to build all supported platforms.

### 3. Provide local settings

Create a local `.env` file. Never commit this file or paste its token into an issue or CI log.

```dotenv
PRTG_PATH=prtg.example.com
PRTG_API=replace-with-your-prtg-api-token
PRTG_CACHE_TIME=6000
TIMEZONE=Europe/Berlin
GRAFANA_USERNAME=admin
GRAFANA_PASSWORD=admin
```

`PRTG_PATH` must contain the host name only—do not include `https://`.

### 4. Start Grafana OSS

```bash
docker compose up --build -d
```

Open [http://localhost:3001](http://localhost:3001). Docker Compose starts Grafana OSS 13.3.0 by default, provisions the PRTG datasource, and permits the unsigned development build. Set `GRAFANA_VERSION` to test another Grafana version.

```bash
# Follow logs
docker compose logs -f grafana

# Stop the environment
docker compose down
```

## Configuration

In Grafana, open **Connections → Data sources → PRTG** and enter:

| Setting     | Required | Example            | Description                                               |
| ----------- | :------: | ------------------ | --------------------------------------------------------- |
| PRTG server |   Yes    | `prtg.example.com` | Host name without `https://`                              |
| API token   |   Yes    | `••••••••`         | Stored in `secureJsonData`; never returned to the browser |
| Cache time  |    No    | `6000`             | Duration in seconds; minimum `10`                         |
| Timezone    |    No    | `Europe/Berlin`    | Used to normalize PRTG timestamps                         |

Select **Save & test** to run the backend health check.

<details>
<summary>Provision with YAML</summary>

```yaml
apiVersion: 1

datasources:
  - name: PRTG
    uid: prtg
    type: maxmarkusprogram-prtg-datasource
    access: proxy
    editable: true
    jsonData:
      path: prtg.example.com
      cacheTime: 6000
      timeZone: Europe/Berlin
    secureJsonData:
      apiKey: your-prtg-api-token
```

For deployments, inject the token through your secret manager instead of committing it in YAML.

</details>

## Query modes

| Mode        | Best for                                   | Inputs                                         |
| ----------- | ------------------------------------------ | ---------------------------------------------- |
| **Metrics** | Time series, stat, gauge, and alert panels | Group → device → sensor → one or more channels |
| **Raw**     | PRTG object metadata and filtered values   | Object selection, property, filter property    |
| **Text**    | Status messages and other text values      | Object selection and text property             |
| **Manual**  | Supported PRTG API operations              | API method and object ID                       |

### Create your first panel

1. Create a dashboard and select **Add visualization**.
2. Choose the **PRTG** datasource.
3. Select **Metrics** as the query type.
4. Select a group, device, sensor, and one or more channels.
5. Choose a visualization and save the dashboard.

<details>
<summary>Configuration and dashboard screenshots</summary>

### Datasource setup

![PRTG datasource configuration](image/README/1739793866048.png)

![PRTG datasource save and test](image/README/1739793921893.png)

### Metrics query

![PRTG metrics query](image/README/1739795402834.png)

![PRTG metrics panel](image/README/1739795452206.png)

### Raw and text queries

![PRTG raw query](image/README/1739796514348.png)

![PRTG text query](image/README/1739796808632.png)

</details>

## Development

| Command             | Purpose                                           |
| ------------------- | ------------------------------------------------- |
| `npm run dev`       | Build the frontend in watch mode                  |
| `npm run build`     | Build the production frontend in `dist/`          |
| `npm run typecheck` | Check TypeScript without emitting files           |
| `npm run lint`      | Run ESLint                                        |
| `npm run test:ci`   | Run the Jest suite once                           |
| `npm run e2e`       | Run Playwright against a running Grafana instance |
| `mage -v coverage`  | Run Go tests with coverage                        |
| `mage -v buildAll`  | Build all supported backend binaries              |

### Recommended validation sequence

```bash
npm run typecheck
npm run lint
npm run test:ci
npm run build
mage -v coverage
```

End-to-end tests use the official `grafana/grafana` OSS image. The CI matrix includes the latest stable 13.3.0 release and other supported versions.

```bash
docker compose up -d
npm exec playwright install chromium
npm run e2e
```

## Build, sign, and release

Grafana requires signed plugins outside development mode. Set `GRAFANA_ACCESS_POLICY_TOKEN` in your local environment or as a GitHub Actions repository secret, then run:

```bash
npm run build
mage -v buildAll
npm run sign
```

Do not add signing tokens to `.env.example`, source control, screenshots, or workflow output. See [Grafana plugin signing](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin) for publisher and signature requirements.

Version tags trigger the repository's release workflow:

```bash
npm version patch
git push origin main --follow-tags
```

## Project structure

```text
.
├── .config/          # Webpack, Docker, and development configuration
├── .github/          # CI, E2E, signing, and release workflows
├── e2e/              # Playwright end-to-end tests
├── pkg/              # Go backend implementation
├── provisioning/     # Local Grafana datasource provisioning
├── src/              # React/TypeScript plugin frontend
├── go.mod             # Go toolchain and backend dependencies
└── package.json       # Node.js scripts and frontend dependencies
```

## Troubleshooting

- **Plugin does not appear:** confirm the backend binary was built for the container platform and restart Grafana.
- **Datasource health check fails:** verify the PRTG host, network/DNS access, TLS, and API token permissions.
- **No groups or sensors appear:** ensure the token can read those PRTG objects and check `docker compose logs grafana`.
- **Unsigned plugin error:** use the provided development Compose setup or sign the production artifact.
- **Port conflict:** local Grafana is exposed on host port `3001`, not `3000`.

## Contributing

Issues and pull requests are welcome.

1. Fork the repository and create a focused branch.
2. Add or update tests with your change.
3. Run the recommended validation sequence.
4. Open a pull request with a clear description and screenshots for UI changes.

Please follow the [Code of Conduct](./CODE_OF_CONDUCT.md), review the [development guide](./DEVELOPMENT.md), and report security issues according to the [security policy](./SECURITY.md).

## Support the project

If this plugin helps you, [sponsor its development via PayPal](https://paypal.me/pehlion).

## Resources

- [Grafana Plugin Tools](https://grafana.com/developers/plugin-tools/)
- [Grafana plugin.json reference](https://grafana.com/developers/plugin-tools/reference/plugin-json)
- [Grafana plugin signing](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin)
- [PRTG HTTP API documentation](https://www.paessler.com/manuals/prtg/http_api)

## License

Copyright © Mustafa Özdemir. Licensed under the [Apache License 2.0](./LICENSE).
