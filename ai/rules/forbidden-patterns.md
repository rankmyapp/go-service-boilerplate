# Padrões Proibidos

Estes padrões nunca devem ser introduzidos sem aprovação humana explícita.

## Código

### Erros ignorados
```go
// PROIBIDO
result, _ := someFunction()

// CORRETO
result, err := someFunction()
if err != nil {
    return fmt.Errorf("context: %w", err)
}
```

### Contexto ignorado
```go
// PROIBIDO — perde cancelamento e timeout
repo.Create(context.Background(), user)

// CORRETO
repo.Create(c.Request.Context(), user)
```

### Lógica de negócio no handler
```go
// PROIBIDO — regra de negócio pertence ao usecase
func (h *UserHandler) CreateUser(c *gin.Context) {
    if user.Email == "" {
        user.Email = "default@example.com" // REGRA DE NEGÓCIO AQUI = ERRADO
    }
    // ...
}
```

### Imports de Gin no usecase
```go
// PROIBIDO — usecase não pode conhecer Gin
package usecase

import "github.com/gin-gonic/gin" // NUNCA
```

### Stack trace em resposta HTTP
```go
// PROIBIDO
c.JSON(500, gin.H{"error": err.Error(), "trace": debug.Stack()})

// CORRETO
slog.Error("unexpected error", "error", err)
c.JSON(500, gin.H{"error": "internal server error"})
```

### Credenciais hardcoded
```go
// PROIBIDO
const jwtSecret = "minha-senha-super-secreta"

// CORRETO
secret := os.Getenv("AUTH_JWT_SECRET")
```

## Arquitetura

### Inversão de dependência violada
```go
// PROIBIDO — usecase importando implementação concreta
import repoMongo "github.com/user/gin-microservice-boilerplate/internal/repository/mongo"

// CORRETO — usecase usa interface
import "github.com/user/gin-microservice-boilerplate/internal/repository"
```

### Acesso direto ao banco fora do repositório
```go
// PROIBIDO — usecase acessando MongoDB diretamente
func (uc *userUsecase) CreateUser(ctx context.Context, user *models.User) (string, error) {
    result, err := uc.db.Collection("users").InsertOne(ctx, user) // ERRADO
```

## Testes

### Mock sem AssertExpectations
```go
// PROIBIDO — mock pode não ter sido chamado como esperado
func TestCreate(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
    mockRepo.On("Create", ...).Return(...)
    // ... execução ...
    assert.NoError(t, err)
    // FALTA: mockRepo.AssertExpectations(t)
}
```

### Teste de integração sem build tag
```go
// PROIBIDO — vai rodar em make test sem Docker
package mongo_test

import "github.com/testcontainers/testcontainers-go" // sem build tag = erro em CI

// CORRETO
//go:build integration

package mongo_test
```

## Dependências e Infraestrutura

### Nova dependência sem justificativa
Adicionar entrada em `go.mod` sem consultar `/ai/standards/libraries-policy.md` e sem aprovação humana.

### Novo banco de dados sem aprovação
Adicionar provider para banco não catalogado em `/ai/standards/database-policy.md`.

### Modificar `.env.example` com valores reais
```bash
# PROIBIDO — nunca valores reais no .env.example
AUTH_JWT_SECRET=minha-senha-real-aqui

# CORRETO
AUTH_JWT_SECRET=your-secret-here
```

## Git e Processo

### Push direto em `main`/`master`
O agente não realiza push em branches protegidas.

### Commit sem revisão humana
O agente não commita código sem que o humano tenha revisado o diff.

### Modificar CI/CD sem aprovação
Alterações em `Dockerfile`, `docker-compose.yaml`, `Makefile` ou pipelines exigem aprovação humana.
