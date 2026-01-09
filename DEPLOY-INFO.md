# Pessoas API - Informações de Deploy no Azure

## ✅ Deploy Realizado com Sucesso!

**Data do Deploy:** 09/01/2026

## 🌐 URLs da Aplicação

- **API Base URL:** https://pessoas-api.politerock-20d39e69.brazilsouth.azurecontainerapps.io
- **Swagger UI:** https://pessoas-api.politerock-20d39e69.brazilsouth.azurecontainerapps.io/swagger/index.html
- **Health Check:** https://pessoas-api.politerock-20d39e69.brazilsouth.azurecontainerapps.io/health

## 🗄️ Database (PostgreSQL)

- **Host:** pessoas-api-db-br.postgres.database.azure.com
- **Database:** postgres
- **Username:** apiuser
- **Password:** `jz7ZmUdRneE+42XDLlLfaqz3ycKmdzD3`
- **Connection String:**
  ```
  postgresql://apiuser:jz7ZmUdRneE+42XDLlLfaqz3ycKmdzD3@pessoas-api-db-br.postgres.database.azure.com/postgres?sslmode=require
  ```

## 🔐 JWT Secret

```
IpnLE98R+8OZQDbrOyMpslGwlPL6SNb7m51UGTpoVTQ=
```

## ☁️ Recursos Azure

| Recurso | Nome | Tipo | Localização |
|---------|------|------|-------------|
| Resource Group | `pessoas-api-rg` | Resource Group | Brazil South |
| Container Registry | `pessoasapiacr` | Azure Container Registry | East US |
| Container App | `pessoas-api` | Container App | Brazil South |
| Database | `pessoas-api-db-br` | PostgreSQL Flexible Server | Brazil South |
| Environment | `pessoas-api-env` | Container Apps Environment | Brazil South |

## 📊 Configuração dos Recursos

### Container App
- **CPU:** 0.5 vCPU
- **Memory:** 1.0 GB
- **Min Replicas:** 1
- **Max Replicas:** 2
- **Port:** 8080

### PostgreSQL
- **SKU:** Standard_B1ms (Burstable)
- **Storage:** 32 GB
- **Version:** 16

## 🔧 Comandos Úteis

### Ver logs em tempo real
```bash
az containerapp logs show \
  --name pessoas-api \
  --resource-group pessoas-api-rg \
  --follow
```

### Verificar status
```bash
az containerapp show \
  --name pessoas-api \
  --resource-group pessoas-api-rg \
  --query properties.runningStatus
```

### Atualizar a aplicação
```bash
# 1. Build nova imagem
docker build -t pessoas-api:latest .

# 2. Login no ACR
az acr login --name pessoasapiacr

# 3. Tag e push
docker tag pessoas-api:latest pessoasapiacr.azurecr.io/pessoas-api:latest
docker push pessoasapiacr.azurecr.io/pessoas-api:latest

# 4. Atualizar container app
az containerapp update \
  --name pessoas-api \
  --resource-group pessoas-api-rg \
  --image pessoasapiacr.azurecr.io/pessoas-api:latest
```

### Escalar aplicação
```bash
az containerapp update \
  --name pessoas-api \
  --resource-group pessoas-api-rg \
  --min-replicas 2 \
  --max-replicas 5
```

### Parar aplicação (economizar custos)
```bash
# Parar container app
az containerapp update \
  --name pessoas-api \
  --resource-group pessoas-api-rg \
  --min-replicas 0 \
  --max-replicas 0

# Parar PostgreSQL
az postgres flexible-server stop \
  --resource-group pessoas-api-rg \
  --name pessoas-api-db-br
```

### Iniciar aplicação novamente
```bash
# Iniciar container app
az containerapp update \
  --name pessoas-api \
  --resource-group pessoas-api-rg \
  --min-replicas 1 \
  --max-replicas 2

# Iniciar PostgreSQL
az postgres flexible-server start \
  --resource-group pessoas-api-rg \
  --name pessoas-api-db-br
```

### Deletar todos os recursos
```bash
az group delete \
  --name pessoas-api-rg \
  --yes --no-wait
```

## 💰 Estimativa de Custos

| Serviço | SKU/Config | Custo Mensal (USD) |
|---------|------------|-------------------|
| Container App | 0.5 vCPU, 1GB RAM | ~$10-15 |
| PostgreSQL | Standard_B1ms | ~$15-20 |
| Container Registry | Basic | ~$5 |
| **TOTAL ESTIMADO** | | **~$30-40/mês** |

**Nota:** Com créditos gratuitos do Azure ($200 USD), você pode usar gratuitamente por vários meses!

## 🧪 Testar a API

### Health Check
```bash
curl https://pessoas-api.politerock-20d39e69.brazilsouth.azurecontainerapps.io/health
```

### Registrar Operador
```bash
curl -X POST https://pessoas-api.politerock-20d39e69.brazilsouth.azurecontainerapps.io/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "Admin123!"
  }'
```

### Login
```bash
curl -X POST https://pessoas-api.politerock-20d39e69.brazilsouth.azurecontainerapps.io/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin123!"
  }'
```

## ⚠️ Troubleshooting

### Container não inicia
1. Ver logs detalhados:
   ```bash
   az containerapp logs show --name pessoas-api --resource-group pessoas-api-rg --type console
   ```

2. Verificar variáveis de ambiente:
   ```bash
   az containerapp show --name pessoas-api --resource-group pessoas-api-rg --query 'properties.template.containers[0].env'
   ```

### Problemas de conexão com banco
1. Verificar firewall do PostgreSQL:
   ```bash
   az postgres flexible-server firewall-rule list \
     --resource-group pessoas-api-rg \
     --name pessoas-api-db-br
   ```

2. Testar conexão ao banco:
   ```bash
   psql "postgresql://apiuser:jz7ZmUdRneE+42XDLlLfaqz3ycKmdzD3@pessoas-api-db-br.postgres.database.azure.com/postgres?sslmode=require"
   ```

## 📚 Próximos Passos

1. **Configurar CI/CD** - Automatizar deploy via GitHub Actions
2. **Custom Domain** - Adicionar domínio personalizado
3. **Monitoring** - Configurar alertas e métricas
4. **Backup** - Configurar backup automático do PostgreSQL
5. **Scaling Rules** - Adicionar regras de auto-scaling baseadas em métricas

## 🔗 Links Úteis

- [Azure Container Apps Documentation](https://learn.microsoft.com/azure/container-apps/)
- [Azure PostgreSQL Flexible Server](https://learn.microsoft.com/azure/postgresql/flexible-server/)
- [Azure CLI Reference](https://learn.microsoft.com/cli/azure/)
- [Swagger UI](https://pessoas-api.politerock-20d39e69.brazilsouth.azurecontainerapps.io/swagger/index.html)
