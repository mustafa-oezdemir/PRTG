<div align="center">

<img src="./src/img/logo.svg" alt="PRTG datasource logo" width="112" />

# PRTG Datasource for Grafana

Bring PRTG Network Monitor metrics, status data, and sensor values into Grafana OSS.

[![CI](https://img.shields.io/github/actions/workflow/status/1DeliDolu/PRTG/ci.yml?branch=main&style=for-the-badge&logo=githubactions&logoColor=white&label=CI)](https://github.com/1DeliDolu/PRTG/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/1DeliDolu/PRTG?style=for-the-badge&logo=github&color=2ea44f)](https://github.com/1DeliDolu/PRTG/releases)
[![Grafana OSS](https://img.shields.io/badge/Grafana_OSS-%E2%89%A5_11.0-F46800?style=for-the-badge&logo=grafana&logoColor=white)](https://grafana.com/oss/grafana/)
[![License](https://img.shields.io/github/license/1DeliDolu/PRTG?style=for-the-badge&color=blue)](./LICENSE)

![Node.js](https://img.shields.io/badge/%E2%AC%A2_Node.js-%E2%89%A5_22-5FA04E?style=flat-square&logo=nodedotjs&logoColor=white)
![npm](https://img.shields.io/badge/%F0%9F%93%A6_npm-11.3.0-CB3837?style=flat-square&logo=npm&logoColor=white)
![Go](https://img.shields.io/badge/%F0%9F%90%B9_Go-1.25.7-00ADD8?style=flat-square&logo=go&logoColor=white)
![React](https://img.shields.io/badge/%E2%9A%9B%EF%B8%8F_React-18.2.0-61DAFB?style=flat-square&logo=react&logoColor=101010)
![TypeScript](https://img.shields.io/badge/%F0%9F%94%B7_TypeScript-5.5.4-3178C6?style=flat-square&logo=typescript&logoColor=white)
![Grafana UI](https://img.shields.io/badge/%F0%9F%8E%A8_Grafana_UI-12.4.5-F46800?style=flat-square&logo=grafana&logoColor=white)

[Features](#-features) · [Quick start](#-quick-start) · [Configuration](#%EF%B8%8F-configuration) · [Development](#-development) · [Contributing](#-contributing)

</div>

---

## ✨ Overview

PRTG Datasource is a backend-enabled Grafana plugin that connects Grafana OSS to a PRTG Network Monitor instance. Build dashboards from groups, devices, sensors, and channels without exposing the PRTG API token to the browser.

| Plugin metadata | Value                              |
| --------------- | ---------------------------------- |
| 🧩 Plugin ID    | `maxmarkusprogram-prtg-datasource` |
| 📦 Version      | `1.0.0`                            |
| 📊 Type         | Backend datasource                 |
| 🟠 Grafana      | OSS `>= 11.0.0`                    |
| ⚙️ Backend      | Go + Grafana Plugin SDK `v0.292.0` |
| 🖥️ Frontend     | React + TypeScript + Grafana UI    |
| 📄 License      | Apache 2.0                         |

## 🚀 Features

- 📈 **Metrics queries** — select groups, devices, sensors, and multiple channels.
- 🧾 **Raw queries** — retrieve PRTG properties with optional filters.
- 💬 **Text queries** — display text and status values from PRTG.
- 🛠️ **Manual API queries** — call supported endpoints such as `getsensordetails.json` and `getstatus.htm`.
- 🏷️ **Flexible series labels** — include group, device, and sensor names.
- ⏱️ **Cache and timezone controls** — tune request caching and timestamps per datasource.
- ❤️ **Backend health check** — validate connectivity from Grafana's datasource settings.
- 🔔 **Grafana integrations** — supports metrics, logs, annotations, alerting, and streaming.
- 🧪 **Automated quality checks** — TypeScript, ESLint, Jest, Go tests, and Playwright against Grafana OSS.

## ✅ Requirements

### Run the plugin

![Grafana](https://img.shields.io/badge/%F0%9F%9F%A0_Grafana_OSS-%E2%89%A5_11.0-F46800?style=for-the-badge&logo=grafana&logoColor=white)
![PRTG](https://img.shields.io/badge/%F0%9F%93%A1_PRTG-API_access-0096D6?style=for-the-badge)

- A reachable PRTG Network Monitor server.
- A PRTG API token with access to the objects you want to query.
- Grafana OSS `11.0.0` or newer.

### Build and develop

![Node.js](https://img.shields.io/badge/Node.js-%E2%89%A5_22-5FA04E?style=for-the-badge&logo=nodedotjs&logoColor=white)
![npm](https://img.shields.io/badge/npm-11.3.0-CB3837?style=for-the-badge&logo=npm&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.25.7-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose_v2-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Mage](https://img.shields.io/badge/Mage-1.17.2-7B42BC?style=for-the-badge)

The authoritative versions are defined in [`.nvmrc`](./.nvmrc), [`package.json`](./package.json), and [`go.mod`](./go.mod).

<details>
<summary><strong>🧱 Core dependency versions</strong></summary>

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

## ⚡ Quick start

### 1. Clone and install

```bash
git clone https://github.com/1DeliDolu/PRTG.git
cd PRTG
npm ci
```

### 2. Build the frontend and backend

```bash
npm run build
mage -v build:linux
```

> On Windows, run `mage -v build:windows`. Use `mage -v buildAll` to build all supported platforms.

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

Open [http://localhost:3001](http://localhost:3001). The local environment provisions the PRTG datasource automatically and permits this unsigned development build.

```bash
# Follow logs
docker compose logs -f grafana

# Stop the environment
docker compose down
```

## ⚙️ Configuration

In Grafana, open **Connections → Data sources → PRTG** and enter:

| Setting        | Required | Example            | Description                                               |
| -------------- | :------: | ------------------ | --------------------------------------------------------- |
| 🌐 PRTG server |    ✅    | `prtg.example.com` | Host name without `https://`                              |
| 🔑 API token   |    ✅    | `••••••••`         | Stored in `secureJsonData`; never returned to the browser |
| ⏲️ Cache time  |    —     | `6000`             | Cache duration in seconds; minimum accepted value is `10` |
| 🌍 Timezone    |    —     | `Europe/Berlin`    | Timezone used to normalize PRTG timestamps                |

Select **Save & test** to run the backend health check.

<details>
<summary><strong>📄 Provision with YAML</strong></summary>

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

## 🔎 Query modes

| Mode           | Best for                                   | Inputs                                         |
| -------------- | ------------------------------------------ | ---------------------------------------------- |
| 📈 **Metrics** | Time series, stat, gauge, and alert panels | Group → device → sensor → one or more channels |
| 🧾 **Raw**     | PRTG object metadata and filtered values   | Object selection, property, filter property    |
| 💬 **Text**    | Status messages and other text values      | Object selection and text property             |
| 🛠️ **Manual**  | Supported PRTG API operations              | API method and object ID                       |

### Create your first panel

1. Create a dashboard and select **Add visualization**.
2. Choose the **PRTG** datasource.
3. Select **Metrics** as the query type.
4. Select a group, device, sensor, and one or more channels.
5. Choose a visualization and save the dashboard.

<details>
<summary><strong>🖼️ Configuration and dashboard screenshots</strong></summary>

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

## 🧑‍💻 Development

| Command             | Purpose                                        |
| ------------------- | ---------------------------------------------- |
| `npm run dev`       | Build the frontend in watch mode               |
| `npm run build`     | Create a production frontend build in `dist/`  |
| `npm run typecheck` | Run TypeScript checks without emitting files   |
| `npm run lint`      | Run ESLint                                     |
| `npm run lint:fix`  | Fix lint and formatting issues where possible  |
| `npm run test:ci`   | Run Jest once for CI                           |
| `npm run e2e`       | Run Playwright end-to-end tests                |
| `npm run server`    | Build and start the Docker Compose environment |
| `mage -v coverage`  | Run Go backend tests with coverage             |
| `mage -v buildAll`  | Build backend binaries for all platforms       |

### Recommended validation sequence

```bash
npm run typecheck
npm run lint
npm run test:ci
npm run build
mage -v coverage
```

End-to-end tests use the official `grafana/grafana` OSS image:

```bash
docker compose up -d
npm exec playwright install chromium
npm run e2e
```

## 📦 Build, sign, and release

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

## 🧭 Project structure

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

## 🧯 Troubleshooting

- **Plugin does not appear:** confirm the backend binary was built for the container platform and restart Grafana.
- **Datasource health check fails:** verify the PRTG host, network/DNS access, TLS, and API token permissions.
- **No groups or sensors appear:** ensure the token can read those PRTG objects and check `docker compose logs grafana`.
- **Unsigned plugin error:** use the provided development Compose setup or sign the production artifact.
- **Port conflict:** local Grafana is exposed on host port `3001`, not `3000`.

## 🤝 Contributing

Issues and pull requests are welcome.

1. Fork the repository and create a focused branch.
2. Add or update tests with your change.
3. Run the recommended validation sequence.
4. Open a pull request with a clear description and screenshots for UI changes.

Please follow the [Code of Conduct](./CODE_OF_CONDUCT.md), review the [development guide](./DEVELOPMENT.md), and report security issues according to the [security policy](./SECURITY.md).

## 📚 Resources

- [Grafana Plugin Tools](https://grafana.com/developers/plugin-tools/)
- [Grafana plugin.json reference](https://grafana.com/developers/plugin-tools/reference/plugin-json)
- [Grafana plugin signing](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin)
- [PRTG HTTP API documentation](https://www.paessler.com/manuals/prtg/http_api)

## 📄 License

Copyright © Mustafa Özdemir. Licensed under the [Apache License 2.0](./LICENSE).

---

<div align="center">

Made with ❤️ for the Grafana OSS and PRTG communities.

[![GitHub stars](https://img.shields.io/github/stars/1DeliDolu/PRTG?style=social)](https://github.com/1DeliDolu/PRTG/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/1DeliDolu/PRTG?style=social)](https://github.com/1DeliDolu/PRTG/forks)

</div>
