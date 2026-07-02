package repository

import (
	"context"

	"wms/domain/entity"
)

type StockRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Stock, error)
	FindByItemAndLocation(ctx context.Context, itemID, locationID string) (*entity.Stock, error)
	FindByItem(ctx context.Context, itemID string) ([]*entity.Stock, error)
	Create(ctx context.Context, stock *entity.Stock) error
	UpdateQuantity(ctx context.Context, id string, quantity float64) error
	Reserve(ctx context.Context, itemID, locationID string, quantity float64) error
	Release(ctx context.Context, itemID, locationID string, quantity float64) error
	List(ctx context.Context, page, limit int) ([]*entity.Stock, int, error)
	ListWithDetails(ctx context.Context, page, limit int) ([]map[string]interface{}, int, error)
	GetTotalStockByItem(ctx context.Context, itemID string) (float64, error)
}

type StockMovementRepository interface {
	Create(ctx context.Context, movement *entity.StockMovement) error
	ListByItem(ctx context.Context, itemID string, page, limit int) ([]*entity.StockMovement, int, error)
	List(ctx context.Context, page, limit int) ([]*entity.StockMovement, int, error)
}
