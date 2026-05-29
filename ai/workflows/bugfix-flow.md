# Workflow: Correção de Bug

Fluxo para investigação e correção de bugs com agente de IA.

## Pré-requisitos

- Bug reportado com descrição clara do comportamento esperado vs. atual
- Ambiente de reprodução identificado (local / staging / prod)
- Severidade definida pelo time

## Etapas

### 1. Carregar Contexto
```
Leia antes de iniciar:
- CLAUDE.md
- /ai/context/architecture.md
- /ai/context/domains.md
- /ai/state/STATE.template.md
```

Se a investigação exigir mais de uma rodada, copie `/ai/state/STATE.template.md` para `/ai/state/STATE.md` e registre evidências, hipóteses e próximos passos.

### 2. Criar Branch
```bash
git checkout master && git pull
git checkout -b fix/PROJ-123-descricao-do-bug
```
Convenção: `fix/PROJ-NNN-descricao-curta`

### 3. Reproduzir o Bug
- Identificar o menor input que reproduz o problema
- Escrever um teste que **falha** antes da correção (TDD)
- Confirmar a reprodução com o humano antes de corrigir
- Registrar a evidência de reprodução em `/ai/state/STATE.md`

### 4. Identificar a Causa Raiz
Investigar na camada correta:

| Sintoma | Onde investigar primeiro |
|---|---|
| Status HTTP errado | `internal/handlers/` |
| Lógica incorreta | `internal/usecase/` |
| Dado incorreto do banco | `internal/repository/mongo/` |
| Configuração errada | `config/` |
| Problema de auth | `internal/middleware/` |

### 5. Planejar a Correção

Antes de modificar código, descrever:
- Causa raiz identificada
- Arquivo(s) a serem modificados
- Risco de regressão em outras áreas
- Testes que serão adicionados/modificados

**Apresentar para aprovação humana.**

Após aprovação, registrar em `/ai/state/STATE.md` a causa raiz, plano aprovado, arquivos afetados e riscos de regressão.

### 6. Implementar a Correção

Regras:
- Corrigir apenas o que está no escopo do bug
- Não refatorar código não relacionado no mesmo commit
- Não adicionar features junto com o bugfix

### 7. Verificar que o Teste Passa
```bash
make test
```

O teste escrito no passo 3 deve passar após a correção.

Registrar em `/ai/state/STATE.md` o teste que falhava, o resultado após a correção e qualquer observação relevante.

### 8. Verificar Regressões
```bash
make test
make test-integration  # se a correção tocou repositórios
```

Todos os testes existentes devem continuar passando.

### 9. Checklist de Segurança para Bugfixes

Verificar antes do commit:
- [ ] A correção não introduziu nova vulnerabilidade
- [ ] Dados sensíveis não foram expostos na investigação (logs de debug removidos)
- [ ] Se o bug era de segurança: verificar se há outros pontos similares no código
- [ ] Nenhum secret no código

### 10. Commit
```bash
git add <arquivos específicos>
git commit -m "fix(PROJ-123): descrição clara do que foi corrigido"
```

Mensagem de commit deve explicar **o que foi corrigido e por quê era um bug**.

### 11. PR
Descrição deve incluir:
- Link para o ticket/reporte
- Causa raiz encontrada
- O que foi corrigido
- Como verificar a correção
- Testes adicionados

## Bugs de Segurança

Para bugs com impacto em segurança:
1. **Não criar issue pública** até correção deployada
2. Notificar o time de segurança imediatamente
3. Prioridade máxima independente do sprint atual
4. Post-mortem após resolução

## Critérios de Conclusão

- [ ] Teste de regressão escrito e passando
- [ ] `make test` e `make test-integration` passando
- [ ] Causa raiz documentada no PR
- [ ] Revisão humana aprovada
- [ ] Sem logs de debug no código
- [ ] `/ai/state/STATE.md` atualizado com resumo de handoff
