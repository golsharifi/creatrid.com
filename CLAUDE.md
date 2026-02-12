# Creatrid — Creator Passport

Verified digital identity platform for creators.

## Architecture

- **Backend**: Go (chi router, pgx, JWT auth) — `backend/`
- **Frontend**: Next.js 16 static export — `frontend/`
- **Database**: PostgreSQL (Docker Compose for local dev)
- **Auth**: Google OAuth → Go backend issues JWT → httpOnly cookie

## Commands

### Backend (Go)
- `cd backend && go run ./cmd/server` — Start Go API on :8080
- `go build ./...` — Verify compilation

### Frontend (Next.js)
- `cd frontend && pnpm dev` — Start dev server on :3000
- `cd frontend && pnpm build` — Static export to `out/`

### Infrastructure
- `docker compose up db -d` — Start PostgreSQL
- `docker compose up -d` — Start PostgreSQL + Go backend

## Project Structure

```
backend/
├── cmd/server/main.go       # Entry point
├── internal/
│   ├── config/              # Env loading
│   ├── auth/                # Google OAuth + JWT
│   ├── handler/             # HTTP handlers
│   ├── middleware/           # CORS + auth
│   ├── model/               # Data types
│   └── store/               # PostgreSQL queries (pgx)
└── migrations/              # SQL migrations

frontend/
├── src/app/                 # Pages (all client components)
├── src/components/          # Layout, UI components
└── src/lib/                 # API client, auth context, utils
```

## API Endpoints

- `GET /api/auth/google` — Redirect to Google OAuth
- `GET /api/auth/google/callback` — Handle OAuth callback
- `GET /api/auth/me` — Current user (auth required)
- `POST /api/auth/logout` — Clear session
- `POST /api/users/onboard` — Set username (auth required)
- `PATCH /api/users/profile` — Update profile (auth required)
- `GET /api/users/{username}` — Public profile
