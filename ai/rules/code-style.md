# Estilo de Código

## Go — Convenções Gerais

### Nomenclatura
- Nomes em **inglês** sempre
- Tipos e funções exportadas: `PascalCase`
- Variáveis e funções privadas: `camelCase`
- Constantes: `PascalCase` (exported) ou `camelCase` (private)
- Acrônimos em maiúsculas: `userID`, `httpClient`, `JSONResponse`
- Evitar abreviações não-óbvias: prefira `userRepository` a `userRepo` em assinaturas públicas

### Estrutura de Arquivo
- Um tipo principal por arquivo (ex: `user_handler.go` contém `UserHandler`)
- Nome do arquivo: `snake_case.go`
- Testes: `snake_case_test.go` no mesmo pacote

### Erros
- Sempre verificar e propagar erros — nunca ignorar com `_`
- Usar `fmt.Errorf("context: %w", err)` para wrapping
- Não logar e retornar ao mesmo tempo — escolha um
- Mensagens de erro em minúsculo, sem ponto final: `"user not found"`, não `"User not found."`

### Comentários
- Comentar apenas o **porquê**, nunca o **o quê**
- Funções exportadas recebem doc comment (`// FuncName does X`)
- Não comentar código morto — delete
- Sem comentários de TODO sem issue associada

## Padrões por Camada

### Handlers (`internal/handlers/`)
- Recebe `*gin.Context`, sem lógica de negócio
- Sempre usar `c.Request.Context()` nas chamadas downstream
- Binding: `ShouldBindJSON` para body, `ShouldBindQuery` para query params
- Respostas de erro: `gin.H{"error": "mensagem amigável"}` — sem stack trace
- Anotações Swagger obrigatórias em todos os handlers públicos:

```go
// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with name and email
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User object"
// @Security     BearerAuth
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
```

### Usecases (`internal/usecase/`)
- Sem importação de `gin` ou `net/http`
- Recebe e retorna tipos de `models/` ou primitivos
- Interface definida no mesmo arquivo: `UserUsecase interface { ... }`
- Struct interna: `type userUsecase struct { repo repository.UserRepository }`
- Construtor retorna a interface: `func NewUserUsecase(...) UserUsecase`

### Repositories (`internal/repository/`)
- Interface em `internal/repository/nome_repository.go`
- Implementação em `internal/repository/mongo/nome_repo_mongo.go`
- Construtor retorna a interface: `func NewUserRepo(db *mongo.Database) repository.UserRepository`
- Sempre usar `ctx` em todas as operações do driver MongoDB

### Models (`models/`)
- Structs puras sem lógica
- Tags obrigatórias para MongoDB: `bson:"field_name"`
- Tags para JSON: `json:"field_name"`
- Campos opcionais: `bson:"field,omitempty"` e `json:"field,omitempty"`
- IDs: `bson:"_id,omitempty"` e `json:"id"`
- Campos sensíveis: `json:"-"` para não serializar

## Testes

### Unitários
- Arquivo: `*_test.go` no mesmo pacote
- Mock via `testify/mock` — arquivo de mock gerado em `mocks/`
- Pattern obrigatório:
  1. Cria mock da dependência
  2. Injeta via construtor
  3. Configura expectativas com `.On(...).Return(...)`
  4. Executa
  5. Asserta resultado
  6. `mock.AssertExpectations(t)` ao final

### Integração
- Build tag obrigatória na primeira linha: `//go:build integration`
- Linha em branco após a build tag
- Usa Testcontainers para subir MongoDB
- Função helper `setupMongoContainer(t, ctx)` retorna container e database

## Formatação

- `gofmt` ou `goimports` antes de commitar
- Sem trailing whitespace
- Uma linha em branco entre funções
- Imports agrupados: stdlib → externos → internos
