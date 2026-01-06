# Load Tests - Pessoas API

Scripts para teste de carga da API usando Apache Bench (ab).

## Pré-requisitos

- Apache Bench instalado (já disponível no macOS)
- API rodando em `http://localhost:8080`
- Para testes de busca por CPF, é necessário ter dados no banco

## Scripts Disponíveis

### 1. test-health.sh
Testa o endpoint de health check com diferentes níveis de carga:
- 1000 requisições com 10 conexões concorrentes
- 5000 requisições com 50 conexões concorrentes
- 10000 requisições com 100 conexões concorrentes

```bash
bash scripts/load-tests/test-health.sh
```

### 2. test-list-persons.sh
Testa o endpoint de listagem de pessoas com paginação:
- 500 requisições com 10 conexões concorrentes (página 1, tamanho padrão)
- 1000 requisições com 25 conexões concorrentes (página 1, 20 itens)
- 2000 requisições com 50 conexões concorrentes (ordenado por nome)

```bash
bash scripts/load-tests/test-list-persons.sh
```

### 3. test-find-by-cpf.sh
Testa o endpoint de busca por CPF:
- 1000 requisições com 10 conexões concorrentes
- 3000 requisições com 30 conexões concorrentes
- 5000 requisições com 50 conexões concorrentes

**Importante:** Certifique-se de que o CPF `11144477735` existe no banco antes de rodar este teste.

```bash
bash scripts/load-tests/test-find-by-cpf.sh
```

### 4. test-create-person.sh
Testa o endpoint de criação de pessoa:
- 100 requisições com 5 conexões concorrentes
- 200 requisições com 10 conexões concorrentes

**Nota:** Este teste criará entradas duplicadas e falhará devido à constraint de CPF único. Use com cautela.

```bash
bash scripts/load-tests/test-create-person.sh
```

### 5. run-all-tests.sh
Executa todos os testes de leitura (health, list, find by CPF) em sequência:

```bash
bash scripts/load-tests/run-all-tests.sh
```

## Métricas Importantes

O Apache Bench fornece várias métricas importantes:

- **Requests per second**: Número de requisições processadas por segundo
- **Time per request (mean)**: Tempo médio por requisição
- **Time per request (mean, across all concurrent requests)**: Tempo médio considerando concorrência
- **Transfer rate**: Taxa de transferência de dados
- **Connection Times (min/mean/median/max)**: Tempos de conexão
- **Percentage of requests served within a certain time**: Distribuição percentil dos tempos de resposta

## Exemplo de Uso Completo

```bash
# 1. Inicie a API
go run cmd/api/main.go

# 2. Em outro terminal, popule o banco com alguns dados
curl -X POST http://localhost:8080/api/v1/persons \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "cpf": "11144477735",
    "birth_date": "1990-01-01T00:00:00Z",
    "phone": "81 91234-5678",
    "email": "test@example.com"
  }'

# 3. Execute os testes
bash scripts/load-tests/run-all-tests.sh
```

## Interpretando os Resultados

### Bons Resultados
- **Requests per second**: > 1000 rps para endpoints simples (health)
- **Mean time per request**: < 100ms para a maioria das requisições
- **Failed requests**: 0 (exceto para testes de criação com CPF duplicado)
- **50% of requests served within**: < 50ms
- **95% of requests served within**: < 200ms

### Sinais de Problemas
- Alta taxa de requisições falhadas
- Tempo médio de resposta crescendo significativamente
- Muitas requisições levando mais de 1 segundo
- Erros de timeout ou conexão

## Customizando os Testes

Você pode modificar os scripts para testar diferentes cenários:

```bash
# Aumentar o número de requisições
ab -n 20000 -c 100 http://localhost:8080/health

# Testar com timeout customizado
ab -n 1000 -c 50 -s 30 http://localhost:8080/api/v1/persons

# Salvar resultados em arquivo
ab -n 1000 -c 50 http://localhost:8080/health > results.txt

# Formato verboso com detalhes de percentis
ab -n 1000 -c 50 -v 2 http://localhost:8080/health
```

## Dicas de Performance

1. **Aqueça a aplicação** antes de medir: Execute alguns requests para garantir que o JIT e caches estejam prontos
2. **Execute múltiplas vezes**: Execute os testes várias vezes para obter médias mais precisas
3. **Monitore recursos**: Use `top` ou `htop` para monitorar CPU e memória durante os testes
4. **Monitore o banco**: Observe as queries no PostgreSQL durante os testes
5. **Varie a carga**: Teste com diferentes níveis de concorrência para encontrar o ponto ideal

## Ferramentas Alternativas

Se quiser ferramentas mais avançadas, considere instalar:

```bash
# hey - ferramenta moderna em Go
brew install hey

# wrk - benchmarking HTTP altamente performático
brew install wrk

# k6 - testes de carga com scripts em JavaScript
brew install k6
```
