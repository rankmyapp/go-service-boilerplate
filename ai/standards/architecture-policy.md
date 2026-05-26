# Política de Arquitetura

## Princípio Fundamental

**Dependências apontam para dentro.** Camadas externas conhecem as internas, nunca o contrário.

```
cmd/api/main.go (wiring)
    │
    ▼
internal/handlers/   ──► internal/usecase/   ──► internal/repository/ (interface)
                                                          │
                                                          ▼
                                               internal/repository/mongo/ (impl)
                                                          │
                                                          ▼
                                                       pkg/db/
```

## Responsabilidade de Cada Camada

### `cmd/api/main.go`
- Único ponto de wiring de dependências
- Inicialização de config, logger, DB, usecases, handlers
- Registro de routes
- Graceful shutdown
- **Sem lógica de negócio**

### `config/`
- Carrega e valida variáveis de ambiente
- Retorna struct tipada com todas as configurações
- Falha fast se configuração obrigatória estiver ausente
- **Sem acesso ao banco ou frameworks**

### `internal/handlers/`
- Traduz `*gin.Context` ↔ DTO ↔ usecase
- Validação de input (binding), sem regra de negócio
- Retorna status HTTP e corpo padronizado
- Propaga contexto via `c.Request.Context()`
- **Sem importação de `internal/repository/` ou drivers de banco**

### `internal/usecase/`
- Toda a lógica de negócio reside aqui
- Depende de interfaces de repositório (`internal/repository/`)
- Não sabe que existe Gin, HTTP ou MongoDB
- Recebe e retorna tipos de `models/`
- **Sem importação de `gin`, `net/http`, ou drivers de banco**

### `internal/repository/`
- Define apenas interfaces (contratos)
- Sem implementação concreta
- Importado por usecases

### `internal/repository/mongo/`
- Implementações concretas para MongoDB
- Importa `go.mongodb.org/mongo-driver`
- Implementa as interfaces de `internal/repository/`
- **Sem lógica de negócio**

### `models/`
- Structs puras de domínio
- Tags `json` e `bson`
- Sem importação de pacotes internos ou frameworks

### `pkg/`
- Pacotes reutilizáveis e independentes de negócio
- `pkg/db/` — ConnectionManager genérico
- `pkg/web/` — Setup do Gin e Swagger
- `pkg/logging/` — Logger estruturado
- `pkg/export/` — Estratégias de exportação

## Adicionando Nova Entidade

Siga esta ordem obrigatória:

1. `models/entidade.go` — struct com tags
2. `internal/repository/entidade_repository.go` — interface
3. `internal/repository/mongo/entidade_repo_mongo.go` — implementação
4. `internal/usecase/entidade_usecase.go` — lógica
5. `internal/handlers/entidade_handler.go` — HTTP
6. `cmd/api/main.go` — wiring
7. Testes em cada camada
8. `make swagger`

## Adicionando Novo Provider de Banco

Hoje, apenas o provider MongoDB está registrado no `ConnectionManager`.
Introduzir qualquer outro provider de banco exige aprovação humana — ver
`/ai/standards/restricted-technologies.md` e `/ai/standards/database-policy.md`.

Após aprovação formal, o padrão arquitetural a seguir é:

1. Criar `pkg/db/<tipo>/provider.go` implementando `db.ProviderRegistration`
2. Registrar em `cmd/api/main.go`: `mgr.RegisterProvider("tipo", tipoProvider.Registration())`
3. Configurar via env: `DB_<NOME>_KIND=tipo`, `DB_<NOME>_URI=...`, `DB_<NOME>_DATABASE=...`

## Padrões para Workers/Processos Assíncronos

Se o projeto evoluir para workers, cada worker deve ter:

- Message handler (recebe a mensagem)
- Validation (valida o payload)
- Business logic (orquestra usecases existentes)
- Retry policy (configurável via env)
- Idempotency (evitar processamento duplicado)
- Dead-letter handling (mensagens que falharam após retries)
- Logging e observabilidade

## O que Exige Aprovação Antes de Implementar

- Mudança de padrão de wiring (ex: introduzir DI framework)
- Mudança de padrão de erro (ex: introduzir tipos de erro customizados)
- Adicionar nova camada à arquitetura
- Remover ou fundir camadas existentes
- Mover lógica entre camadas (ex: lógica do usecase para o handler)
- Adicionar pacote em `pkg/` que depende de `internal/`
