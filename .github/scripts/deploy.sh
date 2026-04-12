#!/bin/bash
set -e
ENVIRONMENT=$1
if [ "$ENVIRONMENT" = "prod" ]; then
    DEPLOY_DIR="/opt/antipratik"
else
    DEPLOY_DIR="/opt/antipratik-${ENVIRONMENT}"
fi
mkdir -p "$DEPLOY_DIR"
cp "/tmp/docker-compose.${ENVIRONMENT}.yml" "${DEPLOY_DIR}/docker-compose.yml"
echo "$GHCR_TOKEN" | docker login ghcr.io -u "$GHCR_USER" --password-stdin
docker compose -f "${DEPLOY_DIR}/docker-compose.yml" pull
docker compose -f "${DEPLOY_DIR}/docker-compose.yml" up -d --remove-orphans
docker image prune -f
echo "✓ ${ENVIRONMENT} deploy complete"