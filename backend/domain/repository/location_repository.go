package repository

import (
	"context"

	"wms/domain/entity"
)

type LocationRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Location, error)
	FindByCode(ctx context.Context, code string) (*entity.Location, error)
	Create(ctx context.Context, location *entity.Location) error
	Update(ctx context.Context, location *entity.Location) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*entity.Location, int, error)
	ListByZone(ctx context.Context, zone string) ([]*entity.Location, error)
}
