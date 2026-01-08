# 🔐 Sistema de Autenticação com Operadores

## Resumo da Implementação

Foi implementado um sistema completo de autenticação baseado em **operadores** (usuários do sistema que gerenciam as pessoas cadastradas), substituindo o modelo anterior que usava dados das próprias pessoas para autenticação.

## Arquitetura

### Modelo de Domínio

**Operator** (`internal/domain/operator/model/operator.go`):
- ID, Username, Email, Password (bcrypt), Active, CreatedAt, UpdatedAt
- Validações: username (3-50 chars), email válido, senha (8-72 chars)
- Método `ValidatePassword()` para verificar senha com bcrypt
- Método `UpdatePassword()` para atualizar senha

### Segurança

✅ **Senhas hasheadas com bcrypt** (custo padrão: 10)  
✅ **Validação de credenciais única** ("invalid credentials" para username/password incorretos)  
✅ **Verificação de conta ativa** antes do login  
✅ **JWT gerado com dados do operador** (operator_id, username)  
✅ **Verificação de duplicação** (username e email únicos)  

### Endpoints Criados

#### 1. POST /api/v1/auth/register
Registra um novo operador no sistema.

**Request:**
```json
{
  "username": "john.doe",
  "email": "john.doe@company.com",
  "password": "SecurePass123!"
}
```

**Responses:**
- **201 Created**: Operador criado com sucesso
```json
{
  "id": 1,
  "message": "Operator registered successfully"
}
```

- **400 Bad Request**: Dados inválidos
- **409 Conflict**: Username ou email já existe
- **422 Unprocessable Entity**: Erro de validação

#### 2. POST /api/v1/auth/login
Autentica um operador e retorna JWT.

**Request:**
```json
{
  "username": "john.doe",
  "password": "SecurePass123!"
}
```

**Responses:**
- **200 OK**: Login bem-sucedido
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Login successful"
}
```

- **400 Bad Request**: Dados inválidos
- **401 Unauthorized**: Credenciais inválidas ou conta inativa

### Estrutura da Tabela

```sql
CREATE TABLE operators (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,  -- Bcrypt hash
    active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
```

### Fluxo de Autenticação

```
1. Operador se registra      → POST /auth/register
2. Senha é hasheada (bcrypt) → Salva no banco
3. Operador faz login        → POST /auth/login
4. Sistema valida credenciais → Verifica hash
5. JWT é gerado               → Token com 24h de validade
6. Token é usado nas requests → Header: Authorization: Bearer <token>
7. Middleware valida JWT      → Extrai operator_id e username
```

### Rotas Protegidas

Todas as rotas `/api/v1/persons/*` agora requerem autenticação:

```
GET    /api/v1/persons          ← JWT obrigatório
POST   /api/v1/persons          ← JWT obrigatório
GET    /api/v1/persons/cpf/:cpf ← JWT obrigatório
```

### Rotas Públicas

```
POST   /api/v1/auth/register    ← Sem JWT
POST   /api/v1/auth/login       ← Sem JWT
GET    /health                  ← Sem JWT
GET    /swagger/*any            ← Sem JWT
```

## Como Usar

### 1. Criar a tabela no banco

```bash
psql -U postgres -d postgres -f scripts/create_operators_table.sql
```

### 2. Registrar um operador

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@company.com",
    "password": "Admin@123456"
  }'
```

### 3. Fazer login e obter token

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin@123456"
  }' | jq -r '.token')

echo $TOKEN
```

### 4. Usar o token para acessar rotas protegidas

```bash
# Listar pessoas
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/persons

# Criar pessoa
curl -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "cpf": "111.444.777-35",
    "birth_date": "1990-01-15T00:00:00Z",
    "phone": "81912345678",
    "email": "joao.silva@email.com"
  }' \
  http://localhost:8080/api/v1/persons
```

## Arquivos Criados

```
internal/domain/operator/
├── model/
│   └── operator.go                    # Modelo de domínio com bcrypt
├── ports/
│   ├── repository.go                  # Interface do repositório
│   └── service.go                     # Interface do serviço de auth
└── service/
    └── auth_service.go                # Lógica de registro e login

internal/infrastructure/persistence/operator/
├── operator_entity.go                 # Entidade GORM
└── operator_repository_impl.go        # Implementação do repositório

internal/infrastructure/http/handler/
└── auth_handler.go                    # Handler HTTP para auth

internal/contract/auth/
├── register_dto.go                    # DTO de registro
└── login_dto.go                       # DTO de login/resposta

scripts/
└── create_operators_table.sql         # Script de criação da tabela
```

## Segurança Implementada

- ✅ **Bcrypt** para hash de senhas (salt automático)
- ✅ **Validação forte** de senhas (mínimo 8 caracteres)
- ✅ **JWT com expiração** (24 horas)
- ✅ **Validação de unicidade** (username e email)
- ✅ **Conta ativa** verificada no login
- ✅ **Mensagens genéricas** de erro (não revela se username existe)
- ✅ **Rate limiting** aplicado (60 req/min)
- ✅ **CORS** configurável
- ✅ **Security headers** aplicados

## Testes

Todos os 72 testes anteriores continuam passando:
```bash
go test ./...
# PASS: 72/72 tests
```

## Próximos Passos Sugeridos

1. **Refresh tokens** para renovar sessão
2. **Recuperação de senha** via email
3. **2FA (Two-Factor Authentication)**
4. **Roles e permissões** (admin, operator, viewer)
5. **Auditoria** de ações dos operadores
6. **Limite de tentativas de login**
7. **Sessões ativas** (logout de todas as sessões)

---

**Implementado por:** Claude Sonnet 4.5  
**Data:** 2026-01-08
