# WMS Project Session Notes

## Session Info
- **Session ID**: `ses_0e216e6faffex9jesqSCspbcTb`
- **Session memory**: `/home/mukhlisadigunawan/.local/share/mimocode/memory/sessions/ses_0e216e6faffex9jesqSCspbcTb/`
- **Git branch**: `main`

## Progress
- [x] User Registration Page (admin/manager only)
- [x] User Management Page (list, create, edit, deactivate users)
- [x] Toggle status dengan Popconfirm
- [x] Auth middleware cek `is_active` di DB → auto logout jika deactivate
- [x] Docker setup (dev + prod)
- [x] Seed data (3 users, 10 items, 10 locations, 13 stock, 3 PO, 2 SO, 2 transfers)
- [x] Frontend proxy Vite → backend Docker
- [x] Migrate router: chi → gin-gonic/gin
- [x] UsersPage: table with list, create, edit, deactivate, toggle status
- [x] Backend: ListUsers, UpdateUser, DeleteUser endpoints

## Cara Run
```bash
# Dev (hot-reload)
# 1. Start PostgreSQL & Redis
docker compose -f docker-compose.prod.yml up postgres redis -d

# 2. Start backend
cd backend && make dev

# 3. Start frontend
cd frontend && npm run dev

# 4. Seed database
cd backend && make seed

# Akses
# Frontend: http://localhost:3001
# Backend:  http://localhost:8080

# Login
# superadmin/admin123, manager/manager123, operator/operator123
```

## File Structure
```
wms/
├── backend/
│   ├── cmd/wms/main.go          # Main app
│   ├── cmd/seed/main.go         # Seed script
│   ├── config.docker.json       # Docker config (host=postgres)
│   ├── delivery/http/
│   │   ├── handler/auth_handler.go    # Login, Register, ListUsers, UpdateUser, DeleteUser
│   │   ├── middleware/auth_middleware.go  # JWT + cek is_active
│   │   └── router/router.go          # Routes
│   ├── application/usecase/auth_usecase.go  # Auth + user management
│   └── infrastructure/
│       ├── database/postgres.go  # Migrate + Seed
│       └── repository/postgres_user_repository.go
├── frontend/
│   ├── src/
│   │   ├── features/auth/
│   │   │   ├── LoginPage.tsx
│   │   │   └── UsersPage.tsx     # User Management (CRUD)
│   │   ├── api/auth.ts           # API functions (login, register, list, update, delete)
│   │   ├── stores/authStore.ts   # Persist user + token ke localStorage
│   │   └── App.tsx               # Routes (/users for admin/manager)
│   └── vite.config.ts            # Proxy /api → backend
└── docker-compose.prod.yml       # Production (nginx + binary)
```

## Design Decisions
- **Registration accessible to admin AND manager**: User initially said admin-only, then corrected to include manager
- **Seed code separate from main app**: `cmd/seed/main.go` as standalone entry point
- **Docker for production only**: Local dev runs natively, Docker only for production deployment
- **CONFIG_FILE env var for Docker**: Backend config supports env var override
- **authStore localStorage persistence**: User data (including role) must survive page reloads

## Errors Fixed
- User registration menu not showing → Zustand memory-only state, fixed by persisting to localStorage
- Docker DB fresh → Added comprehensive seed function
- Backend Docker DNS error → Fixed by full `docker compose down && up --build`
- Port conflicts → Killed old processes, changed ports
- air v1.65.3 requires Go >= 1.25 → Pinned to v1.61.7, later upgraded to golang:1.25-alpine
- Frontend blank page → Backend response double-nested, fixed with `JSONWithMeta`
- Vite not accessible in Docker → Added `host: '0.0.0.0'` + COPY . . in Dockerfile

## Pending
- [ ] Commit changes (user management + auth middleware)
- [ ] Push ke remote
- [ ] Test production docker setup
