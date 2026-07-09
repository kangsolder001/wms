package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wms/domain/entity"
	"wms/pkg/logger"
)

type postgresCategoryRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewPostgresCategoryRepository(db *sql.DB, log logger.Logger) *postgresCategoryRepository {
	return &postgresCategoryRepository{db: db, log: log}
}

func (r *postgresCategoryRepository) FindByID(ctx context.Context, id string) (*entity.Category, error) {
	c := &entity.Category{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, abbreviation, description, is_active, created_at FROM categories WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Abbreviation, &c.Description, &c.IsActive, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return c, nil
}

func (r *postgresCategoryRepository) FindByAbbreviation(ctx context.Context, abbreviation string) (*entity.Category, error) {
	c := &entity.Category{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, abbreviation, description, is_active, created_at FROM categories WHERE abbreviation = $1`, abbreviation,
	).Scan(&c.ID, &c.Name, &c.Abbreviation, &c.Description, &c.IsActive, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return c, nil
}

func (r *postgresCategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO categories (name, abbreviation, description, is_active, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		category.Name, category.Abbreviation, category.Description, category.IsActive, category.CreatedAt,
	).Scan(&category.ID)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *postgresCategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE categories SET name=$1, abbreviation=$2, description=$3, is_active=$4 WHERE id=$5`,
		category.Name, category.Abbreviation, category.Description, category.IsActive, category.ID,
	)
	return err
}

func (r *postgresCategoryRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE categories SET is_active = false WHERE id = $1", id)
	return err
}

func (r *postgresCategoryRepository) List(ctx context.Context, page, limit int) ([]*entity.Category, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM categories WHERE is_active = true").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, abbreviation, description, is_active, created_at FROM categories WHERE is_active = true ORDER BY name ASC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		c := &entity.Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Abbreviation, &c.Description, &c.IsActive, &c.CreatedAt); err != nil {
			continue
		}
		categories = append(categories, c)
	}

	return categories, total, nil
}

func (r *postgresCategoryRepository) ListAll(ctx context.Context) ([]*entity.Category, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, abbreviation, description, is_active, created_at FROM categories WHERE is_active = true ORDER BY name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		c := &entity.Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Abbreviation, &c.Description, &c.IsActive, &c.CreatedAt); err != nil {
			continue
		}
		categories = append(categories, c)
	}

	return categories, nil
}
