package repository

import (
	"context"

	"wms/domain/entity"
)

type StockTransferRepository interface {
	FindByID(ctx context.Context, id string) (*entity.StockTransfer, error)
	Create(ctx context.Context, transfer *entity.StockTransfer) error
	Update(ctx context.Context, transfer *entity.StockTransfer) error
	UpdateStatus(ctx context.Context, id, status string) error
	List(ctx context.Context, page, limit int) ([]*entity.StockTransfer, int, error)
}
