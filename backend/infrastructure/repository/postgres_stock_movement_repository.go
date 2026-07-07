package repository

import (
	"context"
	"database/sql"

	"wms/domain/entity"
	"wms/pkg/logger"
)

type postgresStockMovementRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewPostgresStockMovementRepository(db *sql.DB, log logger.Logger) *postgresStockMovementRepository {
	return &postgresStockMovementRepository{db: db, log: log}
}

func (r *postgresStockMovementRepository) Create(ctx context.Context, movement *entity.StockMovement) error {
	var refID sql.NullString
	if movement.ReferenceID != nil {
		refID = sql.NullString{String: *movement.ReferenceID, Valid: true}
	}
	
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO stock_movements (item_id, from_location_id, to_location_id, quantity, movement_type, reference_type, reference_id, notes, created_by, created_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		movement.ItemID, movement.FromLocationID, movement.ToLocationID, movement.Quantity, movement.MovementType, movement.ReferenceType, refID, movement.Notes, movement.CreatedBy, movement.CreatedAt,
	).Scan(&movement.ID)
	return err
}

func (r *postgresStockMovementRepository) ListByItem(ctx context.Context, itemID string, page, limit int) ([]*entity.StockMovement, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stock_movements WHERE item_id = $1", itemID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, item_id, from_location_id, to_location_id, quantity, movement_type, reference_type, reference_id, notes, created_by, created_at 
		 FROM stock_movements WHERE item_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, itemID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movements []*entity.StockMovement
	for rows.Next() {
		m := &entity.StockMovement{}
		var refID sql.NullString
		if err := rows.Scan(&m.ID, &m.ItemID, &m.FromLocationID, &m.ToLocationID, &m.Quantity, &m.MovementType, &m.ReferenceType, &refID, &m.Notes, &m.CreatedBy, &m.CreatedAt); err != nil {
			continue
		}
		if refID.Valid {
			m.ReferenceID = &refID.String
		}
		movements = append(movements, m)
	}

	return movements, total, nil
}

func (r *postgresStockMovementRepository) List(ctx context.Context, page, limit int) ([]*entity.StockMovement, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stock_movements").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, item_id, from_location_id, to_location_id, quantity, movement_type, reference_type, reference_id, notes, created_by, created_at 
		 FROM stock_movements ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movements []*entity.StockMovement
	for rows.Next() {
		m := &entity.StockMovement{}
		var refID sql.NullString
		if err := rows.Scan(&m.ID, &m.ItemID, &m.FromLocationID, &m.ToLocationID, &m.Quantity, &m.MovementType, &m.ReferenceType, &refID, &m.Notes, &m.CreatedBy, &m.CreatedAt); err != nil {
			continue
		}
		if refID.Valid {
			m.ReferenceID = &refID.String
		}
		movements = append(movements, m)
	}

	return movements, total, nil
}
