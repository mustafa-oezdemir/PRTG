<div align="center">

# 🔐 Security Policy

How to report vulnerabilities in the PRTG datasource for Grafana safely and privately.

[![Security policy](https://img.shields.io/badge/Security-Responsible_disclosure-2ea44f?style=for-the-badge&logo=github)](https://github.com/mustafa-oezdemir/PRTG/security)
[![Supported release](https://img.shields.io/badge/Supported-1.0.x-blue?style=for-the-badge)](https://github.com/mustafa-oezdemir/PRTG/releases)

</div>

## ✅ Supported versions

Security fixes are applied to the latest stable release line.

| Version               |                  Status                  |
| --------------------- | :--------------------------------------: |
| `1.0.x`               |               ✅ Supported               |
| `< 1.0.0`             |              ❌ Unsupported              |
| `main` / `main-build` | 🧪 Development preview; best effort only |

Users should upgrade to the newest stable release before reporting an issue that may already be fixed.

## 🚨 Report a vulnerability

**Do not open a public GitHub issue, discussion, or pull request for a suspected vulnerability.**

Use one of these private channels:

1. **Preferred:** [Open a private GitHub security advisory](https://github.com/mustafa-oezdemir/PRTG/security/advisories/new).
2. **Fallback:** Email [mustafa.ozdemir1408@gmail.com](mailto:mustafa.ozdemir1408@gmail.com?subject=%5BSECURITY%5D%20PRTG%20Datasource).

Include enough information to reproduce and assess the issue:

- A short description and the affected plugin version or commit.
- The Grafana OSS and PRTG versions involved.
- The expected behavior and observed security impact.
- Reproduction steps or a minimal proof of concept.
- Relevant logs with tokens, cookies, host names, and personal data removed.
- Any known mitigations or suggested fixes.
- Your preferred name for acknowledgment, or state that you prefer anonymity.

Never send live API tokens, passwords, session cookies, private keys, production data, or an unredacted `.env` file.

## ⏱️ What to expect

These are best-effort response targets, not a service-level agreement:

| Stage                  |                                       Target |
| ---------------------- | -------------------------------------------: |
| Initial acknowledgment |                              Within 72 hours |
| Preliminary triage     |                       Within 7 calendar days |
| Progress updates       | At least every 14 calendar days while active |

After triage, the maintainer will confirm whether the report is accepted, needs more information, is a duplicate, or is out of scope. Valid issues will be coordinated privately until a fix and disclosure plan are ready.

## 🎯 Scope

Examples of in-scope issues include:

- Exposure or improper handling of PRTG API tokens or Grafana secrets.
- Authentication or authorization bypasses in plugin resources.
- Server-side request forgery, injection, path traversal, or unsafe deserialization.
- Cross-site scripting or unsafe rendering caused by the plugin frontend.
- Sensitive information leakage through logs, traces, metrics, build artifacts, or error messages.
- Compromise of the plugin's build, signing, packaging, or release process.

Generally out of scope:

- Vulnerabilities in Grafana, PRTG, Docker, browsers, or other upstream dependencies without a plugin-specific impact.
- Reports that only identify an outdated dependency without a reachable exploit path.
- Denial-of-service testing against public or production systems.
- Social engineering, phishing, physical attacks, or credential stuffing.
- Findings that require access to already-compromised administrator credentials without additional impact.

Report upstream vulnerabilities directly to the affected vendor. If an upstream issue becomes exploitable through this plugin, explain that plugin-specific path in the report.

## 🤝 Coordinated disclosure

Please allow reasonable time to investigate, prepare a fix, test supported Grafana versions, and publish an update before public disclosure. The maintainer will aim to:

- Minimize the collection and retention of reporter data.
- Share report details only with people needed to resolve the issue.
- Credit reporters who request acknowledgment.
- Publish remediation guidance when users need to take action.

This project currently does not operate a paid bug bounty program.

## 🛡️ Deployment guidance

- Store PRTG tokens in Grafana `secureJsonData` or an external secret manager.
- Keep `.env` and signing tokens out of source control and CI logs.
- Grant tokens only the PRTG permissions the datasource needs.
- Use TLS between Grafana and PRTG and restrict network access where possible.
- Run a supported Grafana OSS version and keep the plugin updated.
- Verify release checksums and signatures before production deployment.

For general bugs and feature requests that do not have security impact, use [GitHub Issues](https://github.com/mustafa-oezdemir/PRTG/issues).
