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
	redisInfra "wms/infrastructure/redis"
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

	redisClient, err := redisInfra.NewRedisClient(redisInfra.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		appLogger.Fatal("failed to connect redis", "error", err)
	}
	defer redisClient.Close()

	jwtService := auth.NewJWTService(cfg.Auth)

	userRepo := postgresRepo.NewPostgresUserRepository(db, appLogger)
	itemRepo := postgresRepo.NewPostgresItemRepository(db, appLogger)
	locationRepo := postgresRepo.NewPostgresLocationRepository(db, appLogger)
	categoryRepo := postgresRepo.NewPostgresCategoryRepository(db, appLogger)
	zoneRepo := postgresRepo.NewPostgresZoneRepository(db, appLogger)
	stockRepo := postgresRepo.NewPostgresStockRepository(db, appLogger)
	purchaseOrderRepo := postgresRepo.NewPostgresPurchaseOrderRepository(db, appLogger)
	salesOrderRepo := postgresRepo.NewPostgresSalesOrderRepository(db, appLogger)
	stockTransferRepo := postgresRepo.NewPostgresStockTransferRepository(db, appLogger)
	stockMovementRepo := postgresRepo.NewPostgresStockMovementRepository(db, appLogger)

	authUC := usecase.NewAuthUsecase(userRepo, jwtService, appLogger)
	itemUC := usecase.NewItemUsecase(itemRepo, categoryRepo, appLogger)
	locationUC := usecase.NewLocationUsecase(locationRepo, appLogger)
	categoryUC := usecase.NewCategoryUsecase(categoryRepo, appLogger)
	zoneUC := usecase.NewZoneUsecase(zoneRepo, appLogger)
	inventoryUC := usecase.NewInventoryUsecase(stockRepo, stockMovementRepo, appLogger)
	inboundUC := usecase.NewInboundUsecase(purchaseOrderRepo, stockRepo, stockMovementRepo, appLogger)
	outboundUC := usecase.NewOutboundUsecase(salesOrderRepo, stockRepo, stockMovementRepo, appLogger)
	transferUC := usecase.NewTransferUsecase(stockTransferRepo, stockRepo, stockMovementRepo, appLogger)
	dashboardUC := usecase.NewDashboardUsecase(stockRepo, appLogger)

	authHandler := handler.NewAuthHandler(authUC, appLogger)
	itemHandler := handler.NewItemHandler(itemUC, appLogger)
	locationHandler := handler.NewLocationHandler(locationUC, appLogger)
	categoryHandler := handler.NewCategoryHandler(categoryUC, appLogger)
	zoneHandler := handler.NewZoneHandler(zoneUC, appLogger)
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
		categoryHandler,
		zoneHandler,
		inventoryHandler,
		inboundHandler,
		outboundHandler,
		transferHandler,
		dashboardHandler,
		authMiddleware,
		roleMiddleware,
		loggingMiddleware,
		redisClient,
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
