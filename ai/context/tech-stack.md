# Tech Stack

## Linguagem e Runtime

| Tecnologia | Versão | Uso |
|---|---|---|
| Go | 1.25+ | Linguagem principal |

## Framework e Bibliotecas Core

| Biblioteca | Uso |
|---|---|
| `github.com/gin-gonic/gin` | HTTP framework (router, middleware, binding) |
| `go.mongodb.org/mongo-driver` | Driver oficial MongoDB |
| `github.com/golang-jwt/jwt/v5` | Validação e parsing de JWT |
| `github.com/joho/godotenv` | Carrega `.env` em variáveis de ambiente |
| `log/slog` (stdlib) | Logging estruturado (JSON/text) |

## Documentação de API

| Biblioteca | Uso |
|---|---|
| `github.com/swaggo/swag` | Geração de OpenAPI/Swagger via anotações |
| `github.com/swaggo/gin-swagger` | Endpoint Swagger UI no Gin |
| `github.com/swaggo/files` | Assets estáticos do Swagger |

## Export

| Biblioteca | Uso |
|---|---|
| `github.com/go-pdf/fpdf` | Geração de PDF |
| `golang.org/x/image` | Processamento de imagens (JPEG) |

## Testes

| Biblioteca | Uso |
|---|---|
| `github.com/stretchr/testify` | Asserts e mocks (testify/mock) |
| `github.com/testcontainers/testcontainers-go` | Testes de integração com MongoDB real |
| `go.uber.org/mock` | Geração de mocks (mockgen) |

## Infraestrutura

| Tecnologia | Uso |
|---|---|
| Docker | Containerização (multi-stage build) |
| Docker Compose | Orquestração local (API + MongoDB) |
| MongoDB | Banco de dados principal |

## Comandos do Projeto

```bash
make run              # Gera Swagger + inicia servidor
make build            # Compila binário em bin/api
make test             # Testes unitários
make test-integration # Testes de integração (requer Docker)
make swagger          # Regenera docs Swagger
make docker-build     # Builda imagem Docker
make docker-up        # Sobe API + MongoDB
make docker-down      # Para os containers
make clean            # Remove binários
```

## Variáveis de Ambiente

| Variável | Default | Descrição |
|---|---|---|
| `SERVER_PORT` | `8080` | Porta do servidor |
| `CORS_ALLOWED_ORIGINS` | vazio | Origins permitidas (CSV) |
| `LOG_LEVEL` | `info` | `debug\|info\|warn\|error` |
| `LOG_FORMAT` | `json` | `json\|text` |
| `LOG_ADD_SOURCE` | `false` | Inclui arquivo/linha no log |
| `AUTH_ENABLED` | `false` | Habilita validação JWT |
| `AUTH_JWT_SECRET` | vazio | Segredo JWT (obrigatório se habilitado) |
| `AUTH_JWT_ISSUER` | vazio | Issuer JWT |
| `AUTH_JWT_AUDIENCE` | vazio | Audience JWT |
| `DB_CONNECTIONS` | `primary` | Nomes das conexões (CSV) |
| `DB_PRIMARY_KIND` | `mongodb` | Tipo do banco |
| `DB_PRIMARY_URI` | `mongodb://localhost:27017` | URI de conexão |
| `DB_PRIMARY_DATABASE` | `appdb` | Nome do database |

## Build Tags

- Sem tag → apenas testes unitários (sem Docker)
- `-tags=integration` → inclui testes de integração (Testcontainers)
