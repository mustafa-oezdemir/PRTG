# Changelog

All notable changes to the PRTG datasource for Grafana are documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Backend health checks, structured logging, metrics, and tracing.
- Metrics, raw, text, and manual PRTG query modes.
- Multi-channel metric selection and configurable series labels.
- Grafana annotations and live streaming support.
- Configurable request caching and timezone handling.
- Provisioned local datasource and health dashboard examples.
- Jest, Go, and Playwright coverage for frontend, backend, configuration, and query flows.
- Automated preview releases with archive checksums and build provenance.

### Changed

- Restructured the repository into top-level frontend, backend, provisioning, and test directories.
- Updated the supported runtime to Grafana OSS `>= 11.0.0`, Node.js `>= 22`, Go `1.25.7`, and Grafana packages `12.x`.
- Migrated selection controls to the current Grafana UI Combobox APIs.
- Updated CI and scheduled E2E workflows to test official `grafana/grafana` OSS images.
- Updated release signing to use `GRAFANA_ACCESS_POLICY_TOKEN`.
- Modernized project, development, security, and community documentation.

### Fixed

- Corrected plugin archive layout and backend executable permissions.
- Stabilized datasource health checks, E2E authentication, and dashboard datasource references.
- Improved timezone parsing, cache defaults, query editor state, and multi-query behavior.
- Prevented enterprise Grafana images from being selected by the OSS test matrix.

### Security

- Kept PRTG API credentials in Grafana `secureJsonData` and out of frontend responses.
- Added release validation, signing checks, checksums, and artifact provenance.
- Documented a private vulnerability reporting and coordinated disclosure process.

## [1.0.0] - 2025-02-17

### Added

- Initial PRTG datasource plugin release for Grafana.
- PRTG server and API token configuration.
- Group, device, sensor, and channel-based metric queries.
- Go backend binaries and React/TypeScript query editor.
- Apache 2.0 license and initial project documentation.

[Unreleased]: https://github.com/1DeliDolu/PRTG/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/1DeliDolu/PRTG/releases/tag/v1.0.0
