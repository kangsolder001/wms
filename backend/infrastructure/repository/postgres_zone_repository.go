package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wms/domain/entity"
	"wms/pkg/logger"
)

type postgresZoneRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewPostgresZoneRepository(db *sql.DB, log logger.Logger) *postgresZoneRepository {
	return &postgresZoneRepository{db: db, log: log}
}

func (r *postgresZoneRepository) FindByID(ctx context.Context, id string) (*entity.Zone, error) {
	z := &entity.Zone{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, code, name, description, is_active, created_at FROM zones WHERE id = $1`, id,
	).Scan(&z.ID, &z.Code, &z.Name, &z.Description, &z.IsActive, &z.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("zone not found: %w", err)
	}
	return z, nil
}

func (r *postgresZoneRepository) Create(ctx context.Context, zone *entity.Zone) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO zones (code, name, description, is_active, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		zone.Code, zone.Name, zone.Description, zone.IsActive, zone.CreatedAt,
	).Scan(&zone.ID)
	if err != nil {
		return fmt.Errorf("failed to create zone: %w", err)
	}
	return nil
}

func (r *postgresZoneRepository) Update(ctx context.Context, zone *entity.Zone) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE zones SET code=$1, name=$2, description=$3, is_active=$4 WHERE id=$5`,
		zone.Code, zone.Name, zone.Description, zone.IsActive, zone.ID,
	)
	return err
}

func (r *postgresZoneRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE zones SET is_active = false WHERE id = $1", id)
	return err
}

func (r *postgresZoneRepository) List(ctx context.Context, page, limit int) ([]*entity.Zone, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM zones WHERE is_active = true").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, code, name, description, is_active, created_at FROM zones WHERE is_active = true ORDER BY code ASC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var zones []*entity.Zone
	for rows.Next() {
		z := &entity.Zone{}
		if err := rows.Scan(&z.ID, &z.Code, &z.Name, &z.Description, &z.IsActive, &z.CreatedAt); err != nil {
			continue
		}
		zones = append(zones, z)
	}

	return zones, total, nil
}

func (r *postgresZoneRepository) ListAll(ctx context.Context) ([]*entity.Zone, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, code, name, description, is_active, created_at FROM zones WHERE is_active = true ORDER BY code ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zones []*entity.Zone
	for rows.Next() {
		z := &entity.Zone{}
		if err := rows.Scan(&z.ID, &z.Code, &z.Name, &z.Description, &z.IsActive, &z.CreatedAt); err != nil {
			continue
		}
		zones = append(zones, z)
	}

	return zones, nil
}
