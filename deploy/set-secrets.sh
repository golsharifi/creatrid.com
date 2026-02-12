#!/usr/bin/env bash
set -euo pipefail

# Sets Container App secrets from a .env.production.local file
# Usage: ./deploy/set-secrets.sh [path-to-env-file]

RESOURCE_GROUP="creatrid-rg"
ACA_APP="creatrid-api"
ENV_FILE="${1:-deploy/.env.production.local}"

if [ ! -f "$ENV_FILE" ]; then
  echo "Error: $ENV_FILE not found."
  echo "Copy deploy/.env.production to deploy/.env.production.local and fill in values."
  exit 1
fi

# Read values from env file
source "$ENV_FILE"

echo "==> Setting secrets on Container App: $ACA_APP"

az containerapp secret set \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_APP" \
  --secrets \
    jwt-secret="$JWT_SECRET" \
    google-client-id="$GOOGLE_CLIENT_ID" \
    google-client-secret="$GOOGLE_CLIENT_SECRET" \
  -o none

echo "==> Updating environment variables..."

ENV_VARS=(
  "PORT=8080"
  "DATABASE_URL=$DATABASE_URL"
  "JWT_SECRET=secretref:jwt-secret"
  "GOOGLE_CLIENT_ID=secretref:google-client-id"
  "GOOGLE_CLIENT_SECRET=secretref:google-client-secret"
  "BACKEND_URL=$BACKEND_URL"
  "FRONTEND_URL=$FRONTEND_URL"
  "COOKIE_DOMAIN=$COOKIE_DOMAIN"
  "COOKIE_SECURE=true"
)

if [ -n "${GITHUB_CLIENT_ID:-}" ] && [ -n "${GITHUB_CLIENT_SECRET:-}" ]; then
  az containerapp secret set \
    --resource-group "$RESOURCE_GROUP" \
    --name "$ACA_APP" \
    --secrets \
      github-client-id="$GITHUB_CLIENT_ID" \
      github-client-secret="$GITHUB_CLIENT_SECRET" \
    -o none
  ENV_VARS+=("GITHUB_CLIENT_ID=secretref:github-client-id")
  ENV_VARS+=("GITHUB_CLIENT_SECRET=secretref:github-client-secret")
fi

az containerapp update \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_APP" \
  --set-env-vars "${ENV_VARS[@]}" \
  -o none

echo "==> Done! Secrets and env vars updated."
