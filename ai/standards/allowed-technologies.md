# Tecnologias Permitidas

Catálogo de linguagens, frameworks e ferramentas aprovados para este repositório.

## Linguagem

| Tecnologia | Status | Uso |
|---|---|---|
| Go 1.25+ | **Aprovado** | Linguagem principal do serviço |

O agente não deve sugerir troca de linguagem sem justificativa técnica e aprovação humana.

## HTTP Framework

| Tecnologia | Status | Uso |
|---|---|---|
| Gin (`github.com/gin-gonic/gin`) | **Aprovado** | Router HTTP, middleware, binding |

Não introduzir outros frameworks HTTP (Echo, Fiber, Chi) sem aprovação.

## Banco de Dados

| Tecnologia | Status | Uso |
|---|---|---|
| MongoDB (`go.mongodb.org/mongo-driver`) | **Aprovado** | Banco principal e único banco em uso |

Qualquer banco diferente de MongoDB (PostgreSQL, Redis, MySQL etc.) não está no catálogo deste repositório e exige justificativa técnica e aprovação humana antes de ser introduzido — ver `/ai/standards/database-policy.md` e `/ai/standards/restricted-technologies.md`.

## Autenticação

| Tecnologia | Status | Uso |
|---|---|---|
| JWT (`github.com/golang-jwt/jwt/v5`) | **Aprovado** | Autenticação stateless |

## Documentação de API

| Tecnologia | Status | Uso |
|---|---|---|
| swaggo/swag | **Aprovado** | Geração OpenAPI via anotações |
| swaggo/gin-swagger | **Aprovado** | Swagger UI integrado ao Gin |

## Configuração

| Tecnologia | Status | Uso |
|---|---|---|
| godotenv | **Aprovado** | Carrega `.env` em desenvolvimento |
| Variáveis de ambiente (stdlib) | **Aprovado** | Configuração principal |

## Testes

| Tecnologia | Status | Uso |
|---|---|---|
| testify | **Aprovado** | Asserts e mocks |
| Testcontainers (MongoDB) | **Aprovado** | Testes de integração |
| go.uber.org/mock | **Aprovado** | Geração de mocks (mockgen) |

## Observabilidade

| Tecnologia | Status | Uso |
|---|---|---|
| log/slog (stdlib) | **Aprovado** | Logging estruturado |

## Export / Utilitários

| Tecnologia | Status | Uso |
|---|---|---|
| go-pdf/fpdf | **Aprovado** | Geração de PDF |
| golang.org/x/image | **Aprovado** | Processamento de imagens |

## Infraestrutura

| Tecnologia | Status | Uso |
|---|---|---|
| Docker | **Aprovado** | Containerização |
| Docker Compose | **Aprovado** | Orquestração local |

## Regra para Tecnologias Não Listadas

Se uma tecnologia não está neste catálogo, o agente deve:

1. Justificar a necessidade
2. Apresentar alternativas já existentes no catálogo
3. Descrever trade-offs
4. Avaliar riscos de segurança e licença
5. Avaliar impacto operacional e em manutenção
6. **Solicitar aprovação humana antes de prosseguir**
