# paddle-mock-api

A lightweight Go mock server that simulates the [Paddle Billing API](https://developer.paddle.com/api-reference/overview) for local development and E2E testing.

No dependencies. No database. Single binary.

## Quick Start

```bash
# Download from releases (macOS ARM64 example)
curl -L -o paddle-mock-api \
  https://github.com/vlah-software-house/paddle-api-mock/releases/latest/download/paddle-mock-api-darwin-arm64
chmod +x paddle-mock-api

# Run with defaults (port 8081, auth enabled, seed data loaded)
./paddle-mock-api

# Or disable auth for quick local testing
./paddle-mock-api -no-auth
```

```bash
# Build from source
go build -o paddle-mock-api ./cmd/server
./paddle-mock-api
```

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `8081` | Server port |
| `-no-auth` | `false` | Disable API key validation |
| `-no-seed` | `false` | Start with an empty store |
| `-webhook-url` | — | Register a webhook URL on startup |
| `-signing-secret` | `pdl_test_signing_secret` | Webhook signing secret |
| `-api-key` | `test_paddle_api_key` | API key for Bearer auth |

## Authentication

All `/v1/` endpoints require a `Authorization: Bearer <api_key>` header (default key: `test_paddle_api_key`). Admin and health endpoints are unauthenticated.

Disable auth entirely with `-no-auth`.

## API Endpoints

### Products & Prices (read-only, seeded)

```
GET  /v1/products
GET  /v1/products/{id}
GET  /v1/prices
GET  /v1/prices/{id}
```

### Customers

```
POST  /v1/customers
GET   /v1/customers
GET   /v1/customers/{id}
PATCH /v1/customers/{id}
```

### Subscriptions

```
POST  /v1/subscriptions
GET   /v1/subscriptions
GET   /v1/subscriptions/{id}
PATCH /v1/subscriptions/{id}
POST  /v1/subscriptions/{id}/activate
POST  /v1/subscriptions/{id}/charge
```

### Transactions

```
GET /v1/transactions
GET /v1/transactions/{id}
```

### Events & Notification Settings

```
GET  /v1/events
GET  /v1/notification-settings
POST /v1/notification-settings
```

### Admin (test helpers)

```
POST /admin/reset                          # Reset to seed state
POST /admin/advance-time/{subscription_id} # Simulate time passing
POST /admin/trigger-webhook/{event_type}   # Manually fire a webhook
GET  /ping                                 # Health check
```

## Seed Data

On startup (unless `-no-seed`), the store is populated with:

| Type | ID | Details |
|------|----|---------|
| Product | `prod_yieldly_base` | Yieldly Base Plan |
| Price | `pri_yieldly_monthly` | $5.00/month, 3-month trial |
| Customer | `ctm_test_alice` | alice@test.com, has trialing subscription |
| Customer | `ctm_test_bob` | bob@test.com, no subscription |
| Subscription | `sub_test_alice` | trialing, trial ends in 90 days |

## Subscription Lifecycle

The mock tracks subscription state through:

```
trialing → active → past_due → canceled
                  → paused → active
```

Use `POST /admin/advance-time/{id}` to move a subscription to its next state:

- **trialing** → activates (trial ends, first billing)
- **active** → next billing cycle (or `?fail=true` for payment failure → past_due)
- **past_due** → canceled
- **paused** → resumed (active)

## Webhooks

Register a webhook URL via the API or the `-webhook-url` flag. When subscription or transaction state changes, the mock POSTs to all registered URLs with:

- Paddle's event payload format
- `Paddle-Signature` header (`ts=...;h1=...`) signed with HMAC-SHA256 using the configured signing secret

Event types fired: `subscription.created`, `subscription.updated`, `subscription.activated`, `subscription.canceled`, `subscription.past_due`, `transaction.completed`, `transaction.payment_failed`.

## Response Format

All responses use Paddle's standard envelope:

```json
{
  "data": { ... },
  "meta": { "request_id": "req_00000001" }
}
```

List endpoints include pagination metadata:

```json
{
  "data": [ ... ],
  "meta": { "request_id": "req_00000002" },
  "pagination": { "per_page": 50, "has_more": false, "estimated_total": 2 }
}
```

## Disclaimer

**This is not an official Paddle product.** This project is an independent mock server created for local development and testing purposes. It is not affiliated with, endorsed by, or supported by [Paddle.com](https://www.paddle.com).

The user is fully responsible for how this tool is used. The authors make no warranties regarding accuracy, completeness, or fitness for any particular purpose. Always test against Paddle's real sandbox environment before going to production.

## License

MIT License. See [LICENSE](LICENSE).
