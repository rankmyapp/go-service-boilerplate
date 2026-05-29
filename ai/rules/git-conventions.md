# Convenções de Git — Commits, Branches e PRs

Padrão único de histórico para este serviço, baseado em **Conventional Commits**. Habilita changelog e versionamento semântico automatizados e mantém o histórico legível por humanos e agentes.

> Commit, push e merge **nunca** são autônomos — exigem revisão humana (Human-in-the-Loop). Não há push direto em `master`/`main`/`production`.

## Conventional Commits

### Formato

```
<tipo>(<escopo opcional>)<!>: <descrição>

<corpo opcional — explica o porquê>

<footers opcionais>
```

### Tipos

| Tipo | Uso |
|------|-----|
| `feat` | Nova funcionalidade |
| `fix` | Correção de bug |
| `docs` | Documentação (inclui Swagger/OpenAPI) |
| `style` | Formatação, sem mudança de comportamento (`gofmt`, imports) |
| `refactor` | Refatoração sem mudança de comportamento |
| `perf` | Melhoria de performance |
| `test` | Testes (unitários/integração) |
| `build` | Build, `go.mod`/dependências, Dockerfile |
| `ci` | Pipeline/CI-CD |
| `chore` | Tarefas auxiliares sem impacto em código de produção |
| `revert` | Reversão de commit anterior |

### Escopo

Opcional, entre parênteses — indica a **camada ou módulo** afetado. Use os nomes reais do projeto:

- `handler` — `internal/handlers`
- `usecase` — `internal/usecase`
- `repository` / `mongo` — `internal/repository`, `internal/repository/mongo`
- `middleware` — `internal/middleware`
- Domínios: nome da entidade afetada (ex.: `user`, `auth`)

Exemplos: `feat(handler):`, `fix(repository):`, `refactor(usecase):`.

### Regras de mensagem

- Descrição no **imperativo**, em **minúscula**, **sem ponto final**, assunto ≤ ~72 caracteres.
- O **corpo** (opcional) explica o **porquê** da mudança, não o "como".
- Linhas em branco separam assunto, corpo e footers.
- Mensagem em português (idioma do repositório).

### Breaking Changes

Sinalizar incompatibilidade com `!` após o tipo/escopo **e/ou** footer `BREAKING CHANGE:`:

```
feat(handler)!: remove campo legado do payload de resposta

BREAKING CHANGE: consumidores que liam `legacy_id` devem migrar para `id`.
```

### Vínculo com Ticket

Referencie o card (Jira/Linear) no footer:

```
fix(repository): corrige timeout em consultas longas ao MongoDB

Adiciona context com deadline na camada de repositório para evitar
conexões penduradas sob carga.

Refs: PROJ-456
```

Use `Closes: PROJ-123` quando o commit encerra o ticket.

### Exemplos completos

```
feat(usecase): adiciona validação de permissões no fluxo de criação de usuário
fix(middleware): trata header Authorization ausente sem panic
docs(api): regenera Swagger após novo endpoint de health
test(repository): cobre cenário de documento não encontrado
chore(deps): atualiza driver oficial do MongoDB
```

## Convenção de Branches

Formato: `<tipo>/<TICKET>-<nome-curto-kebab>`, usando os mesmos tipos do commit.

- `feat/PROJ-123-login-google`
- `fix/PROJ-456-timeout-mongo`
- `chore/PROJ-789-bump-go-1-25`

Branch base: `master`. Mantenha a branch curta e focada em um único ticket.

## Pull Requests

- **Título**: no padrão Conventional Commits (ex.: `feat(handler): endpoint de criação de usuário`).
- **Descrição**: derivada do ticket — o que muda, por quê, como testar, e o link para o card.
- **Revisão humana obrigatória**: nenhum merge sem aprovação. O agente atua como segundo revisor, nunca como aprovador.
- Não fazer push direto nem merge em `master`/`main`/`production`.
- Garantir que testes e lint passam antes de abrir o PR (ver `/ai/rules/code-style.md` e `/ai/standards/`).

## Proibido

Ver também `/ai/rules/forbidden-patterns.md` (seção *Git e Processo*):

- Commit fora do padrão Conventional Commits.
- Commit sem revisão ou push autônomo.
- Push/merge direto em branch protegida.
- Commit de secrets, `.env`, credenciais ou `ai/state/STATE.md`.
