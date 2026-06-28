package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wms/domain/entity"
	"wms/pkg/logger"
)

type postgresLocationRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewPostgresLocationRepository(db *sql.DB, log logger.Logger) *postgresLocationRepository {
	return &postgresLocationRepository{db: db, log: log}
}

func (r *postgresLocationRepository) FindByID(ctx context.Context, id string) (*entity.Location, error) {
	loc := &entity.Location{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, code, name, zone, aisle, rack, level, bin, type, capacity, is_active, created_at 
		 FROM locations WHERE id = $1`, id,
	).Scan(&loc.ID, &loc.Code, &loc.Name, &loc.Zone, &loc.Aisle, &loc.Rack, &loc.Level, &loc.Bin, &loc.Type, &loc.Capacity, &loc.IsActive, &loc.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("location not found: %w", err)
	}
	return loc, nil
}

func (r *postgresLocationRepository) FindByCode(ctx context.Context, code string) (*entity.Location, error) {
	loc := &entity.Location{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, code, name, zone, aisle, rack, level, bin, type, capacity, is_active, created_at 
		 FROM locations WHERE code = $1`, code,
	).Scan(&loc.ID, &loc.Code, &loc.Name, &loc.Zone, &loc.Aisle, &loc.Rack, &loc.Level, &loc.Bin, &loc.Type, &loc.Capacity, &loc.IsActive, &loc.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("location not found: %w", err)
	}
	return loc, nil
}

func (r *postgresLocationRepository) Create(ctx context.Context, loc *entity.Location) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO locations (code, name, zone, aisle, rack, level, bin, type, capacity, is_active, created_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		loc.Code, loc.Name, loc.Zone, loc.Aisle, loc.Rack, loc.Level, loc.Bin, loc.Type, loc.Capacity, loc.IsActive, loc.CreatedAt,
	).Scan(&loc.ID)
	if err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}
	return nil
}

func (r *postgresLocationRepository) Update(ctx context.Context, loc *entity.Location) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE locations SET name=$1, zone=$2, aisle=$3, rack=$4, level=$5, bin=$6, type=$7, capacity=$8, is_active=$9 WHERE id=$10`,
		loc.Name, loc.Zone, loc.Aisle, loc.Rack, loc.Level, loc.Bin, loc.Type, loc.Capacity, loc.IsActive, loc.ID,
	)
	return err
}

func (r *postgresLocationRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE locations SET is_active = false WHERE id = $1", id)
	return err
}

func (r *postgresLocationRepository) List(ctx context.Context, page, limit int) ([]*entity.Location, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM locations WHERE is_active = true").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, code, name, zone, aisle, rack, level, bin, type, capacity, is_active, created_at 
		 FROM locations WHERE is_active = true ORDER BY code ASC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var locations []*entity.Location
	for rows.Next() {
		l := &entity.Location{}
		if err := rows.Scan(&l.ID, &l.Code, &l.Name, &l.Zone, &l.Aisle, &l.Rack, &l.Level, &l.Bin, &l.Type, &l.Capacity, &l.IsActive, &l.CreatedAt); err != nil {
			continue
		}
		locations = append(locations, l)
	}

	return locations, total, nil
}

func (r *postgresLocationRepository) ListByZone(ctx context.Context, zone string) ([]*entity.Location, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, code, name, zone, aisle, rack, level, bin, type, capacity, is_active, created_at 
		 FROM locations WHERE zone = $1 AND is_active = true ORDER BY code ASC`, zone,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*entity.Location
	for rows.Next() {
		l := &entity.Location{}
		if err := rows.Scan(&l.ID, &l.Code, &l.Name, &l.Zone, &l.Aisle, &l.Rack, &l.Level, &l.Bin, &l.Type, &l.Capacity, &l.IsActive, &l.CreatedAt); err != nil {
			continue
		}
		locations = append(locations, l)
	}

	return locations, nil
}
