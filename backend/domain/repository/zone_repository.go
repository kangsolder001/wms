package repository

import (
	"context"

	"wms/domain/entity"
)

type ZoneRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Zone, error)
	Create(ctx context.Context, zone *entity.Zone) error
	Update(ctx context.Context, zone *entity.Zone) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*entity.Zone, int, error)
	ListAll(ctx context.Context) ([]*entity.Zone, error)
}
