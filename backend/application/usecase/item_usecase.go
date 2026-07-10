package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"wms/application/dto"
	"wms/domain/entity"
	"wms/domain/repository"
	"wms/pkg/logger"
)

type ItemUsecase interface {
	GenerateSKU(ctx context.Context, req *dto.GenerateSKURequest) (*dto.GenerateSKUResponse, error)
	CreateItem(ctx context.Context, req *dto.CreateItemRequest) (*dto.ItemResponse, error)
	GetItem(ctx context.Context, id string) (*dto.ItemResponse, error)
	UpdateItem(ctx context.Context, id string, req *dto.UpdateItemRequest) (*dto.ItemResponse, error)
	DeleteItem(ctx context.Context, id string) error
	ListItems(ctx context.Context, page, limit int) ([]*dto.ItemResponse, int, error)
}

type itemUsecase struct {
	itemRepo     repository.ItemRepository
	categoryRepo repository.CategoryRepository
	log          logger.Logger
}

func NewItemUsecase(itemRepo repository.ItemRepository, categoryRepo repository.CategoryRepository, log logger.Logger) ItemUsecase {
	return &itemUsecase{itemRepo: itemRepo, categoryRepo: categoryRepo, log: log}
}

func (uc *itemUsecase) GenerateSKU(ctx context.Context, req *dto.GenerateSKURequest) (*dto.GenerateSKUResponse, error) {
	category, err := uc.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	seq, err := uc.itemRepo.GetNextSKUSequence(ctx, category.Abbreviation)
	if err != nil {
		return nil, err
	}

	sku := fmt.Sprintf("%s-%04d", category.Abbreviation, seq)
	return &dto.GenerateSKUResponse{SKU: sku}, nil
}

func (uc *itemUsecase) CreateItem(ctx context.Context, req *dto.CreateItemRequest) (*dto.ItemResponse, error) {
	category, err := uc.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	seq, err := uc.itemRepo.GetNextSKUSequence(ctx, category.Abbreviation)
	if err != nil {
		return nil, err
	}

	sku := fmt.Sprintf("%s-%04d", category.Abbreviation, seq)

	uc.log.Info("creating item", "sku", sku, "name", req.Name)

	item := &entity.Item{
		SKU:           sku,
		Name:          req.Name,
		Description:   req.Description,
		CategoryID:    req.CategoryID,
		Category:      category.Name,
		Barcode:       entity.GenerateBarcode(sku, req.Name, category.Name),
		UnitOfMeasure: req.UnitOfMeasure,
		Weight:        req.Weight,
		Length:        req.Length,
		Width:         req.Width,
		Height:        req.Height,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.itemRepo.Create(ctx, item); err != nil {
		uc.log.Error("failed to create item", "sku", sku, "error", err)
		return nil, err
	}

	uc.log.Info("item created successfully", "id", item.ID, "sku", item.SKU)

	return &dto.ItemResponse{
		ID:            item.ID,
		SKU:           item.SKU,
		Name:          item.Name,
		Description:   item.Description,
		Category:      item.Category,
		CategoryID:    item.CategoryID,
		Barcode:       item.Barcode,
		UnitOfMeasure: item.UnitOfMeasure,
		Weight:        item.Weight,
		Length:        item.Length,
		Width:         item.Width,
		Height:        item.Height,
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
		CategoryID:    item.CategoryID,
		Barcode:       item.Barcode,
		UnitOfMeasure: item.UnitOfMeasure,
		Weight:        item.Weight,
		Length:        item.Length,
		Width:         item.Width,
		Height:        item.Height,
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
	if req.CategoryID != nil {
		category, err := uc.categoryRepo.FindByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		item.CategoryID = *req.CategoryID
		item.Category = category.Name
	}
	if req.UnitOfMeasure != nil {
		item.UnitOfMeasure = *req.UnitOfMeasure
	}
	if req.Weight != nil {
		item.Weight = *req.Weight
	}
	if req.Length != nil {
		item.Length = *req.Length
	}
	if req.Width != nil {
		item.Width = *req.Width
	}
	if req.Height != nil {
		item.Height = *req.Height
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	item.Barcode = entity.GenerateBarcode(item.SKU, item.Name, item.Category)
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
		CategoryID:    item.CategoryID,
		Barcode:       item.Barcode,
		UnitOfMeasure: item.UnitOfMeasure,
		Weight:        item.Weight,
		Length:        item.Length,
		Width:         item.Width,
		Height:        item.Height,
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
			CategoryID:    item.CategoryID,
			Barcode:       item.Barcode,
			UnitOfMeasure: item.UnitOfMeasure,
			Weight:        item.Weight,
			Length:        item.Length,
			Width:         item.Width,
			Height:        item.Height,
			IsActive:      item.IsActive,
		})
	}

	return result, total, nil
}
