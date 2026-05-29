# Política de Banco de Dados

## Banco em Uso

| Banco | Uso Recomendado | Cuidados |
|---|---|---|
| MongoDB | Documentos flexíveis, dados semi-estruturados, entidades com campos variáveis | Evitar uso como banco relacional disfarçado; controlar tamanho de documentos e arrays |

MongoDB é o **único** banco em uso neste repositório. Introduzir qualquer outro
banco (PostgreSQL, Redis, MySQL, SQLite, banco vetorial, storage externo etc.)
exige justificativa técnica e aprovação humana — ver
`/ai/standards/restricted-technologies.md`.

## MongoDB — Padrões Obrigatórios

### Driver
- Usar exclusivamente o driver oficial: `go.mongodb.org/mongo-driver`
- Nunca construir queries com interpolação de strings (risco de injeção)
- Usar tipos BSON para todas as operações

### Contexto
```go
// OBRIGATÓRIO — sempre usar ctx com timeout
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()

result, err := col.InsertOne(ctx, doc)
```

### Estrutura do Repositório
```go
// Implementação sempre via interface
type userRepoMongo struct {
    collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) repository.UserRepository {
    return &userRepoMongo{collection: db.Collection("users")}
}
```

### Indexes
- Criar indexes para campos de busca frequente
- Criar index único para campos que precisam de unicidade (ex: email)
- Documentar indexes criados nos comentários do repositório

### IDs
- Usar `bson:"_id,omitempty"` para deixar MongoDB gerar o ObjectID
- Retornar o ID gerado como string para as camadas superiores

## ConnectionManager

O `pkg/db/manager.go` foi desenhado para suportar múltiplas conexões nomeadas
via `ProviderRegistration`. Hoje, apenas o provider MongoDB está registrado.

Introduzir um novo provider (qualquer banco diferente de MongoDB) exige
aprovação humana — ver `/ai/standards/restricted-technologies.md`. Após
aprovação, o procedimento envolve criar o pacote provider em `pkg/db/<tipo>/`,
registrá-lo no `cmd/api/main.go` via `mgr.RegisterProvider(...)` e definir as
variáveis de ambiente `DB_<NOME>_KIND`, `DB_<NOME>_URI`, `DB_<NOME>_DATABASE`.

## Regras de Segurança para Bancos

- URIs de conexão: **somente via variável de ambiente**, nunca hardcoded
- Usuário de banco: mínimo de privilégios necessários (não usar root/admin em produção)
- Em produção: TLS/SSL obrigatório nas conexões
- Não logar URIs de conexão nos logs de startup
- Produção deve ter proteção contra DROP, TRUNCATE e DELETE sem WHERE

## Proibido sem Aprovação

- Introduzir novo banco não listado acima
- Usar banco vetorial sem análise de privacidade dos dados indexados
- Usar banco em memória (SQLite, H2) como substituto de banco real em testes de integração
- Conectar a banco de produção a partir de ambiente de desenvolvimento
