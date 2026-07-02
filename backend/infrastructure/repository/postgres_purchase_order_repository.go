package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wms/domain/entity"
	"wms/pkg/logger"
)

type postgresPurchaseOrderRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewPostgresPurchaseOrderRepository(db *sql.DB, log logger.Logger) *postgresPurchaseOrderRepository {
	return &postgresPurchaseOrderRepository{db: db, log: log}
}

func (r *postgresPurchaseOrderRepository) FindByID(ctx context.Context, id string) (*entity.PurchaseOrder, error) {
	po := &entity.PurchaseOrder{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, po_number, supplier_name, status, expected_date, storage_location_id, notes, created_by, u.full_name, created_at, updated_at 
		 FROM purchase_orders po
		 LEFT JOIN users u ON po.created_by = u.id
		 WHERE po.id = $1`, id,
	).Scan(&po.ID, &po.PONumber, &po.SupplierName, &po.Status, &po.ExpectedDate, &po.StorageLocationID, &po.Notes, &po.CreatedBy, &po.CreatedByName, &po.CreatedAt, &po.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}
	return po, nil
}

func (r *postgresPurchaseOrderRepository) Create(ctx context.Context, po *entity.PurchaseOrder) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO purchase_orders (po_number, supplier_name, status, expected_date, storage_location_id, notes, created_by, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		po.PONumber, po.SupplierName, po.Status, po.ExpectedDate, po.StorageLocationID, po.Notes, po.CreatedBy, po.CreatedAt, po.UpdatedAt,
	).Scan(&po.ID)
	if err != nil {
		return fmt.Errorf("failed to create purchase order: %w", err)
	}
	return nil
}

func (r *postgresPurchaseOrderRepository) Update(ctx context.Context, po *entity.PurchaseOrder) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE purchase_orders SET supplier_name=$1, status=$2, expected_date=$3, notes=$4, updated_at=$5 WHERE id=$6`,
		po.SupplierName, po.Status, po.ExpectedDate, po.Notes, po.UpdatedAt, po.ID,
	)
	return err
}

func (r *postgresPurchaseOrderRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE purchase_orders SET status = $1, updated_at = NOW() WHERE id = $2", status, id,
	)
	return err
}

func (r *postgresPurchaseOrderRepository) List(ctx context.Context, page, limit int) ([]*entity.PurchaseOrder, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM purchase_orders").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT po.id, po.po_number, po.supplier_name, po.status, po.expected_date, po.storage_location_id, po.notes, 
		        po.created_by, u.full_name, po.created_at, po.updated_at 
		 FROM purchase_orders po
		 LEFT JOIN users u ON po.created_by = u.id
		 ORDER BY po.created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var pos []*entity.PurchaseOrder
	for rows.Next() {
		po := &entity.PurchaseOrder{}
		var storageLocationID, createdByName sql.NullString
		if err := rows.Scan(&po.ID, &po.PONumber, &po.SupplierName, &po.Status, &po.ExpectedDate, &storageLocationID, &po.Notes, &po.CreatedBy, &createdByName, &po.CreatedAt, &po.UpdatedAt); err != nil {
			continue
		}
		if storageLocationID.Valid {
			po.StorageLocationID = storageLocationID.String
		}
		if createdByName.Valid {
			po.CreatedByName = createdByName.String
		}
		pos = append(pos, po)
	}

	return pos, total, nil
}

func (r *postgresPurchaseOrderRepository) CreateItem(ctx context.Context, item *entity.PurchaseOrderItem) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO purchase_order_items (po_id, item_id, expected_quantity, received_quantity, unit_price) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		item.POID, item.ItemID, item.ExpectedQuantity, item.ReceivedQuantity, item.UnitPrice,
	).Scan(&item.ID)
	return err
}

func (r *postgresPurchaseOrderRepository) UpdateReceivedQuantity(ctx context.Context, id string, receivedQty float64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE purchase_order_items SET received_quantity = $1 WHERE id = $2", receivedQty, id,
	)
	return err
}

func (r *postgresPurchaseOrderRepository) FindItemsByPOID(ctx context.Context, poID string) ([]*entity.PurchaseOrderItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, po_id, item_id, expected_quantity, received_quantity, unit_price 
		 FROM purchase_order_items WHERE po_id = $1`, poID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entity.PurchaseOrderItem
	for rows.Next() {
		item := &entity.PurchaseOrderItem{}
		if err := rows.Scan(&item.ID, &item.POID, &item.ItemID, &item.ExpectedQuantity, &item.ReceivedQuantity, &item.UnitPrice); err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, nil
}
