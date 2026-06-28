package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wms/domain/entity"
	"wms/pkg/logger"
)

type postgresSalesOrderRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewPostgresSalesOrderRepository(db *sql.DB, log logger.Logger) *postgresSalesOrderRepository {
	return &postgresSalesOrderRepository{db: db, log: log}
}

func (r *postgresSalesOrderRepository) FindByID(ctx context.Context, id string) (*entity.SalesOrder, error) {
	so := &entity.SalesOrder{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, so_number, customer_name, status, notes, created_by, created_at 
		 FROM sales_orders WHERE id = $1`, id,
	).Scan(&so.ID, &so.SONumber, &so.CustomerName, &so.Status, &so.Notes, &so.CreatedBy, &so.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("sales order not found: %w", err)
	}
	return so, nil
}

func (r *postgresSalesOrderRepository) Create(ctx context.Context, so *entity.SalesOrder) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO sales_orders (so_number, customer_name, status, notes, created_by, created_at) 
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		so.SONumber, so.CustomerName, so.Status, so.Notes, so.CreatedBy, so.CreatedAt,
	).Scan(&so.ID)
	if err != nil {
		return fmt.Errorf("failed to create sales order: %w", err)
	}
	return nil
}

func (r *postgresSalesOrderRepository) Update(ctx context.Context, so *entity.SalesOrder) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sales_orders SET customer_name=$1, status=$2, notes=$3 WHERE id=$4`,
		so.CustomerName, so.Status, so.Notes, so.ID,
	)
	return err
}

func (r *postgresSalesOrderRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sales_orders SET status = $1 WHERE id = $2", status, id,
	)
	return err
}

func (r *postgresSalesOrderRepository) List(ctx context.Context, page, limit int) ([]*entity.SalesOrder, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sales_orders").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, so_number, customer_name, status, notes, created_by, created_at 
		 FROM sales_orders ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sos []*entity.SalesOrder
	for rows.Next() {
		so := &entity.SalesOrder{}
		if err := rows.Scan(&so.ID, &so.SONumber, &so.CustomerName, &so.Status, &so.Notes, &so.CreatedBy, &so.CreatedAt); err != nil {
			continue
		}
		sos = append(sos, so)
	}

	return sos, total, nil
}

func (r *postgresSalesOrderRepository) CreateItem(ctx context.Context, item *entity.SalesOrderItem) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO sales_order_items (so_id, item_id, ordered_quantity, picked_quantity) 
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		item.SOID, item.ItemID, item.OrderedQuantity, item.PickedQuantity,
	).Scan(&item.ID)
	return err
}

func (r *postgresSalesOrderRepository) UpdatePickedQuantity(ctx context.Context, id string, pickedQty float64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sales_order_items SET picked_quantity = $1 WHERE id = $2", pickedQty, id,
	)
	return err
}

func (r *postgresSalesOrderRepository) CreatePickList(ctx context.Context, pl *entity.PickList) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO pick_lists (so_id, status, picked_by, picked_at, created_at) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		pl.SOID, pl.Status, pl.PickedBy, pl.PickedAt, pl.CreatedAt,
	).Scan(&pl.ID)
	return err
}

func (r *postgresSalesOrderRepository) CreateShipment(ctx context.Context, s *entity.Shipment) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO shipments (shipment_number, so_id, carrier, tracking_number, status, shipped_at, created_by) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		s.ShipmentNumber, s.SOID, s.Carrier, s.TrackingNumber, s.Status, s.ShippedAt, s.CreatedBy,
	).Scan(&s.ID)
	return err
}
