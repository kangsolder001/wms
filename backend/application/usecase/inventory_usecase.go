package usecase

import (
	"context"
	"time"

	"wms/application/dto"
	"wms/domain/entity"
	"wms/domain/repository"
	"wms/pkg/logger"
)

type InventoryUsecase interface {
	GetStock(ctx context.Context, page, limit int, itemID, locationID, search string) ([]*dto.StockResponse, int, error)
	GetStockByItem(ctx context.Context, itemID string) ([]*dto.StockResponse, error)
	AdjustStock(ctx context.Context, req *dto.AdjustStockRequest, userID string) error
	GetStockMovements(ctx context.Context, page, limit int) ([]*dto.StockMovementResponse, int, error)
}

type inventoryUsecase struct {
	stockRepo         repository.StockRepository
	stockMovementRepo repository.StockMovementRepository
	log               logger.Logger
}

func NewInventoryUsecase(
	stockRepo repository.StockRepository,
	stockMovementRepo repository.StockMovementRepository,
	log logger.Logger,
) InventoryUsecase {
	return &inventoryUsecase{
		stockRepo:         stockRepo,
		stockMovementRepo: stockMovementRepo,
		log:               log,
	}
}

func (uc *inventoryUsecase) GetStock(ctx context.Context, page, limit int, itemID, locationID, search string) ([]*dto.StockResponse, int, error) {
	results, total, err := uc.stockRepo.ListWithDetails(ctx, page, limit, itemID, locationID, search)
	if err != nil {
		uc.log.Error("failed to list stock", "error", err)
		return nil, 0, err
	}

	var stockList []*dto.StockResponse
	for _, r := range results {
		stockList = append(stockList, &dto.StockResponse{
			ID:               r["id"].(string),
			ItemID:           r["item_id"].(string),
			LocationID:       r["location_id"].(string),
			Quantity:         r["quantity"].(float64),
			ReservedQuantity: r["reserved_quantity"].(float64),
			BatchNumber:      r["batch_number"].(string),
			ItemSKU:          r["item_sku"].(string),
			ItemName:         r["item_name"].(string),
			LocationCode:     r["location_code"].(string),
		})
	}

	return stockList, total, nil
}

func (uc *inventoryUsecase) GetStockByItem(ctx context.Context, itemID string) ([]*dto.StockResponse, error) {
	stocks, err := uc.stockRepo.FindByItem(ctx, itemID)
	if err != nil {
		uc.log.Error("failed to find stock by item", "item_id", itemID, "error", err)
		return nil, err
	}

	var result []*dto.StockResponse
	for _, s := range stocks {
		result = append(result, &dto.StockResponse{
			ID:               s.ID,
			ItemID:           s.ItemID,
			LocationID:       s.LocationID,
			Quantity:         s.Quantity,
			ReservedQuantity: s.ReservedQuantity,
			BatchNumber:      s.BatchNumber,
		})
	}

	return result, nil
}

func (uc *inventoryUsecase) AdjustStock(ctx context.Context, req *dto.AdjustStockRequest, userID string) error {
	uc.log.Info("adjusting stock", "item_id", req.ItemID, "location_id", req.LocationID, "quantity", req.Quantity)

	existing, _ := uc.stockRepo.FindByItemAndLocation(ctx, req.ItemID, req.LocationID)
	if existing == nil {
		stock := &entity.Stock{
			ItemID:     req.ItemID,
			LocationID: req.LocationID,
			Quantity:   req.Quantity,
			UpdatedAt:  time.Now(),
		}
		if err := uc.stockRepo.Create(ctx, stock); err != nil {
			uc.log.Error("failed to create stock", "error", err)
			return err
		}
	} else {
		newQty := existing.Quantity + req.Quantity
		if newQty < 0 {
			newQty = 0
		}
		if err := uc.stockRepo.UpdateQuantity(ctx, existing.ID, newQty); err != nil {
			uc.log.Error("failed to update stock", "error", err)
			return err
		}
	}

	movement := &entity.StockMovement{
		ItemID:        req.ItemID,
		ToLocationID:  &req.LocationID,
		Quantity:      req.Quantity,
		MovementType:  "adjustment",
		ReferenceType: "manual",
		Notes:         req.Notes,
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
	}

	if err := uc.stockMovementRepo.Create(ctx, movement); err != nil {
		uc.log.Error("failed to create stock movement", "error", err)
		return err
	}

	uc.log.Info("stock adjusted successfully", "item_id", req.ItemID, "location_id", req.LocationID)
	return nil
}

func (uc *inventoryUsecase) GetStockMovements(ctx context.Context, page, limit int) ([]*dto.StockMovementResponse, int, error) {
	movements, total, err := uc.stockMovementRepo.List(ctx, page, limit)
	if err != nil {
		uc.log.Error("failed to list stock movements", "error", err)
		return nil, 0, err
	}

	var result []*dto.StockMovementResponse
	for _, m := range movements {
		result = append(result, &dto.StockMovementResponse{
			ID:             m.ID,
			ItemID:         m.ItemID,
			FromLocationID: m.FromLocationID,
			ToLocationID:   m.ToLocationID,
			Quantity:       m.Quantity,
			MovementType:   m.MovementType,
			ReferenceType:  m.ReferenceType,
			ReferenceID:    m.ReferenceID,
			Notes:          m.Notes,
			CreatedBy:      m.CreatedBy,
		})
	}

	return result, total, nil
}
