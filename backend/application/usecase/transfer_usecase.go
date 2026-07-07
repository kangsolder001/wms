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

type TransferUsecase interface {
	CreateTransfer(ctx context.Context, req *dto.CreateTransferRequest, userID string) (*dto.TransferResponse, error)
	CompleteTransfer(ctx context.Context, id string) error
	ListTransfers(ctx context.Context, page, limit int) ([]*dto.TransferResponse, int, error)
}

type transferUsecase struct {
	transferRepo     repository.StockTransferRepository
	stockRepo        repository.StockRepository
	stockMovementRepo repository.StockMovementRepository
	log              logger.Logger
}

func NewTransferUsecase(
	transferRepo repository.StockTransferRepository,
	stockRepo repository.StockRepository,
	stockMovementRepo repository.StockMovementRepository,
	log logger.Logger,
) TransferUsecase {
	return &transferUsecase{
		transferRepo:     transferRepo,
		stockRepo:        stockRepo,
		stockMovementRepo: stockMovementRepo,
		log:              log,
	}
}

func (uc *transferUsecase) CreateTransfer(ctx context.Context, req *dto.CreateTransferRequest, userID string) (*dto.TransferResponse, error) {
	uc.log.Info("creating stock transfer", "item_id", req.ItemID, "from", req.FromLocationID, "to", req.ToLocationID)

	stock, err := uc.stockRepo.FindByItemAndLocation(ctx, req.ItemID, req.FromLocationID)
	if err != nil || stock == nil {
		return nil, fmt.Errorf("no stock found at source location")
	}

	if stock.Quantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock at source location")
	}

	transfer := &entity.StockTransfer{
		TransferNumber: fmt.Sprintf("TRF-%s", time.Now().Format("20060102150405")),
		FromLocationID: req.FromLocationID,
		ToLocationID:   req.ToLocationID,
		ItemID:         req.ItemID,
		Quantity:       req.Quantity,
		Status:         "pending",
		Notes:          req.Notes,
		CreatedBy:      userID,
		CreatedAt:      time.Now(),
	}

	if err := uc.transferRepo.Create(ctx, transfer); err != nil {
		uc.log.Error("failed to create transfer", "error", err)
		return nil, err
	}

	uc.log.Info("transfer created", "id", transfer.ID, "transfer_number", transfer.TransferNumber)

	return &dto.TransferResponse{
		ID:             transfer.ID,
		TransferNumber: transfer.TransferNumber,
		FromLocationID: transfer.FromLocationID,
		ToLocationID:   transfer.ToLocationID,
		ItemID:         transfer.ItemID,
		Quantity:       transfer.Quantity,
		Status:         transfer.Status,
		Notes:          transfer.Notes,
	}, nil
}

func (uc *transferUsecase) CompleteTransfer(ctx context.Context, id string) error {
	uc.log.Info("completing transfer", "id", id)

	transfer, err := uc.transferRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("transfer not found")
	}

	if transfer.Status != "pending" {
		return fmt.Errorf("transfer cannot be completed, current status: %s", transfer.Status)
	}

	fromStock, err := uc.stockRepo.FindByItemAndLocation(ctx, transfer.ItemID, transfer.FromLocationID)
	if err != nil || fromStock == nil || fromStock.Quantity < transfer.Quantity {
		return fmt.Errorf("insufficient stock at source location")
	}

	uc.stockRepo.UpdateQuantity(ctx, fromStock.ID, fromStock.Quantity-transfer.Quantity)

	toStock, _ := uc.stockRepo.FindByItemAndLocation(ctx, transfer.ItemID, transfer.ToLocationID)
	if toStock == nil {
		newStock := &entity.Stock{
			ItemID:     transfer.ItemID,
			LocationID: transfer.ToLocationID,
			Quantity:   transfer.Quantity,
			UpdatedAt:  time.Now(),
		}
		uc.stockRepo.Create(ctx, newStock)
	} else {
		uc.stockRepo.UpdateQuantity(ctx, toStock.ID, toStock.Quantity+transfer.Quantity)
	}

	fromMovement := &entity.StockMovement{
		ItemID:         transfer.ItemID,
		FromLocationID: &transfer.FromLocationID,
		Quantity:       transfer.Quantity,
		MovementType:   "transfer_out",
		ReferenceType:  "stock_transfer",
		ReferenceID:    &transfer.ID,
		CreatedBy:      transfer.CreatedBy,
		CreatedAt:      time.Now(),
	}
	uc.stockMovementRepo.Create(ctx, fromMovement)

	toMovement := &entity.StockMovement{
		ItemID:        transfer.ItemID,
		ToLocationID:  &transfer.ToLocationID,
		Quantity:      transfer.Quantity,
		MovementType:  "transfer_in",
		ReferenceType: "stock_transfer",
		ReferenceID:   &transfer.ID,
		CreatedBy:     transfer.CreatedBy,
		CreatedAt:     time.Now(),
	}
	uc.stockMovementRepo.Create(ctx, toMovement)

	uc.transferRepo.UpdateStatus(ctx, id, "completed")
	uc.log.Info("transfer completed", "id", id)
	return nil
}

func (uc *transferUsecase) ListTransfers(ctx context.Context, page, limit int) ([]*dto.TransferResponse, int, error) {
	transfers, total, err := uc.transferRepo.List(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var result []*dto.TransferResponse
	for _, t := range transfers {
		result = append(result, &dto.TransferResponse{
			ID:             t.ID,
			TransferNumber: t.TransferNumber,
			FromLocationID: t.FromLocationID,
			ToLocationID:   t.ToLocationID,
			ItemID:         t.ItemID,
			Quantity:       t.Quantity,
			Status:         t.Status,
			Notes:          t.Notes,
		})
	}

	return result, total, nil
}
