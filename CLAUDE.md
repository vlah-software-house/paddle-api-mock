**Project:** `paddle-mock-api` — A lightweight Go mock server that simulates the Paddle Billing API for local development and E2E testing.

**Reference:** This follows the same architecture as [etsy-mock-api](https://github.com/vlah-software-house/etsy-mock-api) — standalone Go binary, in-memory store, seed data, optional auth, admin reset endpoint.

**Purpose:** The consuming application (Yieldly) integrates Paddle for subscription billing. This mock lets us run E2E tests and develop locally without hitting Paddle's sandbox. The mock must be able to:
1. Simulate the Paddle checkout flow (return a transaction with subscription)
2. Accept and respond to Paddle Billing API calls
3. Fire webhook payloads to a configurable callback URL when state changes
4. Track subscription lifecycle (trialing → active → past_due → canceled)

---

### API Endpoints to Implement

Base all endpoints on the [Paddle Billing API](https://developer.paddle.com/api-reference/overview). Use the `/v1/` prefix like real Paddle.

**Products & Prices (read-only, seeded):**
- `GET /v1/products` — list products
- `GET /v1/products/{id}` — get product
- `GET /v1/prices` — list prices
- `GET /v1/prices/{id}` — get price

**Customers:**
- `POST /v1/customers` — create customer
- `GET /v1/customers` — list customers
- `GET /v1/customers/{id}` — get customer
- `PATCH /v1/customers/{id}` — update customer

**Subscriptions:**
- `POST /v1/subscriptions` — create subscription (simulate checkout result)
- `GET /v1/subscriptions` — list subscriptions
- `GET /v1/subscriptions/{id}` — get subscription
- `PATCH /v1/subscriptions/{id}` — update subscription (cancel, pause, change price)
- `POST /v1/subscriptions/{id}/activate` — activate a trialing subscription
- `POST /v1/subscriptions/{id}/charge` — create a one-time charge (usage extras)

**Transactions:**
- `GET /v1/transactions` — list transactions
- `GET /v1/transactions/{id}` — get transaction

**Events / Webhooks:**
- `GET /v1/events` — list fired events
- `GET /v1/notification-settings` — list webhook endpoints
- `POST /v1/notification-settings` — register a webhook URL
- `POST /admin/trigger-webhook/{event_type}` — manually trigger a webhook (test helper)

**Admin (test helpers, not part of real Paddle API):**
- `POST /admin/reset` — reset all data to seed state
- `POST /admin/advance-time/{subscription_id}` — simulate time passing (trigger trial end, billing cycle, payment failure)
- `GET /ping` — health check

---

### Seed Data

Pre-seed these on startup (matching Yieldly's pricing model):

**Product:**
- `prod_yieldly_base` — "Yieldly Base Plan"

**Price:**
- `pri_yieldly_monthly` — $5.00/month, 3-month trial, linked to `prod_yieldly_base`

**Customers (for testing):**
- `ctm_test_alice` — Alice (alice@test.com), active subscription, trial period
- `ctm_test_bob` — Bob (bob@test.com), no subscription

**Subscriptions:**
- `sub_test_alice` — trialing, linked to `ctm_test_alice` and `pri_yieldly_monthly`, trial ends in 90 days

---

### Webhook Behavior

When subscription state changes (via API calls or `admin/advance-time`), the mock should:

1. Create an event record in the in-memory store
2. POST the event payload to all registered webhook URLs
3. Use Paddle's real webhook payload format with these event types:
   - `subscription.created`
   - `subscription.updated`
   - `subscription.activated` (trial → active)
   - `subscription.canceled`
   - `subscription.past_due`
   - `transaction.completed`
   - `transaction.payment_failed`

Include Paddle's `Paddle-Signature` header using a configurable signing secret (default: `pdl_test_signing_secret`) so the consuming app can verify signatures in tests.

---

### CLI Flags

```
-port          Server port (default: 8081)
-no-auth       Disable API key validation (default: false)
-no-seed       Start with empty store (default: false)
-webhook-url   Default webhook URL to register on startup (e.g., http://localhost:8080/api/paddle/webhooks)
-signing-secret Webhook signing secret (default: pdl_test_signing_secret)
```

---

### Auth

When auth is enabled, require `Authorization: Bearer {api_key}` header on all `/v1/` endpoints. Pre-seed a test API key: `test_paddle_api_key`. Skip auth for `/admin/`, `/ping`, and webhook notification endpoints.

---

### Key Design Decisions

- **In-memory store only** — no database, data resets on restart
- **Synchronous webhooks** — when a state change happens, POST to webhook URLs before returning the API response (simpler for testing, no missed events)
- **Time simulation** — the `admin/advance-time` endpoint is critical for testing trial expiry and billing cycles without waiting 90 days
- **Paddle response format** — all responses must use Paddle's envelope: `{ "data": { ... }, "meta": { "request_id": "..." } }`
- **Subscription status machine**: `trialing → active → past_due → canceled` (and `paused` if needed)

---

That's it. Create the project repo, drop this into `CLAUDE.md`, and point Claude at it. The mock will be self-contained and ready for Yieldly's E2E tests.
