#!/bin/bash
set -e
ENVIRONMENT=$1
DEPLOY_DIR="/opt/antipratik-${ENVIRONMENT}"
mkdir -p "$DEPLOY_DIR"
cp "/tmp/docker-compose.${ENVIRONMENT}.yml" "${DEPLOY_DIR}/docker-compose.yml"
cp "/tmp/config.${ENVIRONMENT}.yaml" "${DEPLOY_DIR}/config.yaml"
echo "$GHCR_TOKEN" | docker login ghcr.io -u "$GHCR_USER" --password-stdin
docker compose -f "${DEPLOY_DIR}/docker-compose.yml" pull
docker compose -f "${DEPLOY_DIR}/docker-compose.yml" up -d --remove-orphans
docker image prune -f
echo "✓ ${ENVIRONMENT} deploy complete"