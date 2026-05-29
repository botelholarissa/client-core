# Client Core

API for managing clients and simulating integration with Pipefy.

## Local Execution

**Requirements:** Go 1.21+.

```bash
# Copy environment file
cp .env.example .env

# Start API (port 8080)
go run cmd/main.go
```

On the first execution, the application creates the `clientcore.db` file with the `clients` and `processed_events` tables.

## Tests

```bash
go test ./...
```

Tests use in-memory SQLite with `modernc.org/sqlite`, fully isolated from each other. They do not touch `clientcore.db` or depend on external files.

## Endpoints

* Create client: `POST /clientes`

```bash
curl -X POST http://localhost:8080/clientes \
	-H 'Content-Type: application/json' \
	-d '{"cliente_nome":"João Silva","cliente_email":"joao.silva@email.com","tipo_solicitacao":"Profile update","valor_patrimonio":250000}'
```

* Webhook: `POST /webhooks/pipefy/card-updated`

```bash
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
	-H 'Content-Type: application/json' \
	-d '{"event_id":"evt_123","card_id":"card_456","cliente_email":"joao.silva@email.com","timestamp":"2026-05-18T12:00:00Z"}'
```

* Get client: `GET /clientes/:email`

```bash
curl http://localhost:8080/clientes/joao.silva@email.com
```

## Simulated GraphQL

Pipefy mutations are located in `internal/pipefy/client.go`:

* `createCard` — creates a card with client data
* `updateFieldsValues` — updates status and priority

The mutations are assembled as strings and returned in the response. No external request is actually executed.

## Note About `tipo_solicitacao`

This field is received during client creation but is not persisted in the database. It is only used in the Pipefy `createCard` mutation. Since it is not part of any business rule, I decided not to store it. Pipefy-specific formatting remains isolated within the integration layer.

## Production Perspective

If this application were deployed to production, the current structure would adapt well because the layers are already separated. The idea would be to replace local components with managed services without changing the business logic.

### What Would Change

* **API Gateway** would replace Gin as the entry point. It would receive HTTP requests, validate authentication, and forward them to Lambda. It would scale automatically and charge per request.

* **AWS Lambda** would run the Go code currently located in `cmd/main.go`. Each request would trigger an execution managed by AWS. There would be no cost for idle infrastructure, and scaling would happen automatically during traffic spikes.

* **RDS (PostgreSQL)** or **DynamoDB** would replace SQLite. With RDS, backups, restoration, and failover would be simpler. DynamoDB would provide horizontal scalability with minimal operational effort. For this use case, both would work well: I would choose RDS for relational querying and joins, or DynamoDB for maximum scalability.

### What Would Stay the Same

* The `processed_events` table with `event_id` as a key already handles idempotency. In production, the database would guarantee that duplicated webhook events are ignored after the first successful processing.

### For High Traffic Scenarios

* **SQS** would work as a buffering layer. API Gateway would enqueue the event and immediately return HTTP 200 to Pipefy. Lambda would process messages asynchronously at a sustainable rate. Without SQS, slow processing combined with traffic spikes could lead to timeouts.

* **CloudWatch** or **Datadog** would collect logs and metrics such as execution time, errors, and request volume. Alerts could be configured to notify abnormal behavior. With Datadog, centralized dashboards could also be created to monitor application health.
