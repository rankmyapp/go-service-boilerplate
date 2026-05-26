# Domínios e Entidades

## Entidades de Domínio

### User (`models/user.go`)

Entidade principal do boilerplate. Representa um usuário do sistema.

```go
type User struct {
    ID    string `json:"id" bson:"_id,omitempty"`
    Name  string `json:"name" bson:"name"`
    Email string `json:"email" bson:"email"`
}
```

**Operações disponíveis:**
- `CreateUser` — cria usuário, retorna ID gerado
- `GetUserByID` — busca por ID
- `GetAllUsers` — lista todos
- `UpdateUser` — atualiza por ID
- `DeleteUser` — remove por ID

### Export (`models/export.go`)

Representa uma solicitação de exportação de dados. Suporta múltiplos formatos e fontes.

**Formatos suportados:** `CSV`, `JPEG`, `PDF`  
**Fontes suportadas:** `chart`, `table`

**Combinações implementadas:**
| Formato | Fonte | Estratégia |
|---|---|---|
| CSV | chart | `pkg/export/csv/chart_strategy.go` |
| CSV | table | `pkg/export/csv/table_strategy.go` |
| JPEG | chart | `pkg/export/jpeg/chart_strategy.go` |
| PDF | table | `pkg/export/pdf/table_strategy.go` |

## Autenticação e Permissões

### JWT Auth (`internal/middleware/jwt_auth.go`)

Middleware opcional (controlado por `AUTH_ENABLED`). Valida Bearer token no header `Authorization`.

Configurações via env:
- `AUTH_JWT_SECRET` — segredo HMAC (obrigatório quando habilitado)
- `AUTH_JWT_ISSUER` — issuer esperado
- `AUTH_JWT_AUDIENCE` — audience esperado
- `AUTH_TOKEN_COOKIE_NAME` — aceita token via cookie também

### Permissions Guard (`internal/middleware/permissions_guard.go`)

Middleware de autorização baseado em permissões numéricas. Cada rota define quais permission IDs são aceitos.

```go
users.POST("", middleware.RequirePermissions(perms.Create...), h.CreateUser)
```

Permissões são configuradas via env (`AUTH_PERMISSIONS_*`) e injetadas pelo `main.go`.

## Contratos das Interfaces

### UserRepository

```go
type UserRepository interface {
    Create(ctx context.Context, user *models.User) (string, error)
    GetByID(ctx context.Context, id string) (*models.User, error)
    GetAll(ctx context.Context) ([]*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id string) error
}
```

### ExportRepository

Interface de acesso a dados para exportações (definida em `internal/repository/export_repository.go`).

### ExportStrategy

```go
type ExportStrategy interface {
    Export(ctx context.Context, data interface{}) ([]byte, error)
}
```

## Regras de Negócio Importantes

1. **ID gerado pelo banco**: MongoDB gera o `_id`; o usecase retorna o ID para o caller
2. **Permissões são opcionais**: quando `AUTH_ENABLED=false`, todas as rotas são públicas
3. **Exportação sem storage**: o endpoint retorna o arquivo diretamente na resposta (sem persistência)
4. **Estratégia obrigatória**: combinação formato+fonte sem estratégia registrada retorna erro 400
