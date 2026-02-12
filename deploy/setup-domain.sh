#!/usr/bin/env bash
set -euo pipefail

# Sets up custom domains for Creatrid
# Prerequisites:
#   1. You own the domain and can add DNS records
#   2. Azure CLI is authenticated (az login)
#   3. Resources are already deployed (deploy.sh was run)
#
# DNS records to add BEFORE running this script:
#   creatrid.com       → CNAME → lemon-plant-0ae85c203.2.azurestaticapps.net
#   api.creatrid.com   → CNAME → creatrid-api.bluewave-351c9003.westeurope.azurecontainerapps.io
#
# Usage: ./deploy/setup-domain.sh

RESOURCE_GROUP="creatrid-rg"
ACA_APP="creatrid-api"
ACA_ENV="creatrid-env"
SWA_NAME="creatrid-web"

FRONTEND_DOMAIN="creatrid.com"
API_DOMAIN="api.creatrid.com"

echo "============================================"
echo "  Creatrid Custom Domain Setup"
echo "============================================"
echo ""
echo "Make sure these DNS records exist:"
echo "  $FRONTEND_DOMAIN  → CNAME → lemon-plant-0ae85c203.2.azurestaticapps.net"
echo "  $API_DOMAIN        → CNAME → creatrid-api.bluewave-351c9003.westeurope.azurecontainerapps.io"
echo ""
read -p "DNS records configured? (y/N) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
  echo "Please configure DNS records first."
  exit 1
fi

# --- Frontend: Static Web Apps custom domain ---
echo ""
echo "==> Adding custom domain to Static Web App: $FRONTEND_DOMAIN"
az staticwebapp hostname set \
  --resource-group "$RESOURCE_GROUP" \
  --name "$SWA_NAME" \
  --hostname "$FRONTEND_DOMAIN" \
  -o none

echo "  Managed SSL certificate will be provisioned automatically."

# --- Backend: Container Apps custom domain ---
echo ""
echo "==> Adding custom domain to Container App: $API_DOMAIN"

# Create a managed certificate
echo "  Creating managed certificate..."
az containerapp env certificate create \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_ENV" \
  --hostname "$API_DOMAIN" \
  --validation-method CNAME \
  -o none 2>/dev/null || echo "  (Certificate may already exist)"

# Get certificate ID
CERT_ID=$(az containerapp env certificate list \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_ENV" \
  --query "[?properties.subjectName=='$API_DOMAIN'].id" \
  -o tsv 2>/dev/null || true)

if [ -z "$CERT_ID" ]; then
  echo "  Warning: Could not find certificate for $API_DOMAIN."
  echo "  You may need to wait for DNS propagation and retry."
  echo "  Certificate provisioning can take a few minutes."
else
  echo "  Binding domain with certificate..."
  az containerapp hostname bind \
    --resource-group "$RESOURCE_GROUP" \
    --name "$ACA_APP" \
    --hostname "$API_DOMAIN" \
    --certificate "$CERT_ID" \
    -o none
fi

# --- Update secrets/env vars for the new domains ---
echo ""
echo "==> Updating backend URLs to use custom domains..."
echo "  Run: ./deploy/set-secrets.sh with updated .env.production.local"
echo "  Set these values:"
echo "    BACKEND_URL=https://$API_DOMAIN"
echo "    FRONTEND_URL=https://$FRONTEND_DOMAIN"
echo "    COOKIE_DOMAIN=$FRONTEND_DOMAIN"
echo ""
echo "  Also update GitHub Actions vars:"
echo "    NEXT_PUBLIC_API_URL=https://$API_DOMAIN"
echo "    NEXT_PUBLIC_PROFILE_URL=https://$FRONTEND_DOMAIN"

echo ""
echo "==> Done! Custom domains configured."
echo "  SSL certificates are managed automatically by Azure."
echo "  It may take a few minutes for certificates to be issued."
