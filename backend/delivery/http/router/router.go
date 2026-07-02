package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wms/config"
	"wms/delivery/http/handler"
	"wms/delivery/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     corsConfig.AllowedOrigins,
		AllowMethods:     corsConfig.AllowedMethods,
		AllowHeaders:     corsConfig.AllowedHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           12 * time.Hour,
	}))

	adminRoles := []string{"admin"}
	managerRoles := []string{"admin", "manager"}
	operatorRoles := []string{"admin", "manager", "operator"}

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
				stock.POST("/opname", roleMiddleware.RequireRole(operatorRoles...), inventoryHandler.StockOpname)
			}

			po := auth.Group("/purchase-orders")
			{
				po.GET("", inboundHandler.ListPurchaseOrders)
				po.POST("", roleMiddleware.RequireRole(managerRoles...), inboundHandler.CreatePurchaseOrder)
				po.GET("/:id", inboundHandler.GetPurchaseOrder)
				po.POST("/:id/approve", roleMiddleware.RequireRole(managerRoles...), inboundHandler.ApprovePurchaseOrder)
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
