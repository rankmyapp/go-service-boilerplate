# Prompt: Revisar Código

## Objetivo

Realizar revisão técnica do diff atual ou de um PR específico, focando em correção, segurança e aderência aos padrões do projeto.

## Instruções

Revise o código considerando os seguintes aspectos, na ordem de prioridade:

### 1. Corretude
- A lógica implementada está correta para o requisito?
- Há casos de borda não tratados?
- Erros são propagados corretamente?
- Contexto (`ctx`) é passado em todas as chamadas assíncronas?

### 2. Segurança
Consulte `/ai/rules/security.md` e verifique:
- [ ] Dados sensíveis não expostos em respostas de erro
- [ ] Input validado antes de uso
- [ ] Tokens e secrets não logados
- [ ] Authorization headers não expostos em logs
- [ ] JWT validado corretamente (exp, iss, aud)
- [ ] Queries MongoDB seguras (sem injeção)

### 3. Aderência à Arquitetura
Consulte `/ai/context/architecture.md` e verifique:
- [ ] Handler não contém lógica de negócio
- [ ] Usecase não conhece Gin ou HTTP
- [ ] Repository implementa a interface definida
- [ ] Wiring correto no `main.go`
- [ ] Dependências apontam para dentro (Clean Architecture)

### 4. Padrões de Código
Consulte `/ai/rules/code-style.md` e verifique:
- [ ] Nomenclatura em inglês, camelCase/PascalCase correto
- [ ] Erros verificados e não ignorados com `_`
- [ ] Sem comentários desnecessários (apenas WHY quando não óbvio)
- [ ] Anotações Swagger presentes nos handlers públicos
- [ ] Build tag `//go:build integration` nos testes de integração

### 5. Testes
- [ ] Testes unitários cobrem os caminhos happy path e error
- [ ] Mocks configurados com `AssertExpectations`
- [ ] Testes de integração usam `//go:build integration`
- [ ] Testes passam (`make test`)

### 6. Dependências
Consulte `/ai/standards/libraries-policy.md`:
- [ ] Nenhuma biblioteca nova foi adicionada sem justificativa
- [ ] Bibliotecas novas têm manutenção ativa e licença compatível

## Output Esperado

Para cada problema encontrado, forneça:

```
**Severidade:** [Critical / High / Medium / Low]
**Arquivo:** `internal/handlers/user_handler.go:45`
**Problema:** [descrição clara do problema]
**Sugestão:** [como corrigir]
```

Finalize com um resumo:
```
## Resumo da Revisão
- Aprovado / Aprovado com ressalvas / Bloqueado
- [N] problemas críticos
- [N] problemas altos
- [N] sugestões de melhoria
```

## Restrições

- Não fazer alterações diretas — apenas comentários e sugestões
- Problemas críticos ou altos de segurança bloqueiam o merge
- Solicitar aprovação humana para decisões de arquitetura
