package usecase

import (
	"context"
	"fmt"
	"time"

	"wms/application/dto"
	"wms/domain/entity"
	"wms/domain/repository"
	"wms/pkg/logger"
)

type OutboundUsecase interface {
	CreateSalesOrder(ctx context.Context, req *dto.CreateSalesOrderRequest, userID string) (*dto.SalesOrderResponse, error)
	GetSalesOrder(ctx context.Context, id string) (*dto.SalesOrderResponse, error)
	ListSalesOrders(ctx context.Context, page, limit int) ([]*dto.SalesOrderResponse, int, error)
	PickOrder(ctx context.Context, soID string, req *dto.PickListRequest, userID string) error
	ShipOrder(ctx context.Context, soID string, req *dto.ShipRequest, userID string) (*dto.ShipmentResponse, error)
}

type outboundUsecase struct {
	soRepo            repository.SalesOrderRepository
	stockRepo         repository.StockRepository
	stockMovementRepo repository.StockMovementRepository
	log               logger.Logger
}

func NewOutboundUsecase(
	soRepo repository.SalesOrderRepository,
	stockRepo repository.StockRepository,
	stockMovementRepo repository.StockMovementRepository,
	log logger.Logger,
) OutboundUsecase {
	return &outboundUsecase{
		soRepo:            soRepo,
		stockRepo:         stockRepo,
		stockMovementRepo: stockMovementRepo,
		log:               log,
	}
}

func (uc *outboundUsecase) CreateSalesOrder(ctx context.Context, req *dto.CreateSalesOrderRequest, userID string) (*dto.SalesOrderResponse, error) {
	uc.log.Info("creating sales order", "customer", req.CustomerName)

	so := &entity.SalesOrder{
		SONumber:     fmt.Sprintf("SO-%s", time.Now().Format("20060102150405")),
		CustomerName: req.CustomerName,
		Status:       "pending",
		Notes:        req.Notes,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
	}

	if err := uc.soRepo.Create(ctx, so); err != nil {
		uc.log.Error("failed to create sales order", "error", err)
		return nil, err
	}

	for _, item := range req.Items {
		soItem := &entity.SalesOrderItem{
			SOID:             so.ID,
			ItemID:           item.ItemID,
			OrderedQuantity: item.OrderedQuantity,
		}
		if err := uc.soRepo.CreateItem(ctx, soItem); err != nil {
			uc.log.Error("failed to create SO item", "error", err)
		}
	}

	uc.log.Info("sales order created", "id", so.ID, "so_number", so.SONumber)

	return &dto.SalesOrderResponse{
		ID:           so.ID,
		SONumber:     so.SONumber,
		CustomerName: so.CustomerName,
		Status:       so.Status,
		Notes:        so.Notes,
		CreatedAt:    so.CreatedAt,
	}, nil
}

func (uc *outboundUsecase) GetSalesOrder(ctx context.Context, id string) (*dto.SalesOrderResponse, error) {
	so, err := uc.soRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("sales order not found")
	}

	return &dto.SalesOrderResponse{
		ID:           so.ID,
		SONumber:     so.SONumber,
		CustomerName: so.CustomerName,
		Status:       so.Status,
		Notes:        so.Notes,
		CreatedAt:    so.CreatedAt,
	}, nil
}

func (uc *outboundUsecase) ListSalesOrders(ctx context.Context, page, limit int) ([]*dto.SalesOrderResponse, int, error) {
	sos, total, err := uc.soRepo.List(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var result []*dto.SalesOrderResponse
	for _, so := range sos {
		result = append(result, &dto.SalesOrderResponse{
			ID:           so.ID,
			SONumber:     so.SONumber,
			CustomerName: so.CustomerName,
			Status:       so.Status,
			Notes:        so.Notes,
			CreatedAt:    so.CreatedAt,
		})
	}

	return result, total, nil
}

func (uc *outboundUsecase) PickOrder(ctx context.Context, soID string, req *dto.PickListRequest, userID string) error {
	uc.log.Info("picking order", "so_id", soID)

	_, err := uc.soRepo.FindByID(ctx, soID)
	if err != nil {
		return fmt.Errorf("sales order not found")
	}

	for _, item := range req.Items {
		stock, err := uc.stockRepo.FindByItem(ctx, item.ItemID)
		if err != nil || len(stock) == 0 {
			return fmt.Errorf("insufficient stock for item %s", item.ItemID)
		}

		totalAvailable := 0.0
		for _, s := range stock {
			totalAvailable += s.Quantity - s.ReservedQuantity
		}

		if totalAvailable < item.Quantity {
			return fmt.Errorf("insufficient stock for item %s", item.ItemID)
		}

		remaining := item.Quantity
		for _, s := range stock {
			if remaining <= 0 {
				break
			}
			available := s.Quantity - s.ReservedQuantity
			if available <= 0 {
				continue
			}
			deduct := remaining
			if deduct > available {
				deduct = available
			}
			uc.stockRepo.UpdateQuantity(ctx, s.ID, s.Quantity-deduct)
			remaining -= deduct

			movement := &entity.StockMovement{
				ItemID:         item.ItemID,
				FromLocationID: &s.LocationID,
				Quantity:       deduct,
				MovementType:   "pick",
				ReferenceType:  "sales_order",
				ReferenceID:    &soID,
				CreatedBy:      userID,
				CreatedAt:      time.Now(),
			}
			uc.stockMovementRepo.Create(ctx, movement)
		}
	}

	uc.soRepo.UpdateStatus(ctx, soID, "picked")
	uc.log.Info("order picked successfully", "so_id", soID)
	return nil
}

func (uc *outboundUsecase) ShipOrder(ctx context.Context, soID string, req *dto.ShipRequest, userID string) (*dto.ShipmentResponse, error) {
	uc.log.Info("shipping order", "so_id", soID)

	shipment := &entity.Shipment{
		ShipmentNumber: fmt.Sprintf("SHIP-%s", time.Now().Format("20060102150405")),
		SOID:           soID,
		Carrier:        req.Carrier,
		TrackingNumber: req.TrackingNumber,
		Status:         "shipped",
		CreatedBy:      userID,
	}

	now := time.Now()
	shipment.ShippedAt = &now

	if err := uc.soRepo.CreateShipment(ctx, shipment); err != nil {
		uc.log.Error("failed to create shipment", "error", err)
		return nil, err
	}

	uc.soRepo.UpdateStatus(ctx, soID, "shipped")

	uc.log.Info("order shipped successfully", "so_id", soID, "shipment_number", shipment.ShipmentNumber)

	return &dto.ShipmentResponse{
		ID:             shipment.ID,
		ShipmentNumber: shipment.ShipmentNumber,
		SOID:           shipment.SOID,
		Carrier:        shipment.Carrier,
		TrackingNumber: shipment.TrackingNumber,
		Status:         shipment.Status,
		ShippedAt:      shipment.ShippedAt,
	}, nil
}
