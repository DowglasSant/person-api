# Pessoas API

API REST para gerenciamento de pessoas, construída com Go seguindo os princípios de **Hexagonal Architecture** (Ports & Adapters).

## Arquitetura

Esta aplicação implementa **Hexagonal Architecture**, separando claramente o domínio da infraestrutura através de portas (interfaces) e adaptadores (implementações).

```
.
├── cmd/
│   └── api/
│       └── main.go                    # Entry point da aplicação
├── internal/
│   ├── contract/                      # DTOs e contratos de API
│   │   └── person/
│   │       ├── new_person_dto.go
│   │       ├── person_response_dto.go
│   │       └── pagination_dto.go
│   │
│   ├── domain/                        # 🔵 HEXÁGONO (Núcleo da Aplicação)
│   │   └── person/
│   │       ├── model/                 # Entidades de domínio
│   │       │   └── person.go
│   │       ├── ports/                 # 🔌 PORTAS (Interfaces)
│   │       │   ├── service.go         # PersonService interface
│   │       │   └── repository.go      # PersonRepository interface
│   │       ├── service/               # Lógica de negócio (implementa porta)
│   │       │   └── person_service.go  # PersonServiceImpl
│   │       ├── error/                 # Erros de domínio
│   │       └── utils/                 # Utilitários de domínio
│   │
│   └── infrastructure/                # ⚙️ ADAPTADORES (Camada Externa)
│       ├── database/                  # Configuração de banco de dados
│       ├── persistence/               # Adapter de persistência
│       │   └── person/
│       │       ├── person_entity.go   # Entidade GORM
│       │       └── person_repository_impl.go  # Implementa porta
│       └── http/                      # Adapter HTTP
│           ├── handler/               # HTTP handlers
│           ├── router/                # Configuração de rotas
│           └── middleware/            # Middlewares
│
├── scripts/
│   ├── create_schema.sql              # Schema SQL (documentação)
│   └── load-tests/                    # Scripts de teste de carga
└── .env                               # Variáveis de ambiente
```

## Tecnologias

- **Go 1.25.3**
- **Gin** - Framework web
- **GORM** - ORM
- **PostgreSQL** - Banco de dados
- **Testify** - Biblioteca de testes
- **Swagger/OpenAPI 3.0** - Documentação da API

## Configuração

### Pré-requisitos

- Go 1.25+
- PostgreSQL 14+

### Variáveis de Ambiente

Copie o arquivo `.env.example` para `.env` e configure:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=@Pass2025
DB_NAME=postgres
DB_SCHEMA=people
DB_SSLMODE=disable
```

### Instalação

```bash
# Instalar dependências
go mod download

# Build da aplicação
go build -o bin/api cmd/api/main.go

# Executar
./bin/api
```

## Rodando a aplicação

```bash
go run cmd/api/main.go
```

A API estará disponível em `http://localhost:8080`

## Documentação da API

### Swagger UI

A documentação completa da API está disponível através do **Swagger UI** quando a aplicação está rodando:

🔗 **http://localhost:8080/swagger/index.html**

A interface Swagger fornece:
- ✅ Documentação interativa de todos os endpoints
- ✅ Descrição detalhada de parâmetros e respostas
- ✅ Exemplos de requisições e respostas
- ✅ Possibilidade de testar endpoints diretamente no navegador
- ✅ Schemas de dados (DTOs) com validações

### Gerando Documentação

Se você modificar as anotações Swagger nos handlers ou DTOs, regenere a documentação:

```bash
# Instalar swag CLI (uma vez)
go install github.com/swaggo/swag/cmd/swag@latest

# Gerar documentação
swag init -g cmd/api/main.go -o docs
```

Os arquivos gerados em `docs/` devem ser commitados no repositório.

## Endpoints

### Health Check

```bash
GET /health
```

**Resposta (200):**
```json
{
  "status": "ok"
}
```

### Criar Pessoa

```bash
POST /api/v1/persons
Content-Type: application/json

{
  "name": "John Doe",
  "cpf": "111.444.777-35",
  "birth_date": "1990-01-01T00:00:00Z",
  "phone": "81 91234-5678",
  "email": "john.doe@example.com"
}
```

**Resposta de sucesso (201):**
```json
{
  "id": 1,
  "message": "Person created successfully"
}
```

**Resposta de erro (422):**
```json
{
  "error": "validation_error",
  "message": "CPF inválido"
}
```

### Listar Pessoas (com paginação)

```bash
GET /api/v1/persons?page=1&page_size=10&sort=name&order=asc
```

**Parâmetros de query (opcionais):**
- `page` - Número da página (default: 1, mínimo: 1)
- `page_size` - Itens por página (default: 10, máximo: 100)
- `sort` - Campo para ordenação (default: id)
  - Valores válidos: `id`, `name`, `cpf`, `email`, `created_at`, `updated_at`
- `order` - Direção da ordenação (default: desc)
  - Valores válidos: `asc`, `desc`

**Resposta de sucesso (200):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Alice Silva",
      "cpf": "11144477735",
      "birth_date": "1990-01-01T00:00:00Z",
      "phone_number": "81912345678",
      "email": "alice@example.com",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    {
      "id": 2,
      "name": "Bob Santos",
      "cpf": "22233344405",
      "birth_date": "1985-03-15T00:00:00Z",
      "phone_number": "11987654321",
      "email": "bob@example.com",
      "created_at": "2024-01-01T11:00:00Z",
      "updated_at": "2024-01-01T11:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_items": 2,
  "total_pages": 1
}
```

**Exemplos de uso:**
```bash
# Primeira página com 10 itens
GET /api/v1/persons

# Segunda página com 20 itens
GET /api/v1/persons?page=2&page_size=20

# Ordenar por nome em ordem crescente
GET /api/v1/persons?sort=name&order=asc

# Ordenar por data de criação (mais recentes primeiro)
GET /api/v1/persons?sort=created_at&order=desc
```

### Buscar Pessoa por CPF

```bash
GET /api/v1/persons/cpf/:cpf
```

**Parâmetros:**
- `cpf` - CPF da pessoa (pode ser formatado ou apenas números)

**Resposta de sucesso (200):**
```json
{
  "id": 1,
  "name": "John Doe",
  "cpf": "11144477735",
  "birth_date": "1990-01-01T00:00:00Z",
  "phone_number": "81912345678",
  "email": "john.doe@example.com",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

**Resposta de não encontrado (404):**
```json
{
  "error": "not_found",
  "message": "Person not found with the provided CPF"
}
```

**Exemplos de uso:**
```bash
# CPF sem formatação
GET /api/v1/persons/cpf/11144477735

# CPF com formatação (será automaticamente convertido)
GET /api/v1/persons/cpf/111.444.777-35
```

## Testes

### Testes Unitários

```bash
# Rodar todos os testes
go test ./...

# Rodar testes com verbose
go test -v ./...

# Rodar testes de um pacote específico
go test -v ./internal/domain/person/model/
go test -v ./internal/infrastructure/persistence/person/
```

### Testes de Carga

A aplicação inclui scripts de teste de carga usando Apache Bench para avaliar performance sob diferentes cenários:

```bash
# Teste rápido (1000 requisições, 20 conexões concorrentes)
bash scripts/load-tests/quick-test.sh

# Rodar todos os testes de carga
bash scripts/load-tests/run-all-tests.sh

# Testes individuais
bash scripts/load-tests/test-health.sh           # Testa endpoint /health
bash scripts/load-tests/test-list-persons.sh     # Testa listagem com paginação
bash scripts/load-tests/test-find-by-cpf.sh      # Testa busca por CPF
bash scripts/load-tests/test-create-person.sh    # Testa criação de pessoa
```

**Documentação completa:** Ver [scripts/load-tests/README.md](scripts/load-tests/README.md) para detalhes sobre os testes, métricas e interpretação de resultados.

## Estrutura do Banco de Dados

### Schema: people

**Tabela: person**

| Campo        | Tipo         | Descrição                    |
|--------------|--------------|------------------------------|
| id           | SERIAL4      | Chave primária (autogerado)  |
| name         | VARCHAR(255) | Nome completo                |
| cpf          | VARCHAR(11)  | CPF (apenas números, único)  |
| birth_date   | DATE         | Data de nascimento           |
| phone_number | VARCHAR(11)  | Telefone (apenas números)    |
| email        | VARCHAR(255) | Email                        |
| created_at   | TIMESTAMP    | Data de criação              |
| updated_at   | TIMESTAMP    | Data de atualização          |

### SQL de criação

```sql
CREATE SCHEMA IF NOT EXISTS people;

CREATE TABLE people.person (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cpf VARCHAR(11) NOT NULL UNIQUE,
    birth_date DATE NOT NULL,
    phone_number VARCHAR(11) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_person_cpf ON people.person(cpf);
```

## Validações

O domínio aplica as seguintes validações:

- **Nome**: obrigatório, não pode ser vazio
- **CPF**: obrigatório, deve ser válido (algoritmo de validação de CPF)
- **Email**: obrigatório, deve ter formato válido
- **Telefone**: obrigatório, 10 ou 11 dígitos
- **Data de nascimento**: obrigatória, não pode ser futura

## Logging

A aplicação possui um sistema de logging abrangente que registra:

### Logs de Requisição HTTP

O middleware de logging registra automaticamente:
- **Início da requisição**: Método HTTP, path e IP do cliente
- **Fim da requisição**: Status HTTP, duração e IP do cliente
- **Erros**: Requisições com status >= 400 são marcadas como `[REQUEST ERROR]`

**Formato dos logs:**
```
[REQUEST START] POST /api/v1/persons - Client: 127.0.0.1
[REQUEST END] POST /api/v1/persons - Status: 201 - Duration: 15.4ms - Client: 127.0.0.1
[REQUEST ERROR] POST /api/v1/persons - Status: 422 - Duration: 2.1ms - Client: 127.0.0.1
```

### Logs de Handler

Cada handler registra:
- **[INFO]**: Parâmetros recebidos
- **[SUCCESS]**: Operações bem-sucedidas com detalhes
- **[ERROR]**: Erros de validação ou processamento
- **[WARN]**: Situações de atenção (ex: recurso não encontrado)

**Exemplos:**
```
[INFO] ListPersons - Fetching page: 1, pageSize: 10, sort: name, order: asc
[SUCCESS] CreatePerson - Person created with ID: 1, Name: João Silva
[ERROR] CreatePerson - Validation error for CPF 12345678900: cpf is invalid
[WARN] FindPersonByCPF - Person not found with CPF: 99999999999
```

### Logs de Repository

A camada de persistência registra:
- **[REPO]**: Início de operações no banco
- **[REPO SUCCESS]**: Operações bem-sucedidas
- **[REPO ERROR]**: Erros de banco de dados
- **[REPO WARN]**: Registros não encontrados

**Exemplos:**
```
[REPO] Save - Attempting to save person with CPF: 11144477735
[REPO SUCCESS] Save - Person saved with ID: 1, CPF: 11144477735
[REPO] FindAll - Querying page: 1, pageSize: 10, sort: name asc
[REPO SUCCESS] FindAll - Retrieved 10 persons out of 25 total
[REPO WARN] FindByCPF - Person not found with CPF: 12345678900
```

### Níveis de Log

- **[INFO]**: Informações gerais sobre o fluxo da aplicação
- **[SUCCESS]**: Operações concluídas com sucesso
- **[WARN]**: Situações que merecem atenção mas não são erros
- **[ERROR]**: Erros que impediram a conclusão de uma operação
- **[REPO]**: Operações específicas do repositório

## Princípios de Arquitetura Hexagonal

Este projeto segue os princípios de **Hexagonal Architecture** (também conhecida como Ports & Adapters):

### Conceitos Principais

#### 🔵 Hexágono (Núcleo)
O domínio da aplicação, contendo a lógica de negócio pura, sem dependências externas.

- **Entidades**: `domain/person/model/` - Objetos de negócio
- **Serviços**: `domain/person/service/` - Lógica de negócio
- **Portas**: `domain/person/ports/` - Interfaces que definem contratos

#### 🔌 Portas (Interfaces)
Definem os contratos entre o hexágono e o mundo externo.

**Portas Primárias (Driving Ports)**: Comandam o hexágono
- `PersonService` - Interface que os handlers HTTP usam para executar operações de negócio

**Portas Secundárias (Driven Ports)**: São comandadas pelo hexágono
- `PersonRepository` - Interface que o domínio usa para persistência

#### ⚙️ Adaptadores (Implementações)
Conectam o hexágono ao mundo externo.

**Adaptadores Primários (Driving Adapters)**: Iniciam interações
- `PersonHandler` - Adapter HTTP que implementa endpoints REST

**Adaptadores Secundários (Driven Adapters)**: Respondem a solicitações
- `PersonRepositoryImpl` - Adapter PostgreSQL que implementa persistência

### Fluxo de Dependências

```
┌─────────────────────────────────────────────────────────────┐
│                    ADAPTADORES PRIMÁRIOS                     │
│                     (Driving Adapters)                       │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │          PersonHandler (HTTP Adapter)                │   │
│  │     infrastructure/http/handler/person_handler.go    │   │
│  └────────────────────┬─────────────────────────────────┘   │
│                       │ depende de                           │
│                       ▼                                      │
│  ┌──────────────────────────────────────────────────────┐   │
│  │         ports.PersonService (interface)              │   │
│  │         domain/person/ports/service.go               │   │
│  └────────────────────┬─────────────────────────────────┘   │
└───────────────────────┼──────────────────────────────────────┘
                        │ implementada por
┌───────────────────────┼──────────────────────────────────────┐
│                       ▼                HEXÁGONO              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │         PersonServiceImpl                            │   │
│  │      domain/person/service/person_service.go         │   │
│  └────────────────────┬─────────────────────────────────┘   │
│                       │ depende de                           │
│                       ▼                                      │
│  ┌──────────────────────────────────────────────────────┐   │
│  │      ports.PersonRepository (interface)              │   │
│  │      domain/person/ports/repository.go               │   │
│  └────────────────────┬─────────────────────────────────┘   │
└───────────────────────┼──────────────────────────────────────┘
                        │ implementada por
┌───────────────────────┼──────────────────────────────────────┐
│                       ▼        ADAPTADORES SECUNDÁRIOS       │
│  ┌──────────────────────────────────────────────────────┐   │
│  │       PersonRepositoryImpl (PostgreSQL Adapter)      │   │
│  │  infrastructure/persistence/person/repository.go     │   │
│  └──────────────────────────────────────────────────────┘   │
│                     (Driven Adapters)                        │
└──────────────────────────────────────────────────────────────┘
```

### Benefícios da Arquitetura

1. **Inversão de Dependência**: Infraestrutura depende do domínio, não o contrário
2. **Testabilidade**: Lógica de negócio pode ser testada sem banco ou HTTP
3. **Substituibilidade**: Troque PostgreSQL por MongoDB sem tocar no domínio
4. **Isolamento**: Mudanças em frameworks não afetam regras de negócio
5. **Clareza**: Separação explícita entre portas (contratos) e adaptadores (implementações)

### Regras de Dependência

- ✅ Adaptadores **podem** depender de Portas
- ✅ Hexágono **pode** definir Portas
- ✅ Hexágono **pode** depender apenas de suas próprias Portas
- ❌ Hexágono **nunca** depende de Adaptadores
- ❌ Portas **nunca** dependem de Adaptadores
