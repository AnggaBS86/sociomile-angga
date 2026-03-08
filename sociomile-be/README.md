# sociomile-be

Backend service for the Sociomile, built with Go, Echo, MySQL, Redis, and Asynq.

## Overview

This project implements the main flow:

`channel -> conversation -> assignment -> reply -> escalate -> ticket`

Implemented features:

- JWT login with `admin` and `agent` roles
- Role-based route protection
- Multi-tenant isolation using `tenant_id`
- Channel webhook simulation
- Conversation list, detail, and agent reply
- Ticket escalation and ticket listing
- Admin-only ticket status update
- Redis rate limiting for `POST /channel/webhook`
- Event persistence to `activity_logs`
- Async job processing with Asynq worker
- Service-layer unit tests

## Architecture

- `cmd/api`: HTTP API entrypoint
- `cmd/worker`: Asynq worker entrypoint
- `internal/http/handler`: Echo handlers
- `internal/http/request`: request DTOs
- `internal/http/response`: success and error response helpers
- `internal/http/middleware`: auth, validation, global error handling
- `internal/service`: business logic
- `internal/service/events`: event dispatching
- `internal/service/ratelimiter`: Redis rate limiting
- `internal/repository`: MySQL repository implementations
- `internal/domain/model`: domain models
- `internal/domain/repository_interface`: repository contracts
- `internal/worker/asynq`: async job handlers
- `migrations`: SQL migrations

## Requirements

- Go `1.24+`
- MySQL `8+`
- Redis `6+`

## Environment Variables

See `.env.example`.

- `APP_PORT`
- `DB_DSN`
- `JWT_SECRET`
- `JWT_TTL_MINUTES`
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `REDIS_DB`
- `WEBHOOK_RATE_LIMIT_PER_MINUTE`
- `ASYNQ_QUEUE`
- `ASYNQ_CONCURRENCY`

Example:

```env
APP_PORT=8080
DB_DSN=root:123456@tcp(127.0.0.1:3306)/sociomile?parseTime=true
JWT_SECRET=replace-with-strong-secret
JWT_TTL_MINUTES=120
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0
WEBHOOK_RATE_LIMIT_PER_MINUTE=60
ASYNQ_QUEUE=events
ASYNQ_CONCURRENCY=10
```

## Database Migration

Create database first:

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS sociomile;"
```

Run migrations:

```bash
cd "sociomile-angga/sociomile-be"
make migrate-up MIGRATE_DATABASE_URL='mysql://root:123456@tcp(localhost:3306)/sociomile'
```

Rollback one version:

```bash
make migrate-down MIGRATE_DATABASE_URL='mysql://root:123456@tcp(localhost:3306)/sociomile'
```

## Running the App

Start API:

```bash
make run
```

OR 

```bash
air
```

Start worker in another terminal:

```bash
make worker
```

Development mode with hot reload:

```bash
make dev
```

## Seed Data

At minimum, seed:

- one record in `tenants`
- one `admin` user
- one `agent` user

Example:

```sql
INSERT INTO tenants (id, name) VALUES (1, 'Tenant A');

INSERT INTO users (tenant_id, email, password, role)
VALUES
  (1, 'admin@email.com', '123456', 'admin'),
  (1, 'angga@email.com', '123456', 'agent');
```

## Response Format

Success response:

```json
{
  "status": "OK",
  "message": "operation succeeded",
  "data": {}
}
```

Error response:

```json
{
  "status": "ERROR",
  "message": "something went wrong"
}
```

Validation error response:

```json
{
  "status": "ERROR",
  "message": "validation failed",
  "errors": [
    {
      "field": "email",
      "rule": "email",
      "message": "email must be a valid email"
    }
  ]
}
```

## Authentication

Use Bearer token for protected endpoints:

```http
Authorization: Bearer <token>
```

## API Docs

### `GET /health`

Health check endpoint.

Response:

```json
{
  "status": "OK",
  "message": "operation succeeded",
  "data": {
    "status": "ok"
  }
}
```

### `POST /auth/login`

Login using email and password.

Request:

```json
{
  "email": "angga@email.com",
  "password": "123456"
}
```

Response:

```json
{
  "status": "OK",
  "message": "operation succeeded",
  "data": {
    "token": "jwt-token",
    "user": {
      "id": 2,
      "tenant_id": 1,
      "email": "angga@email.com",
      "role": "agent"
    }
  }
}
```

### `POST /channel/webhook`

Simulates customer message from external channel.

Request:

```json
{
  "tenant_id": 1,
  "customer_external_id": "cust-001",
  "message": "Hello"
}
```

Response:

```json
{
  "status": "OK",
  "message": "operation succeeded",
  "data": {
    "conversation_id": 10,
    "status": "open"
  }
}
```

Possible errors:

- `400` invalid payload or validation failure
- `429` webhook rate limit exceeded

### `GET /conversations`

Protected roles: `admin`, `agent`

Query params:

- `status`
- `assigned_agent`
- `limit`
- `offset`

Example:

```bash
GET /conversations?status=open&assigned_agent=2&limit=20&offset=0
```

### `GET /conversations/:id`

Protected roles: `admin`, `agent`

Returns conversation detail and messages.

### `POST /conversations/:id/messages`

Protected roles: `admin`, `agent`

Request:

```json
{
  "message": "We are checking your issue"
}
```

Notes:

- If conversation has not been assigned yet, the service dispatches an async assignment event.
- Agent reply authorization is checked in service layer.

### `POST /conversations/:id/escalate`

Protected roles: `agent`

Request:

```json
{
  "title": "Escalated issue",
  "description": "Need internal follow up",
  "priority": "high"
}
```

Current behavior:

- API returns `202 Accepted`
- Ticket creation is processed asynchronously by worker

Response:

```json
{
  "status": "OK",
  "message": "event queued",
  "data": {
    "tenant_id": 1,
    "conversation_id": 10,
    "title": "Escalated issue",
    "description": "Need internal follow up",
    "status": "queued",
    "priority": "high"
  }
}
```

### `GET /tickets`

Protected roles: `admin`, `agent`

Query params:

- `status`
- `assigned_agent`
- `limit`
- `offset`

### `PATCH /tickets/:id/status`

Protected roles: `admin`

Request:

```json
{
  "status": "open"
}
```

Validation tag currently allows:

- `open`
- `in_progress`
- `resolved`
- `closed`

## Async Events

Emitted events:

- `conversation.assigned`
- `conversation.escalated`
- `ticket.created`

Flow:

1. API dispatches domain event
2. Event is persisted into `activity_logs`
3. Event is enqueued to Asynq
4. Worker consumes task and performs async processing

Worker command:

```bash
make worker
```

## Multi-Tenancy

Isolation is enforced by:

- `tenant_id` stored in JWT claims
- handlers passing `claims.TenantID` into service layer
- repositories filtering tenant-sensitive queries by `tenant_id`

Important example:

- conversation lookup by ID uses `tenant_id` and `conversation_id`
- ticket lookup by conversation uses `tenant_id` and `conversation_id`

## Testing

Run all tests:

```bash
go test ./...
```

Current unit tests focus on:

- auth service login logic
- conversation lifecycle
- assignment event dispatching
- escalation rule: one conversation cannot be escalated twice

## Trade-offs
- Passwords are plain text for this exercise; production should use bcrypt.
- Some business operations are now asynchronous, so API may return queued state before worker finishes.
- No Swagger/OpenAPI yet.
- No Docker Compose yet.
