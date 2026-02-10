# Go Gin Clean Architecture Boilerplate

Base para APIs em Go usando **Gin**, **Clean Architecture**, **MongoDB** e **Swagger**, com testes unitarios e de integracao via **Testcontainers**.

---

## Recursos

- HTTP API com [Gin](https://github.com/gin-gonic/gin)
- Clean Architecture pragmatica (Handlers -> Usecases -> Repositories)
- Suporte a multiplos bancos de dados via ConnectionManager
- Swagger (OpenAPI) gerado automaticamente a partir de anotacoes
- Configuracao via `.env` (variaveis de ambiente)
- Testes unitarios com mocks (testify)
- Testes de integracao com MongoDB real via [Testcontainers](https://github.com/testcontainers/testcontainers-go)
- Containerizacao com Docker e Docker Compose
- Graceful shutdown

---

## Requisitos

- Go 1.25+
- Docker (para rodar via container e/ou testes de integracao)
- [swag](https://github.com/swaggo/swag) (para gerar docs do Swagger)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

---

## Quick Start

### Rodando localmente

```bash
# 1. Clone o repositorio
git clone https://github.com/user/gin-microservice-boilerplate.git
cd gin-microservice-boilerplate

# 2. Copie o arquivo de configuracao
cp .env.example .env

# 3. Suba um MongoDB local (ou ajuste a URI no .env)
docker run -d -p 27017:27017 mongo:7

# 4. Rode a aplicacao
make run
```

A API estara disponivel em `http://localhost:8080` e o Swagger em `http://localhost:8080/swagger/index.html`.

### Rodando via Docker Compose

```bash
# Sobe a API + MongoDB com um unico comando
make docker-up

# Para parar
make docker-down
```

O Docker Compose cuida de tudo: builda a imagem da API, sobe o MongoDB com healthcheck e conecta os dois.

---

## Configuracao

A aplicacao usa variaveis de ambiente carregadas via `.env`. O arquivo `.env.example` serve como referencia:

| Variavel              | Default                       | Descricao            |
|-----------------------|-------------------------------|----------------------|
| `SERVER_PORT`         | `8080`                        | Porta do servidor    |
| `DB_PRIMARY_KIND`     | `mongodb`                     | Tipo do banco        |
| `DB_PRIMARY_URI`      | `mongodb://localhost:27017`   | URI de conexao       |
| `DB_PRIMARY_DATABASE` | `appdb`                       | Nome do database     |

O `.env` esta no `.gitignore` para nao expor credenciais. Sempre commite apenas o `.env.example`.

---

## Endpoints da API

| Metodo   | Rota                  | Descricao         |
|----------|-----------------------|--------------------|
| `POST`   | `/api/v1/users`       | Criar usuario      |
| `GET`    | `/api/v1/users`       | Listar todos       |
| `GET`    | `/api/v1/users/:id`   | Buscar por ID      |
| `PUT`    | `/api/v1/users/:id`   | Atualizar usuario  |
| `DELETE` | `/api/v1/users/:id`   | Deletar usuario    |
| `GET`    | `/swagger/index.html` | Swagger UI         |

---

## Comandos (Makefile)

| Comando                | Descricao                                        |
|------------------------|--------------------------------------------------|
| `make run`             | Gera swagger e inicia o servidor                 |
| `make build`           | Compila binario em `bin/api`                     |
| `make test`            | Roda testes unitarios                            |
| `make test-integration`| Roda testes de integracao (precisa de Docker)     |
| `make swagger`         | Regenera documentacao Swagger                    |
| `make docker-build`    | Builda imagem Docker                             |
| `make docker-up`       | Sobe API + MongoDB via Docker Compose            |
| `make docker-down`     | Para os containers                               |
| `make clean`           | Remove binarios                                  |

---

## Estrutura de Pastas

```
.
├── cmd/
│   └── api/
│       └── main.go                    # Entry point e wiring de dependencias
├── config/
│   └── config.go                      # Carrega configuracao do .env
├── docs/                              # Swagger gerado (docs.go, swagger.json, swagger.yaml)
├── internal/
│   ├── handlers/
│   │   ├── user_handler.go            # Handlers HTTP (traduz HTTP <-> DTO)
│   │   └── user_handler_test.go       # Testes unitarios dos handlers
│   ├── usecase/
│   │   ├── user_usecase.go            # Regras de negocio
│   │   └── user_usecase_test.go       # Testes unitarios dos usecases
│   └── repository/
│       ├── user_repository.go         # Interface do repositorio (contrato)
│       └── mongo/
│           ├── user_repo_mongo.go     # Implementacao MongoDB
│           └── user_repo_mongo_it_test.go  # Testes de integracao
├── models/
│   └── user.go                        # Entidades de dominio
├── pkg/
│   ├── db/
│   │   ├── manager.go                 # ConnectionManager generico (multi-db)
│   │   └── mongo/
│   │       └── provider.go            # Provider MongoDB
│   └── web/
│       └── server.go                  # Setup do Gin + Swagger
├── .env.example                       # Exemplo de variaveis de ambiente
├── .gitignore
├── .dockerignore
├── Dockerfile                         # Multi-stage build
├── docker-compose.yaml                # API + MongoDB
├── Makefile
├── go.mod
└── go.sum
```

### O que cada camada faz

| Camada                  | Pasta                      | Responsabilidade                                      |
|-------------------------|----------------------------|-------------------------------------------------------|
| **Entry Point**         | `cmd/api/`                 | Inicializa a aplicacao, faz o wiring de dependencias   |
| **Config**              | `config/`                  | Carrega variaveis de ambiente                          |
| **Handlers**            | `internal/handlers/`       | Traduz HTTP para chamadas de usecase, sem regra de negocio |
| **Usecases**            | `internal/usecase/`        | Contem a logica de negocio, orquestra repositorios     |
| **Repository Interface**| `internal/repository/`     | Define contratos (interfaces) de acesso a dados        |
| **Repository Impl**     | `internal/repository/mongo/`| Implementacao concreta para MongoDB                  |
| **Models**              | `models/`                  | Entidades de dominio, compartilhadas entre camadas     |
| **DB Manager**          | `pkg/db/`                  | Gerencia conexoes com multiplos bancos                 |
| **DB Provider**         | `pkg/db/mongo/`            | Abre/fecha conexao com MongoDB                         |
| **Web**                 | `pkg/web/`                 | Configura o Gin e Swagger                              |

---

## Arquitetura

```
Request HTTP
    │
    ▼
┌──────────┐     ┌──────────┐     ┌──────────────┐     ┌─────────┐
│  Handler  │────>│  Usecase │────>│  Repository   │────>│ MongoDB │
│  (Gin)    │     │ (Logica) │     │  (Interface)  │     │         │
└──────────┘     └──────────┘     └──────────────┘     └─────────┘
```

**Regra de dependencia:** cada camada so conhece a camada imediatamente abaixo, e sempre via interface. O Handler nao sabe qual banco existe. O Usecase nao sabe que usa Gin.

---

## Testes

### Testes Unitarios

Rodam sem banco de dados real. Cada camada testa de forma isolada usando mocks (testify/mock).

```bash
make test
```

**O que e testado:**
- **Usecases** (`internal/usecase/user_usecase_test.go`): mock do repositorio, valida logica de negocio
- **Handlers** (`internal/handlers/user_handler_test.go`): mock do usecase, valida status codes e respostas HTTP usando `httptest.NewRecorder`

**Padrao de cada teste:**
1. Cria um mock da dependencia
2. Injeta no componente via construtor
3. Configura expectativas: `mock.On("Method", args).Return(values)`
4. Executa a chamada
5. Valida com `assert` e `mock.AssertExpectations`

### Testes de Integracao

Usam Testcontainers para subir um MongoDB real dentro do Docker, por demanda. Nao precisam de `docker-compose` -- basta ter o Docker rodando.

```bash
make test-integration
```

**O que e testado:**
- **Repositorios MongoDB** (`internal/repository/mongo/user_repo_mongo_it_test.go`): CRUD real contra um container MongoDB

Os testes de integracao usam a build tag `//go:build integration`, por isso **nao rodam** com `go test ./...` (apenas unitarios). Para roda-los, use `-tags=integration`.

---

## Como usar este boilerplate em outro projeto

### 1. Clone e renomeie o modulo

```bash
git clone https://github.com/user/gin-microservice-boilerplate.git meu-servico
cd meu-servico
rm -rf .git
git init
```

Altere o module path no `go.mod`:

```
module github.com/sua-org/meu-servico
```

Depois faca um find-and-replace em todos os arquivos `.go`:

```bash
# Linux/Mac
grep -rl "github.com/user/gin-microservice-boilerplate" --include="*.go" | xargs sed -i 's|github.com/user/gin-microservice-boilerplate|github.com/sua-org/meu-servico|g'

# Windows (PowerShell)
Get-ChildItem -Recurse -Filter *.go | ForEach-Object { (Get-Content $_.FullName) -replace 'github.com/user/gin-microservice-boilerplate','github.com/sua-org/meu-servico' | Set-Content $_.FullName }
```

### 2. Configure o `.env`

```bash
cp .env.example .env
# edite o .env com suas credenciais
```

### 3. Rode

```bash
make run
# ou via Docker
make docker-up
```

---

## Como adicionar novas features

### Adicionar uma nova entidade (ex: `Product`)

Siga estes passos na ordem. A ideia e ir da camada mais interna para a mais externa.

#### 1. Criar o model

```go
// models/product.go
package models

type Product struct {
    ID    string  `json:"id" bson:"_id,omitempty"`
    Name  string  `json:"name" bson:"name"`
    Price float64 `json:"price" bson:"price"`
}
```

#### 2. Criar a interface do repositorio

```go
// internal/repository/product_repository.go
package repository

import (
    "context"
    "github.com/sua-org/meu-servico/models"
)

type ProductRepository interface {
    Create(ctx context.Context, product *models.Product) (string, error)
    GetByID(ctx context.Context, id string) (*models.Product, error)
    GetAll(ctx context.Context) ([]*models.Product, error)
}
```

#### 3. Criar a implementacao MongoDB

```go
// internal/repository/mongo/product_repo_mongo.go
package mongo

import (
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/sua-org/meu-servico/internal/repository"
)

type productRepoMongo struct {
    collection *mongo.Collection
}

func NewProductRepo(db *mongo.Database) repository.ProductRepository {
    return &productRepoMongo{collection: db.Collection("products")}
}

// Implemente os metodos da interface...
```

#### 4. Criar o usecase

```go
// internal/usecase/product_usecase.go
package usecase

import "github.com/sua-org/meu-servico/internal/repository"

type ProductUsecase interface {
    // defina os metodos
}

type productUsecase struct {
    productRepo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
    return &productUsecase{productRepo: repo}
}
```

#### 5. Criar o handler

```go
// internal/handlers/product_handler.go
package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/sua-org/meu-servico/internal/usecase"
)

type ProductHandler struct {
    usecase usecase.ProductUsecase
}

func NewProductHandler(uc usecase.ProductUsecase) *ProductHandler {
    return &ProductHandler{usecase: uc}
}

func (h *ProductHandler) RegisterRoutes(rg *gin.RouterGroup) {
    products := rg.Group("/products")
    {
        products.POST("", h.CreateProduct)
        products.GET("", h.GetAllProducts)
        products.GET("/:id", h.GetProductByID)
    }
}

// Implemente os handler methods com anotacoes Swagger...
```

#### 6. Fazer o wiring no main.go

Adicione as linhas de wiring junto das existentes em `cmd/api/main.go`:

```go
productRepo := repoMongo.NewProductRepo(primaryDB)
productUC := usecase.NewProductUsecase(productRepo)
productHandler := handlers.NewProductHandler(productUC)
productHandler.RegisterRoutes(api)
```

#### 7. Regenerar o Swagger

```bash
make swagger
```

#### 8. Criar os testes

- Teste unitario do usecase com mock do repositorio
- Teste unitario do handler com mock do usecase
- Teste de integracao do repositorio com `//go:build integration`

---

### Adicionar um novo banco de dados (ex: PostgreSQL)

O ConnectionManager suporta multiplos providers. Para adicionar Postgres:

#### 1. Criar o provider

```go
// pkg/db/postgres/provider.go
package postgres

import (
    "context"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/sua-org/meu-servico/pkg/db"
)

func Registration() db.ProviderRegistration {
    return db.ProviderRegistration{
        Open:  openPostgres,
        Close: closePostgres,
    }
}

func openPostgres(ctx context.Context, cfg map[string]string) (interface{}, error) {
    conn, err := sql.Open("postgres", cfg["uri"])
    if err != nil {
        return nil, err
    }
    return conn, conn.Ping()
}

func closePostgres(ctx context.Context, conn interface{}) error {
    return conn.(*sql.DB).Close()
}
```

#### 2. Registrar no main.go

```go
import dbPostgres "github.com/sua-org/meu-servico/pkg/db/postgres"

mgr.RegisterProvider("postgres", dbPostgres.Registration())
```

#### 3. Adicionar ao .env

```env
DB_ANALYTICS_KIND=postgres
DB_ANALYTICS_URI=postgres://user:pass@localhost:5432/analytics?sslmode=disable
DB_ANALYTICS_DATABASE=analytics
```

#### 4. Adicionar ao config.go

Inclua a nova instancia no mapa `Databases` dentro de `config.Load()`.

---

### Adicionar middleware (ex: CORS, Auth)

Middlewares sao adicionados no `pkg/web/server.go` ou diretamente em `main.go`:

```go
// pkg/web/server.go - middleware global
r := gin.Default()
r.Use(corsMiddleware())

// cmd/api/main.go - middleware por grupo
api := router.Group("/api/v1")
api.Use(authMiddleware())
```

---

## Conceitos Chave

### Clean Architecture

O codigo e organizado em camadas concentricas. A regra fundamental: **dependencias apontam para dentro**.

- **Models** (centro): entidades puras, sem dependencia de framework
- **Repository interface**: contrato de acesso a dados
- **Usecase**: regras de negocio, depende apenas de interfaces
- **Handler**: traduz HTTP, depende do usecase via interface
- **Infra** (pkg): implementacoes concretas (MongoDB, Gin, etc.)

Beneficio: trocar o banco de MongoDB para Postgres so requer uma nova implementacao em `internal/repository/postgres/` -- nenhuma outra camada muda.

### ConnectionManager

Gerencia multiplas conexoes de banco nomeadas (`primary`, `analytics`, etc.). Cada tipo de banco registra um provider (open/close). Thread-safe via `sync.RWMutex`.

### Inversao de Dependencia

Toda comunicacao entre camadas usa interfaces. Isso permite:
- Testes unitarios com mocks (sem banco real)
- Trocar implementacoes sem afetar a logica de negocio
- Wiring explicito no `main.go` (sem framework de DI)

### Build Tags para Testes

- `go test ./...` roda apenas testes unitarios (rapidos, sem Docker)
- `go test -tags=integration ./...` roda testes de integracao (Testcontainers sobe MongoDB automaticamente)

---

## Stack

| Tecnologia                  | Uso                            |
|-----------------------------|--------------------------------|
| [Gin](https://github.com/gin-gonic/gin) | HTTP framework      |
| [MongoDB Driver](https://go.mongodb.org/mongo-driver) | Acesso ao MongoDB |
| [swaggo/swag](https://github.com/swaggo/swag) | Geracao de Swagger  |
| [godotenv](https://github.com/joho/godotenv) | Carrega `.env`      |
| [testify](https://github.com/stretchr/testify) | Asserts e mocks    |
| [Testcontainers](https://github.com/testcontainers/testcontainers-go) | Testes de integracao |
