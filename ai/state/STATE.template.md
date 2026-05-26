# STATE.md

Estado vivo da tarefa atual para agentes de IA.

Use este arquivo para continuidade entre sessões, retomada após compactação de contexto e handoff entre humano e agente. Este arquivo não substitui `CLAUDE.md`, `/ai/context/`, `/ai/rules/` ou `/ai/standards/`.

## Regras de Uso

- Copie este template para `ai/state/STATE.md` ao iniciar uma tarefa longa.
- Atualize o estado sempre que houver plano aprovado, edição relevante, teste executado, bloqueio ou decisão humana.
- Registre fatos verificáveis, não suposições soltas.
- Se uma decisão precisar virar regra permanente, mova para o documento correto em `/ai/`.
- Não registre secrets, tokens, credenciais, dados sensíveis ou dumps de ambiente.
- Não commite `ai/state/STATE.md`; apenas este template deve ser versionado.

## Identificação

- Ticket/issue:
- Tipo de tarefa: feature | bugfix | refactor | docs | investigação | outro
- Branch:
- Responsável humano:
- Agente/ferramenta:
- Data de início:
- Última atualização:

## Objetivo

Descreva em uma ou duas frases o resultado esperado da tarefa.

## Escopo

### Dentro do Escopo

- 

### Fora do Escopo

- 

## Contexto Carregado

Liste apenas os documentos e arquivos realmente lidos para esta tarefa.

- 

## Plano Atual

- [ ] 

## Estado da Implementação

### Arquivos Lidos

- 

### Arquivos Alterados

- 

### Decisões Tomadas

| Decisão | Motivo | Aprovada por |
|---|---|---|
| | | |

### Suposições

| Suposição | Como validar | Status |
|---|---|---|
| | | |

### Bloqueios / Dúvidas

- 

## Verificação

| Comando/Teste | Quando executar | Resultado | Observações |
|---|---|---|---|
| `make test` | Antes de solicitar revisão | pendente | |
| `make test-integration` | Se tocar repository, MongoDB ou fluxo com Docker | pendente | |
| `make swagger` | Se alterar handlers, rotas ou contratos HTTP | pendente | |

## Evidências

Registre outputs importantes de forma resumida. Não cole logs longos.

- 

## Próximo Passo

Uma frase objetiva para retomada.

## Resumo de Handoff

Resumo curto para outro agente ou humano continuar sem depender do histórico do chat.
