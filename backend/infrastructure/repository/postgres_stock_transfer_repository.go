package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wms/domain/entity"
	"wms/pkg/logger"
)

type postgresStockTransferRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewPostgresStockTransferRepository(db *sql.DB, log logger.Logger) *postgresStockTransferRepository {
	return &postgresStockTransferRepository{db: db, log: log}
}

func (r *postgresStockTransferRepository) FindByID(ctx context.Context, id string) (*entity.StockTransfer, error) {
	t := &entity.StockTransfer{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, transfer_number, from_location_id, to_location_id, item_id, quantity, status, notes, created_by, created_at 
		 FROM stock_transfers WHERE id = $1`, id,
	).Scan(&t.ID, &t.TransferNumber, &t.FromLocationID, &t.ToLocationID, &t.ItemID, &t.Quantity, &t.Status, &t.Notes, &t.CreatedBy, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("transfer not found: %w", err)
	}
	return t, nil
}

func (r *postgresStockTransferRepository) Create(ctx context.Context, t *entity.StockTransfer) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO stock_transfers (transfer_number, from_location_id, to_location_id, item_id, quantity, status, notes, created_by, created_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		t.TransferNumber, t.FromLocationID, t.ToLocationID, t.ItemID, t.Quantity, t.Status, t.Notes, t.CreatedBy, t.CreatedAt,
	).Scan(&t.ID)
	if err != nil {
		return fmt.Errorf("failed to create transfer: %w", err)
	}
	return nil
}

func (r *postgresStockTransferRepository) Update(ctx context.Context, t *entity.StockTransfer) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE stock_transfers SET status=$1, notes=$2 WHERE id=$3`,
		t.Status, t.Notes, t.ID,
	)
	return err
}

func (r *postgresStockTransferRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE stock_transfers SET status = $1 WHERE id = $2", status, id,
	)
	return err
}

func (r *postgresStockTransferRepository) List(ctx context.Context, page, limit int) ([]*entity.StockTransfer, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stock_transfers").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, transfer_number, from_location_id, to_location_id, item_id, quantity, status, notes, created_by, created_at 
		 FROM stock_transfers ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transfers []*entity.StockTransfer
	for rows.Next() {
		t := &entity.StockTransfer{}
		if err := rows.Scan(&t.ID, &t.TransferNumber, &t.FromLocationID, &t.ToLocationID, &t.ItemID, &t.Quantity, &t.Status, &t.Notes, &t.CreatedBy, &t.CreatedAt); err != nil {
			continue
		}
		transfers = append(transfers, t)
	}

	return transfers, total, nil
}
