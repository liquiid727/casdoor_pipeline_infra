#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REMOTE_HOST="${REMOTE_HOST:-realdesk-dev}"
REMOTE_DIR="${REMOTE_DIR:-/opt/casdoor}"
SERVICE_NAME="${SERVICE_NAME:-casdoor}"
LOCAL_CACHE_DIR="${ROOT_DIR}/.cache"
YARN_CACHE_DIR="${YARN_CACHE_DIR:-${LOCAL_CACHE_DIR}/yarn}"
TMP_ARCHIVE="${ROOT_DIR}/.release/casdoor-test-src.tar.gz"
REMOTE_BUILD_DIR="/tmp/casdoor-build"

mkdir -p "${YARN_CACHE_DIR}" "${ROOT_DIR}/.release"

echo "==> Building frontend locally"
(
  cd "${ROOT_DIR}/web"
  CYPRESS_INSTALL_BINARY=0 YARN_CACHE_FOLDER="${YARN_CACHE_DIR}" yarn install --frozen-lockfile
  yarn build
)

echo "==> Packing source bundle for remote build"
tar --exclude='.git' \
  --exclude='.cache' \
  --exclude='.release' \
  --exclude='web/node_modules' \
  -czf "${TMP_ARCHIVE}" \
  -C "${ROOT_DIR}" .

echo "==> Uploading source bundle to ${REMOTE_HOST}"
cat "${TMP_ARCHIVE}" | ssh "${REMOTE_HOST}" "mkdir -p '${REMOTE_BUILD_DIR}' && tar -xzf - -C '${REMOTE_BUILD_DIR}'"

echo "==> Building backend on ${REMOTE_HOST}"
ssh "${REMOTE_HOST}" "
set -e
cd '${REMOTE_BUILD_DIR}'
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o casdoor .
sudo mkdir -p '${REMOTE_DIR}'
sudo mkdir -p '${REMOTE_DIR}/conf' '${REMOTE_DIR}/web/build' '${REMOTE_DIR}/web/public' '${REMOTE_DIR}/deploy/systemd'
sudo install -m 755 casdoor '${REMOTE_DIR}/casdoor'
sudo rsync -a --delete conf/ '${REMOTE_DIR}/conf/'
sudo rsync -a --delete web/build/ '${REMOTE_DIR}/web/build/'
sudo rsync -a --delete web/public/ '${REMOTE_DIR}/web/public/'
sudo install -m 644 deploy/systemd/casdoor.service '${REMOTE_DIR}/deploy/systemd/casdoor.service'
sudo install -m 644 deploy/systemd/casdoor.env.example '${REMOTE_DIR}/deploy/systemd/casdoor.env.example'
"

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
sed "s|__INSTALL_DIR__|${REMOTE_DIR}|g" "${ROOT_DIR}/deploy/systemd/casdoor.service" | ssh "${REMOTE_HOST}" "sudo tee '/etc/systemd/system/${SERVICE_NAME}.service' >/dev/null"

echo "==> Restarting ${SERVICE_NAME}.service"
ssh "${REMOTE_HOST}" "sudo systemctl daemon-reload && sudo systemctl enable '${SERVICE_NAME}.service' && sudo systemctl restart '${SERVICE_NAME}.service' && sudo systemctl --no-pager --full status '${SERVICE_NAME}.service'"
