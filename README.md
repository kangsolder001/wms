# Warehouse Management System (WMS)

Web-based Warehouse Management System built with Go backend (Clean Architecture) and React frontend. Designed for medium warehouses (5-10 users) with modular feature packages.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.25+, chi router, pgx (PostgreSQL), zerolog, JWT |
| Frontend | React 18, TypeScript, Vite, Ant Design, Zustand, TanStack Query |
| Database | PostgreSQL 16 |
| Infra | Docker Compose, Makefile |

## Features

- **Auth & RBAC** — Login, JWT, 4 roles (superadmin, manager, operator, viewer)
- **User Management** — List, create, edit, deactivate users (admin/manager only)
- **Master Data** — Item/product catalog with SKU, categories, UoM
- **Location Management** — Warehouse zones, aisles, racks, bins hierarchy
- **Inventory** — Stock levels per location, adjustments, stock movements
- **Inbound** — Purchase orders, goods receipt (GRN)
- **Outbound** — Sales orders, pick list, packing & dispatch
- **Stock Transfer** — Move stock between locations
- **Dashboard** — Stock summary, recent activities, KPI cards
- **Dark Mode** — Catppuccin Mocha palette with light/dark toggle

## Quick Start

### Local Development

```bash
# 1. Start PostgreSQL & Redis (via Docker)
docker compose -f docker-compose.prod.yml up postgres redis -d

# 2. Start backend (port 8080)
cd backend && make dev

# 3. Start frontend (dev mode, port 3001)
cd frontend && npm run dev

# 4. Seed database
cd backend && make seed
```

### Default Login

| Username | Password | Role |
|----------|----------|------|
| superadmin | admin123 | superadmin |
| manager | manager123 | manager |
| operator | operator123 | operator |

### Access URLs

| Service | URL |
|---------|-----|
| Frontend (Vite dev) | http://localhost:3001 |
| Backend API | http://localhost:8080 |

## Project Structure

```
wms/
├── backend/                          # Go Clean Architecture
│   ├── cmd/
│   │   ├── wms/main.go              # Entry point, DI wiring
│   │   └── seed/main.go            # Database seed script
│   ├── config/config.go             # Config loader (supports CONFIG_FILE env)
│   ├── domain/                       # Entities + repository interfaces
│   ├── application/                  # Use cases + DTOs
│   ├── infrastructure/              # DB, logger, JWT, repo implementations
│   ├── delivery/http/               # Handlers, middleware, router
│   └── pkg/logger/                  # Logger interface
├── frontend/                         # React + Vite + TypeScript
│   └── src/
│       ├── features/auth/           # Login, User Management
│       ├── features/dashboard/      # Dashboard page
│       ├── features/masterdata/     # Items management
│       ├── features/location/       # Location management
│       ├── features/inventory/      # Stock management
│       ├── features/inbound/        # Purchase orders
│       ├── features/outbound/       # Sales orders
│       ├── features/transfer/       # Stock transfers
│       ├── api/                     # API client layer
│       ├── stores/                  # Zustand state (auth, theme)
│       └── shared/components/       # Layout, sidebar
├── docker-compose.prod.yml          # Production setup (Docker)
└── Makefile
```

## Docker (Production)

Uses multi-stage builds, static frontend via Nginx, Go binary:

```bash
# Build and start
docker compose -f docker-compose.prod.yml up --build -d

# Access on port 80
curl http://localhost
```

## API

Base URL: `http://localhost:8080/api/v1`

| Method | Endpoint | Description | Role |
|--------|----------|-------------|------|
| POST | `/auth/login` | Login | Public |
| GET | `/auth/me` | Get profile | All |
| GET | `/auth/users` | List users | Admin, Manager |
| PUT | `/auth/users/:id` | Update user | Admin, Manager |
| DELETE | `/auth/users/:id` | Deactivate user | Admin, Manager |
| GET/POST | `/items` | Item CRUD | Manager+ |
| GET/POST | `/locations` | Location CRUD | Manager+ |
| GET | `/stock` | Stock levels | All |
| GET | `/stock/movements` | Stock movements | All |
| POST | `/stock/adjust` | Adjust stock | Manager+ |
| GET/POST | `/purchase-orders` | PO management | Manager+ |
| POST | `/purchase-orders/:id/receive` | Receive goods | Operator+ |
| GET/POST | `/sales-orders` | SO management | Manager+ |
| POST | `/sales-orders/:id/pick` | Pick order | Operator+ |
| POST | `/sales-orders/:id/ship` | Ship order | Operator+ |
| GET/POST | `/transfers` | Stock transfers | Manager+ |
| PUT | `/transfers/:id/complete` | Complete transfer | Operator+ |
| GET | `/dashboard/summary` | Dashboard data | All |

## Build & Run

```bash
# Backend
cd backend
make build          # Build binary
make dev            # Run in dev mode

# Frontend
cd frontend
npm run build       # Production build
npm run dev         # Dev server

# Seed database
cd backend && make seed  # or: go run ./cmd/seed/main.go
```

## License

Private
