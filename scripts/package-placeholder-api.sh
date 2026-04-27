#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
ARTIFACT_DIR="${ROOT_DIR}/.artifacts/api"
BUILD_DIR="${ARTIFACT_DIR}/build"
BIN_PATH="${BUILD_DIR}/bootstrap"
ZIP_PATH="${ARTIFACT_DIR}/bootstrap.zip"

rm -rf "${BUILD_DIR}" "${ZIP_PATH}"
mkdir -p "${BUILD_DIR}"

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
  go build -tags lambda.norpc -o "${BIN_PATH}" "${ROOT_DIR}/cmd/placeholder-api"

(
  cd "${BUILD_DIR}"
  zip -q "${ZIP_PATH}" bootstrap
)
