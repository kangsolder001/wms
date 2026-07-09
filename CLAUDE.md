# WMS Project Guidance

## Overview

Warehouse Management System (WMS) — Go backend (Clean Architecture) + React frontend (TypeScript).

**Ports:**
- Backend: `localhost:8080`
- Frontend: `localhost:3001`
- PostgreSQL: `localhost:5334` (mapped from container 5432)
- Redis: `localhost:6380` (mapped from container 6379)

**Default Login:** `superadmin / admin123`

---

## Tech Stack

| Layer | Tech |
|-------|------|
| Backend | Go 1.25, Gin, pgx (PostgreSQL), zerolog, JWT |
| Frontend | React 19, TypeScript, Vite, Ant Design, Zustand, TanStack Query |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Infra | Docker Compose |

---

## Architecture — Backend (Clean Architecture)

```
backend/
├── cmd/wms/main.go              ← Entry point, DI wiring
├── domain/
│   ├── entity/                  ← Business entities (struct + methods)
│   └── repository/              ← Repository interfaces
├── application/
│   ├── dto/                     ← Request/Response DTOs
│   └── usecase/                 ← Business logic (implements interfaces)
├── infrastructure/
│   ├── repository/              ← PostgreSQL repository implementations
│   ├── auth/                    ← JWT service
│   ├── database/                ← Migration + seed
│   ├── logger/                  ← Zerolog implementation
│   └── redis/                   ← Redis client
├── delivery/http/
│   ├── handler/                 ← HTTP handlers (Gin)
│   ├── middleware/              ← Auth, Role, Logging, RateLimiter
│   ├── router/                  ← Route registration
│   └── response/                ← Standard JSON response helpers
└── pkg/logger/                  ← Logger interface
```

### Flow for each request

```
HTTP Request → Middleware (auth/role) → Handler → Usecase → Repository → PostgreSQL
```

### Adding a New Master Data (e.g., Supplier)

**Step 1: Entity** — `domain/entity/supplier.go`
```go
package entity

import "time"

type Supplier struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Code        string    `json:"code"`
    Description string    `json:"description"`
    IsActive    bool      `json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
}
```

**Step 2: Repository interface** — `domain/repository/supplier_repository.go`
```go
package repository

import (
    "context"
    "wms/domain/entity"
)

type SupplierRepository interface {
    FindByID(ctx context.Context, id string) (*entity.Supplier, error)
    Create(ctx context.Context, supplier *entity.Supplier) error
    Update(ctx context.Context, supplier *entity.Supplier) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, page, limit int) ([]*entity.Supplier, int, error)
    ListAll(ctx context.Context) ([]*entity.Supplier, error)
}
```

**Step 3: DTO** — `application/dto/supplier_dto.go`
```go
package dto

type CreateSupplierRequest struct {
    Name        string `json:"name" validate:"required"`
    Code        string `json:"code" validate:"required"`
    Description string `json:"description"`
}

type UpdateSupplierRequest struct {
    Name        *string `json:"name"`
    Code        *string `json:"code"`
    Description *string `json:"description"`
    IsActive    *bool   `json:"is_active"`
}

type SupplierResponse struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Code        string `json:"code"`
    Description string `json:"description"`
    IsActive    bool   `json:"is_active"`
}
```

**Step 4: Usecase** — `application/usecase/supplier_usecase.go`
```go
package usecase

// Interface + implementation, following item_usecase.go pattern
// Constructor: NewSupplierUsecase(supplierRepo repository.SupplierRepository, log logger.Logger) SupplierUsecase
```

**Step 5: Repository impl** — `infrastructure/repository/postgres_supplier_repository.go`
```go
package repository

// PostgreSQL implementation, following postgres_item_repository.go pattern
// Constructor: NewPostgresSupplierRepository(db *sql.DB, log logger.Logger) *postgresSupplierRepository
```

**Step 6: Handler** — `delivery/http/handler/supplier_handler.go`
```go
package handler

// Gin handlers, following item_handler.go pattern
// Constructor: NewSupplierHandler(supplierUC usecase.SupplierUsecase, log logger.Logger) *SupplierHandler
```

**Step 7: Register in main.go** — Add:
```go
// Repository
supplierRepo := postgresRepo.NewPostgresSupplierRepository(db, appLogger)
// Usecase
supplierUC := usecase.NewSupplierUsecase(supplierRepo, appLogger)
// Handler
supplierHandler := handler.NewSupplierHandler(supplierUC, appLogger)
```

**Step 8: Register in router.go** — Add to `NewRouter()` params + routes:
```go
// Add param
supplierHandler *handler.SupplierHandler,
// Add routes
suppliers := auth.Group("/suppliers")
{
    suppliers.GET("", supplierHandler.List)
    suppliers.GET("/all", supplierHandler.ListAll)
    suppliers.POST("", roleMiddleware.RequireRole(managerRoles...), supplierHandler.Create)
    suppliers.GET("/:id", supplierHandler.Get)
    suppliers.PUT("/:id", roleMiddleware.RequireRole(managerRoles...), supplierHandler.Update)
    suppliers.DELETE("/:id", roleMiddleware.RequireRole(adminRoles...), supplierHandler.Delete)
}
```

**Step 9: Database** — Add to `postgres.go`:
- Migration: `CREATE TABLE IF NOT EXISTS suppliers (...)`
- Alterations (if needed): `ALTER TABLE ... ADD COLUMN ...`
- Seed data (if needed)

### DB Patterns

- **Soft delete:** `UPDATE ... SET is_active = false WHERE id = $1` (never hard delete)
- **Pagination:** `SELECT ... LIMIT $1 OFFSET $2`, total via `SELECT COUNT(*)`
- **Auto-increment SKU:** `MAX(CAST(SUBSTRING(sku FROM LENGTH(prefix)+1) AS INTEGER))` → +1
- **NULL handling:** Use `COALESCE(column, 0)` in SELECT for nullable numeric columns
- **UUID PK:** `DEFAULT gen_random_uuid()`
- **Timestamps:** `created_at TIMESTAMPTZ DEFAULT NOW()`, `updated_at` manual set

### API Response Format

```go
// Success (single)
response.JSON(c, http.StatusOK, data)
// → {"success": true, "data": {...}}

// Success (paginated)
response.JSONWithMeta(c, http.StatusOK, data, &response.Meta{Page, Limit, Total})
// → {"success": true, "data": [...], "meta": {"page":1, "limit":10, "total":100}}

// Error
response.Error(c, http.StatusBadRequest, "message")
// → {"success": false, "error": "message"}
```

### Middleware Chain

```
Logging → CORS → RateLimiter (login only) → AuthMiddleware (JWT) → RoleMiddleware → Handler
```

- **AuthMiddleware:** Validates JWT, checks `is_active`, sets `user_id` and `role` in context
- **RoleMiddleware:** `RequireRole("superadmin", "manager", ...)`
- **Role hierarchy:** `superadmin > manager > operator > viewer`

---

## Architecture — Frontend

```
frontend/src/
├── api/                        ← Axios API clients (one per module)
│   ├── client.ts              ← Base axios instance + interceptors
│   ├── items.ts               ← Item API functions
│   ├── categories.ts          ← Category API
│   └── ...
├── features/                   ← Page components (one dir per module)
│   ├── auth/                  ← LoginPage, UsersPage
│   ├── masterdata/            ← ItemsPage, CategoriesPage, ZonesPage
│   │   ├── components/        ← QRLabel, QRLabelModal, LabelSizeSelector
│   │   └── styles/
│   ├── location/              ← LocationsPage
│   ├── inventory/             ← StockPage
│   ├── inbound/               ← PurchaseOrdersPage
│   ├── outbound/              ← SalesOrdersPage
│   ├── transfer/              ← TransfersPage
│   └── dashboard/             ← DashboardPage
├── shared/components/          ← MainLayout (sidebar + header)
├── stores/                     ← Zustand stores
│   ├── authStore.ts           ← Auth state (token, user, login/logout)
│   └── themeStore.ts          ← Dark/light theme toggle
├── App.tsx                     ← Routes + theme config
└── main.tsx                    ← React entry point
```

### Adding a New Page (e.g., Suppliers)

**Step 1: API client** — `api/suppliers.ts`
```typescript
import api from './client';

export interface Supplier {
  id: string;
  name: string;
  code: string;
  description: string;
  is_active: boolean;
}

export interface CreateSupplierRequest {
  name: string;
  code: string;
  description?: string;
}

export const supplierApi = {
  list: (page = 1, limit = 10) =>
    api.get('/suppliers', { params: { page, limit } }),
  listAll: () => api.get('/suppliers/all'),
  get: (id: string) => api.get<Supplier>(`/suppliers/${id}`),
  create: (data: CreateSupplierRequest) => api.post<Supplier>('/suppliers', data),
  update: (id: string, data: Partial<Supplier>) => api.put<Supplier>(`/suppliers/${id}`, data),
  delete: (id: string) => api.delete(`/suppliers/${id}`),
};
```

**Step 2: Page component** — `features/masterdata/SuppliersPage.tsx`
```tsx
// Follow CategoriesPage.tsx pattern:
// - useQuery for data fetching
// - useMutation for create/delete
// - Table with columns
// - Modal with Form for create/edit
```

**Step 3: Register route** — `App.tsx`
```tsx
import SuppliersPage from './features/masterdata/SuppliersPage';
// Add inside <Route path="/" ...>:
<Route path="suppliers" element={<SuppliersPage />} />
```

**Step 4: Add to sidebar** — `shared/components/MainLayout.tsx`
```tsx
// Add to baseMenuItems:
{ key: '/suppliers', icon: <TeamOutlined />, label: 'Suppliers' },
```

### Data Fetching Pattern

```tsx
// List (paginated) — used in table
const { data, isLoading } = useQuery({
  queryKey: ['items'],
  queryFn: () => itemApi.list().then((res) => res.data),
});
// Table: dataSource={data?.data || []}

// List All — used in dropdowns
const { data: categoriesData } = useQuery({
  queryKey: ['categories-all'],
  queryFn: () => categoryApi.listAll().then((res) => res.data.data), // ← note .data.data
});
// Select: options={(categoriesData || []).map(...)}
```

**Important:** The `/all` endpoint wraps response in `{ success, data }`, so you need `.data.data` to get the array.

### Auth Flow

1. Login → POST `/api/v1/auth/login` → returns `{ token, user }`
2. Store in `authStore` + `localStorage`
3. Axios interceptor adds `Authorization: Bearer {token}` to all requests
4. 401 response → auto-logout + redirect to `/login`
5. `ProtectedRoute` wrapper checks `isAuthenticated` before rendering

### Theme

- Zustand store (`themeStore.ts`) toggles `isDark`
- Ant Design `ConfigProvider` with dark/light token configs
- Toggle switch in header

---

## Docker

### Development

```bash
docker compose up -d              # Start all (postgres, redis, backend, frontend)
docker compose up -d --build      # Rebuild and start
docker compose logs -f backend    # View backend logs
docker compose logs -f frontend   # View frontend logs
docker compose down               # Stop all
```

### Production

```bash
docker compose -f docker-compose.prod.yml up --build -d
```

### Services

| Service | Container | Port | Notes |
|---------|-----------|------|-------|
| PostgreSQL | wms-postgres | 5434→5432 | Data in volume `postgres_data` |
| Redis | wms-redis | 6380→6379 | Data in volume `redis_data` |
| Backend | wms-backend | 8080 | Hot-reload with `air` |
| Frontend | wms-frontend | 3001 | Vite dev server |

---

## Database Schema (Master Data)

### categories
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | Auto |
| name | VARCHAR(100) UNIQUE | Required |
| abbreviation | VARCHAR(10) UNIQUE | Required — used for SKU prefix |
| description | TEXT | Optional |
| is_active | BOOLEAN | Default true |
| created_at | TIMESTAMPTZ | Default now() |

### zones
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | Auto |
| code | VARCHAR(50) UNIQUE | Required |
| name | VARCHAR(100) | Required |
| description | TEXT | Optional |
| is_active | BOOLEAN | Default true |
| created_at | TIMESTAMPTZ | Default now() |

### items
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | Auto |
| sku | VARCHAR(50) UNIQUE | Auto-generated: `{abbreviation}-{sequence}` |
| name | VARCHAR(255) | Required |
| description | TEXT | Optional |
| category_id | UUID FK | → categories.id |
| category | VARCHAR(100) | Denormalized category name |
| barcode | VARCHAR(500) | Auto-generated: `{sku}\|{name}\|{category}` |
| unit_of_measure | VARCHAR(20) | Required (pcs/kg/box) |
| weight | DECIMAL(10,2) | Optional |
| length | DECIMAL(10,2) | Optional (cm) |
| width | DECIMAL(10,2) | Optional (cm) |
| height | DECIMAL(10,2) | Optional (cm) |
| is_active | BOOLEAN | Default true |
| created_at | TIMESTAMPTZ | Default now() |
| updated_at | TIMESTAMPTZ | Manual |

### locations
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | Auto |
| code | VARCHAR(50) UNIQUE | Required |
| name | VARCHAR(100) | Required |
| zone | VARCHAR(50) | FK → zones.name (text) |
| aisle, rack, level, bin | VARCHAR(20) | Hierarchy |
| type | VARCHAR(20) | storage/receiving/shipping/staging |
| capacity | DECIMAL(10,2) | Optional |
| is_active | BOOLEAN | Default true |

### Other tables
- `users` — Auth + RBAC (4 roles)
- `stock` — Item qty per location (unique: item_id + location_id + batch_number)
- `stock_movements` — Audit log for stock changes
- `purchase_orders` + `purchase_order_items` — Inbound flow
- `goods_receipts` + `goods_receipt_items` — GRN
- `sales_orders` + `sales_order_items` — Outbound flow
- `pick_lists` — Pick assignment
- `shipments` — Shipping tracking
- `stock_transfers` — Inter-location transfers

---

## API Endpoints

Base: `http://localhost:8080/api/v1`

### Auth
| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| POST | `/auth/login` | Public | Login (rate-limited) |
| GET | `/auth/me` | All | Get profile |
| GET | `/auth/users` | Manager+ | List users |
| PUT | `/auth/users/:id` | Manager+ | Update user |
| DELETE | `/auth/users/:id` | Manager+ | Deactivate user |

### Master Data
| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/items` | All | List items (paginated) |
| GET | `/items/:id` | All | Get item |
| POST | `/items` | Manager+ | Create item (auto SKU) |
| POST | `/items/generate-sku` | Manager+ | Preview SKU for category |
| PUT | `/items/:id` | Manager+ | Update item |
| DELETE | `/items/:id` | Superadmin | Soft delete item |
| GET | `/categories` | All | List categories (paginated) |
| GET | `/categories/all` | All | List all categories (no pagination) |
| POST | `/categories` | Manager+ | Create category |
| PUT | `/categories/:id` | Manager+ | Update category |
| DELETE | `/categories/:id` | Superadmin | Soft delete category |
| GET | `/locations` | All | List locations (paginated) |
| POST | `/locations` | Manager+ | Create location |
| PUT | `/locations/:id` | Manager+ | Update location |
| DELETE | `/locations/:id` | Superadmin | Soft delete location |
| GET | `/zones` | All | List zones (paginated) |
| GET | `/zones/all` | All | List all zones (no pagination) |
| POST | `/zones` | Manager+ | Create zone |
| PUT | `/zones/:id` | Manager+ | Update zone |
| DELETE | `/zones/:id` | Superadmin | Soft delete zone |

### Inventory
| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/stock` | All | List stock levels |
| GET | `/stock/movements` | All | List movements |
| POST | `/stock/adjust` | Manager+ | Adjust stock |
| POST | `/stock/opname` | Operator+ | Stock opname |

### Inbound
| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/purchase-orders` | All | List POs |
| POST | `/purchase-orders` | Manager+ | Create PO |
| GET | `/purchase-orders/:id` | All | Get PO |
| POST | `/purchase-orders/:id/approve` | Manager+ | Approve PO |
| POST | `/purchase-orders/:id/receive` | Operator+ | Receive goods (GRN) |

### Outbound
| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/sales-orders` | All | List SOs |
| POST | `/sales-orders` | Manager+ | Create SO |
| GET | `/sales-orders/:id` | All | Get SO |
| POST | `/sales-orders/:id/pick` | Operator+ | Pick order |
| POST | `/sales-orders/:id/ship` | Operator+ | Ship order |

### Transfer
| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/transfers` | All | List transfers |
| POST | `/transfers` | Manager+ | Create transfer |
| PUT | `/transfers/:id/complete` | Operator+ | Complete transfer |

### Dashboard
| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/dashboard/summary` | All | Dashboard stats |

---

## Code Conventions

### Backend (Go)
- **Naming:** PascalCase for exported, camelCase for unexported
- **Files:** `snake_case.go`
- **Packages:** lowercase, single word
- **Errors:** Wrap with `fmt.Errorf("context: %w", err)`
- **Logging:** `uc.log.Info("message", "key", value)` (zerolog structured)
- **No comments** unless asked
- **Constructor pattern:** `NewXxxHandler(uc, log)`, `NewXxxUsecase(repo, log)`, `NewXxxRepository(db, log)`

### Frontend (TypeScript)
- **Files:** `PascalCase.tsx` for components, `camelCase.ts` for utilities
- **Components:** Export default function component
- **Hooks:** `use` prefix
- **API types:** Defined in `api/*.ts`, co-located with API functions
- **State:** Zustand for global (auth, theme), React Query for server state
- **Styling:** Inline styles (Ant Design pattern), no CSS modules
- **No comments** unless asked

---

## Commands

```bash
# Backend
cd backend
make dev                  # Run with hot-reload (air)
make build                # Build binary
make test                 # Run tests
make lint                 # Lint (golangci-lint)
make fmt                  # Format code
go build ./...            # Verify compilation

# Frontend
cd frontend
npm run dev               # Dev server (port 3001)
npm run build             # Production build
npm run lint              # Lint (oxlint)
npx tsc --noEmit          # Type check

# Docker
docker compose up -d --build        # Rebuild all
docker compose up -d --build backend   # Rebuild backend only
docker compose up -d --build frontend  # Rebuild frontend only
docker compose restart backend      # Restart backend (no rebuild)
docker compose logs -f backend      # Tail backend logs
```

---

## Common Pitfalls

1. **`/all` endpoint response:** Always `.data.data` to get array (backend wraps in `{success, data}`)
2. **NULL columns:** Use `COALESCE(col, default)` in SQL SELECT for nullable numeric columns
3. **SKU auto-gen:** Must call `generate-sku` endpoint first to get next SKU, then submit with `category_id`
4. **Soft delete:** Never hard delete. Set `is_active = false`
5. **Barcode:** Auto-generated from `{sku}|{name}|{category}`, not user input
6. **Category in items:** Stored as both `category_id` (FK) and `category` (denormalized name)
7. **Zone in locations:** Stored as text `zone` (references zone name, not FK)
