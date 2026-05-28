# Client Core

API para gerenciar clientes e simular integração com Pipefy.

## Execução local

**Pré-requisitos:** Go 1.21+.

```bash
# Copia o arquivo de ambiente
cp .env.example .env

# Sobe a API (porta 8080)
go run cmd/main.go
```

Na primeira execução a aplicação cria o arquivo `clientcore.db` com as tabelas `clients` e `processed_events`. 

## Testes

```bash
go test ./...
```
Os testes usam SQLite em memória com `modernc.org/sqlite`, isolados entre si. Não tocam no `clientcore.db` nem dependem de arquivo externo.

## Endpoints

- Criar cliente: `POST /clientes`

```bash
curl -X POST http://localhost:8080/clientes \
	-H 'Content-Type: application/json' \
	-d '{"cliente_nome":"João Silva","cliente_email":"joao.silva@email.com","tipo_solicitacao":"Atualização cadastral","valor_patrimonio":250000}'
```

- Webhook: `POST /webhooks/pipefy/card-updated`

```bash
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
	-H 'Content-Type: application/json' \
	-d '{"event_id":"evt_123","card_id":"card_456","cliente_email":"joao.silva@email.com","timestamp":"2026-05-18T12:00:00Z"}'
```

- Consultar cliente: `GET /clientes/:email`

```bash
curl http://localhost:8080/clientes/joao.silva@email.com
```

## GraphQL (simulado)

As mutations do Pipefy ficam em `internal/pipefy/client.go`:

- `createCard` — cria card com os dados do cliente
- `updateFieldsValues` — atualiza status e prioridade

Elas são montadas como string e retornadas na resposta não é feita chamada externa de fato.

## Nota sobre `tipo_solicitacao`

O campo é recebido na criação do cliente, mas não vai pro banco. Ele só é usado na mutation `createCard` do Pipefy. Como não participa de nenhuma regra de negócio, optei por não persistir. A adaptação pro formato do Pipefy fica isolada na camada de integração.

## Visão de Produção

Se fosse pra produção, a estrutura atual se adaptaria bem porque as camadas já são separadas. A ideia é trocar o que é local por serviços gerenciados, sem mexer na lógica.

**O que subiria:**

- **API Gateway** no lugar do Gin como porta de entrada. Receberia o HTTP, validaria autenticação e repassaria pro Lambda. Custaria por requisição e escalaria automaticamente.

- **Lambda** rodaria o código Go que hoje está no `cmd/main.go`. Cada requisição viraria uma execução e a AWS alocaria, rodaria, desligaria. Não pagaria por máquina ociosa. Se ninguém chamasse a API por uma hora, não geraria custo. Se houvesse pico, escalaria sozinho.

- **RDS (PostgreSQL)** ou **DynamoDB** no lugar do SQLite. Com RDS o backup, restauração e failover seriam simplificados. O DynamoDB escalaria horizontal sem muito esforço. Pra esse caso, ambos serviriam e eu usaria o RDS se a prioridade fosse SQL com joins e o DynamoDB se a prioridade fosse máxima escalabilidade.

**O que ficaria como está:**

- A tabela `processed_events` com `event_id` como chave já resolveria idempotência. Em produção o banco garantiria que mesmo que o webhook chegasse duas vezes, o segundo processamento seria ignorado.

**Pra volume alto:**

- **SQS** entraria como buffer. O API Gateway jogaria o evento na fila e voltaria 200 pro Pipefy na hora. O Lambda processaria depois no ritmo que aguentasse. Sem SQS, se o Lambda demorasse e chegassem mais requisições do que ele conseguisse atender, começaria a dar timeout.

- **CloudWatch** ou **Datadog** coletaria logs e métricas como tempo de execução, erros e volume de requisições. Daria pra configurar alarmes para acionar quando algo saísse do normal. Com Datadog também seria possível fazer dashboards para visualizar a saúde da aplicação de forma centralizada.
