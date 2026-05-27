# Roteiro para Gravação do Vídeo de Defesa

**Duração total: até 7 minutos**

---

## Antes de gravar: preparar o VS Code

Deixe aberto (nesta ordem):
- `internal/pipefy/client.go`
- `internal/service/webhook_service.go`
- `cmd/main.go`

E dois terminais lado a lado: um com `go run cmd/main.go` pronto, outro com os curls.

---

## 1. Abertura — 15s

"Olá, sou [seu nome] e este é meu teste técnico. Uma API em Go que gerencia clientes e simula integração com Pipefy via GraphQL."

[fala do seu background aqui, ~5s]

---

## 2. Arquitetura — 90s

Mostre a estrutura de pastas:

```
cmd/
  main.go              ← entrada, rotas
internal/
  handlers/            ← HTTP (só orquestra)
  service/             ← regras de negócio
  repository/          ← SQL
  pipefy/              ← GraphQL
  database/            ← conexão SQLite
  models/              ← DTOs
```

**Desenho (mostre com o mouse ou num slide):**

```
HTTP request
    │
    ▼
┌──────────┐      ┌──────────┐     ┌────────────┐
│ handler   │────▶│ service  │────▶│ repository │────▶ SQLite
│ (gin)     │     │ (regras) │     │ (SQL puro) │
└──────────┘      └────┬─────┘     └────────────┘
                      │
                      ▼
               ┌────────────┐
               │ pipefy/    │
               │ (GraphQL   │
               │  strings)  │
               └────────────┘
```

**Explique cada camada e JUSTIFIQUE:**

- **Handler** — só pega o JSON, valida o formato e chama o service. Zero lógica de negócio. Se trocar de framework (gin → chi, echo), só mexe aqui.
- **Service** — toda regra de negócio mora aqui. Validações, idempotência, cálculo de prioridade. Testável sem HTTP.
- **Repository** — SQL puro. Se trocar SQLite por PostgreSQL, só mexe aqui. O service não sabe que banco é.
- **Pipefy** — isolado porque o Pipefy tem formato próprio de GraphQL. O service não monta string GraphQL, só chama esse pacote.
- **Models** — structs compartilhadas entre as camadas.

"Essa separação não é overengineering. Cada camada tem um motivo claro pra existir e pode ser trocada sem quebrar as outras."

**Mostre também o diagrama de produção:**

```
         ┌─────────────┐
         │   Pipefy     │
         └──────┬───────┘
                │ webhook
                ▼
┌──────────────────────┐     ┌────────────┐
│    API Gateway       │────▶│    SQS     │
│  POST /clientes      │     │  (fila)    │
│  GET /clientes/:email│     └──────┬─────┘
│  POST /webhooks/...  │            │
└─────────┬────────────┘            │
          │                         │
          ▼                         ▼
┌──────────────────────────────────────┐
│          Lambda (Go)                 │
│  handler → service → repository      │
└──────┬───────────────────────────┬───┘
       │                           │
       ▼                           ▼
┌──────────────┐          ┌──────────────┐
│  RDS (SQL)   │          │  DynamoDB    │
│  PostgreSQL  │          │  (NoSQL)     │
└──────────────┘          └──────────────┘
       │
       ▼
┌──────────────────────┐
│  CloudWatch / Datadog │
└───────────────────────┘
```

Diga: "Em produção, cada peça vira um serviço AWS. API Gateway no lugar do Gin. Lambda no lugar do servidor fixo. RDS ou DynamoDB no lugar do SQLite. Se o volume crescer, SQS desacopla o webhook do processamento. Tudo monitorado por CloudWatch ou Datadog."

---

## 3. O ponto alto: GraphQL — 100s

Abra `internal/pipefy/client.go`.

"Diferente de integrar com REST, o Pipefy usa GraphQL. Precisei pesquisar a documentação oficial em developers.pipefy.com pra entender as mutations exatas."

### createCard 
https://developers.pipefy.com/reference/create-a-card-with-the-required-fields-fulfilled

Mostre o método `BuildCreateCardMutation`:

```go
mutation { createCard(input: {
  pipe_id: 34,
  fields_attributes: [
    {field_id: "cliente_nome", field_value: "João Silva"},
    {field_id: "cliente_email", field_value: "joao@example.com"},
    {field_id: "valor_patrimonio", field_value: "250000"},
    {field_id: "tipo_solicitacao", field_value: "Atualização cadastral"}
  ]
}) { card { id } } }
```

"A documentação do Pipefy mostra que o `createCard` recebe um `pipe_id` e um array `fields_attributes`. Mapeei cada campo do payload da API para um field_id. A string é montada como template e retornada — numa integração real, seria enviada via HTTP pro GraphQL do Pipefy."

### updateFieldsValues
https://developers.pipefy.com/reference/fields#updating-fields-values

Mostre o método `BuildUpdateFieldsValuesMutation`:

```go
mutation { updateFieldsValues(input: {
  nodeId: 456,
  values: [
    {fieldId: "status", value: "Processado"},
    {fieldId: "priority", value: "prioridade_alta"}
  ]
}) { success } }
}

"Pra atualizar múltiplos campos de uma vez, a doc recomenda `updateFieldsValues`. Descobri que ele usa `nodeId` (o ID do card) e um array `values`. Optei por ele em vez de chamar `updateCardField` duas vezes — uma chamada só, mais eficiente."

"E tem a função `escapeString`: sem ela, um nome com aspas quebra a mutation."

**Fale sobre a decisão de não persistir `tipo_solicitacao`:**

"O campo `tipo_solicitacao` é recebido e usado apenas nessa mutation do Pipefy. Não particpa de nenhuma regra depois — nem prioridade, nem webhook. Por isso não persisto. A adaptação pro formato Pipefy fica toda aqui na camada de integração."

---

## 4. Endpoints e regras de negócio — 60s

Abra `cmd/main.go`, aponte as rotas:

```go
router.POST("/clientes", clientHandler.Create)
router.POST("/webhooks/pipefy/card-updated", webhookHandler.CardUpdated)
```

Vá para `internal/service/webhook_service.go` e destaque dois pontos:

**Idempotência** — linhas 39-45:

"Verifico se o `event_id` já foi processado. Se sim, retorno mutation vazia. Isso garante que o mesmo webhook do Pipefy não processe duas vezes o mesmo evento, mesmo que chegue repetido."

**Regra de prioridade** — linhas 52-55:

"Se `valor_patrimonio >= 200000`, prioridade alta. Senão, normal. Tudo testado."

---

## 5. Testes — 60s

Rode no terminal:

```bash
go test ./... -v
```

Mostre os 10 testes passando. Aponte os três grupos:

1. **ClientService** — criação, validações (email, nome), duplicidade
2. **WebhookService** — prioridade alta, prioridade normal, idempotência, event_id vazio, cliente inexistente
3. **Helper** — banco in-memory com `modernc.org/sqlite` (pure Go, sem CGO, sem arquivo físico)

"Usei banco in-memory nos testes. Isso isola cada teste, não deixa lixo e não precisa de CGO. Em produção uso o `mattn/go-sqlite3` com arquivo."

---

## 6. Demonstração ao vivo — 60s

Mude pro terminal com o servidor rodando.

### 1. Criar cliente

```bash
curl -X POST http://localhost:8080/clientes \
  -H 'Content-Type: application/json' \
  -d '{"cliente_nome":"João Silva","cliente_email":"joao.silva@example.com","tipo_solicitacao":"Atualização cadastral","valor_patrimonio":250000}'
```

Mostre a resposta com a mutation `createCard`.

### 2. Consultar cliente (mostra o banco)

```bash
curl http://localhost:8080/clientes/joao.silva@example.com
```

Mostre o JSON: `"status": "Aguardando Análise"`, sem prioridade.

### 3. Webhook

```bash
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H 'Content-Type: application/json' \
  -d '{"event_id":"evt_123","card_id":"card_456","cliente_email":"joao.silva@example.com","timestamp":"2026-05-18T12:00:00Z"}'
```

Mostre a mutation `updateFieldsValues` com status "Processado" e prioridade "prioridade_alta".

### 4. Consultar de novo (estado mudou)

```bash
curl http://localhost:8080/clientes/joao.silva@example.com
```

Mostre: `"status": "Processado"`, `"prioridade": "prioridade_alta"`. O dado mudou — o webhook funcionou.

### Idempotência (opcional, se der tempo)

Repita o mesmo curl — a mutation volta vazia.

---

## 7. Fechamento — 15s

"Projeto completo: duas rotas, regras de negócio testadas, mutations GraphQL seguindo a doc oficial do Pipefy, banco SQLite, e separação clara entre camadas. Código disponível no repositório. Obrigado!"

---

## Dicas de gravação

- **VS Code:** fonte 17-18pt, abra os arquivos do roteiro antes de gravar
- **Terminal:** tenha os comandos copiados pra colar rápido
- **Não leia:** roteiro é pra guiar, não pra ler palavra por palavra
- **Se errar:** respira e repete a frase — edita depois
- **Áudio:** ambiente silencioso, microfone próximo

## Checklist pré-gravação

- [ ] VS Code aberto com os 3 arquivos principais
- [ ] Terminal 1: `go run cmd/main.go` rodando
- [ ] Terminal 2: pronto com os curls
- [ ] Fonte legível (17pt+)
- [ ] Áudio testado
