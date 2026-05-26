# Workflow: Desenvolvimento de Feature

Fluxo completo end-to-end para implementação de uma nova feature com agente de IA.

## Pré-requisitos

- Ticket criado e detalhado no sistema de gestão (Jira/Linear)
- Aceitação do ticket revisada e aprovada pelo time
- Branch `master` atualizada localmente

## Etapas

### 1. Carregar Contexto
```
Leia antes de iniciar qualquer tarefa:
- CLAUDE.md (entrypoint)
- /ai/context/architecture.md
- /ai/context/domains.md
- /ai/context/tech-stack.md
```

### 2. Carregar Padrões Downstream
```
Consulte antes de propor qualquer decisão técnica:
- /ai/standards/allowed-technologies.md
- /ai/standards/restricted-technologies.md
- /ai/standards/libraries-policy.md
- /ai/standards/database-policy.md
- /ai/standards/architecture-policy.md
- /ai/standards/api-policy.md
- /ai/standards/security-iso27001.md
```

### 3. Criar Branch
```bash
git checkout master && git pull
git checkout -b feat/PROJ-123-nome-da-feature
```
Convenção: `feat/PROJ-NNN-descricao-curta`

### 4. Planejar Implementação
Usar o prompt `/ai/prompts/plan-from-ticket.md`.

Entregar ao humano:
- Lista de arquivos afetados (criar/modificar/ler)
- Ordem de implementação
- Riscos e dependências identificados
- Testes planejados
- Dúvidas que bloqueiam implementação

**Aguardar aprovação humana antes de continuar.**

### 5. Validar Aderência Técnica
Antes de implementar, verificar se o plano respeita:
- [ ] Catálogo de tecnologias (`/ai/standards/allowed-technologies.md`)
- [ ] Padrões de arquitetura (`/ai/standards/architecture-policy.md`)
- [ ] Política de dependências (`/ai/standards/libraries-policy.md`)
- [ ] Política de banco (`/ai/standards/database-policy.md`)
- [ ] Segurança (`/ai/standards/security-iso27001.md`)

### 6. Implementar
Usar o prompt `/ai/prompts/implement.md`.

Seguir a ordem Clean Architecture:
1. Model
2. Interface do repositório
3. Implementação MongoDB
4. Usecase
5. Handler
6. Wiring no `main.go`

### 7. Escrever Testes
- Testes unitários do usecase (mock do repository)
- Testes unitários do handler (mock do usecase)
- Testes de integração do repository (`//go:build integration`)

Rodar e confirmar:
```bash
make test
make test-integration  # se aplicável
```

### 8. Regenerar Swagger
```bash
make swagger
```

### 9. Revisão de Código
- Agente faz auto-revisão usando `/ai/prompts/review-code.md`
- Humano revisa o diff completo

**Nenhum commit sem revisão humana.**

### 10. Validar Segurança
Checar antes do PR:
- [ ] Nenhum secret no código ou testes
- [ ] Logs sem dados sensíveis
- [ ] Respostas de erro sem stack trace
- [ ] Novas dependências aprovadas

### 11. Commit e PR
Após aprovação humana:
```bash
git add <arquivos específicos>
git commit -m "feat(PROJ-123): descrição do que foi implementado"
git push -u origin feat/PROJ-123-nome-da-feature
```

Descrição do PR deve incluir:
- Link para o ticket
- O que foi implementado
- Como testar
- Checklist de segurança

## Critérios de Conclusão

- [ ] `make test` passa
- [ ] `make test-integration` passa (se aplicável)
- [ ] Swagger atualizado
- [ ] Revisão humana aprovada
- [ ] Nenhum secret exposto
- [ ] PR criado com descrição completa
