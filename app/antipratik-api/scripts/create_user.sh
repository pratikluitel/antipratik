#!/usr/bin/env bash
set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
go run "$SCRIPT_DIR/create_user/main.go" "$@"
