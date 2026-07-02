package repository

import (
	"context"

	"wms/domain/entity"
)

type PurchaseOrderRepository interface {
	FindByID(ctx context.Context, id string) (*entity.PurchaseOrder, error)
	Create(ctx context.Context, po *entity.PurchaseOrder) error
	Update(ctx context.Context, po *entity.PurchaseOrder) error
	UpdateStatus(ctx context.Context, id, status string) error
	List(ctx context.Context, page, limit int) ([]*entity.PurchaseOrder, int, error)
	CreateItem(ctx context.Context, item *entity.PurchaseOrderItem) error
	FindItemsByPOID(ctx context.Context, poID string) ([]*entity.PurchaseOrderItem, error)
	UpdateReceivedQuantity(ctx context.Context, id string, receivedQty float64) error
}

type SalesOrderRepository interface {
	FindByID(ctx context.Context, id string) (*entity.SalesOrder, error)
	Create(ctx context.Context, so *entity.SalesOrder) error
	Update(ctx context.Context, so *entity.SalesOrder) error
	UpdateStatus(ctx context.Context, id, status string) error
	List(ctx context.Context, page, limit int) ([]*entity.SalesOrder, int, error)
	CreateItem(ctx context.Context, item *entity.SalesOrderItem) error
	UpdatePickedQuantity(ctx context.Context, id string, pickedQty float64) error
	CreatePickList(ctx context.Context, pl *entity.PickList) error
	CreateShipment(ctx context.Context, s *entity.Shipment) error
}
