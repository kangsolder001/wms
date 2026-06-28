package repository

import (
	"context"

	"wms/domain/entity"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	List(ctx context.Context, page, limit int) ([]*entity.User, int, error)
}
