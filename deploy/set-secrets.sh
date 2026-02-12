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

echo "==> Setting core secrets on Container App: $ACA_APP"

az containerapp secret set \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_APP" \
  --secrets \
    jwt-secret="$JWT_SECRET" \
    google-client-id="$GOOGLE_CLIENT_ID" \
    google-client-secret="$GOOGLE_CLIENT_SECRET" \
  -o none

echo "==> Building environment variables..."

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

# Helper: conditionally add an OAuth platform's secrets
add_platform_secrets() {
  local name="$1"
  local id_var="$2"
  local secret_var="$3"

  local id_val="${!id_var:-}"
  local secret_val="${!secret_var:-}"

  if [ -n "$id_val" ] && [ -n "$secret_val" ]; then
    local lower_name
    lower_name=$(echo "$name" | tr '[:upper:]' '[:lower:]')
    echo "  Adding $name OAuth secrets..."
    az containerapp secret set \
      --resource-group "$RESOURCE_GROUP" \
      --name "$ACA_APP" \
      --secrets \
        "${lower_name}-client-id=${id_val}" \
        "${lower_name}-client-secret=${secret_val}" \
      -o none
    ENV_VARS+=("${id_var}=secretref:${lower_name}-client-id")
    ENV_VARS+=("${secret_var}=secretref:${lower_name}-client-secret")
  fi
}

add_platform_secrets "GitHub" "GITHUB_CLIENT_ID" "GITHUB_CLIENT_SECRET"
add_platform_secrets "Twitter" "TWITTER_CLIENT_ID" "TWITTER_CLIENT_SECRET"
add_platform_secrets "LinkedIn" "LINKEDIN_CLIENT_ID" "LINKEDIN_CLIENT_SECRET"
add_platform_secrets "Instagram" "INSTAGRAM_CLIENT_ID" "INSTAGRAM_CLIENT_SECRET"
add_platform_secrets "Dribbble" "DRIBBBLE_CLIENT_ID" "DRIBBBLE_CLIENT_SECRET"
add_platform_secrets "Behance" "BEHANCE_CLIENT_ID" "BEHANCE_CLIENT_SECRET"

# Azure Storage
if [ -n "${AZURE_STORAGE_ACCOUNT:-}" ] && [ -n "${AZURE_STORAGE_KEY:-}" ]; then
  echo "  Adding Azure Storage secrets..."
  az containerapp secret set \
    --resource-group "$RESOURCE_GROUP" \
    --name "$ACA_APP" \
    --secrets \
      azure-storage-key="$AZURE_STORAGE_KEY" \
    -o none
  ENV_VARS+=("AZURE_STORAGE_ACCOUNT=$AZURE_STORAGE_ACCOUNT")
  ENV_VARS+=("AZURE_STORAGE_KEY=secretref:azure-storage-key")
  ENV_VARS+=("AZURE_STORAGE_CONTAINER=${AZURE_STORAGE_CONTAINER:-avatars}")
fi

# SMTP
if [ -n "${SMTP_HOST:-}" ] && [ -n "${SMTP_USERNAME:-}" ] && [ -n "${SMTP_PASSWORD:-}" ]; then
  echo "  Adding SMTP secrets..."
  az containerapp secret set \
    --resource-group "$RESOURCE_GROUP" \
    --name "$ACA_APP" \
    --secrets \
      smtp-password="$SMTP_PASSWORD" \
    -o none
  ENV_VARS+=("SMTP_HOST=$SMTP_HOST")
  ENV_VARS+=("SMTP_PORT=${SMTP_PORT:-587}")
  ENV_VARS+=("SMTP_USERNAME=$SMTP_USERNAME")
  ENV_VARS+=("SMTP_PASSWORD=secretref:smtp-password")
  ENV_VARS+=("SMTP_FROM=${SMTP_FROM:-noreply@creatrid.com}")
fi

# Sentry
if [ -n "${SENTRY_DSN:-}" ]; then
  echo "  Adding Sentry DSN..."
  ENV_VARS+=("SENTRY_DSN=$SENTRY_DSN")
fi

# Refresh interval
if [ -n "${REFRESH_INTERVAL:-}" ]; then
  ENV_VARS+=("REFRESH_INTERVAL=$REFRESH_INTERVAL")
fi

echo "==> Updating environment variables..."

az containerapp update \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_APP" \
  --set-env-vars "${ENV_VARS[@]}" \
  -o none

echo "==> Done! Secrets and env vars updated."
