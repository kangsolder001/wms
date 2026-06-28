# Warehouse Management System (WMS)

Web-based Warehouse Management System built with Go backend (Clean Architecture) and React frontend. Designed for medium warehouses (5-10 users) with modular feature packages.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22+, chi router, pgx (PostgreSQL), zerolog, JWT |
| Frontend | React 18, TypeScript, Vite, Ant Design, Zustand, TanStack Query |
| Database | PostgreSQL 16 |
| Infra | Docker Compose, Makefile |

## Features

- **Auth & RBAC** — Login, JWT, 4 roles (admin, manager, operator, viewer)
- **Master Data** — Item/product catalog with SKU, categories, UoM
- **Location Management** — Warehouse zones, aisles, racks, bins hierarchy
- **Inventory** — Stock levels per location, adjustments, stock movements
- **Inbound** — Purchase orders, goods receipt (GRN)
- **Outbound** — Sales orders, pick list, packing & dispatch
- **Stock Transfer** — Move stock between locations
- **Dashboard** — Stock summary, recent activities, KPI cards
- **Dark Mode** — Catppuccin Mocha palette with light/dark toggle

## Quick Start

```bash
# 1. Start PostgreSQL
make docker-up

# 2. Start backend (port 8080)
cd backend && make dev

# 3. Start frontend (dev mode, port 5173)
cd frontend && npm run dev
```

Default login: `admin` / `admin123`

## Project Structure

```
wms/
├── backend/                          # Go Clean Architecture
│   ├── cmd/wms/main.go              # Entry point, DI wiring
│   ├── config/config.go             # Config loader
│   ├── domain/                       # Entities + repository interfaces
│   ├── application/                  # Use cases + DTOs
│   ├── infrastructure/              # DB, logger, JWT, repo implementations
│   ├── delivery/http/               # Handlers, middleware, router
│   └── pkg/logger/                  # Logger interface
├── frontend/                         # React + Vite + TypeScript
│   └── src/
│       ├── features/                # Auth, Dashboard, Items, Locations, Stock, etc.
│       ├── api/                     # API client layer
│       ├── stores/                  # Zustand state (auth, theme)
│       └── shared/components/       # Layout, sidebar
├── docker-compose.yml
└── Makefile
```

## API

Base URL: `http://localhost:8080/api/v1`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | Login |
| GET/POST | `/items` | Item CRUD |
| GET/POST | `/locations` | Location CRUD |
| GET | `/stock` | Stock levels |
| GET/POST | `/purchase-orders` | PO management |
| POST | `/purchase-orders/:id/receive` | Receive goods |
| GET/POST | `/sales-orders` | SO management |
| POST | `/sales-orders/:id/pick` | Pick order |
| POST | `/sales-orders/:id/ship` | Ship order |
| GET/POST | `/transfers` | Stock transfers |
| GET | `/dashboard/summary` | Dashboard data |

## Build & Run

```bash
# Backend
cd backend
make build          # Build binary
make dev            # Run in dev mode
make test           # Run tests
make lint           # Lint with golangci-lint

# Frontend
cd frontend
npm run build       # Production build
npm run dev         # Dev server
```

## License

Private
