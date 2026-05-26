# Arquitetura do Projeto

## Visão Geral

API REST em Go com Clean Architecture pragmática. Camadas desacopladas via interfaces — trocar banco ou framework não afeta regras de negócio.

## Camadas e Responsabilidades

```
Request HTTP
    │
    ▼
┌──────────┐     ┌──────────┐     ┌──────────────┐     ┌─────────┐
│  Handler  │────>│  Usecase │────>│  Repository   │────>│   DB    │
│  (Gin)    │     │ (Lógica) │     │  (Interface)  │     │         │
└──────────┘     └──────────┘     └──────────────┘     └─────────┘
```

| Camada | Pasta | Responsabilidade |
|---|---|---|
| Entry Point | `cmd/api/` | Inicializa app e faz wiring de dependências |
| Config | `config/` | Carrega variáveis de ambiente |
| Handlers | `internal/handlers/` | Traduz HTTP ↔ DTO, sem regra de negócio |
| Usecases | `internal/usecase/` | Lógica de negócio, orquestra repositórios |
| Repository Interface | `internal/repository/` | Contratos (interfaces) de acesso a dados |
| Repository Impl | `internal/repository/mongo/` | Implementação concreta para MongoDB |
| Models | `models/` | Entidades de domínio, compartilhadas entre camadas |
| DB Manager | `pkg/db/` | Gerencia conexões com múltiplos bancos |
| Web | `pkg/web/` | Configura Gin, CORS e Swagger |
| Logging | `pkg/logging/` | Logger estruturado (slog) |
| Export | `pkg/export/` | Estratégias de exportação (CSV, JPEG, PDF) |

## Regra de Dependência

**Dependências apontam para dentro.** Cada camada só conhece a imediatamente abaixo, sempre via interface.

- `Handler` não sabe qual banco existe
- `Usecase` não sabe que usa Gin
- `Repository interface` é definida em `internal/repository/`, implementada em `internal/repository/mongo/`

## Padrões de Wiring

O `main.go` faz o wiring explícito de todas as dependências (sem framework de DI):

```go
userRepo := repoMongo.NewUserRepo(primaryDB)
userUC := usecase.NewUserUsecase(userRepo)
userHandler := handlers.NewUserHandler(userUC)
userHandler.RegisterRoutes(api, perms)
```

## Estratégia de Export

Usa Strategy Pattern para formatos de exportação. Cada formato/fonte é uma estratégia independente registrada no `main.go`:

```go
exportStrategies := map[usecase.ExportStrategyKey]usecase.ExportStrategy{
    usecase.NewExportStrategyKey(models.ExportFormatCSV, models.ExportSourceChart): csvExport.NewChartStrategy(),
    ...
}
```

## ConnectionManager

`pkg/db/manager.go` gerencia múltiplas conexões nomeadas (`primary`, `analytics`, etc.). Cada tipo de banco registra um provider. Thread-safe via `sync.RWMutex`.

## Graceful Shutdown

O servidor captura `SIGINT`/`SIGTERM` e aguarda requests em andamento antes de encerrar (timeout de 15s).
