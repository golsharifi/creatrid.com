#!/usr/bin/env bash
set -euo pipefail

# Deploy frontend static export to Azure Static Web Apps
# Usage: ./deploy/deploy-frontend.sh
#
# Required env vars:
#   SWA_DEPLOYMENT_TOKEN  — from Azure portal or deploy.sh output
#   NEXT_PUBLIC_API_URL   — backend URL (e.g., https://creatrid-api.*.azurecontainerapps.io)
#   NEXT_PUBLIC_PROFILE_URL — frontend URL (e.g., https://creatrid-web.azurestaticapps.net)

RESOURCE_GROUP="creatrid-rg"
SWA_NAME="creatrid-web"

if [ -z "${SWA_DEPLOYMENT_TOKEN:-}" ]; then
  echo "==> Fetching SWA deployment token..."
  SWA_DEPLOYMENT_TOKEN=$(az staticwebapp secrets list \
    --resource-group "$RESOURCE_GROUP" \
    --name "$SWA_NAME" \
    --query "properties.apiKey" -o tsv)
fi

if [ -z "${NEXT_PUBLIC_API_URL:-}" ]; then
  echo "Error: NEXT_PUBLIC_API_URL is required"
  echo "Example: NEXT_PUBLIC_API_URL=https://creatrid-api.xxx.azurecontainerapps.io ./deploy/deploy-frontend.sh"
  exit 1
fi

echo "==> Building frontend..."
echo "    API URL: $NEXT_PUBLIC_API_URL"
echo "    Profile URL: ${NEXT_PUBLIC_PROFILE_URL:-$NEXT_PUBLIC_API_URL}"

cd frontend
pnpm build

echo "==> Copying staticwebapp.config.json to output..."
cp staticwebapp.config.json out/

echo "==> Deploying to Azure Static Web Apps..."
npx --yes @azure/static-web-apps-cli deploy out \
  --deployment-token "$SWA_DEPLOYMENT_TOKEN" \
  --env production

echo "==> Frontend deployed!"
echo "    URL: https://$(az staticwebapp show --resource-group $RESOURCE_GROUP --name $SWA_NAME --query defaultHostname -o tsv)"
