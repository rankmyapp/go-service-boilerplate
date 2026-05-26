# Tecnologias Restritas ou Proibidas

## Ações que Exigem Aprovação Humana Explícita

O agente **não pode** realizar as seguintes ações de forma autônoma:

| Ação | Motivo |
|---|---|
| Introduzir novo banco de dados | Impacto em infraestrutura e operação |
| Introduzir novo broker/mensageria | Impacto arquitetural significativo |
| Trocar linguagem do serviço | Impacto total no repositório |
| Trocar framework HTTP principal | Impacto em toda a camada de handlers |
| Adicionar dependência sem manutenção ativa | Risco de segurança e abandono |
| Usar biblioteca sem licença clara | Risco jurídico |
| Usar pacote com vulnerabilidade crítica conhecida | Risco de segurança |
| Usar serviço externo sem avaliação | Risco de privacidade e compliance |
| Enviar dados para fornecedor externo | Risco de vazamento de dados |
| Remover autenticação ou autorização existente | Risco de segurança crítico |
| Reduzir cobertura de testes | Redução de qualidade |
| Desabilitar lint, CI ou verificações de segurança | Risco de qualidade e segurança |
| Alterar Dockerfile, docker-compose, Makefile, CI/CD | Impacto em infraestrutura |

## Tecnologias Proibidas (sem exceção)

| Tecnologia | Motivo |
|---|---|
| Frameworks de DI automático (wire, dig, fx) | Contradiz wiring explícito no `main.go` |
| ORM completo (gorm) | Contradiz uso direto do driver MongoDB |
| `gorilla/mux` | Substituído pelo Gin já adotado |
| Qualquer lib que acesse `.env` diretamente além do godotenv | Risco de duplicidade e inconsistência |

## Padrões Proibidos (independente de tecnologia)

- Wildcard CORS (`*`) em produção
- Secrets hardcoded em qualquer arquivo Go ou de configuração
- Logs com tokens, passwords ou dados sensíveis
- `context.TODO()` ou `context.Background()` em handlers (usar `c.Request.Context()`)
- Goroutines sem controle de ciclo de vida (sem WaitGroup ou canal de shutdown)
- `panic()` fora de `main.go` (sem recover)

## Processo para Solicitar Exceção

Quando uma exceção for necessária, o agente deve apresentar:

1. **Justificativa**: por que é necessário neste caso
2. **Alternativas avaliadas**: o que foi considerado antes
3. **Riscos**: o que pode dar errado
4. **Impacto em segurança**: vulnerabilidades potenciais
5. **Impacto operacional**: o que muda na operação
6. **Impacto em manutenção**: o que fica mais difícil
7. **Plano de rollback**: como desfazer se der errado
8. **Aprovação necessária**: de quem (Tech Lead, Security, etc.)

O agente não pode tratar exceções como decisão aprovada sem confirmação humana explícita.
