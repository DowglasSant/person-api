#!/bin/bash

# ========================================
# Script de Deploy para Azure
# ========================================
# Este script automatiza o deploy da API no Azure Container Apps
# com PostgreSQL Flexible Server

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Deploy Pessoas API - Azure${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# ========================================
# 1. VARIÁVEIS DE CONFIGURAÇÃO
# ========================================
echo -e "${YELLOW}Configurando variáveis...${NC}"

RESOURCE_GROUP="pessoas-api-rg"
LOCATION="eastus"
ACR_NAME="pessoasapiacr"
POSTGRES_SERVER="pessoas-api-db"
POSTGRES_DB="postgres"
POSTGRES_USER="apiuser"
POSTGRES_PASSWORD="$(openssl rand -base64 24)"
JWT_SECRET="$(openssl rand -base64 32)"
CONTAINER_APP_NAME="pessoas-api"
CONTAINER_APP_ENV="pessoas-api-env"

echo -e "${GREEN}✓ Variáveis configuradas${NC}"
echo ""

# ========================================
# 2. LOGIN NO AZURE
# ========================================
echo -e "${YELLOW}Fazendo login no Azure...${NC}"
az login
echo -e "${GREEN}✓ Login realizado${NC}"
echo ""

# ========================================
# 3. CRIAR RESOURCE GROUP
# ========================================
echo -e "${YELLOW}Criando Resource Group...${NC}"
az group create \
  --name $RESOURCE_GROUP \
  --location $LOCATION
echo -e "${GREEN}✓ Resource Group criado: $RESOURCE_GROUP${NC}"
echo ""

# ========================================
# 4. CRIAR AZURE CONTAINER REGISTRY
# ========================================
echo -e "${YELLOW}Criando Azure Container Registry...${NC}"
az acr create \
  --resource-group $RESOURCE_GROUP \
  --name $ACR_NAME \
  --sku Basic \
  --admin-enabled true
echo -e "${GREEN}✓ ACR criado: $ACR_NAME${NC}"
echo ""

# ========================================
# 5. BUILD E PUSH DA IMAGEM DOCKER
# ========================================
echo -e "${YELLOW}Fazendo build e push da imagem Docker...${NC}"
az acr build \
  --registry $ACR_NAME \
  --image pessoas-api:latest \
  --file Dockerfile \
  .
echo -e "${GREEN}✓ Imagem Docker criada e enviada${NC}"
echo ""

# ========================================
# 6. CRIAR POSTGRESQL FLEXIBLE SERVER
# ========================================
echo -e "${YELLOW}Criando PostgreSQL Flexible Server...${NC}"
az postgres flexible-server create \
  --resource-group $RESOURCE_GROUP \
  --name $POSTGRES_SERVER \
  --location $LOCATION \
  --admin-user $POSTGRES_USER \
  --admin-password $POSTGRES_PASSWORD \
  --sku-name Standard_B1ms \
  --tier Burstable \
  --storage-size 32 \
  --version 16 \
  --public-access 0.0.0.0
echo -e "${GREEN}✓ PostgreSQL criado: $POSTGRES_SERVER${NC}"
echo ""

# ========================================
# 7. CRIAR DATABASE
# ========================================
echo -e "${YELLOW}Criando database...${NC}"
az postgres flexible-server db create \
  --resource-group $RESOURCE_GROUP \
  --server-name $POSTGRES_SERVER \
  --database-name $POSTGRES_DB
echo -e "${GREEN}✓ Database criado: $POSTGRES_DB${NC}"
echo ""

# ========================================
# 8. CRIAR CONTAINER APPS ENVIRONMENT
# ========================================
echo -e "${YELLOW}Criando Container Apps Environment...${NC}"
az containerapp env create \
  --name $CONTAINER_APP_ENV \
  --resource-group $RESOURCE_GROUP \
  --location $LOCATION
echo -e "${GREEN}✓ Environment criado: $CONTAINER_APP_ENV${NC}"
echo ""

# ========================================
# 9. OBTER CREDENCIAIS DO ACR
# ========================================
echo -e "${YELLOW}Obtendo credenciais do ACR...${NC}"
ACR_USERNAME=$(az acr credential show --name $ACR_NAME --query username -o tsv)
ACR_PASSWORD=$(az acr credential show --name $ACR_NAME --query passwords[0].value -o tsv)
ACR_LOGIN_SERVER=$(az acr show --name $ACR_NAME --query loginServer -o tsv)
echo -e "${GREEN}✓ Credenciais obtidas${NC}"
echo ""

# ========================================
# 10. CRIAR CONTAINER APP
# ========================================
echo -e "${YELLOW}Criando Container App...${NC}"

DB_HOST="${POSTGRES_SERVER}.postgres.database.azure.com"
DB_CONNECTION_STRING="host=${DB_HOST} port=5432 user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=require"

az containerapp create \
  --name $CONTAINER_APP_NAME \
  --resource-group $RESOURCE_GROUP \
  --environment $CONTAINER_APP_ENV \
  --image "${ACR_LOGIN_SERVER}/pessoas-api:latest" \
  --target-port 8080 \
  --ingress external \
  --registry-server $ACR_LOGIN_SERVER \
  --registry-username $ACR_USERNAME \
  --registry-password $ACR_PASSWORD \
  --cpu 0.5 \
  --memory 1.0Gi \
  --min-replicas 1 \
  --max-replicas 2 \
  --env-vars \
    DB_HOST=$DB_HOST \
    DB_PORT=5432 \
    DB_USER=$POSTGRES_USER \
    DB_PASSWORD=$POSTGRES_PASSWORD \
    DB_NAME=$POSTGRES_DB \
    DB_SCHEMA=people \
    DB_SSLMODE=require \
    JWT_SECRET=$JWT_SECRET \
    CORS_ALLOWED_ORIGINS="*"

echo -e "${GREEN}✓ Container App criado: $CONTAINER_APP_NAME${NC}"
echo ""

# ========================================
# 11. OBTER URL DA APLICAÇÃO
# ========================================
APP_URL=$(az containerapp show \
  --name $CONTAINER_APP_NAME \
  --resource-group $RESOURCE_GROUP \
  --query properties.configuration.ingress.fqdn \
  -o tsv)

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  DEPLOY CONCLUÍDO COM SUCESSO!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}Informações da Aplicação:${NC}"
echo -e "URL: ${GREEN}https://$APP_URL${NC}"
echo -e "Swagger: ${GREEN}https://$APP_URL/swagger/index.html${NC}"
echo ""
echo -e "${YELLOW}Credenciais do Banco de Dados:${NC}"
echo -e "Host: ${GREEN}$DB_HOST${NC}"
echo -e "Database: ${GREEN}$POSTGRES_DB${NC}"
echo -e "User: ${GREEN}$POSTGRES_USER${NC}"
echo -e "Password: ${GREEN}$POSTGRES_PASSWORD${NC}"
echo ""
echo -e "${YELLOW}JWT Secret:${NC}"
echo -e "${GREEN}$JWT_SECRET${NC}"
echo ""
echo -e "${RED}IMPORTANTE: Salve essas credenciais em um local seguro!${NC}"
echo ""
