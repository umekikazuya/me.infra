#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

cd "${ROOT_DIR}"
sh ./scripts/package-placeholder-api.sh
go mod download
go run ./cmd/app
