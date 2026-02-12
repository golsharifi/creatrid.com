# Creatrid — Creator Passport

Verified digital identity platform for creators. Connect social accounts, build a Creator Score, and share a public profile.

## Architecture

- **Backend**: Go (chi/v5 router, pgx/v5, JWT auth) — `backend/`
- **Frontend**: Next.js 16 static export (`output: 'export'`) — `frontend/`
- **Database**: PostgreSQL 16 (Docker Compose local, Azure Flexible Server production)
- **Auth**: Google OAuth → Go backend issues JWT → httpOnly cookie (`SameSite=None` in production for cross-origin)
- **Hosting**: Azure Container Apps (backend), Azure Static Web Apps (frontend)
- **CI/CD**: GitHub Actions — auto-deploy on push to `main`

## Commands

### Backend (Go)
- `cd backend && go run ./cmd/server` — Start API on :8080
- `cd backend && go build ./...` — Verify compilation

### Frontend (Next.js)
- `cd frontend && pnpm dev` — Dev server on :3000
- `cd frontend && pnpm build` — Static export to `out/`

### Infrastructure
- `docker compose up db -d` — Start local PostgreSQL
- `./deploy/deploy.sh` — Provision all Azure resources (one-time)
- `./deploy/deploy-backend.sh` — Rebuild and redeploy backend
- `./deploy/deploy-frontend.sh` — Build and deploy frontend
- `./deploy/set-secrets.sh` — Update Container App secrets

## Production URLs

| Resource | URL |
|----------|-----|
| Frontend | https://lemon-plant-0ae85c203.2.azurestaticapps.net |
| Backend API | https://creatrid-api.bluewave-351c9003.westeurope.azurecontainerapps.io |
| Database | creatrid-pg.postgres.database.azure.com |
| Container Registry | creatridacr.azurecr.io |

## Project Structure

```
backend/
├── cmd/server/main.go          # Entry point, routes, migration runner
├── internal/
│   ├── auth/                   # Google OAuth + JWT service
│   ├── config/                 # Env var loading
│   ├── handler/                # HTTP handlers (auth, user, connection, og, score, analytics, admin)
│   ├── middleware/              # CORS, RequireAuth, RequireAuthRedirect, RateLimit
│   ├── model/                  # User, Account, Connection structs
│   ├── platform/               # Provider interface + YouTube, GitHub, Twitter, LinkedIn, Instagram, Dribbble, Behance
│   ├── score/                  # Creator Score calculation engine
│   ├── storage/                # Azure Blob Storage client (image uploads)
│   └── store/                  # PostgreSQL queries (pgx)
├── migrations/                 # 001_initial, 002_connections, 003_profile_customization, 004_analytics

frontend/
├── src/app/                    # Pages (all "use client" components)
│   ├── admin/                  # Admin dashboard (stats, user management, verification)
│   ├── connections/            # Connect/disconnect social accounts
│   ├── dashboard/              # Main dashboard with score, analytics, share, QR
│   ├── onboarding/             # First-time username setup
│   ├── profile/                # Public profile with themes, links, connections
│   ├── settings/               # Profile editor, avatar upload, theme picker, custom links
│   └── sign-in/                # Google OAuth sign-in
├── src/components/             # Header, footer, providers, share buttons
└── src/lib/                    # API client, auth context, types

deploy/
├── deploy.sh                   # Full Azure provisioning script
├── deploy-backend.sh           # Backend redeploy
├── deploy-frontend.sh          # Frontend redeploy
├── set-secrets.sh              # Update Container App secrets
└── .env.production             # Production env template

.github/workflows/
├── backend.yml                 # Go build + ACR + Container Apps deploy
└── frontend.yml                # pnpm build + Static Web Apps deploy
```

## API Endpoints

### Auth
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/auth/google` | No | Redirect to Google OAuth |
| GET | `/api/auth/google/callback` | No | Handle OAuth callback |
| GET | `/api/auth/me` | Yes | Current user |
| POST | `/api/auth/logout` | Yes | Clear session |

### Users
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/users/onboard` | Yes | Set username + name |
| PATCH | `/api/users/profile` | Yes | Update profile, theme, custom links |
| POST | `/api/users/profile/image` | Yes | Upload profile photo (multipart) |
| GET | `/api/users/{username}` | No | Public profile |
| GET | `/api/users/{username}/connections` | No | Public connections |

### Connections
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/connections` | Yes | List user's connections |
| GET | `/api/connections/{platform}/connect` | Yes (redirect) | Start platform OAuth |
| GET | `/api/connections/{platform}/callback` | Yes (redirect) | Handle platform OAuth callback |
| DELETE | `/api/connections/{platform}` | Yes | Remove connection |
| POST | `/api/connections/{platform}/refresh` | Yes | Refresh token and re-fetch profile |

### Analytics
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/analytics` | Yes | Analytics summary for authenticated user |
| POST | `/api/users/{username}/view` | No | Track profile view |
| POST | `/api/users/{username}/click` | No | Track link click |

### Admin
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/admin/stats` | Admin | Platform-wide stats |
| GET | `/api/admin/users` | Admin | List all users (paginated) |
| POST | `/api/admin/users/verify` | Admin | Set/unset verification badge |

### Other
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/health` | No | Health check |
| GET | `/p/{username}` | No | OG meta tags page (social previews) |

## Roadmap

### Completed

- [x] **Phase 1: Go Backend** — chi router, pgx, JWT, Google OAuth, user CRUD, migrations
- [x] **Phase 2: Next.js Frontend** — static export, auth context, sign-in, onboarding, dashboard, settings, public profile
- [x] **Phase 3: Integration** — API client, cookie auth, CORS, end-to-end flows
- [x] **Phase 4: Social Connections** — YouTube + GitHub OAuth providers, connections table, connect/disconnect flow
- [x] **Phase 5: Creator Score** — score engine (profile 20 + email 10 + connections 50 + followers 20), auto-recalc on changes
- [x] **Phase 6: Public Profiles & Sharing** — OG meta tags, QR codes, copy link, share to Twitter/LinkedIn
- [x] **Phase 7: All Social Platforms** — Twitter/X, LinkedIn, Instagram, Behance, Dribbble providers (enabled via env vars)
- [x] **Phase 8: Profile Customization** — 6 themes, custom links (up to 10), settings UI
- [x] **Phase 9: Azure Deployment** — Container Apps, PostgreSQL Flexible Server, Static Web Apps, ACR
- [x] **Phase 10: CI/CD** — GitHub Actions workflows, Azure service principal, auto-deploy on push
- [x] **Phase 11: Profile Analytics** — profile_views + link_clicks tables, view/click tracking, analytics dashboard cards
- [x] **Phase 12: Image Upload** — Azure Blob Storage, multipart upload endpoint, avatar upload UI in settings
- [x] **Phase 13: API Rate Limiting** — Token bucket per IP (20 req/s, burst 40), auto-cleanup of stale visitors
- [x] **Phase 14: Admin Dashboard** — `/admin` page, platform stats, user table with pagination, verify/unverify
- [x] **Phase 15: Verification System** — Admin-managed verification badges, badge display on public profile
- [x] **Phase 16: Token Refresh** — RefreshToken on all 7 providers, refresh endpoint, auto-update stored tokens

### Upcoming

- [ ] **Custom Domain** — Point `creatrid.com` → Static Web Apps, `api.creatrid.com` → Container Apps
- [ ] **Email Notifications** — Welcome email, connection alerts, weekly digest
- [ ] **B2B Features** — Brand accounts, creator discovery, collaboration requests

## Key Design Decisions

- **Go backend over Next.js API routes**: Better performance, simpler deployment, explicit control over auth flow
- **Static export**: Frontend is pure client-side SPA, no server-side rendering needed
- **SameSite=None cookies**: Required for cross-origin auth between different Azure domains (frontend ≠ backend domain)
- **Provider interface pattern**: All social platforms implement the same `Provider` interface for generic connect/callback handlers
- **RequireAuthRedirect middleware**: Browser OAuth flows redirect to sign-in (not 401 JSON) when unauthenticated
- **Logarithmic follower bonus**: Score is fair across audience sizes (100 followers → 5 pts, 100K → 20 pts)
- **Conditional platform providers**: Each platform is enabled only if its env vars are set (no crash on missing credentials)

## Environment Variables

See `.env.example` for the full list. Required: `DATABASE_URL`, `JWT_SECRET`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`. All social platform credentials are optional. Image upload requires `AZURE_STORAGE_ACCOUNT`, `AZURE_STORAGE_KEY`, `AZURE_STORAGE_CONTAINER`.
