# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Casdoor is an open-source Identity and Access Management (IAM) platform and MCP (Model Context Protocol) Gateway. It provides authentication, authorization, and user management with a web UI. Backend is Go (Beego v2 framework, Casbin for RBAC, xorm ORM). Frontend is React 18 with Ant Design 6.

## Common Commands

### Backend (Go)
```bash
make backend          # Build to bin/manager (runs fmt + vet first)
make run              # Run with go run ./main.go
make ut               # Run Go tests with coverage
make lint             # Run golangci-lint (requires: make lint-install)
make fmt              # go fmt ./...
make vet              # go vet ./...
```

### Frontend (React)
```bash
cd web && yarn install && yarn build   # Build frontend
make frontend                          # Same as above
cd web && yarn start                   # Dev server on port 7001
cd web && yarn fix                     # ESLint fix
cd web && yarn lint:css                # Stylelint fix
```

### Docker
```bash
docker compose up     # Start Casdoor + PostgreSQL
```

### Deploy
```bash
make deploy           # Helm upgrade/install to Kubernetes
make dry-run          # Helm dry-run preview
make undeploy         # Helm delete
make package-systemd  # Build a linux/amd64 systemd release bundle
make casdoor-test     # Build and deploy to ssh realdesk-dev:/opt/casdoor
```

## Architecture

### Request Flow
```
HTTP Request → routers/ (Beego router + middleware/filters) → controllers/ (handlers) → object/ (business logic + data access) → Database
```

### Key Backend Directories
- **`controllers/`** — HTTP handlers (one file per domain: users, applications, organizations, auth, LDAP, MCP, etc.)
- **`object/`** — Core business logic and ORM models (~160 files). This is the data access layer using xorm.
- **`routers/`** — Beego router setup, middleware (CORS, auth checks, Prometheus metrics, timeout, static serving)
- **`authz/`** — Casbin authorization policy enforcement
- **`conf/`** — Configuration loading from `app.conf` (database, ports, LDAP/RADIUS settings, session config)
- **`idp/`** — Identity provider implementations (Google, GitHub, Azure AD, WeChat, Telegram, etc.)
- **`notification/`** — Notification providers (Slack, Discord, Telegram, email, web push)
- **`storage/`** — File storage providers (S3, Azure, MinIO, local filesystem)
- **`pp/`** — Payment providers (Stripe, PayPal, Alipay, WeChat Pay)
- **`mcp/`** — MCP server utilities
- **`ldap/`** — Built-in LDAP server
- **`scim/`** — SCIM 2.0 provisioning endpoint
- **`cred/`** — Credential/password hashing strategies

### Frontend Structure (`web/src/`)
- Pages organized by entity type (users, applications, providers, roles, permissions, etc.)
- Ant Design 6 component library
- React Router 5 for client-side routing
- i18next for internationalization
- CRACO overrides Create React App config

### Database
Default is PostgreSQL (configured in `conf/app.conf`). Also supports MySQL, SQLite, SQL Server. CI tests run against MySQL 5.7.

### Server Ports
- HTTP: 8000 (default)
- LDAP: 389 / LDAPS: 636
- RADIUS: 1812
- Frontend dev server: 7001

## Linting Rules

- **Go**: golangci-lint with gofumpt formatter (see `.golangci.yml`)
- **JS/CSS**: ESLint + Stylelint with Husky pre-commit hooks via lint-staged

## Conventions

- Go code uses Beego controller patterns — new endpoints follow the `controllers/` file naming convention
- ORM models live in `object/` with corresponding database operations
- Provider integrations (identity, notification, storage, payment) follow a common interface pattern in their respective directories
- Commit messages follow conventional commits (used by semantic-release for automated versioning)

## Deployment Agent Knowledge

- The canonical SSH/systemd target is `realdesk-dev`.
- Install path is `/opt/casdoor`, and the service unit is `/etc/systemd/system/casdoor.service`.
- Reusable release commands are `make package-systemd` and `make casdoor-test`.
- The release bundle must include the backend binary plus frontend runtime assets (`web/build` and `web/public`), otherwise Casdoor cannot serve the UI correctly.
- Runtime overrides live in `/opt/casdoor/.env`; keep secrets there rather than committing them.
- If `/opt/casdoor/.env` is absent, Casdoor falls back to the bundled `conf/app.conf`.
