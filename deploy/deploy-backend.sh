#!/usr/bin/env bash
set -euo pipefail

# Build and deploy backend to Azure Container Apps
# Usage: ./deploy/deploy-backend.sh

RESOURCE_GROUP="creatrid-rg"
ACR_NAME="creatridacr"
ACA_APP="creatrid-api"

echo "==> Building and pushing backend image to ACR..."
az acr build \
  --registry "$ACR_NAME" \
  --image creatrid-api:latest \
  --file backend/Dockerfile \
  backend/

echo "==> Updating Container App..."
ACR_SERVER=$(az acr show \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACR_NAME" \
  --query loginServer -o tsv)

az containerapp update \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_APP" \
  --image "${ACR_SERVER}/creatrid-api:latest" \
  -o none

echo "==> Backend deployed!"
ACA_URL=$(az containerapp show \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_APP" \
  --query "properties.configuration.ingress.fqdn" -o tsv)
echo "    URL: https://$ACA_URL"
