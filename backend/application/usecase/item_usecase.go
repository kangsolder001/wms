package usecase

import (
	"context"
	"errors"
	"time"

	"wms/application/dto"
	"wms/domain/entity"
	"wms/domain/repository"
	"wms/pkg/logger"
)

type ItemUsecase interface {
	CreateItem(ctx context.Context, req *dto.CreateItemRequest) (*dto.ItemResponse, error)
	GetItem(ctx context.Context, id string) (*dto.ItemResponse, error)
	UpdateItem(ctx context.Context, id string, req *dto.UpdateItemRequest) (*dto.ItemResponse, error)
	DeleteItem(ctx context.Context, id string) error
	ListItems(ctx context.Context, page, limit int) ([]*dto.ItemResponse, int, error)
}

type itemUsecase struct {
	itemRepo repository.ItemRepository
	log      logger.Logger
}

func NewItemUsecase(itemRepo repository.ItemRepository, log logger.Logger) ItemUsecase {
	return &itemUsecase{itemRepo: itemRepo, log: log}
}

func (uc *itemUsecase) CreateItem(ctx context.Context, req *dto.CreateItemRequest) (*dto.ItemResponse, error) {
	uc.log.Info("creating item", "sku", req.SKU, "name", req.Name)

	existing, _ := uc.itemRepo.FindBySKU(ctx, req.SKU)
	if existing != nil {
		return nil, errors.New("SKU already exists")
	}

	item := &entity.Item{
		SKU:           req.SKU,
		Name:          req.Name,
		Description:   req.Description,
		Category:      req.Category,
		UnitOfMeasure: req.UnitOfMeasure,
		Weight:        req.Weight,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.itemRepo.Create(ctx, item); err != nil {
		uc.log.Error("failed to create item", "sku", req.SKU, "error", err)
		return nil, err
	}

	uc.log.Info("item created successfully", "id", item.ID, "sku", item.SKU)

	return &dto.ItemResponse{
		ID:            item.ID,
		SKU:           item.SKU,
		Name:          item.Name,
		Description:   item.Description,
		Category:      item.Category,
		UnitOfMeasure: item.UnitOfMeasure,
		Weight:        item.Weight,
		IsActive:      item.IsActive,
	}, nil
}

func (uc *itemUsecase) GetItem(ctx context.Context, id string) (*dto.ItemResponse, error) {
	item, err := uc.itemRepo.FindByID(ctx, id)
	if err != nil {
		uc.log.Error("item not found", "id", id, "error", err)
		return nil, errors.New("item not found")
	}

	return &dto.ItemResponse{
		ID:            item.ID,
		SKU:           item.SKU,
		Name:          item.Name,
		Description:   item.Description,
		Category:      item.Category,
		UnitOfMeasure: item.UnitOfMeasure,
		Weight:        item.Weight,
		IsActive:      item.IsActive,
	}, nil
}

func (uc *itemUsecase) UpdateItem(ctx context.Context, id string, req *dto.UpdateItemRequest) (*dto.ItemResponse, error) {
	uc.log.Info("updating item", "id", id)

	item, err := uc.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("item not found")
	}

	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.Category != nil {
		item.Category = *req.Category
	}
	if req.UnitOfMeasure != nil {
		item.UnitOfMeasure = *req.UnitOfMeasure
	}
	if req.Weight != nil {
		item.Weight = *req.Weight
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	item.UpdatedAt = time.Now()

	if err := uc.itemRepo.Update(ctx, item); err != nil {
		uc.log.Error("failed to update item", "id", id, "error", err)
		return nil, err
	}

	uc.log.Info("item updated successfully", "id", item.ID)

	return &dto.ItemResponse{
		ID:            item.ID,
		SKU:           item.SKU,
		Name:          item.Name,
		Description:   item.Description,
		Category:      item.Category,
		UnitOfMeasure: item.UnitOfMeasure,
		Weight:        item.Weight,
		IsActive:      item.IsActive,
	}, nil
}

func (uc *itemUsecase) DeleteItem(ctx context.Context, id string) error {
	uc.log.Info("deleting item", "id", id)

	if err := uc.itemRepo.Delete(ctx, id); err != nil {
		uc.log.Error("failed to delete item", "id", id, "error", err)
		return err
	}

	uc.log.Info("item deleted successfully", "id", id)
	return nil
}

func (uc *itemUsecase) ListItems(ctx context.Context, page, limit int) ([]*dto.ItemResponse, int, error) {
	items, total, err := uc.itemRepo.List(ctx, page, limit)
	if err != nil {
		uc.log.Error("failed to list items", "error", err)
		return nil, 0, err
	}

	var result []*dto.ItemResponse
	for _, item := range items {
		result = append(result, &dto.ItemResponse{
			ID:            item.ID,
			SKU:           item.SKU,
			Name:          item.Name,
			Description:   item.Description,
			Category:      item.Category,
			UnitOfMeasure: item.UnitOfMeasure,
			Weight:        item.Weight,
			IsActive:      item.IsActive,
		})
	}

	return result, total, nil
}
