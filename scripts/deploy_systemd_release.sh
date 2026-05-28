#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ARCHIVE_PATH="${SYSTEMD_RELEASE_ARCHIVE:-${SYSTEMD_RELEASE_DIR:-$ROOT_DIR/.release/casdoor-systemd}.tar.gz}"
REMOTE_HOST="${REMOTE_HOST:-realdesk-dev}"
REMOTE_DIR="${REMOTE_DIR:-/opt/casdoor}"
SERVICE_NAME="${SERVICE_NAME:-casdoor}"
UNIT_TEMPLATE="${ROOT_DIR}/deploy/systemd/casdoor.service"

if [[ ! -f "${ARCHIVE_PATH}" ]]; then
	echo "Release archive not found at ${ARCHIVE_PATH}, building it first."
	"${ROOT_DIR}/scripts/package_systemd_release.sh"
fi

echo "==> Ensuring remote install directory exists: ${REMOTE_HOST}:${REMOTE_DIR}"
ssh "${REMOTE_HOST}" "sudo mkdir -p '${REMOTE_DIR}'"

echo "==> Uploading release bundle"
cat "${ARCHIVE_PATH}" | ssh "${REMOTE_HOST}" "sudo tar -xzf - -C '${REMOTE_DIR}'"

echo "==> Bootstrapping PostgreSQL and runtime env"
ssh "${REMOTE_HOST}" "
set -e
sudo -u postgres psql -v ON_ERROR_STOP=1 <<'SQL'
DO \$\$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'casdoor') THEN
    CREATE ROLE casdoor LOGIN PASSWORD 'casdoor';
  ELSE
    ALTER ROLE casdoor WITH LOGIN PASSWORD 'casdoor';
  END IF;
END
\$\$;
SELECT 'CREATE DATABASE casdoor OWNER casdoor'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'casdoor')\gexec
SQL

if [ ! -f '${REMOTE_DIR}/.env' ]; then
  tmp_env=\$(mktemp)
  cat > \"\${tmp_env}\" <<'ENV'
runmode=prod
httpport=8000
driverName=postgres
dataSourceName=user=casdoor password=casdoor host=127.0.0.1 port=5432 sslmode=disable dbname=casdoor
dbName=casdoor
ENV
  sudo install -m 600 \"\${tmp_env}\" '${REMOTE_DIR}/.env'
  rm -f \"\${tmp_env}\"
fi
"

echo "==> Installing systemd unit"
sed "s|__INSTALL_DIR__|${REMOTE_DIR}|g" "${UNIT_TEMPLATE}" | ssh "${REMOTE_HOST}" "sudo tee '/etc/systemd/system/${SERVICE_NAME}.service' >/dev/null"

echo "==> Reloading and restarting ${SERVICE_NAME}.service"
ssh "${REMOTE_HOST}" "sudo systemctl daemon-reload && sudo systemctl enable '${SERVICE_NAME}.service' && sudo systemctl restart '${SERVICE_NAME}.service' && sudo systemctl --no-pager --full status '${SERVICE_NAME}.service'"
