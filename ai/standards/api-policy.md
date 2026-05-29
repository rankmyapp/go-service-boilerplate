# Política de API e Integração

## Padrões Obrigatórios para Novos Endpoints

Toda API ou integração proposta deve contemplar:

- [ ] Autenticação (JWT quando `AUTH_ENABLED=true`)
- [ ] Autorização (permissões via `RequirePermissions`)
- [ ] Validação de entrada (`ShouldBindJSON` / `ShouldBindQuery`)
- [ ] Contrato de request/response documentado com Swagger
- [ ] Respostas de erro padronizadas (ver abaixo)
- [ ] Paginação quando o endpoint retorna listas
- [ ] Timeout de contexto em operações de banco
- [ ] Logs seguros (sem dados sensíveis)
- [ ] Propagação de `requestID` / `correlationID`

## Contratos HTTP

### Prefixo de rota
Todos os endpoints da API REST ficam sob `/api/v1/`.

### Verbos HTTP
| Ação | Verbo | Exemplo |
|---|---|---|
| Criar | POST | `POST /api/v1/users` |
| Listar | GET | `GET /api/v1/users` |
| Buscar por ID | GET | `GET /api/v1/users/:id` |
| Atualizar | PUT | `PUT /api/v1/users/:id` |
| Deletar | DELETE | `DELETE /api/v1/users/:id` |
| Exportar | POST | `POST /api/v1/exports` |

### Status Codes Padrão
| Situação | Status |
|---|---|
| Criação bem-sucedida | `201 Created` |
| Leitura / update bem-sucedido | `200 OK` |
| Deleção bem-sucedida | `200 OK` com `{"message": "..."}` |
| Input inválido | `400 Bad Request` |
| Não autenticado | `401 Unauthorized` |
| Sem permissão | `403 Forbidden` |
| Não encontrado | `404 Not Found` |
| Erro interno | `500 Internal Server Error` |

### Formato de Resposta de Erro
```json
{
  "error": "mensagem amigável ao usuário"
}
```

Nunca expor:
- Stack traces
- Queries de banco
- Paths internos
- Versões de libs
- Detalhes de infraestrutura

## Swagger / OpenAPI

Toda função handler pública deve ter anotações completas:

```go
// FunctionName godoc
// @Summary      Resumo em uma linha
// @Description  Descrição detalhada
// @Tags         nome-do-grupo
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "ID"
// @Param        body  body      models.Type  true  "Body"
// @Security     BearerAuth
// @Success      200   {object}  models.Type
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /resource/{id} [get]
```

Regenerar docs após qualquer alteração de rota: `make swagger`

## Paginação

Endpoints que retornam listas devem implementar paginação quando o volume pode crescer:

```go
// Query params padrão
type PaginationQuery struct {
    Page  int `form:"page,default=1"  binding:"min=1"`
    Limit int `form:"limit,default=20" binding:"min=1,max=100"`
}

// Response padrão
type PaginatedResponse struct {
    Data  interface{} `json:"data"`
    Total int64       `json:"total"`
    Page  int         `json:"page"`
    Limit int         `json:"limit"`
}
```

## Rate Limiting

Rate limit não está implementado hoje no boilerplate. Quando necessário para
endpoints críticos ou públicos, avaliar primeiro uma implementação em memória.
A introdução de qualquer store externo (Redis ou equivalente) exige aprovação
humana — ver `/ai/standards/database-policy.md` e
`/ai/standards/restricted-technologies.md`.

## Versionamento de API

- Versão atual: `/api/v1`
- Nova versão major: `/api/v2` (manter v1 funcionando durante migração)
- Não quebrar contratos existentes sem período de deprecação

## Integração com Serviços Externos

O boilerplate atual **não consome APIs externas**. Esta seção define os padrões
a aplicar quando integrações externas forem introduzidas. A inclusão de uma
nova dependência externa exige aprovação humana — ver
`/ai/standards/restricted-technologies.md`.

Quando o serviço precisar consumir APIs externas:

- [ ] Timeout configurável via env
- [ ] Retry com backoff exponencial para erros transientes
- [ ] Circuit breaker para dependências críticas
- [ ] Fallback definido
- [ ] Logs de chamadas externas (sem dados sensíveis)
- [ ] Health check da dependência
- [ ] Aprovação humana para adicionar nova dependência externa
