# Agente — Entrypoint

## Projeto

Go Gin Microservice Boilerplate. API REST com Clean Architecture, Gin, MongoDB e JWT.  
Usado como template para novos serviços Go na organização.

## Contexto

Leia os arquivos em `/ai/context/` antes de iniciar qualquer tarefa:

- `/ai/context/architecture.md` — camadas, responsabilidades e regra de dependência
- `/ai/context/domains.md` — entidades de domínio, interfaces e regras de negócio
- `/ai/context/tech-stack.md` — stack completa, dependências e variáveis de ambiente

## Regras

Siga todas as diretrizes em `/ai/rules/`:

- `/ai/rules/code-style.md` — convenções Go, padrões por camada, estrutura de testes
- `/ai/rules/security.md` — JWT, logs, erros HTTP, validação de input
- `/ai/rules/forbidden-patterns.md` — o que nunca fazer

## Padrões Técnicos Downstream

Antes de sugerir arquitetura, biblioteca, framework, banco ou integração, leia:

- `/ai/standards/allowed-technologies.md`
- `/ai/standards/restricted-technologies.md`
- `/ai/standards/libraries-policy.md`
- `/ai/standards/database-policy.md`
- `/ai/standards/architecture-policy.md`
- `/ai/standards/api-policy.md`
- `/ai/standards/security-iso27001.md`
- `/ai/standards/logging-observability.md`

Não introduza tecnologia nova sem justificar e solicitar aprovação humana.

## Prompts Disponíveis

- Planejar: `/ai/prompts/plan-from-ticket.md`
- Implementar: `/ai/prompts/implement.md`
- Revisar código: `/ai/prompts/review-code.md`

## Workflows

- Feature completa: `/ai/workflows/feature-flow.md`
- Correção de bug: `/ai/workflows/bugfix-flow.md`

## Comandos Úteis

```bash
make run              # Gera Swagger + inicia servidor
make build            # Compila binário
make test             # Testes unitários
make test-integration # Testes de integração (requer Docker)
make swagger          # Regenera docs Swagger
make docker-up        # Sobe API + MongoDB
make docker-down      # Para containers
```

## Princípios de Governança

1. **Human-in-the-Loop**: planos aprovados antes de implementar; commits revisados antes de criar
2. **Menor Privilégio**: o agente opera com o mínimo necessário para a tarefa
3. **Contexto Suficiente**: ler `/ai/context/` e `/ai/standards/` antes de qualquer decisão técnica
4. **Sem Secrets**: nunca acessar, copiar ou expor `.env`, tokens ou credenciais
5. **Sem Push Autônomo**: o agente não commita nem faz push sem revisão humana
