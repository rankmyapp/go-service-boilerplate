# Política de Bibliotecas e Dependências

## Antes de Sugerir uma Nova Biblioteca

O agente deve verificar **todos** os itens abaixo antes de recomendar qualquer nova dependência:

- [ ] O projeto já possui solução equivalente?
- [ ] A stdlib do Go já oferece este recurso nativamente?
- [ ] A biblioteca possui manutenção ativa (commits recentes, issues respondidas)?
- [ ] A licença é compatível com uso comercial? (MIT, Apache 2.0, BSD são aceitas)
- [ ] Há vulnerabilidades críticas conhecidas (CVE)?
- [ ] O ganho real compensa o aumento de complexidade?
- [ ] A biblioteca processa ou envia dados para terceiros?
- [ ] A inclusão exige aprovação humana?

Se qualquer item resultar em dúvida, **não adicionar** e apresentar para aprovação humana.

## Checklist de Avaliação de Licença

| Licença | Status |
|---|---|
| MIT | Aprovada |
| Apache 2.0 | Aprovada |
| BSD (2-clause, 3-clause) | Aprovada |
| ISC | Aprovada |
| GPL (qualquer versão) | **Requer análise jurídica** |
| AGPL | **Requer análise jurídica** |
| LGPL | **Requer análise jurídica** |
| Sem licença | **Proibida** |
| Licença proprietária | **Requer aprovação** |

## O que Evitar

| Padrão | Motivo |
|---|---|
| Biblioteca abandonada (último commit > 2 anos) | Risco de CVEs sem patch |
| Pacote sem licença clara | Risco jurídico |
| Dependência desnecessária (stdlib resolve) | Aumenta superfície de ataque |
| Pacote com CVE crítico não patcheado | Risco direto de segurança |
| Lib que dificulta testes (singletons globais, init()) | Compromete testabilidade |
| Dependências com muitas dependências transitivas | Aumenta risco de supply chain |

## Dependências Diretas Atuais (go.mod)

As seguintes dependências são aprovadas e fazem parte da stack oficial:

```
github.com/gin-gonic/gin         — HTTP framework
github.com/go-pdf/fpdf           — PDF generation
github.com/golang-jwt/jwt/v5     — JWT validation
github.com/joho/godotenv         — .env loader
github.com/stretchr/testify      — test assertions and mocks
github.com/swaggo/files          — Swagger static assets
github.com/swaggo/gin-swagger    — Swagger UI for Gin
github.com/swaggo/swag           — OpenAPI generation
github.com/testcontainers/...    — integration test containers
go.mongodb.org/mongo-driver      — MongoDB driver
golang.org/x/image               — image processing
go.uber.org/mock                 — mock generation
```

## Processo para Adicionar Nova Dependência

1. Verificar o checklist acima
2. Verificar se não existe solução já no projeto
3. Documentar justificativa no PR
4. Obter aprovação do Tech Lead antes de adicionar ao `go.mod`
5. Verificar vulnerabilidades: `go list -m -json all | nancy` ou equivalente

## Atualização de Dependências

- Atualizações de patch (ex: `v1.2.3` → `v1.2.4`): podem ser feitas com análise de changelog
- Atualizações de minor (ex: `v1.2.x` → `v1.3.x`): verificar breaking changes, testar
- Atualizações de major (ex: `v1.x` → `v2.x`): exigem aprovação humana e plano de migração
