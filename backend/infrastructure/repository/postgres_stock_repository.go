package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wms/domain/entity"
	"wms/pkg/logger"
)

type postgresStockRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewPostgresStockRepository(db *sql.DB, log logger.Logger) *postgresStockRepository {
	return &postgresStockRepository{db: db, log: log}
}

func (r *postgresStockRepository) FindByID(ctx context.Context, id string) (*entity.Stock, error) {
	s := &entity.Stock{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, item_id, location_id, quantity, reserved_quantity, batch_number, expiry_date, updated_at 
		 FROM stock WHERE id = $1`, id,
	).Scan(&s.ID, &s.ItemID, &s.LocationID, &s.Quantity, &s.ReservedQuantity, &s.BatchNumber, &s.ExpiryDate, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("stock not found: %w", err)
	}
	return s, nil
}

func (r *postgresStockRepository) FindByItemAndLocation(ctx context.Context, itemID, locationID string) (*entity.Stock, error) {
	s := &entity.Stock{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, item_id, location_id, quantity, reserved_quantity, batch_number, expiry_date, updated_at 
		 FROM stock WHERE item_id = $1 AND location_id = $2`, itemID, locationID,
	).Scan(&s.ID, &s.ItemID, &s.LocationID, &s.Quantity, &s.ReservedQuantity, &s.BatchNumber, &s.ExpiryDate, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("stock not found: %w", err)
	}
	return s, nil
}

func (r *postgresStockRepository) FindByItem(ctx context.Context, itemID string) ([]*entity.Stock, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, item_id, location_id, quantity, reserved_quantity, batch_number, expiry_date, updated_at 
		 FROM stock WHERE item_id = $1 AND quantity > 0`, itemID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []*entity.Stock
	for rows.Next() {
		s := &entity.Stock{}
		if err := rows.Scan(&s.ID, &s.ItemID, &s.LocationID, &s.Quantity, &s.ReservedQuantity, &s.BatchNumber, &s.ExpiryDate, &s.UpdatedAt); err != nil {
			continue
		}
		stocks = append(stocks, s)
	}

	return stocks, nil
}

func (r *postgresStockRepository) Create(ctx context.Context, stock *entity.Stock) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO stock (item_id, location_id, quantity, reserved_quantity, batch_number, expiry_date, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		stock.ItemID, stock.LocationID, stock.Quantity, stock.ReservedQuantity, stock.BatchNumber, stock.ExpiryDate, stock.UpdatedAt,
	).Scan(&stock.ID)
	if err != nil {
		return fmt.Errorf("failed to create stock: %w", err)
	}
	return nil
}

func (r *postgresStockRepository) UpdateQuantity(ctx context.Context, id string, quantity float64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE stock SET quantity = $1, updated_at = NOW() WHERE id = $2", quantity, id,
	)
	return err
}

func (r *postgresStockRepository) Reserve(ctx context.Context, itemID, locationID string, quantity float64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE stock SET reserved_quantity = reserved_quantity + $1 WHERE item_id = $2 AND location_id = $3",
		quantity, itemID, locationID,
	)
	return err
}

func (r *postgresStockRepository) Release(ctx context.Context, itemID, locationID string, quantity float64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE stock SET reserved_quantity = reserved_quantity - $1 WHERE item_id = $2 AND location_id = $3",
		quantity, itemID, locationID,
	)
	return err
}

func (r *postgresStockRepository) List(ctx context.Context, page, limit int) ([]*entity.Stock, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stock WHERE quantity > 0").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT s.id, s.item_id, s.location_id, s.quantity, s.reserved_quantity, s.batch_number, s.expiry_date, s.updated_at 
		 FROM stock s WHERE s.quantity > 0 ORDER BY s.updated_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var stocks []*entity.Stock
	for rows.Next() {
		s := &entity.Stock{}
		if err := rows.Scan(&s.ID, &s.ItemID, &s.LocationID, &s.Quantity, &s.ReservedQuantity, &s.BatchNumber, &s.ExpiryDate, &s.UpdatedAt); err != nil {
			continue
		}
		stocks = append(stocks, s)
	}

	return stocks, total, nil
}

func (r *postgresStockRepository) ListWithDetails(ctx context.Context, page, limit int, itemID, locationID, search string) ([]map[string]interface{}, int, error) {
	offset := (page - 1) * limit

	where := "WHERE s.quantity > 0"
	args := []interface{}{}
	argIdx := 1

	if itemID != "" {
		where += fmt.Sprintf(" AND s.item_id = $%d", argIdx)
		args = append(args, itemID)
		argIdx++
	}
	if locationID != "" {
		where += fmt.Sprintf(" AND s.location_id = $%d", argIdx)
		args = append(args, locationID)
		argIdx++
	}
	if search != "" {
		where += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.sku ILIKE $%d OR l.code ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM stock s LEFT JOIN items i ON s.item_id = i.id LEFT JOIN locations l ON s.location_id = l.id %s", where)
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`SELECT s.id, s.item_id, s.location_id, s.quantity, s.reserved_quantity, s.batch_number,
	        i.sku, i.name AS item_name, l.code AS location_code, l.name AS location_name
	 FROM stock s
	 LEFT JOIN items i ON s.item_id = i.id
	 LEFT JOIN locations l ON s.location_id = l.id
	 %s ORDER BY s.updated_at DESC LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, itemIDVal, locationIDVal string
		var quantity, reservedQuantity float64
		var batchNumber, itemSKU, itemName, locationCode, locationName string

		if err := rows.Scan(&id, &itemIDVal, &locationIDVal, &quantity, &reservedQuantity, &batchNumber, &itemSKU, &itemName, &locationCode, &locationName); err != nil {
			continue
		}
		results = append(results, map[string]interface{}{
			"id":                id,
			"item_id":           itemIDVal,
			"location_id":       locationIDVal,
			"quantity":          quantity,
			"reserved_quantity": reservedQuantity,
			"batch_number":      batchNumber,
			"item_sku":          itemSKU,
			"item_name":         itemName,
			"location_code":     locationCode,
			"location_name":     locationName,
		})
	}

	return results, total, nil
}

func (r *postgresStockRepository) GetTotalStockByItem(ctx context.Context, itemID string) (float64, error) {
	var total float64
	err := r.db.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(quantity), 0) FROM stock WHERE item_id = $1", itemID,
	).Scan(&total)
	return total, err
}
