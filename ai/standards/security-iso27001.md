# Segurança — Alinhamento ISO 27001

Este documento orienta decisões técnicas com base em controles de segurança inspirados na ISO 27001. Não substitui as políticas formais de Segurança/Compliance da organização.

## Checklist para Toda Especificação Técnica

Antes de propor ou implementar qualquer feature, o agente deve considerar:

- [ ] Quais dados serão processados?
- [ ] Há dados pessoais, sensíveis, financeiros ou regulados?
- [ ] Quem pode acessar esses dados?
- [ ] Como a autenticação será feita?
- [ ] Como a autorização será validada?
- [ ] Onde secrets serão armazenados?
- [ ] Haverá criptografia em trânsito?
- [ ] Haverá criptografia em repouso para dados sensíveis?
- [ ] Quais eventos precisam de auditoria?
- [ ] Existe risco de vazamento em logs?
- [ ] Existe plano de rollback?

## Controles por Tema

### Controle de Acesso
- Autenticação via JWT com validação de `exp`, `iss`, `aud`
- Autorização granular por permissões numéricas por rota
- Menor privilégio: cada endpoint solicita apenas as permissões necessárias
- Sem compartilhamento de secrets entre ambientes

### Gestão de Ativos
- Identificar e documentar: APIs, coleções MongoDB, serviços e integrações impactadas
- Classificar dados: públicos, internos, confidenciais, sensíveis
- Manter inventário de dependências atualizado (`go.mod`)

### Segurança no Desenvolvimento
- Validação de toda entrada de usuário antes de processar
- Nunca concatenar queries — usar driver BSON (proteção contra injeção)
- Code review obrigatório antes de merge
- Testes unitários e de integração como requisito de entrega
- `make test` deve passar antes de qualquer PR

### Gestão de Vulnerabilidades
- Não introduzir dependência com CVE crítico não patcheado
- Atualizações de segurança devem ser aplicadas com urgência
- Verificar: `govulncheck ./...` antes de releases

### Logs e Monitoramento
- Logs estruturados em JSON para facilitar análise (`LOG_FORMAT=json`)
- Incluir `requestID` ou `correlationID` em todas as operações
- Nível de log configurável por ambiente (`LOG_LEVEL`)
- Alertas para erros 5xx e comportamentos anômalos (configurado fora do boilerplate)
- Ver `/ai/standards/logging-observability.md` para detalhes

### Segurança em Nuvem
- Secrets em AWS Secrets Manager ou Parameter Store criptografado (produção)
- IAM com menor privilégio para a aplicação
- Redes segregadas por ambiente (dev, staging, prod)
- TLS em todas as conexões de banco em produção
- Variáveis de ambiente injetadas pelo orquestrador — não no container

### Continuidade
- Graceful shutdown implementado (SIGINT/SIGTERM com timeout de 15s)
- Plano de rollback para toda mudança de schema de banco
- Retry e timeout em operações de banco e HTTP externo

### Gestão de Incidentes
- Logs estruturados permitem rastreamento de incidentes por requestID
- Erros críticos logados com contexto suficiente para investigação
- Nunca deletar logs de produção sem política de retenção

## Dados Sensíveis — Tratamento por Tipo (quando aplicável)

O boilerplate atual processa apenas as entidades `User` (nome, e-mail) e
`Export` (artefatos gerados). Outros tipos de dado aparecem na tabela abaixo
como referência para evoluções futuras — sempre que um novo tipo for
introduzido, aplicar os controles correspondentes.

| Tipo de Dado | Controles Obrigatórios |
|---|---|
| Senhas | Hash (bcrypt/argon2), nunca armazenar em plain text |
| Tokens JWT | TTL configurado, rotação planejada |
| Chaves de API | Secrets Manager, nunca em env do container |
| Dados pessoais (CPF, email) | Verificar necessidade de criptografia em repouso |
| Dados financeiros | Criptografia em repouso + auditoria de acesso |

## Ações de Alto Risco — Aprovação Obrigatória

| Ação | Por quê |
|---|---|
| Desabilitar `AUTH_ENABLED` em produção | Remove toda autenticação |
| Adicionar campo no modelo sem validação | Risco de injeção ou dados inconsistentes |
| Expor novo endpoint sem permissão | Surface de ataque expandida |
| Alterar algoritmo de hash/JWT | Pode invalidar sessões existentes |
| Adicionar log com payload completo | Risco de exposição de dados |
| Conectar a banco de produção em dev | Risco de corrupção de dados |
