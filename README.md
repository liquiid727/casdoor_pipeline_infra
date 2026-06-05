<div align="center">
  <a href="https://github.com/liquiid727/pipeline-auth">
    <img src="https://raw.githubusercontent.com/liquiid727/pipeline-auth/main/web/public/img/logo.png" alt="Pipeline Auth" width="500">
  </a>

  <h3>Pipeline Auth: Independent IAM / MCP Gateway Fork</h3>

  <p align="center">
    <strong>An open-source, AI-first IAM / MCP gateway and authentication server with a web UI.</strong><br>
    Supporting MCP, A2A, OAuth&nbsp;2.0, OIDC (OAuth&nbsp;2.x), SAML, CAS, LDAP, SCIM, WebAuthn, TOTP, MFA, Face ID,<br>
    Google Workspace, Azure AD, and more.
  </p>

  <p align="center">
    <a href="https://github.com/liquiid727/pipeline-auth"><strong>Repository: liquiid727/pipeline-auth</strong></a>
  </p>

  <p>
    <a href="https://github.com/liquiid727/pipeline-auth">
      <img src="https://img.shields.io/badge/repository-pipeline--auth-1890ff?style=flat-square&logo=github&logoColor=white" alt="Repository">
    </a>
    <a href="https://github.com/liquiid727/pipeline-auth/releases/latest">
      <img src="https://img.shields.io/github/v/release/casdoor/casdoor?style=flat-square&color=blue" alt="GitHub Release">
    </a>
    <a href="https://github.com/liquiid727/pipeline-auth/actions/workflows/build.yml">
      <img src="https://img.shields.io/github/actions/workflow/status/liquiid727/pipeline-auth/build.yml?style=flat-square&label=build" alt="Build Status">
    </a>
    <a href="https://goreportcard.com/report/github.com/liquiid727/pipeline-auth">
      <img src="https://goreportcard.com/badge/github.com/liquiid727/pipeline-auth?style=flat-square" alt="Go Report Card">
    </a>
    <a href="https://github.com/liquiid727/pipeline-auth/blob/master/LICENSE">
      <img src="https://img.shields.io/github/license/liquiid727/pipeline-auth?style=flat-square&color=orange" alt="License">
    </a>
  </p>

  <p>
    <a href="https://github.com/liquiid727/pipeline-auth/stargazers">
      <img src="https://img.shields.io/github/stars/liquiid727/pipeline-auth?style=flat-square&color=yellow" alt="GitHub Stars">
    </a>
    <a href="https://github.com/liquiid727/pipeline-auth/network/members">
      <img src="https://img.shields.io/github/forks/liquiid727/pipeline-auth?style=flat-square" alt="GitHub Forks">
    </a>
    <a href="https://github.com/liquiid727/pipeline-auth/issues">
      <img src="https://img.shields.io/github/issues/liquiid727/pipeline-auth?style=flat-square&color=red" alt="GitHub Issues">
    </a>
  </p>

  <p align="center">
    <a href="https://github.com/liquiid727/pipeline-auth"><strong>Repository</strong></a> ·
    <a href="https://github.com/liquiid727/pipeline-auth/issues"><strong>Issues</strong></a>
  </p>
</div>

---

## Table of contents

- [Why Pipeline Auth](#why-pipeline-auth)
- [Quick start](#quick-start)
- [Features](#features)
- [Technology stack](#technology-stack)
- [Documentation](#documentation)
- [Integrations](#integrations)
- [Security](#security)
- [Community and support](#community-and-support)
- [Contributing](#contributing)
- [License](#license)

---

<a id="why-pipeline-auth"></a>
## Why Pipeline Auth

Pipeline Auth is an independent fork of Casdoor focused on running as a self-contained IAM / MCP gateway codebase without depending on a sibling `casdoor` checkout or upstream Casdoor self-branding.

---

<a id="quick-start"></a>
## 🚀 Quick start

Pick one deployment method below. To keep behavior consistent with upstream, the steps are aligned with official docs.

### 🛠️ Source code (default)

1. Install dependencies: **Go 1.25** (follow `go.mod`), **Node.js LTS (20)**, **Yarn 1.x**, and a supported database.
2. Clone the repository:

```bash
git clone https://github.com/liquiid727/pipeline-auth.git
cd casdoor
```

3. Configure database in `conf/app.conf` (at minimum set `driverName`, `dataSourceName`, and `dbName`; for MySQL create database `casdoor` first).
4. Build frontend and start backend:

```bash
cd web
yarn install
yarn build
cd ..
go run main.go
```

5. Open [http://localhost:8000](http://localhost:8000) and sign in with `built-in/admin` / `123` on a fresh install (change password immediately in production).

### 🐳 Docker

Use one of the official Docker paths:

- **All-in-one (SQLite quick trial)**:

```bash
docker run -p 8000:8000 casbin/casdoor-all-in-one
```

- **Docker Compose** (with your `conf/app.conf` next to `docker-compose.yml`):

```bash
docker compose up
```

Then open [http://localhost:8000](http://localhost:8000) and sign in with `built-in/admin` / `123` on a fresh install.

### ☸️ Kubernetes Helm

With Helm v3 and a running Kubernetes cluster:

```bash
helm install casdoor oci://registry-1.docker.io/casbin/casdoor-helm-charts
```

After installation, access Pipeline Auth through your cluster service/ingress.

---

<a id="features"></a>
## ✨ Features

<table>
<tr>
<td width="50%">

### 🔐 Authentication

- **OAuth 2.0 / OIDC** — OpenID Connect and OAuth 2.x authorization
- **SAML 2.0** — Enterprise SSO integration
- **CAS** — Central Authentication Service
- **LDAP** — Directory service integration
- **WebAuthn / Passkeys** — Passwordless authentication
- **TOTP / MFA** — Multi-factor authentication
- **Face ID** — Biometric authentication

</td>
<td width="50%">

### 🏢 Enterprise

- **SCIM 2.0** — User provisioning
- **RBAC** — Role-based access control
- **Social Login** — Google, GitHub, Azure AD, and more
- **Custom providers** — Extensible identity providers
- **User management** — Web UI for administration
- **Audit logs** — Comprehensive logging
- **Multi-tenancy** — Organization support

</td>
</tr>
<tr>
<td width="50%">

### 🤖 AI & MCP

- **MCP Gateway** — Model Context Protocol support
- **A2A Protocol** — Agent-to-Agent communication
- **AI-First Design** — Built for AI applications

</td>
<td width="50%">

### 🛠️ Developer Experience

- **RESTful API** — Complete API coverage
- **SDKs** — Go, Java, Python, Node.js, and more
- **Swagger UI** — Interactive API documentation
- **Webhooks** — Event-driven integrations
- **Customizable UI** — Brand theming support

</td>
</tr>
</table>

---

<a id="technology-stack"></a>
## Technology stack

Casdoor is built as a **frontend–backend separated** project:

- **Web UI**: JavaScript and **React** ([`web/`](https://github.com/liquiid727/pipeline-auth/tree/master/web))
- **API server**: **Go** with **Beego**, RESTful APIs ([repository root](https://github.com/liquiid727/pipeline-auth))
- **Data**: mainstream databases including **MySQL**, **PostgreSQL**, and others ([overview](https://casdoor.ai/docs/overview))
- **Cache**: optional **Redis** for session/cache-style deployments (configure as needed)

---

<a id="documentation"></a>
## 📖 Documentation

**All product documentation, installation, and tutorials live at [casdoor.ai/docs/overview](https://casdoor.ai/docs/overview).** Start here, then use the sections below.

**Install**

- [Install from source](https://casdoor.ai/docs/basic/server-installation)
- [Install with Docker](https://casdoor.ai/docs/basic/try-with-docker)
- [Install with Kubernetes Helm](https://casdoor.ai/docs/basic/try-with-helm)

**Connect applications**

- [How to connect to Casdoor](https://casdoor.ai/docs/how-to-connect/overview)

**APIs**

- [Public API](https://casdoor.ai/docs/basic/public-api)
- [Swagger UI](https://door.casdoor.com/swagger) (live API explorer)

---

<a id="integrations"></a>
## 🔌 Integrations

Casdoor integrates with common languages and frameworks:

<p align="center">
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg" width="40" alt="Go">
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/java/java-original.svg" width="40" alt="Java">
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/python/python-original.svg" width="40" alt="Python">
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg" width="40" alt="Node.js">
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/react/react-original.svg" width="40" alt="React">
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/vuejs/vuejs-original.svg" width="40" alt="Vue">
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/angularjs/angularjs-original.svg" width="40" alt="Angular">
</p>

Browse the full list: [Integrations](https://casdoor.ai/docs/category/integrations).

---

<a id="community-and-support"></a>
## 🤝 Community and support

- **Discord**: [Join our community](https://discord.gg/5rPsrAzK7S)
- **Contact**: [casdoor.ai/help](https://casdoor.ai/help)
- **Issues**: [GitHub Issues](https://github.com/liquiid727/pipeline-auth/issues)
- **Discussions**: [GitHub Discussions](https://github.com/liquiid727/pipeline-auth/discussions)

---

<a id="contributing"></a>
## 🌍 Contributing

If you have questions about Casdoor, you can **[open an issue](https://github.com/liquiid727/pipeline-auth/issues)**. Pull requests are welcome; **we recommend opening an issue first** so you can align with maintainers and the community before larger changes.

Please also read our [contribution guidelines](https://casdoor.ai/docs/contributing/) before contributing.

### Translation and i18n

- **Crowdin** is used for translation workflows: [casdoor-site on Crowdin](https://crowdin.com/project/casdoor-site).
- The web app uses **i18next**. When you add or change user-visible strings under [`web/`](https://github.com/liquiid727/pipeline-auth/tree/master/web), update the English catalog at [`web/src/locales/en/data.json`](web/src/locales/en/data.json) accordingly.

---

<a id="license"></a>
## 📄 License

Casdoor is licensed under the [Apache License 2.0](https://github.com/liquiid727/pipeline-auth/blob/master/LICENSE).

---

<div align="center">

[![Made with ❤️](https://img.shields.io/badge/Made_with-%E2%9D%A4%EF%B8%8F-ff6b6b?style=flat-square&logoColor=white)](https://casdoor.ai) [![By Casdoor](https://img.shields.io/badge/by-Casdoor-4ecdc4?style=flat-square)](https://casdoor.ai)

<a href="https://github.com/liquiid727/pipeline-auth/stargazers"><img src="https://img.shields.io/github/stars/casdoor/casdoor?style=social&logo=github&label=Star" alt="GitHub Stars"></a>

<sub>© 2026 <a href="https://casdoor.ai">Casdoor</a>. Licensed under <a href="https://github.com/liquiid727/pipeline-auth/blob/master/LICENSE">Apache License 2.0</a>.</sub>

</div>
