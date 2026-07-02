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

type InboundUsecase interface {
	CreatePurchaseOrder(ctx context.Context, req *dto.CreatePurchaseOrderRequest, userID string) (*dto.PurchaseOrderResponse, error)
	GetPurchaseOrder(ctx context.Context, id string) (*dto.PurchaseOrderResponse, error)
	ListPurchaseOrders(ctx context.Context, page, limit int) ([]*dto.PurchaseOrderResponse, int, error)
	ReceiveGoods(ctx context.Context, poID string, req *dto.ReceiveGoodsRequest, userID string) (*dto.GoodsReceiptResponse, error)
}

type inboundUsecase struct {
	poRepo           repository.PurchaseOrderRepository
	stockRepo        repository.StockRepository
	stockMovementRepo repository.StockMovementRepository
	log              logger.Logger
}

func NewInboundUsecase(
	poRepo repository.PurchaseOrderRepository,
	stockRepo repository.StockRepository,
	stockMovementRepo repository.StockMovementRepository,
	log logger.Logger,
) InboundUsecase {
	return &inboundUsecase{
		poRepo:            poRepo,
		stockRepo:         stockRepo,
		stockMovementRepo: stockMovementRepo,
		log:               log,
	}
}

func (uc *inboundUsecase) CreatePurchaseOrder(ctx context.Context, req *dto.CreatePurchaseOrderRequest, userID string) (*dto.PurchaseOrderResponse, error) {
	uc.log.Info("creating purchase order", "supplier", req.SupplierName)

	var expectedDate *time.Time
	if req.ExpectedDate != "" {
		t, err := time.Parse("2006-01-02", req.ExpectedDate)
		if err != nil {
			t, err = time.Parse(time.RFC3339, req.ExpectedDate)
		}
		if err == nil {
			expectedDate = &t
		}
	}

	po := &entity.PurchaseOrder{
		PONumber:     fmt.Sprintf("PO-%s", time.Now().Format("20060102150405")),
		SupplierName: req.SupplierName,
		Status:       "pending",
		ExpectedDate: expectedDate,
		Notes:        req.Notes,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.poRepo.Create(ctx, po); err != nil {
		uc.log.Error("failed to create purchase order", "error", err)
		return nil, err
	}

	for _, item := range req.Items {
		poItem := &entity.PurchaseOrderItem{
			POID:             po.ID,
			ItemID:           item.ItemID,
			ExpectedQuantity: item.ExpectedQuantity,
			UnitPrice:        item.UnitPrice,
		}
		if err := uc.poRepo.CreateItem(ctx, poItem); err != nil {
			uc.log.Error("failed to create PO item", "error", err)
		}
	}

	uc.log.Info("purchase order created", "id", po.ID, "po_number", po.PONumber)

	return &dto.PurchaseOrderResponse{
		ID:           po.ID,
		PONumber:     po.PONumber,
		SupplierName: po.SupplierName,
		Status:       po.Status,
		ExpectedDate: po.ExpectedDate,
		Notes:        po.Notes,
		CreatedBy:    po.CreatedBy,
		CreatedByName: po.CreatedByName,
		CreatedAt:    po.CreatedAt,
	}, nil
}

func (uc *inboundUsecase) GetPurchaseOrder(ctx context.Context, id string) (*dto.PurchaseOrderResponse, error) {
	po, err := uc.poRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found")
	}

	return &dto.PurchaseOrderResponse{
		ID:           po.ID,
		PONumber:     po.PONumber,
		SupplierName: po.SupplierName,
		Status:       po.Status,
		ExpectedDate: po.ExpectedDate,
		Notes:        po.Notes,
		CreatedBy:    po.CreatedBy,
		CreatedByName: po.CreatedByName,
		CreatedAt:    po.CreatedAt,
	}, nil
}

func (uc *inboundUsecase) ListPurchaseOrders(ctx context.Context, page, limit int) ([]*dto.PurchaseOrderResponse, int, error) {
	pos, total, err := uc.poRepo.List(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var result []*dto.PurchaseOrderResponse
	for _, po := range pos {
		result = append(result, &dto.PurchaseOrderResponse{
			ID:           po.ID,
			PONumber:     po.PONumber,
			SupplierName: po.SupplierName,
			Status:       po.Status,
			ExpectedDate: po.ExpectedDate,
			Notes:        po.Notes,
			CreatedBy:    po.CreatedBy,
			CreatedByName: po.CreatedByName,
			CreatedAt:    po.CreatedAt,
		})
	}

	return result, total, nil
}

func (uc *inboundUsecase) ReceiveGoods(ctx context.Context, poID string, req *dto.ReceiveGoodsRequest, userID string) (*dto.GoodsReceiptResponse, error) {
	uc.log.Info("receiving goods", "po_id", poID, "grn_number", req.GRNNumber)

	_, err := uc.poRepo.FindByID(ctx, poID)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found")
	}

	grn := &entity.GoodsReceipt{
		GRNNumber:  req.GRNNumber,
		POID:       poID,
		ReceivedBy: userID,
		ReceivedAt: time.Now(),
		Notes:      req.Notes,
	}

	for _, item := range req.Items {
		grnItem := &entity.GoodsReceiptItem{
			ItemID:      item.ItemID,
			Quantity:    item.Quantity,
			BatchNumber: item.BatchNumber,
			LocationID:  item.LocationID,
		}
		grn.Items = append(grn.Items, *grnItem)

		stock, _ := uc.stockRepo.FindByItemAndLocation(ctx, item.ItemID, item.LocationID)
		if stock == nil {
			newStock := &entity.Stock{
				ItemID:     item.ItemID,
				LocationID: item.LocationID,
				Quantity:   item.Quantity,
				UpdatedAt:  time.Now(),
			}
			uc.stockRepo.Create(ctx, newStock)
		} else {
			uc.stockRepo.UpdateQuantity(ctx, stock.ID, stock.Quantity+item.Quantity)
		}

		movement := &entity.StockMovement{
			ItemID:        item.ItemID,
			ToLocationID:  &item.LocationID,
			Quantity:      item.Quantity,
			MovementType:  "receipt",
			ReferenceType: "purchase_order",
			ReferenceID:   poID,
			CreatedBy:     userID,
			CreatedAt:     time.Now(),
		}
		uc.stockMovementRepo.Create(ctx, movement)
	}

	uc.poRepo.UpdateStatus(ctx, poID, "received")

	uc.log.Info("goods received successfully", "grn_number", req.GRNNumber)

	return &dto.GoodsReceiptResponse{
		ID:         grn.ID,
		GRNNumber:  grn.GRNNumber,
		POID:       grn.POID,
		ReceivedBy: grn.ReceivedBy,
		ReceivedAt: grn.ReceivedAt,
		Notes:      grn.Notes,
	}, nil
}
