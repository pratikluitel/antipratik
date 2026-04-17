#!/usr/bin/env bash
set -euo pipefail

# Usage: ./generate.sh <dev|prod>

ENV="${1:-}"
if [[ "$ENV" != "dev" && "$ENV" != "prod" ]]; then
  echo "Usage: $0 <dev|prod>" >&2
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CERT="$SCRIPT_DIR/fullchain-${ENV}.pem"
KEY="$SCRIPT_DIR/privkey-${ENV}.pem"

openssl req -x509 -newkey rsa:4096 -sha256 -days 365 \
  -nodes \
  -keyout "$KEY" \
  -out "$CERT" \
  -subj "/CN=antipratik-${ENV}" \
  -extensions v3_ca \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

echo ""
echo "Generated:"
echo "  $CERT"
echo "  $KEY"
echo ""

# base64 is available on both macOS and Linux; flag differs
if base64 --version 2>/dev/null | grep -q "GNU"; then
  B64="base64 -w 0"   # GNU coreutils (Linux)
else
  B64="base64"        # macOS (no wrap flag needed, outputs single line)
fi

echo "=== fullchain-${ENV}.pem (base64) ==="
$B64 < "$CERT"
echo ""
echo "=== privkey-${ENV}.pem (base64) ==="
$B64 < "$KEY"
echo ""
