# Logging e Observabilidade

## Stack de Logging

Este projeto usa `log/slog` (stdlib) configurado em `pkg/logging/logger.go`.

Configuração via variáveis de ambiente:
- `LOG_LEVEL` — `debug | info | warn | error` (default: `info`)
- `LOG_FORMAT` — `json | text` (default: `json`)
- `LOG_ADD_SOURCE` — `true | false` — inclui arquivo e linha (default: `false`)

## Padrões de Log

### Nível Correto para Cada Situação

| Nível | Quando usar |
|---|---|
| `Debug` | Informação de diagnóstico em desenvolvimento |
| `Info` | Eventos normais do ciclo de vida (startup, requests concluídos) |
| `Warn` | Situação inesperada mas recuperável |
| `Error` | Falha que impede a operação, requer atenção |

### Campos Obrigatórios por Contexto

**Startup:**
```go
slog.Info("server started", "port", cfg.Server.Port)
slog.Info("database connected", "name", name, "kind", dbCfg.Kind)
```

**Request (quando logado no middleware):**
```go
slog.Info("request", "method", r.Method, "path", r.URL.Path, "status", status, "duration_ms", duration)
```

**Erros:**
```go
slog.Error("failed to create user", "error", err, "requestID", requestID)
```

**Shutdown:**
```go
slog.Info("server shutting down")
slog.Info("server stopped")
```

### Campos Proibidos nos Logs

Nunca incluir nos logs:
- Senhas ou hashes de senha
- Tokens JWT (nem parcialmente)
- Valores do header `Authorization`
- Valores de cookies de sessão
- Chaves privadas ou certificados
- CPF, dados financeiros sem mascaramento
- Payloads completos de request/response sem filtro
- URIs de conexão de banco (contêm credenciais)
- Valores de variáveis de ambiente secretas

```go
// PROIBIDO
slog.Info("request headers", "headers", r.Header) // expõe Authorization

// CORRETO
slog.Info("request", "method", r.Method, "path", r.URL.Path)
```

## Correlation ID / Request ID

Para rastreabilidade, propagar um ID único por request:

```go
// No middleware (exemplo simplificado)
requestID := r.Header.Get("X-Request-ID")
if requestID == "" {
    requestID = uuid.New().String()
}
ctx := context.WithValue(r.Context(), "requestID", requestID)
```

Incluir em todos os logs de negócio:
```go
slog.Error("usecase failed", "error", err, "requestID", requestID)
```

## Structured Logging — Regras

1. **Sempre usar key-value pairs** — nunca interpolação de string:
```go
// CORRETO
slog.Info("user created", "userID", id)

// ERRADO
slog.Info(fmt.Sprintf("user %s created", id))
```

2. **Nomes de chave consistentes** — usar snake_case para chaves:
- `user_id`, `request_id`, `error`, `duration_ms`, `status_code`

3. **Sem log de sucesso em hot paths** — evitar log a cada operação de banco de dados normal

4. **Sempre logar erros no handler ou no ponto mais próximo do usuário**, não em todas as camadas:
```go
// No usecase — só retorna o erro, não loga
return nil, fmt.Errorf("create user: %w", err)

// No handler — loga com contexto
slog.Error("failed to create user", "error", err, "requestID", requestID)
```

## Observabilidade Futura (fora do escopo atual)

Quando o projeto crescer, considerar (com aprovação humana):

- **Métricas**: Prometheus + Grafana
- **Tracing**: OpenTelemetry + Jaeger/Zipkin
- **Alertas**: baseados em taxa de erro 5xx e latência P99
- **Health check endpoint**: `GET /healthz` para readiness/liveness

Qualquer adição de biblioteca de observabilidade deve passar pelo processo de `/ai/standards/libraries-policy.md`.
