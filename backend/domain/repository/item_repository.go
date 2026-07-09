package repository

import (
	"context"

	"wms/domain/entity"
)

type ItemRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Item, error)
	FindBySKU(ctx context.Context, sku string) (*entity.Item, error)
	GetNextSKUSequence(ctx context.Context, categoryAbbreviation string) (int, error)
	Create(ctx context.Context, item *entity.Item) error
	Update(ctx context.Context, item *entity.Item) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*entity.Item, int, error)
}
