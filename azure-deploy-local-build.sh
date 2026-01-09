#!/bin/bash

# ========================================
# Script de Deploy para Azure (Build Local)
# ========================================
# Este script faz build local da imagem e deploy no Azure

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Deploy Pessoas API - Azure${NC}"
echo -e "${GREEN}  (Build Local)${NC}"
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
# 2. VERIFICAR RECURSOS EXISTENTES
# ========================================
echo -e "${YELLOW}Verificando recursos existentes...${NC}"

# Verificar se Resource Group existe
RG_EXISTS=$(az group exists --name $RESOURCE_GROUP)

if [ "$RG_EXISTS" = "false" ]; then
    echo -e "${YELLOW}Criando Resource Group...${NC}"
    az group create --name $RESOURCE_GROUP --location $LOCATION
    echo -e "${GREEN}✓ Resource Group criado${NC}"
else
    echo -e "${GREEN}✓ Resource Group já existe${NC}"
fi

# Verificar se ACR existe
ACR_EXISTS=$(az acr show --name $ACR_NAME --resource-group $RESOURCE_GROUP 2>/dev/null || echo "")

if [ -z "$ACR_EXISTS" ]; then
    echo -e "${YELLOW}Criando Azure Container Registry...${NC}"
    az acr create \
      --resource-group $RESOURCE_GROUP \
      --name $ACR_NAME \
      --sku Basic \
      --admin-enabled true
    echo -e "${GREEN}✓ ACR criado${NC}"
else
    echo -e "${GREEN}✓ ACR já existe${NC}"
fi

echo ""

# ========================================
# 3. OBTER CREDENCIAIS DO ACR
# ========================================
echo -e "${YELLOW}Obtendo credenciais do ACR...${NC}"
ACR_USERNAME=$(az acr credential show --name $ACR_NAME --query username -o tsv)
ACR_PASSWORD=$(az acr credential show --name $ACR_NAME --query passwords[0].value -o tsv)
ACR_LOGIN_SERVER=$(az acr show --name $ACR_NAME --query loginServer -o tsv)
echo -e "${GREEN}✓ Credenciais obtidas${NC}"
echo ""

# ========================================
# 4. BUILD LOCAL DA IMAGEM DOCKER
# ========================================
echo -e "${YELLOW}Fazendo build local da imagem Docker...${NC}"
docker build -t pessoas-api:latest .
echo -e "${GREEN}✓ Build concluído${NC}"
echo ""

# ========================================
# 5. LOGIN NO ACR E PUSH DA IMAGEM
# ========================================
echo -e "${YELLOW}Fazendo login no ACR...${NC}"
echo $ACR_PASSWORD | docker login $ACR_LOGIN_SERVER --username $ACR_USERNAME --password-stdin
echo -e "${GREEN}✓ Login realizado${NC}"
echo ""

echo -e "${YELLOW}Tagging e push da imagem...${NC}"
docker tag pessoas-api:latest ${ACR_LOGIN_SERVER}/pessoas-api:latest
docker push ${ACR_LOGIN_SERVER}/pessoas-api:latest
echo -e "${GREEN}✓ Imagem enviada para o ACR${NC}"
echo ""

# ========================================
# 6. CRIAR POSTGRESQL FLEXIBLE SERVER
# ========================================
POSTGRES_EXISTS=$(az postgres flexible-server show --name $POSTGRES_SERVER --resource-group $RESOURCE_GROUP 2>/dev/null || echo "")

if [ -z "$POSTGRES_EXISTS" ]; then
    echo -e "${YELLOW}Criando PostgreSQL Flexible Server...${NC}"
    echo -e "${YELLOW}(Isso pode levar alguns minutos)${NC}"
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
      --public-access 0.0.0.0 \
      --yes
    echo -e "${GREEN}✓ PostgreSQL criado${NC}"

    echo -e "${YELLOW}Criando database...${NC}"
    az postgres flexible-server db create \
      --resource-group $RESOURCE_GROUP \
      --server-name $POSTGRES_SERVER \
      --database-name $POSTGRES_DB
    echo -e "${GREEN}✓ Database criado${NC}"
else
    echo -e "${GREEN}✓ PostgreSQL já existe${NC}"
    # Usar credenciais existentes
    POSTGRES_USER=$(az postgres flexible-server show --name $POSTGRES_SERVER --resource-group $RESOURCE_GROUP --query administratorLogin -o tsv)
    echo -e "${YELLOW}Usando usuário existente: $POSTGRES_USER${NC}"
    echo -e "${RED}ATENÇÃO: Usando senha existente. Se não souber, recrie o servidor.${NC}"
fi
echo ""

# ========================================
# 7. CRIAR CONTAINER APPS ENVIRONMENT
# ========================================
ENV_EXISTS=$(az containerapp env show --name $CONTAINER_APP_ENV --resource-group $RESOURCE_GROUP 2>/dev/null || echo "")

if [ -z "$ENV_EXISTS" ]; then
    echo -e "${YELLOW}Criando Container Apps Environment...${NC}"
    az containerapp env create \
      --name $CONTAINER_APP_ENV \
      --resource-group $RESOURCE_GROUP \
      --location $LOCATION
    echo -e "${GREEN}✓ Environment criado${NC}"
else
    echo -e "${GREEN}✓ Environment já existe${NC}"
fi
echo ""

# ========================================
# 8. CRIAR/ATUALIZAR CONTAINER APP
# ========================================
DB_HOST="${POSTGRES_SERVER}.postgres.database.azure.com"

APP_EXISTS=$(az containerapp show --name $CONTAINER_APP_NAME --resource-group $RESOURCE_GROUP 2>/dev/null || echo "")

if [ -z "$APP_EXISTS" ]; then
    echo -e "${YELLOW}Criando Container App...${NC}"

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

    echo -e "${GREEN}✓ Container App criado${NC}"
else
    echo -e "${YELLOW}Atualizando Container App...${NC}"

    az containerapp update \
      --name $CONTAINER_APP_NAME \
      --resource-group $RESOURCE_GROUP \
      --image "${ACR_LOGIN_SERVER}/pessoas-api:latest"

    echo -e "${GREEN}✓ Container App atualizado${NC}"
fi
echo ""

# ========================================
# 9. OBTER URL DA APLICAÇÃO
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
echo -e "Health: ${GREEN}https://$APP_URL/health${NC}"
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
echo -e "${YELLOW}Azure Container Registry:${NC}"
echo -e "Server: ${GREEN}$ACR_LOGIN_SERVER${NC}"
echo -e "Username: ${GREEN}$ACR_USERNAME${NC}"
echo ""
echo -e "${RED}IMPORTANTE: Salve essas credenciais em um local seguro!${NC}"
echo ""
echo -e "${YELLOW}Comandos úteis:${NC}"
echo -e "Ver logs: ${GREEN}az containerapp logs show --name $CONTAINER_APP_NAME --resource-group $RESOURCE_GROUP --follow${NC}"
echo -e "Deletar tudo: ${GREEN}az group delete --name $RESOURCE_GROUP --yes${NC}"
echo ""
