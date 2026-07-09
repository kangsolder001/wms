package repository

import (
	"context"

	"wms/domain/entity"
)

type CategoryRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Category, error)
	FindByAbbreviation(ctx context.Context, abbreviation string) (*entity.Category, error)
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*entity.Category, int, error)
	ListAll(ctx context.Context) ([]*entity.Category, error)
}
