# Deploy no Azure - Pessoas API

Este guia mostra como fazer deploy da API no Azure usando Container Apps e PostgreSQL.

## 📋 Pré-requisitos

1. **Azure CLI instalado**
   ```bash
   # macOS
   brew install azure-cli

   # Verificar instalação
   az --version
   ```

2. **Conta Azure ativa**
   - Criar conta gratuita: https://azure.microsoft.com/free/
   - Crédito inicial: $200 USD por 30 dias

## 🚀 Opção 1: Deploy Automatizado (Recomendado)

### Passo a Passo:

```bash
# 1. Dar permissão de execução ao script
chmod +x azure-deploy.sh

# 2. Executar o deploy
./azure-deploy.sh
```

O script irá:
- ✅ Criar Resource Group
- ✅ Criar Azure Container Registry (ACR)
- ✅ Fazer build e push da imagem Docker
- ✅ Criar PostgreSQL Flexible Server
- ✅ Criar Container Apps Environment
- ✅ Fazer deploy da aplicação
- ✅ Configurar variáveis de ambiente
- ✅ Exibir URL da aplicação

**Tempo estimado:** 10-15 minutos

## 🛠️ Opção 2: Deploy Manual

### 1. Login no Azure
```bash
az login
```

### 2. Configurar variáveis
```bash
RESOURCE_GROUP="pessoas-api-rg"
LOCATION="eastus"
ACR_NAME="pessoasapiacr"
POSTGRES_SERVER="pessoas-api-db"
POSTGRES_USER="apiuser"
POSTGRES_PASSWORD="SuaSenhaSegura123!"
JWT_SECRET="seu-jwt-secret-minimo-32-caracteres-aqui"
```

### 3. Criar Resource Group
```bash
az group create \
  --name $RESOURCE_GROUP \
  --location $LOCATION
```

### 4. Criar Container Registry
```bash
az acr create \
  --resource-group $RESOURCE_GROUP \
  --name $ACR_NAME \
  --sku Basic \
  --admin-enabled true
```

### 5. Build e Push da Imagem
```bash
az acr build \
  --registry $ACR_NAME \
  --image pessoas-api:latest \
  --file Dockerfile \
  .
```

### 6. Criar PostgreSQL Flexible Server
```bash
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
```

### 7. Criar Database
```bash
az postgres flexible-server db create \
  --resource-group $RESOURCE_GROUP \
  --server-name $POSTGRES_SERVER \
  --database-name postgres
```

### 8. Criar Container Apps Environment
```bash
az containerapp env create \
  --name pessoas-api-env \
  --resource-group $RESOURCE_GROUP \
  --location $LOCATION
```

### 9. Obter credenciais do ACR
```bash
ACR_USERNAME=$(az acr credential show --name $ACR_NAME --query username -o tsv)
ACR_PASSWORD=$(az acr credential show --name $ACR_NAME --query passwords[0].value -o tsv)
ACR_LOGIN_SERVER=$(az acr show --name $ACR_NAME --query loginServer -o tsv)
```

### 10. Fazer Deploy do Container App
```bash
DB_HOST="${POSTGRES_SERVER}.postgres.database.azure.com"

az containerapp create \
  --name pessoas-api \
  --resource-group $RESOURCE_GROUP \
  --environment pessoas-api-env \
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
    DB_NAME=postgres \
    DB_SCHEMA=people \
    DB_SSLMODE=require \
    JWT_SECRET=$JWT_SECRET \
    CORS_ALLOWED_ORIGINS="*"
```

### 11. Obter URL da aplicação
```bash
az containerapp show \
  --name pessoas-api \
  --resource-group $RESOURCE_GROUP \
  --query properties.configuration.ingress.fqdn \
  -o tsv
```

## 💰 Estimativa de Custos (Tier Mais Barato)

| Serviço | SKU | Custo Mensal (USD) |
|---------|-----|-------------------|
| Container Apps | 0.5 vCPU, 1GB RAM | ~$10-15 |
| PostgreSQL Flexible | Standard_B1ms | ~$15-20 |
| Container Registry | Basic | ~$5 |
| **TOTAL** | | **~$30-40/mês** |

### Opções para Reduzir Custos:

1. **Usar créditos gratuitos Azure**: $200 USD por 30 dias
2. **Azure for Students**: $100 USD por ano (sem cartão de crédito)
3. **Parar recursos quando não usar**:
   ```bash
   # Parar Container App
   az containerapp update --name pessoas-api \
     --resource-group $RESOURCE_GROUP \
     --min-replicas 0 --max-replicas 0

   # Parar PostgreSQL
   az postgres flexible-server stop \
     --resource-group $RESOURCE_GROUP \
     --name $POSTGRES_SERVER
   ```

## 🔄 Atualizar a Aplicação

```bash
# 1. Fazer novo build
az acr build \
  --registry $ACR_NAME \
  --image pessoas-api:latest \
  --file Dockerfile \
  .

# 2. Atualizar Container App
az containerapp update \
  --name pessoas-api \
  --resource-group $RESOURCE_GROUP \
  --image "${ACR_LOGIN_SERVER}/pessoas-api:latest"
```

## 📊 Monitoramento

### Ver logs da aplicação
```bash
az containerapp logs show \
  --name pessoas-api \
  --resource-group $RESOURCE_GROUP \
  --follow
```

### Ver status
```bash
az containerapp show \
  --name pessoas-api \
  --resource-group $RESOURCE_GROUP \
  --query properties.runningStatus
```

## 🗑️ Deletar Todos os Recursos

```bash
az group delete \
  --name $RESOURCE_GROUP \
  --yes --no-wait
```

## 🔐 Segurança

### Configurar Custom Domain (Opcional)
```bash
az containerapp hostname add \
  --name pessoas-api \
  --resource-group $RESOURCE_GROUP \
  --hostname api.seudominio.com
```

### Configurar HTTPS (Automático)
O Azure Container Apps já configura HTTPS automaticamente.

## 📝 Testar a API

```bash
# Health check
curl https://seu-app.azurecontainerapps.io/health

# Swagger UI
https://seu-app.azurecontainerapps.io/swagger/index.html

# Registrar operador
curl -X POST https://seu-app.azurecontainerapps.io/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "Admin123!"
  }'
```

## ⚠️ Troubleshooting

### Container não inicia
```bash
# Ver logs detalhados
az containerapp logs show \
  --name pessoas-api \
  --resource-group $RESOURCE_GROUP \
  --tail 100
```

### Problemas de conexão com banco
```bash
# Verificar firewall do PostgreSQL
az postgres flexible-server firewall-rule list \
  --resource-group $RESOURCE_GROUP \
  --name $POSTGRES_SERVER

# Adicionar regra se necessário
az postgres flexible-server firewall-rule create \
  --resource-group $RESOURCE_GROUP \
  --name $POSTGRES_SERVER \
  --rule-name AllowAllAzureIPs \
  --start-ip-address 0.0.0.0 \
  --end-ip-address 0.0.0.0
```

## 📚 Referências

- [Azure Container Apps Documentation](https://learn.microsoft.com/azure/container-apps/)
- [Azure PostgreSQL Flexible Server](https://learn.microsoft.com/azure/postgresql/flexible-server/)
- [Azure CLI Reference](https://learn.microsoft.com/cli/azure/)
