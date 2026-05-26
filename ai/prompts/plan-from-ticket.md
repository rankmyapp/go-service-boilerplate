# Prompt: Planejar Implementação a partir de Ticket

## Instruções

Antes de escrever qualquer código, produza um plano de implementação estruturado seguindo os passos abaixo.

### 1. Entender o Escopo
- Qual é o objetivo da task?
- Quais entidades de domínio são afetadas?
- Há impacto em outras features existentes?

### 2. Listar Arquivos Afetados
Para cada arquivo, indique se será **criado**, **modificado** ou **lido**:

```
CRIAR:
- models/produto.go
- internal/repository/produto_repository.go
- internal/repository/mongo/produto_repo_mongo.go
- internal/usecase/produto_usecase.go
- internal/handlers/produto_handler.go

MODIFICAR:
- cmd/api/main.go (wiring)

LER:
- internal/handlers/user_handler.go (referência de padrão)
```

### 3. Descrever Mudanças por Arquivo
Para cada arquivo, descreva o que será implementado (não como, apenas o quê).

### 4. Identificar Riscos e Dependências
- Há dependências entre os arquivos? (ex: interface antes da implementação)
- Há risco de breaking change em código existente?
- Há necessidade de migration ou mudança de schema?

### 5. Propor Ordem de Implementação
Seguindo a Clean Architecture (de dentro para fora):
1. Model (`models/`)
2. Interface do repositório (`internal/repository/`)
3. Implementação MongoDB (`internal/repository/mongo/`)
4. Usecase (`internal/usecase/`)
5. Handler (`internal/handlers/`)
6. Wiring no `main.go`
7. Testes (unitários + integração)
8. Swagger: `make swagger`

### 6. Sugerir Testes Necessários
- Testes unitários do usecase (mock do repository)
- Testes unitários do handler (mock do usecase)
- Testes de integração do repository (Testcontainers)

## Restrições

- Siga os padrões em `/ai/rules/code-style.md`
- Considere a arquitetura em `/ai/context/architecture.md`
- Não proponha mudanças fora do escopo do ticket
- Consulte `/ai/standards/` antes de sugerir nova tecnologia, biblioteca ou banco
- Submeta o plano para aprovação humana antes de implementar

## Output Esperado

```
## Plano de Implementação: [Nome do Ticket]

### Escopo
[descrição do que será feito]

### Arquivos
[lista com ação para cada arquivo]

### Ordem de Implementação
[lista numerada]

### Riscos
[lista de riscos identificados]

### Testes
[lista de testes a criar]

### Dúvidas / Pontos de Atenção
[qualquer questão que precisa de decisão humana antes de implementar]
```
