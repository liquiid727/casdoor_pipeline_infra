# realdesk-dev systemd deployment

This repo now supports a binary + systemd deploy flow for `realdesk-dev`.

## Layout

- Install dir: `/opt/casdoor`
- Systemd unit: `/etc/systemd/system/casdoor.service`
- Optional runtime overrides: `/opt/casdoor/.env`

## Make targets

- `make package-systemd`
- `make casdoor-test`

## What gets bundled

- `casdoor` Linux amd64 binary
- `web/build`
- `web/public`
- `conf/app.conf`
- `conf/waf.conf`
- `deploy/systemd/casdoor.service`
- `deploy/systemd/casdoor.env.example`

## Runtime notes

- The service uses `/opt/casdoor/.env` when present.
- If `.env` is absent, the bundled `conf/app.conf` defaults are used.
- The deploy script bootstraps a local PostgreSQL role/database pair: `casdoor` / `casdoor`.
- On first deploy, `/opt/casdoor/.env` is created automatically with `httpport=8000`.
- `sub` from Casdoor JWT remains the user identifier for upstream apps.
