# Regras de Segurança

## Princípios Fundamentais

1. **Nunca expor detalhes internos em respostas HTTP de erro**
2. **Nunca logar dados sensíveis**
3. **Nunca acessar ou copiar arquivos `.env`, credentials ou secrets**
4. **Menor privilégio**: solicitar apenas o mínimo necessário

## Autenticação e Autorização

### JWT
- Secret via variável de ambiente `AUTH_JWT_SECRET` — nunca hardcoded
- Validar `exp`, `iss` e `aud` quando configurados
- Token aceito via header `Authorization: Bearer <token>` ou cookie configurável
- Em falha de validação: retornar 401 sem detalhar o motivo exato

```go
// CORRETO
c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})

// ERRADO — expõe detalhes internos
c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
```

### Permissions Guard
- Permissões configuradas via env, não hardcoded
- Retornar 403 em permissão insuficiente, sem listar quais permissões são necessárias

## Logs Seguros

### O que NUNCA logar:
- Senhas e hashes de senha
- Tokens JWT (nem os primeiros caracteres)
- Authorization headers
- Cookies com sessão
- Chaves privadas, certificados
- Dados pessoais sensíveis (CPF, RG, dados financeiros)
- Query params que contenham credentials
- Payloads completos de requests sem filtro

### O que é seguro logar:
- IDs de entidades
- Request IDs / correlation IDs
- Status codes
- Timestamps
- Nomes de operações
- Erros de negócio (sem stack trace completo)

```go
// CORRETO
slog.Info("user created", "userID", id, "requestID", requestID)

// ERRADO
slog.Info("request", "headers", c.Request.Header) // expõe Authorization
slog.Info("user login", "password", req.Password)
```

## Respostas de Erro HTTP

### Nunca expor:
- Stack traces
- Queries MongoDB/SQL
- Paths de arquivos internos
- Versões de bibliotecas
- Detalhes de infraestrutura
- Nomes de tabelas ou coleções

```go
// CORRETO — handler
if err := h.usecase.CreateUser(ctx, &user); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
    return
}

// ERRADO — expõe detalhes internos
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
```

## Validação de Input

- Usar `ShouldBindJSON` — rejeita automaticamente campos inválidos
- Campos obrigatórios definidos com `binding:"required"` nos structs
- Sanitizar strings antes de usar em queries (MongoDB Driver já protege de injeção via BSON)
- Validar IDs antes de passar ao repositório (evitar queries com ID vazio)

```go
// CORRETO
var user models.User
if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
    return
}
```

## CORS

- Origins permitidas via env `CORS_ALLOWED_ORIGINS`
- Em produção: nunca usar `*` (wildcard) para origens
- Configuração em `pkg/web/server.go`

## Secrets e Variáveis de Ambiente

- Secrets **somente** via variáveis de ambiente
- `.env` no `.gitignore` — nunca commitar
- Commitar apenas `.env.example` com valores fictícios
- Em produção: usar AWS Secrets Manager, Parameter Store ou equivalente aprovado
- Nunca imprimir o valor de uma variável de ambiente nos logs de startup

## MongoDB

- Usar sempre o MongoDB Driver oficial com tipos BSON — protegido contra injeção
- Nunca construir queries com interpolação de strings
- Usar `context.WithTimeout` em operações longas
- Não expor URI de conexão nos logs

## Dados Sensíveis em Transit

- HTTPS obrigatório em produção (configurado no load balancer/ingress — fora do escopo do boilerplate)
- Tokens no header, não na query string (evitar aparecimento em logs de proxy)

## Ações que Exigem Aprovação Humana

O agente **não pode** realizar as seguintes ações sem aprovação explícita:
- Remover ou desabilitar autenticação/autorização existente
- Alterar a validação JWT (algoritmo, campos obrigatórios)
- Modificar políticas de CORS para aceitar origens não-confiáveis
- Adicionar endpoint sem proteção de permissão quando `AUTH_ENABLED=true`
- Logar qualquer dado pessoal ou sensível
