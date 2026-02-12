#!/usr/bin/env bash
set -euo pipefail

# ─── Configuration ─────────────────────────────────────────────
RESOURCE_GROUP="creatrid-rg"
LOCATION="westeurope"
PG_SERVER="creatrid-pg"
PG_ADMIN_USER="creatridadmin"
PG_DB_NAME="creatrid"
ACR_NAME="creatridacr"
ACA_ENV="creatrid-env"
ACA_APP="creatrid-api"
SWA_NAME="creatrid-web"

# Prompt for secrets that should not be in scripts
if [ -z "${PG_ADMIN_PASSWORD:-}" ]; then
  echo -n "Enter PostgreSQL admin password: "
  read -rs PG_ADMIN_PASSWORD
  echo
fi

echo "==> Using subscription: $(az account show --query name -o tsv)"
echo "==> Resource group: $RESOURCE_GROUP"
echo "==> Location: $LOCATION"
echo ""

# ─── Resource Group ────────────────────────────────────────────
echo "==> Creating resource group..."
az group create \
  --name "$RESOURCE_GROUP" \
  --location "$LOCATION" \
  -o none

# ─── Azure Database for PostgreSQL Flexible Server ─────────────
echo "==> Creating PostgreSQL Flexible Server..."
az postgres flexible-server create \
  --resource-group "$RESOURCE_GROUP" \
  --name "$PG_SERVER" \
  --location "$LOCATION" \
  --admin-user "$PG_ADMIN_USER" \
  --admin-password "$PG_ADMIN_PASSWORD" \
  --sku-name Standard_B1ms \
  --tier Burstable \
  --storage-size 32 \
  --version 16 \
  --public-access 0.0.0.0 \
  --yes \
  -o none

echo "==> Creating database..."
az postgres flexible-server db create \
  --resource-group "$RESOURCE_GROUP" \
  --server-name "$PG_SERVER" \
  --database-name "$PG_DB_NAME" \
  -o none

# Get the FQDN for the connection string
PG_FQDN=$(az postgres flexible-server show \
  --resource-group "$RESOURCE_GROUP" \
  --name "$PG_SERVER" \
  --query fullyQualifiedDomainName -o tsv)

DATABASE_URL="postgres://${PG_ADMIN_USER}:${PG_ADMIN_PASSWORD}@${PG_FQDN}:5432/${PG_DB_NAME}?sslmode=require"
echo "==> PostgreSQL ready: $PG_FQDN"

# ─── Azure Container Registry ─────────────────────────────────
echo "==> Creating Container Registry..."
az acr create \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACR_NAME" \
  --sku Basic \
  --admin-enabled true \
  -o none

ACR_SERVER=$(az acr show \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACR_NAME" \
  --query loginServer -o tsv)

ACR_PASSWORD=$(az acr credential show \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACR_NAME" \
  --query "passwords[0].value" -o tsv)

echo "==> ACR ready: $ACR_SERVER"

# ─── Build & Push Backend Image ────────────────────────────────
echo "==> Building and pushing backend image..."
az acr build \
  --registry "$ACR_NAME" \
  --image creatrid-api:latest \
  --file backend/Dockerfile \
  backend/

# ─── Container Apps Environment ────────────────────────────────
echo "==> Creating Container Apps environment..."
az containerapp env create \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_ENV" \
  --location "$LOCATION" \
  -o none

# ─── Container App (Backend API) ──────────────────────────────
echo "==> Creating Container App..."
az containerapp create \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_APP" \
  --environment "$ACA_ENV" \
  --image "${ACR_SERVER}/creatrid-api:latest" \
  --registry-server "$ACR_SERVER" \
  --registry-username "$ACR_NAME" \
  --registry-password "$ACR_PASSWORD" \
  --target-port 8080 \
  --ingress external \
  --min-replicas 0 \
  --max-replicas 2 \
  --cpu 0.25 \
  --memory 0.5Gi \
  --env-vars \
    PORT=8080 \
    DATABASE_URL="$DATABASE_URL" \
    JWT_SECRET=secretref:jwt-secret \
    GOOGLE_CLIENT_ID=secretref:google-client-id \
    GOOGLE_CLIENT_SECRET=secretref:google-client-secret \
    COOKIE_SECURE=true \
  -o none

# Get the app URL
ACA_URL=$(az containerapp show \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_APP" \
  --query "properties.configuration.ingress.fqdn" -o tsv)

BACKEND_URL="https://${ACA_URL}"
echo "==> Container App ready: $BACKEND_URL"

# ─── Azure Static Web Apps ─────────────────────────────────────
echo "==> Creating Static Web App..."
az staticwebapp create \
  --resource-group "$RESOURCE_GROUP" \
  --name "$SWA_NAME" \
  --location "$LOCATION" \
  --sku Free \
  -o none

SWA_URL=$(az staticwebapp show \
  --resource-group "$RESOURCE_GROUP" \
  --name "$SWA_NAME" \
  --query defaultHostname -o tsv)

SWA_TOKEN=$(az staticwebapp secrets list \
  --resource-group "$RESOURCE_GROUP" \
  --name "$SWA_NAME" \
  --query "properties.apiKey" -o tsv)

echo "==> Static Web App ready: https://$SWA_URL"

# ─── Update Container App with final URLs ──────────────────────
echo "==> Updating Container App environment variables..."
FRONTEND_URL="https://${SWA_URL}"

az containerapp update \
  --resource-group "$RESOURCE_GROUP" \
  --name "$ACA_APP" \
  --set-env-vars \
    BACKEND_URL="$BACKEND_URL" \
    FRONTEND_URL="$FRONTEND_URL" \
    COOKIE_DOMAIN="" \
  -o none

# ─── Summary ───────────────────────────────────────────────────
echo ""
echo "============================================"
echo "  Deployment Complete!"
echo "============================================"
echo ""
echo "Resources created in: $RESOURCE_GROUP"
echo ""
echo "PostgreSQL:     $PG_FQDN"
echo "Backend API:    $BACKEND_URL"
echo "Frontend:       https://$SWA_URL"
echo "Container Reg:  $ACR_SERVER"
echo ""
echo "─── Next Steps ───"
echo ""
echo "1. Set secrets on the Container App:"
echo "   az containerapp secret set \\"
echo "     --resource-group $RESOURCE_GROUP \\"
echo "     --name $ACA_APP \\"
echo "     --secrets jwt-secret=<YOUR_JWT_SECRET> \\"
echo "              google-client-id=<YOUR_GOOGLE_CLIENT_ID> \\"
echo "              google-client-secret=<YOUR_GOOGLE_CLIENT_SECRET>"
echo ""
echo "2. Update Google Cloud Console OAuth redirect URIs:"
echo "   - ${BACKEND_URL}/api/auth/google/callback"
echo "   - ${BACKEND_URL}/api/connections/youtube/callback"
echo ""
echo "3. Deploy frontend (first time):"
echo "   cd frontend && pnpm build"
echo "   npx @azure/static-web-apps-cli deploy out \\"
echo "     --deployment-token $SWA_TOKEN \\"
echo "     --env production"
echo ""
echo "4. Set frontend environment variables before building:"
echo "   NEXT_PUBLIC_API_URL=$BACKEND_URL"
echo "   NEXT_PUBLIC_PROFILE_URL=https://$SWA_URL"
echo ""
echo "5. (Optional) Add custom domain:"
echo "   az staticwebapp hostname set --name $SWA_NAME --hostname creatrid.com"
echo "   az containerapp hostname add --name $ACA_APP --hostname api.creatrid.com"
echo ""
echo "SWA deployment token (save this):"
echo "  $SWA_TOKEN"
echo ""
