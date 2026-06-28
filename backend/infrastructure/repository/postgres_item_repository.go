package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wms/domain/entity"
	"wms/pkg/logger"
)

type postgresItemRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewPostgresItemRepository(db *sql.DB, log logger.Logger) *postgresItemRepository {
	return &postgresItemRepository{db: db, log: log}
}

func (r *postgresItemRepository) FindByID(ctx context.Context, id string) (*entity.Item, error) {
	item := &entity.Item{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, sku, name, description, category, unit_of_measure, weight, is_active, created_at, updated_at 
		 FROM items WHERE id = $1`, id,
	).Scan(&item.ID, &item.SKU, &item.Name, &item.Description, &item.Category, &item.UnitOfMeasure, &item.Weight, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("item not found: %w", err)
	}
	return item, nil
}

func (r *postgresItemRepository) FindBySKU(ctx context.Context, sku string) (*entity.Item, error) {
	item := &entity.Item{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, sku, name, description, category, unit_of_measure, weight, is_active, created_at, updated_at 
		 FROM items WHERE sku = $1`, sku,
	).Scan(&item.ID, &item.SKU, &item.Name, &item.Description, &item.Category, &item.UnitOfMeasure, &item.Weight, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("item not found: %w", err)
	}
	return item, nil
}

func (r *postgresItemRepository) Create(ctx context.Context, item *entity.Item) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO items (sku, name, description, category, unit_of_measure, weight, is_active, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		item.SKU, item.Name, item.Description, item.Category, item.UnitOfMeasure, item.Weight, item.IsActive, item.CreatedAt, item.UpdatedAt,
	).Scan(&item.ID)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}
	return nil
}

func (r *postgresItemRepository) Update(ctx context.Context, item *entity.Item) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE items SET name=$1, description=$2, category=$3, unit_of_measure=$4, weight=$5, is_active=$6, updated_at=$7 WHERE id=$8`,
		item.Name, item.Description, item.Category, item.UnitOfMeasure, item.Weight, item.IsActive, item.UpdatedAt, item.ID,
	)
	return err
}

func (r *postgresItemRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE items SET is_active = false WHERE id = $1", id)
	return err
}

func (r *postgresItemRepository) List(ctx context.Context, page, limit int) ([]*entity.Item, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM items WHERE is_active = true").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, sku, name, description, category, unit_of_measure, weight, is_active, created_at, updated_at 
		 FROM items WHERE is_active = true ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*entity.Item
	for rows.Next() {
		i := &entity.Item{}
		if err := rows.Scan(&i.ID, &i.SKU, &i.Name, &i.Description, &i.Category, &i.UnitOfMeasure, &i.Weight, &i.IsActive, &i.CreatedAt, &i.UpdatedAt); err != nil {
			continue
		}
		items = append(items, i)
	}

	return items, total, nil
}
