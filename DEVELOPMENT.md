<div align="center">

---

## 🧰 Prerequisites

| Tool    | Required version | Source of truth                                     |
| ------- | ---------------: | --------------------------------------------------- |
| Node.js |        `>= 22` | [`.nvmrc`](./.nvmrc) and `package.json #engines` |
| npm     |       `11.3.0` | `package.json#packageManager`                     |
| Go      |       `1.25.7` | [`go.mod`](./go.mod)                               |
| Mage    |       `1.17.2` | `go.mod` and CI workflows                         |
| Docker  |   Current stable | Docker Compose v2 is required                       |

The plugin targets Grafana OSS `>= 11.0.0`. CI tests the supported Grafana OSS matrix using the official `grafana/grafana` image.

## ⚡ Setup

### 1. Clone the repository

```bash
git clone https://github.com/1DeliDolu/PRTG.git
cd PRTG
```

### 2. Select the Node.js version

```bash
nvm use
node --version
npm --version
```

If your version manager does not read `.nvmrc`, install Node.js 22 manually.

### 3. Install dependencies

```bash
npm ci
go mod download
go install github.com/magefile/mage@v1.17.2
```

Use `npm ci` for reproducible installs. Use `npm install` only when intentionally changing `package.json` or `package-lock.json`.

### 4. Configure the local environment

Create `.env` in the repository root:

```dotenv
PRTG_PATH=prtg.example.com
PRTG_API=replace-with-your-prtg-api-token
PRTG_CACHE_TIME=6000
TIMEZONE=Europe/Berlin
GRAFANA_USERNAME=admin
GRAFANA_PASSWORD=admin
```

`PRTG_PATH` is the host name without `https://`. Never commit `.env`, API tokens, signing credentials, production URLs, or generated authentication state.

## 🔄 Local development

### Build once and start Grafana

The Docker container runs Linux, so build the Linux backend binary before starting it:

```bash
npm run build
mage -v build:linux
docker compose up --build -d
```

Open [http://localhost:3001](http://localhost:3001). Port `2345` is reserved for Delve debugging.

```bash
# Follow Grafana and plugin logs
docker compose logs -f grafana

# Inspect container status
docker compose ps

# Stop and remove local containers
docker compose down
```

### Frontend watch mode

```bash
npm run dev
```

Webpack watches `src/` and writes the updated frontend bundle to `dist/`, which is mounted into the Grafana container.

### Backend development

Rebuild the backend after changes under `pkg/`:

```bash
# Current operating system
mage -v build:backend

# Linux binary used by Docker
mage -v build:linux

# Every supported platform
mage -v buildAll
```

Restart Grafana after rebuilding the backend:

```bash
docker compose restart grafana
```

List all available Mage targets with `mage -l`.

## 🗂️ Project layout

```text
.
├── .config/          # Webpack and Docker development configuration
├── .github/          # CI, E2E, compatibility, and release workflows
├── e2e/              # Plugin E2E helpers
├── pkg/              # Go backend
├── provisioning/     # Local Grafana datasource and dashboards
├── src/              # React and TypeScript frontend
├── tests/            # Playwright test suites
├── go.mod             # Go toolchain and dependencies
└── package.json       # Frontend dependencies and npm scripts
```

## 🧪 Testing

### Frontend quality checks

```bash
npm run typecheck
npm run lint
npm run test:ci
npm run build
```

During active Jest development, use `npm test` for watch mode. Use `npm run lint:fix` to apply supported ESLint and Prettier fixes.

### Backend checks

```bash
mage -v format
mage -v lint
mage -v coverage
```

### Playwright E2E

Install Chromium once, build the plugin, and start Grafana:

```bash
npm exec playwright install chromium
npm run build
mage -v build:linux
docker compose up --build -d
npm run e2e
```

Playwright output is written to `playwright-report/` and `test-results/`. These directories are generated artifacts and must not be committed.

### Full pre-push check

```bash
npm run typecheck
npm run lint
npm run test:ci
npm run build
mage -v coverage
mage -v buildAll
```

## 🧩 Dependency updates

Update dependencies deliberately and keep the lockfiles synchronized.

```bash
# Inspect frontend updates
npm outdated

# Update the Grafana Plugin SDK within the intended compatibility range
go get github.com/grafana/grafana-plugin-sdk-go@latest
go mod tidy
```

After any Grafana package or SDK update, run the full pre-push check and the Playwright suite. Avoid unrelated dependency upgrades in feature pull requests.

## ✅ Pull request checklist

- Keep each pull request focused on one change.
- Add or update Jest, Go, or Playwright tests as appropriate.
- Do not commit `.env`, PRTG credentials, Grafana tokens, reports, or screenshots containing secrets.
- Update [`README.md`](./README.md) when user-facing behavior changes.
- Add a concise entry under **Unreleased** in [`CHANGELOG.md`](./CHANGELOG.md).
- Confirm frontend and backend validation passes locally.
- Include screenshots for visible UI changes.

By participating, you agree to follow the [Code of Conduct](./CODE_OF_CONDUCT.md). Security findings must follow the private process in [`SECURITY.md`](./SECURITY.md), not a public issue.

## 📦 Release process

### Preview builds from `main`

After CI succeeds on a push to `main`, `main-release.yml` publishes or refreshes the `main-build` prerelease. Preview builds may be unsigned when no policy token is configured and are not stable releases.

### Stable releases

Stable release signing requires the GitHub repository secret:

```text
GRAFANA_ACCESS_POLICY_TOKEN
```

The tag must exactly match `package.json`, including the `v` prefix:

```bash
# Creates a commit and tag such as v1.0.1
npm version patch

# Pushes the version commit and tag
git push origin main --follow-tags
```

Pushing `v*` triggers [`.github/workflows/release.yml`](./.github/workflows/release.yml), which verifies the tag, builds frontend and backend artifacts, signs the plugin, validates the archive, creates build provenance, and publishes the GitHub release.

Before tagging:

1. Move entries from **Unreleased** to a dated version in `CHANGELOG.md`.
2. Run the full pre-push check.
3. Verify `package.json` contains the intended version.
4. Confirm `GRAFANA_ACCESS_POLICY_TOKEN` is configured in GitHub Actions secrets.
5. Never print or pass the token as a command-line argument.

## 🔗 References

- [Grafana Plugin Tools](https://grafana.com/developers/plugin-tools/)
- [Grafana backend plugin guide](https://grafana.com/developers/plugin-tools/key-concepts/backend-plugins/)
- [Grafana plugin signing](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin)
- [Playwright documentation](https://playwright.dev/docs/intro)
- [Mage documentation](https://magefile.org/)
