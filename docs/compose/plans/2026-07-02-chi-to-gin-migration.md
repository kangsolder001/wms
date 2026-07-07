# Migrate Router: chi → gin

> **For agentic workers:** REQUIRED SUB-SKILL: Use compose:subagent (recommended) or compose:execute to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace chi router with gin-gonic/gin across all backend HTTP handlers, middleware, and router setup.

**Architecture:** Replace chi-specific patterns (`chi.URLParam`, `chi.Router`, `chi.Middleware`) with gin equivalents (`c.Param`, `*gin.Engine`, `gin.HandlerFunc`). Handler signatures change from `func(w http.ResponseWriter, r *http.Request)` to `func(c *gin.Context)`. Middleware changes from `func(http.Handler) http.Handler` to `gin.HandlerFunc`.

**Tech Stack:** Go 1.25+, gin-gonic/gin v1.10+, remove go-chi/chi/v5

---

## Scope Summary

| Component | chi usage | gin equivalent |
|-----------|-----------|----------------|
| URL params | `chi.URLParam(r, "id")` | `c.Param("id")` |
| Query params | `r.URL.Query().Get("x")` | `c.Query("x")` |
| Body decode | `json.NewDecoder(r.Body).Decode(&v)` | `c.ShouldBindJSON(&v)` |
| Context values | `r.Context().Value("x").(string)` | `c.GetString("x")` |
| Route groups | `r.Route("/path", func(r chi.Router){...})` | `r.Group("/path")` |
| Middleware | `func(http.Handler) http.Handler` | `gin.HandlerFunc` |
| Router | `chi.NewRouter()` | `gin.Default()` |

## Files Affected

### Modify (11 files)
- `backend/go.mod` — replace chi with gin dependency
- `backend/cmd/wms/main.go` — use gin engine, update middleware wiring
- `backend/delivery/http/router/router.go` — rewrite all routes with gin
- `backend/delivery/http/handler/auth_handler.go` — 6 methods → gin.Context
- `backend/delivery/http/handler/item_handler.go` — 5 methods → gin.Context
- `backend/delivery/http/handler/location_handler.go` — 5 methods → gin.Context
- `backend/delivery/http/handler/inventory_handler.go` — 4 methods → gin.Context
- `backend/delivery/http/handler/inbound_handler.go` — 4 methods → gin.Context
- `backend/delivery/http/handler/outbound_handler.go` — 5 methods → gin.Context
- `backend/delivery/http/handler/transfer_handler.go` — 3 methods → gin.Context
- `backend/delivery/http/handler/dashboard_handler.go` — 1 method → gin.Context
- `backend/delivery/http/middleware/auth_middleware.go` — gin.HandlerFunc
- `backend/delivery/http/middleware/role_middleware.go` — gin.HandlerFunc
- `backend/delivery/http/middleware/logging_middleware.go` — gin.HandlerFunc
- `backend/delivery/http/response/response.go` — use gin.Context for response

---

### Task 1: Install gin, remove chi

**Files:**
- Modify: `backend/go.mod`
- Modify: `backend/go.sum`

- [ ] **Step 1: Install gin**

```bash
cd backend && go get github.com/gin-gonic/gin@v1.10.0
```

- [ ] **Step 2: Remove chi**

```bash
cd backend && go get github.com/go-chi/chi/v5@none && go mod tidy
```

- [ ] **Step 3: Commit**

```bash
git add backend/go.mod backend/go.sum
git commit -m "chore: replace chi with gin dependency"
```

---

### Task 2: Convert response helpers to gin

**Files:**
- Modify: `backend/delivery/http/response/response.go`

- [ ] **Step 1: Rewrite response.go**

```go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `total:"total"`
}

func JSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": statusCode >= 200 && statusCode < 300,
		"data":    data,
	})
}

func JSONWithMeta(c *gin.Context, statusCode int, data interface{}, meta *Meta) {
	c.JSON(statusCode, gin.H{
		"success": statusCode >= 200 && statusCode < 300,
		"data":    data,
		"meta":    meta,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error":   message,
	})
}
```

- [ ] **Step 2: Build to verify**

```bash
cd backend && go build ./...
```

- [ ] **Step 3: Commit**

```bash
git add backend/delivery/http/response/response.go
git commit -m "refactor: convert response helpers to use gin.Context"
```

---

### Task 3: Convert auth middleware to gin

**Files:**
- Modify: `backend/delivery/http/middleware/auth_middleware.go`

- [ ] **Step 1: Rewrite auth_middleware.go**

```go
package middleware

import (
	"database/sql"
	"net/http"
	"strings"

	"wms/infrastructure/auth"
	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService auth.JWTService
	db         *sql.DB
	log        logger.Logger
}

func NewAuthMiddleware(jwtService auth.JWTService, db *sql.DB, log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService, db: db, log: log}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		claims, err := m.jwtService.ValidateToken(parts[1])
		if err != nil {
			m.log.Error("invalid token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		var isActive bool
		err = m.db.QueryRow("SELECT is_active FROM users WHERE id = $1", userID).Scan(&isActive)
		if err != nil || !isActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "account is deactivated"})
			return
		}

		role, _ := claims["role"].(string)

		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}
```

- [ ] **Step 2: Build to verify**

```bash
cd backend && go build ./...
```

- [ ] **Step 3: Commit**

```bash
git add backend/delivery/http/middleware/auth_middleware.go
git commit -m "refactor: convert auth middleware to gin.HandlerFunc"
```

---

### Task 4: Convert role & logging middleware to gin

**Files:**
- Modify: `backend/delivery/http/middleware/role_middleware.go`
- Modify: `backend/delivery/http/middleware/logging_middleware.go`

- [ ] **Step 1: Rewrite role_middleware.go**

```go
package middleware

import (
	"net/http"

	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RoleMiddleware struct {
	log logger.Logger
}

func NewRoleMiddleware(log logger.Logger) *RoleMiddleware {
	return &RoleMiddleware{log: log}
}

func (m *RoleMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		allowed := false
		for _, allowedRole := range roles {
			if roleStr == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			m.log.Warn("access denied", "required_roles", roles, "user_role", roleStr)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}
```

- [ ] **Step 2: Rewrite logging_middleware.go**

```go
package middleware

import (
	"time"

	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

type LoggingMiddleware struct {
	log logger.Logger
}

func NewLoggingMiddleware(log logger.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{log: log}
}

func (m *LoggingMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		m.log.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", time.Since(start).String(),
			"remote_addr", c.ClientIP(),
		)
	}
}
```

- [ ] **Step 3: Build to verify**

```bash
cd backend && go build ./...
```

- [ ] **Step 4: Commit**

```bash
git add backend/delivery/http/middleware/role_middleware.go backend/delivery/http/middleware/logging_middleware.go
git commit -m "refactor: convert role and logging middleware to gin.HandlerFunc"
```

---

### Task 5: Convert auth handler to gin

**Files:**
- Modify: `backend/delivery/http/handler/auth_handler.go`

- [ ] **Step 1: Rewrite auth_handler.go**

Key changes:
- Remove `chi` import, add `gin` import
- All methods: `func(h *AuthHandler) Method(c *gin.Context)`
- `chi.URLParam(r, "id")` → `c.Param("id")`
- `json.NewDecoder(r.Body).Decode(&req)` → `c.ShouldBindJSON(&req)`
- `r.Context().Value("user_id").(string)` → `c.GetString("user_id")`
- `response.JSON(w, ...)` → `response.JSON(c, ...)`
- `response.Error(w, ...)` → `response.Error(c, ...)`

```go
package handler

import (
	"net/http"

	"wms/application/dto"
	"wms/application/usecase"
	"wms/delivery/http/response"
	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUC usecase.AuthUsecase
	log    logger.Logger
}

func NewAuthHandler(authUC usecase.AuthUsecase, log logger.Logger) *AuthHandler {
	return &AuthHandler{authUC: authUC, log: log}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authUC.Login(c.Request.Context(), &dto.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, result)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authUC.Register(c.Request.Context(), &dto.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     req.Role,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, result)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.authUC.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, result)
}

func (h *AuthHandler) ListUsers(c *gin.Context) {
	page := 1
	limit := 50

	result, total, err := h.authUC.ListUsers(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(c, http.StatusOK, result, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authUC.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, result)
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.authUC.DeleteUser(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "user deactivated"})
}
```

- [ ] **Step 2: Build to verify**

```bash
cd backend && go build ./...
```

- [ ] **Step 3: Commit**

```bash
git add backend/delivery/http/handler/auth_handler.go
git commit -m "refactor: convert auth handler to use gin.Context"
```

---

### Task 6: Convert item & location handlers to gin

**Files:**
- Modify: `backend/delivery/http/handler/item_handler.go`
- Modify: `backend/delivery/http/handler/location_handler.go`

- [ ] **Step 1: Rewrite item_handler.go**

Same pattern as auth handler. All 5 methods converted. `chi.URLParam` → `c.Param`, `json.Decode` → `ShouldBindJSON`, `response.JSON(w,...)` → `response.JSON(c,...)`.

- [ ] **Step 2: Rewrite location_handler.go**

Same pattern. All 5 methods converted.

- [ ] **Step 3: Build to verify**

```bash
cd backend && go build ./...
```

- [ ] **Step 4: Commit**

```bash
git add backend/delivery/http/handler/item_handler.go backend/delivery/http/handler/location_handler.go
git commit -m "refactor: convert item and location handlers to use gin.Context"
```

---

### Task 7: Convert inventory, inbound, outbound, transfer, dashboard handlers to gin

**Files:**
- Modify: `backend/delivery/http/handler/inventory_handler.go`
- Modify: `backend/delivery/http/handler/inbound_handler.go`
- Modify: `backend/delivery/http/handler/outbound_handler.go`
- Modify: `backend/delivery/http/handler/transfer_handler.go`
- Modify: `backend/delivery/http/handler/dashboard_handler.go`

- [ ] **Step 1: Rewrite inventory_handler.go**

4 methods. `r.URL.Query().Get("x")` → `c.Query("x")`. `r.Context().Value("user_id")` → `c.GetString("user_id")`.

- [ ] **Step 2: Rewrite inbound_handler.go**

4 methods. Same pattern.

- [ ] **Step 3: Rewrite outbound_handler.go**

5 methods. Same pattern.

- [ ] **Step 4: Rewrite transfer_handler.go**

3 methods. Same pattern.

- [ ] **Step 5: Rewrite dashboard_handler.go**

1 method. Simplest — just change signature + response calls.

- [ ] **Step 6: Build to verify**

```bash
cd backend && go build ./...
```

- [ ] **Step 7: Commit**

```bash
git add backend/delivery/http/handler/
git commit -m "refactor: convert remaining handlers to use gin.Context"
```

---

### Task 8: Rewrite router with gin

**Files:**
- Modify: `backend/delivery/http/router/router.go`

- [ ] **Step 1: Rewrite router.go**

Key changes:
- `chi.NewRouter()` → `gin.Default()`
- `r.Route("/path", func(r chi.Router){...})` → `group := r.Group("/path")`
- `r.Use(middleware)` → `r.Use(ginMiddleware)`
- `r.Handle("/path", withRole(rm, roles, handler))` → `group.POST("/path", rm.RequireRole("admin", "manager"), handler)`
- `r.Get("/*", ...)` → `r.NoRoute(...)` for SPA fallback
- Remove `chimw` (chi middleware) import — gin has built-in Recovery, RequestID etc.
- `serveFrontend` uses `r.NoRoute()` for SPA fallback
- CORS handled via `github.com/gin-contrib/cors` or keep `rs/cors` as middleware

```go
package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"wms/config"
	"wms/delivery/http/handler"
	"wms/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	itemHandler *handler.ItemHandler,
	locationHandler *handler.LocationHandler,
	inventoryHandler *handler.InventoryHandler,
	inboundHandler *handler.InboundHandler,
	outboundHandler *handler.OutboundHandler,
	transferHandler *handler.TransferHandler,
	dashboardHandler *handler.DashboardHandler,
	authMiddleware *middleware.AuthMiddleware,
	roleMiddleware *middleware.RoleMiddleware,
	loggingMiddleware *middleware.LoggingMiddleware,
	corsConfig config.CORSConfig,
) *gin.Engine {
	r := gin.Default()

	r.Use(loggingMiddleware.Handle())

	c := cors.New(cors.Options{
		AllowedOrigins:   corsConfig.AllowedOrigins,
		AllowedMethods:   corsConfig.AllowedMethods,
		AllowedHeaders:   corsConfig.AllowedHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
	})
	r.Use(func(ctx *gin.Context) {
		c.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx.Request = r
			ctx.Next()
		})).ServeHTTP(ctx.Writer, ctx.Request)
	})

	adminRoles := []string{"superadmin"}
	managerRoles := []string{"superadmin", "manager"}
	operatorRoles := []string{"superadmin", "manager", "operator"}

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.Login)

		auth := api.Group("")
		auth.Use(authMiddleware.Handle())
		{
			auth.GET("/auth/me", authHandler.GetProfile)

			auth.GET("/auth/users", roleMiddleware.RequireRole(managerRoles...), authHandler.ListUsers)
			auth.PUT("/auth/users/:id", roleMiddleware.RequireRole(managerRoles...), authHandler.UpdateUser)
			auth.DELETE("/auth/users/:id", roleMiddleware.RequireRole(managerRoles...), authHandler.DeleteUser)

			items := auth.Group("/items")
			{
				items.GET("", itemHandler.List)
				items.POST("", roleMiddleware.RequireRole(managerRoles...), itemHandler.Create)
				items.GET("/:id", itemHandler.Get)
				items.PUT("/:id", roleMiddleware.RequireRole(managerRoles...), itemHandler.Update)
				items.DELETE("/:id", roleMiddleware.RequireRole(adminRoles...), itemHandler.Delete)
			}

			locations := auth.Group("/locations")
			{
				locations.GET("", locationHandler.List)
				locations.POST("", roleMiddleware.RequireRole(managerRoles...), locationHandler.Create)
				locations.GET("/:id", locationHandler.Get)
				locations.PUT("/:id", roleMiddleware.RequireRole(managerRoles...), locationHandler.Update)
				locations.DELETE("/:id", roleMiddleware.RequireRole(adminRoles...), locationHandler.Delete)
			}

			stock := auth.Group("/stock")
			{
				stock.GET("", inventoryHandler.ListStock)
				stock.GET("/movements", inventoryHandler.ListMovements)
				stock.POST("/adjust", roleMiddleware.RequireRole(managerRoles...), inventoryHandler.AdjustStock)
			}

			po := auth.Group("/purchase-orders")
			{
				po.GET("", inboundHandler.ListPurchaseOrders)
				po.POST("", roleMiddleware.RequireRole(managerRoles...), inboundHandler.CreatePurchaseOrder)
				po.GET("/:id", inboundHandler.GetPurchaseOrder)
				po.POST("/:id/receive", roleMiddleware.RequireRole(operatorRoles...), inboundHandler.ReceiveGoods)
			}

			so := auth.Group("/sales-orders")
			{
				so.GET("", outboundHandler.ListSalesOrders)
				so.POST("", roleMiddleware.RequireRole(managerRoles...), outboundHandler.CreateSalesOrder)
				so.GET("/:id", outboundHandler.GetSalesOrder)
				so.POST("/:id/pick", roleMiddleware.RequireRole(operatorRoles...), outboundHandler.PickOrder)
				so.POST("/:id/ship", roleMiddleware.RequireRole(operatorRoles...), outboundHandler.ShipOrder)
			}

			transfers := auth.Group("/transfers")
			{
				transfers.GET("", transferHandler.List)
				transfers.POST("", roleMiddleware.RequireRole(managerRoles...), transferHandler.Create)
				transfers.PUT("/:id/complete", roleMiddleware.RequireRole(operatorRoles...), transferHandler.Complete)
			}

			auth.GET("/dashboard/summary", dashboardHandler.GetSummary)
		}
	}

	serveFrontend(r, "../frontend/dist")

	return r
}

func serveFrontend(r *gin.Engine, distPath string) {
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		return
	}

	fileServer := http.FileServer(http.Dir(distPath))

	r.NoRoute(func(c *gin.Context) {
		path := filepath.Join(distPath, c.Request.URL.Path)

		if _, err := os.Stat(path); os.IsNotExist(err) || strings.HasSuffix(path, "/") {
			c.File(filepath.Join(distPath, "index.html"))
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
```

- [ ] **Step 2: Build to verify**

```bash
cd backend && go build ./...
```

- [ ] **Step 3: Commit**

```bash
git add backend/delivery/http/router/router.go
git commit -m "refactor: rewrite router with gin-gonic/gin"
```

---

### Task 9: Update main.go for gin

**Files:**
- Modify: `backend/cmd/wms/main.go`

- [ ] **Step 1: Update main.go**

Key changes:
- Router returns `*gin.Engine` instead of `http.Handler`
- `http.Server{Handler: mux}` stays the same (gin.Engine implements http.Handler)
- Remove `chimw` import
- Update middleware constructor calls

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wms/config"
	"wms/delivery/http/handler"
	"wms/delivery/http/middleware"
	"wms/delivery/http/router"
	"wms/infrastructure/auth"
	"wms/infrastructure/database"
	loggerInfra "wms/infrastructure/logger"
	postgresRepo "wms/infrastructure/repository"
	"wms/application/usecase"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	appLogger := loggerInfra.New(cfg.Logger.Level, cfg.Logger.Output)

	db, err := database.New(cfg.Database)
	if err != nil {
		appLogger.Fatal("failed to connect database", "error", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		appLogger.Fatal("failed to run migrations", "error", err)
	}

	jwtService := auth.NewJWTService(cfg.Auth)

	userRepo := postgresRepo.NewPostgresUserRepository(db, appLogger)
	itemRepo := postgresRepo.NewPostgresItemRepository(db, appLogger)
	locationRepo := postgresRepo.NewPostgresLocationRepository(db, appLogger)
	stockRepo := postgresRepo.NewPostgresStockRepository(db, appLogger)
	purchaseOrderRepo := postgresRepo.NewPostgresPurchaseOrderRepository(db, appLogger)
	salesOrderRepo := postgresRepo.NewPostgresSalesOrderRepository(db, appLogger)
	stockTransferRepo := postgresRepo.NewPostgresStockTransferRepository(db, appLogger)
	stockMovementRepo := postgresRepo.NewPostgresStockMovementRepository(db, appLogger)

	authUC := usecase.NewAuthUsecase(userRepo, jwtService, appLogger)
	itemUC := usecase.NewItemUsecase(itemRepo, appLogger)
	locationUC := usecase.NewLocationUsecase(locationRepo, appLogger)
	inventoryUC := usecase.NewInventoryUsecase(stockRepo, stockMovementRepo, appLogger)
	inboundUC := usecase.NewInboundUsecase(purchaseOrderRepo, stockRepo, stockMovementRepo, appLogger)
	outboundUC := usecase.NewOutboundUsecase(salesOrderRepo, stockRepo, stockMovementRepo, appLogger)
	transferUC := usecase.NewTransferUsecase(stockTransferRepo, stockRepo, stockMovementRepo, appLogger)
	dashboardUC := usecase.NewDashboardUsecase(stockRepo, appLogger)

	authHandler := handler.NewAuthHandler(authUC, appLogger)
	itemHandler := handler.NewItemHandler(itemUC, appLogger)
	locationHandler := handler.NewLocationHandler(locationUC, appLogger)
	inventoryHandler := handler.NewInventoryHandler(inventoryUC, appLogger)
	inboundHandler := handler.NewInboundHandler(inboundUC, appLogger)
	outboundHandler := handler.NewOutboundHandler(outboundUC, appLogger)
	transferHandler := handler.NewTransferHandler(transferUC, appLogger)
	dashboardHandler := handler.NewDashboardHandler(dashboardUC, appLogger)

	authMiddleware := middleware.NewAuthMiddleware(jwtService, db, appLogger)
	roleMiddleware := middleware.NewRoleMiddleware(appLogger)
	loggingMiddleware := middleware.NewLoggingMiddleware(appLogger)

	mux := router.NewRouter(
		authHandler,
		itemHandler,
		locationHandler,
		inventoryHandler,
		inboundHandler,
		outboundHandler,
		transferHandler,
		dashboardHandler,
		authMiddleware,
		roleMiddleware,
		loggingMiddleware,
		cfg.CORS,
	)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		appLogger.Info("server starting", "port", cfg.App.Port, "env", cfg.App.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	appLogger.Info("server stopped")
}
```

- [ ] **Step 2: Build to verify**

```bash
cd backend && go build ./...
```

- [ ] **Step 3: Commit**

```bash
git add backend/cmd/wms/main.go
git commit -m "refactor: update main.go for gin router"
```

---

### Task 10: Full build, test, and clean up

- [ ] **Step 1: Full build**

```bash
cd backend && go build ./...
```

- [ ] **Step 2: Tidy go.mod**

```bash
cd backend && go mod tidy
```

- [ ] **Step 3: Start and test**

```bash
docker compose -f docker-compose.yml down -v && docker compose -f docker-compose.yml up -d
# Wait for backend to start, then test:
curl http://localhost:8080/api/v1/auth/login -X POST -H 'Content-Type: application/json' -d '{"username":"admin","password":"admin123"}'
```

- [ ] **Step 4: Update note.md**

Add migration status to note.md.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "chore: clean up go.mod after chi to gin migration"
```

---

## Migration Summary

| File | Changes |
|------|---------|
| `go.mod` | `chi/v5` → `gin v1.10.0` |
| `response.go` | `http.ResponseWriter` → `*gin.Context` |
| `auth_middleware.go` | `func(http.Handler) http.Handler` → `gin.HandlerFunc` |
| `role_middleware.go` | Same pattern, `RequireRole` returns `gin.HandlerFunc` |
| `logging_middleware.go` | Same pattern, uses `c.Writer.Status()` |
| `auth_handler.go` | 6 methods, `chi.URLParam` → `c.Param`, `ShouldBindJSON` |
| `item_handler.go` | 5 methods, same pattern |
| `location_handler.go` | 5 methods, same pattern |
| `inventory_handler.go` | 4 methods, `c.Query()` |
| `inbound_handler.go` | 4 methods |
| `outbound_handler.go` | 5 methods |
| `transfer_handler.go` | 3 methods |
| `dashboard_handler.go` | 1 method |
| `router.go` | Complete rewrite, `gin.Default()`, `Group()` |
| `main.go` | Minimal changes, router returns `*gin.Engine` |
