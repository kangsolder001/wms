package usecase

import (
	"context"

	"wms/domain/repository"
	"wms/pkg/logger"
)

type DashboardUsecase interface {
	GetSummary(ctx context.Context) (map[string]interface{}, error)
}

type dashboardUsecase struct {
	stockRepo repository.StockRepository
	log       logger.Logger
}

func NewDashboardUsecase(stockRepo repository.StockRepository, log logger.Logger) DashboardUsecase {
	return &dashboardUsecase{stockRepo: stockRepo, log: log}
}

func (uc *dashboardUsecase) GetSummary(ctx context.Context) (map[string]interface{}, error) {
	uc.log.Info("getting dashboard summary")

	stocks, total, err := uc.stockRepo.List(ctx, 1, 1000)
	if err != nil {
		uc.log.Error("failed to get stock for summary", "error", err)
		return nil, err
	}

	totalQuantity := 0.0
	for _, s := range stocks {
		totalQuantity += s.Quantity
	}

	summary := map[string]interface{}{
		"total_stock_items": total,
		"total_quantity":    totalQuantity,
	}

	return summary, nil
}
