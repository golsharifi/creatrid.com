# Creatrid.com — Gap Closure Plan

## Overview

7 remaining gaps to reach full vision completion. Custom domain is already done.

---

## Gap Tracking

| # | Gap | Priority | Status | Notes |
|---|-----|----------|--------|-------|
| 1 | Social OAuth Credentials | Low | ✅ Done | Twitter/X + LinkedIn + Instagram deployed. Behance + Dribbble removed (deprecated/not needed) |
| 2 | SMTP Configuration | Low | ✅ Done | Gmail SMTP via info@creatrid.com — tested and working |
| 3 | Custom Domain | — | ✅ Done | creatrid.com → frontend, api.creatrid.com → backend |
| 4 | Blockchain Anchoring | High | ✅ Done | Backend + frontend + verify page + nav link + vault anchor button |
| 5 | B2B Agency Dashboard | High | ✅ Done | Agency accounts, bulk verification, API billing — 100% complete |
| 6 | Webhook Delivery Worker | High | ✅ Done | Background worker, retry with backoff, delivery logs — 100% complete |
| 7 | Social Token Lite | High | ✅ Done | Tokens, tips, subscriptions — backend + frontend + i18n |

---

## Implementation Order

### Round 1 (parallel) — ✅ Complete
- **Gap 4: Blockchain Anchoring** — ✅ Done
- **Gap 5: B2B Agency Dashboard** — ✅ Done
- **Gap 6: Webhook Delivery Worker** — ✅ Done
- **Gap 7: Social Token Lite** — ✅ Done

### Round 2 (requires user credentials)
- **Gap 1: Social OAuth Credentials** — Configure Twitter, LinkedIn, Instagram, Behance, Dribbble API keys
- **Gap 2: SMTP Configuration** — Configure Gmail App Password or SendGrid API key

---

## Gap 4: Blockchain Anchoring ✅

**Goal**: Hash creator content and anchor proofs on-chain for verifiable ownership timestamps.

### Backend
- [x] Migration: `content_anchors` table (content_id, tx_hash, chain, block_number, anchor_hash, status, created_at)
- [x] `internal/blockchain/anchor.go` — Service to compute SHA-256 hash, submit to Base/Polygon via RPC
- [x] `internal/blockchain/verify.go` — Verify anchor by checking on-chain data (simulated mode)
- [x] `handler/blockchain.go` — POST /api/content/{id}/anchor, GET /api/content/{id}/anchor, GET /api/verify/{hash}, GET /api/anchors
- [x] Config: `BLOCKCHAIN_RPC_URL`, `BLOCKCHAIN_PRIVATE_KEY`, `BLOCKCHAIN_CHAIN_ID`

### Frontend
- [x] Verification public page (`/verify`) with search by hash, proof certificate, QR code
- [x] API client methods (anchor, getAnchor, verify, listAnchors)
- [x] i18n keys (17 translation keys across en/es/fa)
- [ ] Header nav link to `/verify` — in progress
- [ ] "Anchor on Blockchain" button on vault item detail page — in progress

### Smart Contract
- [x] ABI for Go interaction (simulated mode — real contract deployment separate)

---

## Gap 5: B2B Agency Dashboard ✅

**Goal**: Agencies can manage multiple creators, bulk operations, and usage-based API access.

### Backend
- [x] Migration: `agencies`, `agency_creators`, `api_usage` tables
- [x] `handler/agency.go` — 11 endpoints: CRUD, invite, creators, analytics, API usage, bulk verify
- [x] `store/agency.go` — 520-line store with full database layer
- [x] `middleware/apiusage.go` — API usage tracking middleware
- [x] `middleware/apikey.go` — API key authentication

### Frontend
- [x] `/agency` dashboard — 4 tabs (Dashboard, Creators, API Usage, Settings)
- [x] `/agency/invites` page — Creator-side invite management
- [x] API client methods (10 methods)
- [x] i18n keys (30+ keys)
- [x] Header nav link

---

## Gap 6: Webhook Delivery Worker ✅

**Goal**: Reliably deliver webhooks to registered endpoints with retry logic.

### Backend
- [x] Migration: `webhook_deliveries` columns (attempts, max_attempts, next_retry_at, status)
- [x] `internal/webhook/worker.go` — Background goroutine, 5s poll, exponential backoff, 5 max attempts
- [x] `internal/webhook/dispatcher.go` — Queue webhook events with delivery tracking
- [x] Store methods: ListPendingDeliveries, IncrementDeliveryAttempt, MarkDeliveryDead, ResetDeliveryForRetry
- [x] Dead letter handling after max retries
- [x] HMAC-SHA256 signature generation

### Frontend
- [x] Webhook delivery logs (status badges, attempt count, timestamps)
- [x] Manual retry button for failed deliveries
- [x] API client methods (deliveries, retryDelivery)

---

## Gap 7: Social Token Lite ✅

**Goal**: Creator tokens for gated content, tipping, and fan subscriptions.

### Backend
- [x] Migration: `creator_tokens`, `token_balances`, `tips`, `fan_subscriptions`, `gated_content`, `token_transactions` tables
- [x] `handler/tokens.go` — Create/get/update token, holders, transactions, purchase (Stripe + simulated)
- [x] `handler/tips.go` — Send tip, list received/sent, stats
- [x] `handler/subscriptions.go` — Subscribe (3 tiers), list subs/fans, cancel, gate/ungate content, access check
- [x] `store/tokens.go` — 658-line store with full database layer
- [x] Routes wired in main.go (16+ endpoints)

### Frontend
- [x] `/tokens` page — Creator token dashboard with 5 tabs (Token, Holders, Transactions, Tips, Fans)
- [x] API client methods (tokens, tips, subscriptions — 13 methods)
- [x] i18n keys (31 keys across en/es/fa)
- [x] Header nav link (Zap icon, desktop + mobile)

---

## Gap 1: Social OAuth Credentials ⬚

**Requires**: User to provide API keys for each platform.

| Platform | Env Vars Needed | Status |
|----------|----------------|--------|
| Twitter/X | `TWITTER_CLIENT_ID`, `TWITTER_CLIENT_SECRET` | ⬚ Need API key |
| LinkedIn | `LINKEDIN_CLIENT_ID`, `LINKEDIN_CLIENT_SECRET` | ⬚ Need API key |
| Instagram | `INSTAGRAM_CLIENT_ID`, `INSTAGRAM_CLIENT_SECRET` | ⬚ Need API key |
| Behance | `BEHANCE_CLIENT_ID`, `BEHANCE_CLIENT_SECRET` | ⬚ Need API key |
| Dribbble | `DRIBBBLE_CLIENT_ID`, `DRIBBBLE_CLIENT_SECRET` | ⬚ Need API key |

---

## Gap 2: SMTP Configuration ⬚

**Requires**: User to provide email service credentials.

**Option A**: Gmail App Password
- `SMTP_HOST=smtp.gmail.com`
- `SMTP_USERNAME=your@gmail.com`
- `SMTP_PASSWORD=your-app-password`

**Option B**: SendGrid
- `SMTP_HOST=smtp.sendgrid.net`
- `SMTP_USERNAME=apikey`
- `SMTP_PASSWORD=your-sendgrid-api-key`
