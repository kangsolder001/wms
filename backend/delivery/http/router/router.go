package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"wms/config"
	"wms/delivery/http/handler"
	"wms/delivery/http/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func withRole(rm *middleware.RoleMiddleware, roles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rm.RequireRole(roles...).Middleware(next).ServeHTTP(w, r)
	}
}

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
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(loggingMiddleware.Handle)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Heartbeat("/health"))

	c := cors.New(cors.Options{
		AllowedOrigins:   corsConfig.AllowedOrigins,
		AllowedMethods:   corsConfig.AllowedMethods,
		AllowedHeaders:   corsConfig.AllowedHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
	})
	r.Use(c.Handler)

	adminRoles := []string{"admin"}
	managerRoles := []string{"admin", "manager"}
	operatorRoles := []string{"admin", "manager", "operator"}

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Handle)

			r.Handle("/auth/register", withRole(roleMiddleware, adminRoles, authHandler.Register))
			r.Get("/auth/me", authHandler.GetProfile)

			r.Route("/items", func(r chi.Router) {
				r.Get("/", itemHandler.List)
				r.Post("/", withRole(roleMiddleware, managerRoles, itemHandler.Create))
				r.Get("/{id}", itemHandler.Get)
				r.Put("/{id}", withRole(roleMiddleware, managerRoles, itemHandler.Update))
				r.Delete("/{id}", withRole(roleMiddleware, adminRoles, itemHandler.Delete))
			})

			r.Route("/locations", func(r chi.Router) {
				r.Get("/", locationHandler.List)
				r.Post("/", withRole(roleMiddleware, managerRoles, locationHandler.Create))
				r.Get("/{id}", locationHandler.Get)
				r.Put("/{id}", withRole(roleMiddleware, managerRoles, locationHandler.Update))
				r.Delete("/{id}", withRole(roleMiddleware, adminRoles, locationHandler.Delete))
			})

			r.Route("/stock", func(r chi.Router) {
				r.Get("/", inventoryHandler.ListStock)
				r.Get("/movements", inventoryHandler.ListMovements)
				r.Post("/adjust", withRole(roleMiddleware, managerRoles, inventoryHandler.AdjustStock))
			})

			r.Route("/purchase-orders", func(r chi.Router) {
				r.Get("/", inboundHandler.ListPurchaseOrders)
				r.Post("/", withRole(roleMiddleware, managerRoles, inboundHandler.CreatePurchaseOrder))
				r.Get("/{id}", inboundHandler.GetPurchaseOrder)
				r.Post("/{id}/receive", withRole(roleMiddleware, operatorRoles, inboundHandler.ReceiveGoods))
			})

			r.Route("/sales-orders", func(r chi.Router) {
				r.Get("/", outboundHandler.ListSalesOrders)
				r.Post("/", withRole(roleMiddleware, managerRoles, outboundHandler.CreateSalesOrder))
				r.Get("/{id}", outboundHandler.GetSalesOrder)
				r.Post("/{id}/pick", withRole(roleMiddleware, operatorRoles, outboundHandler.PickOrder))
				r.Post("/{id}/ship", withRole(roleMiddleware, operatorRoles, outboundHandler.ShipOrder))
			})

			r.Route("/transfers", func(r chi.Router) {
				r.Get("/", transferHandler.List)
				r.Post("/", withRole(roleMiddleware, managerRoles, transferHandler.Create))
				r.Put("/{id}/complete", withRole(roleMiddleware, operatorRoles, transferHandler.Complete))
			})

			r.Get("/dashboard/summary", dashboardHandler.GetSummary)
		})
	})

	serveFrontend(r, "../frontend/dist")

	return r
}

func serveFrontend(r chi.Router, distPath string) {
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		return
	}

	fileServer := http.FileServer(http.Dir(distPath))

	r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		path := filepath.Join(distPath, req.URL.Path)

		if _, err := os.Stat(path); os.IsNotExist(err) || strings.HasSuffix(path, "/") {
			http.ServeFile(w, req, filepath.Join(distPath, "index.html"))
			return
		}

		fileServer.ServeHTTP(w, req)
	})
}
