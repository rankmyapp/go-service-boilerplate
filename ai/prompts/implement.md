# Prompt: Implementar Feature

## Pré-condições

Antes de implementar, confirme:
- [ ] Plano aprovado pelo humano (ver `/ai/prompts/plan-from-ticket.md`)
- [ ] Branch criada com convenção `feat/PROJ-123-nome-da-task`
- [ ] Padrões downstream consultados em `/ai/standards/`

## Instruções de Implementação

### Ordem Obrigatória (Clean Architecture — de dentro para fora)

1. **Model** — entidade pura em `models/`, sem importar pacotes internos
2. **Interface do repositório** — em `internal/repository/`, apenas o contrato
3. **Implementação MongoDB** — em `internal/repository/mongo/`, implementa a interface
4. **Usecase** — em `internal/usecase/`, lógica de negócio via interface do repositório
5. **Handler** — em `internal/handlers/`, traduz HTTP ↔ DTO, sem lógica de negócio
6. **Wiring** — adicionar ao `cmd/api/main.go`
7. **Testes** — unitários primeiro, integração depois
8. **Swagger** — rodar `make swagger` ao final

### Padrões de Código

Consulte `/ai/rules/code-style.md` para:
- Nomenclatura de funções, tipos e variáveis
- Tratamento de erros
- Estrutura dos testes
- Convenções de comentários Swagger

### Checklist de Segurança por Camada

**Handler:**
- [ ] Validar input com `ShouldBindJSON` ou `ShouldBindQuery`
- [ ] Não expor stack trace ou detalhes internos em respostas de erro
- [ ] Usar `c.Request.Context()` para propagação de contexto

**Usecase:**
- [ ] Validar regras de negócio antes de chamar o repositório
- [ ] Retornar erros semânticos (não de infraestrutura)

**Repository:**
- [ ] Usar `ctx` em todas as operações do MongoDB
- [ ] Não logar dados sensíveis

**Model:**
- [ ] Tags `bson` corretas para MongoDB
- [ ] Tags `json` para serialização HTTP
- [ ] Campos sensíveis com `json:"-"` quando necessário

### Estrutura de Testes

**Teste unitário (usecase):**
```go
func TestCreateUser(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
    uc := usecase.NewUserUsecase(mockRepo)

    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
        Return("generated-id", nil)

    id, err := uc.CreateUser(context.Background(), &models.User{Name: "Test"})

    assert.NoError(t, err)
    assert.Equal(t, "generated-id", id)
    mockRepo.AssertExpectations(t)
}
```

**Teste unitário (handler):**
```go
func TestCreateUser_Success(t *testing.T) {
    mockUC := new(mocks.UserUsecase)
    h := handlers.NewUserHandler(mockUC)

    mockUC.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
        Return("new-id", nil)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    // setup request body...

    h.CreateUser(c)

    assert.Equal(t, http.StatusCreated, w.Code)
    mockUC.AssertExpectations(t)
}
```

**Teste de integração (repository):**
```go
//go:build integration

func TestUserRepo_Create(t *testing.T) {
    ctx := context.Background()
    container, db := setupMongoContainer(t, ctx)
    defer container.Terminate(ctx)

    repo := mongo.NewUserRepo(db)
    id, err := repo.Create(ctx, &models.User{Name: "Test", Email: "test@test.com"})

    assert.NoError(t, err)
    assert.NotEmpty(t, id)
}
```

## Restrições

- Não introduzir nova dependência sem consultar `/ai/standards/libraries-policy.md`
- Não alterar arquitetura sem aprovação humana
- Não commitar — o commit é feito após revisão humana
- Sempre rodar `make test` antes de solicitar revisão
