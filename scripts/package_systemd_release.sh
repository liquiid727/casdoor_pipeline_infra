#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RELEASE_DIR="${SYSTEMD_RELEASE_DIR:-$ROOT_DIR/.release/casdoor-systemd}"
ARCHIVE_PATH="${SYSTEMD_RELEASE_ARCHIVE:-${RELEASE_DIR}.tar.gz}"
STAGING_DIR="${RELEASE_DIR}/staging"
GO_BIN="${GO_BIN:-$(go env GOROOT)/bin/go}"
GOCACHE_DIR="${GOCACHE_DIR:-$ROOT_DIR/.cache/go-build}"
GOMODCACHE_DIR="${GOMODCACHE_DIR:-$ROOT_DIR/.cache/go-mod}"
YARN_CACHE_DIR="${YARN_CACHE_DIR:-$ROOT_DIR/.cache/yarn}"
GOPROXY_URL="${GOPROXY_URL:-https://proxy.golang.org,direct}"

echo "==> Preparing release workspace: ${RELEASE_DIR}"
rm -rf "${RELEASE_DIR}"
mkdir -p "${STAGING_DIR}/conf" "${STAGING_DIR}/web" "${STAGING_DIR}/deploy/systemd" "${STAGING_DIR}/logs" "${STAGING_DIR}/tmp" "${STAGING_DIR}/files" "${GOCACHE_DIR}" "${GOMODCACHE_DIR}" "${YARN_CACHE_DIR}"

echo "==> Building backend binary"
(
	cd "${ROOT_DIR}"
	GOTOOLCHAIN=local GOPROXY="${GOPROXY_URL}" GOCACHE="${GOCACHE_DIR}" GOMODCACHE="${GOMODCACHE_DIR}" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 "${GO_BIN}" build -ldflags="-s -w" -o "${STAGING_DIR}/casdoor" .
)

echo "==> Building frontend assets"
(
	cd "${ROOT_DIR}/web"
	YARN_CACHE_FOLDER="${YARN_CACHE_DIR}" yarn install --frozen-lockfile
	yarn build
)

echo "==> Collecting runtime assets"
cp "${ROOT_DIR}/conf/app.conf" "${STAGING_DIR}/conf/app.conf"
cp "${ROOT_DIR}/conf/waf.conf" "${STAGING_DIR}/conf/waf.conf"
cp -R "${ROOT_DIR}/web/build" "${STAGING_DIR}/web/build"
cp -R "${ROOT_DIR}/web/public" "${STAGING_DIR}/web/public"
cp "${ROOT_DIR}/deploy/systemd/casdoor.service" "${STAGING_DIR}/deploy/systemd/casdoor.service"
cp "${ROOT_DIR}/deploy/systemd/casdoor.env.example" "${STAGING_DIR}/deploy/systemd/casdoor.env.example"

cat > "${STAGING_DIR}/VERSION.txt" <<EOF
commit=$(git -C "${ROOT_DIR}" rev-parse HEAD)
built_at=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
EOF

echo "==> Packaging release archive: ${ARCHIVE_PATH}"
mkdir -p "$(dirname "${ARCHIVE_PATH}")"
rm -f "${ARCHIVE_PATH}"
tar -C "${STAGING_DIR}" -czf "${ARCHIVE_PATH}" .

echo "Systemd release bundle created:"
echo "  staging: ${STAGING_DIR}"
echo "  archive: ${ARCHIVE_PATH}"
