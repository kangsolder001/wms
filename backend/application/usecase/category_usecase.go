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

type CategoryUsecase interface {
	CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategory(ctx context.Context, id string) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id string) error
	ListCategories(ctx context.Context, page, limit int) ([]*dto.CategoryResponse, int, error)
	ListAllCategories(ctx context.Context) ([]*dto.CategoryResponse, error)
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
	log          logger.Logger
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository, log logger.Logger) CategoryUsecase {
	return &categoryUsecase{categoryRepo: categoryRepo, log: log}
}

func (uc *categoryUsecase) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	uc.log.Info("creating category", "name", req.Name, "abbreviation", req.Abbreviation)

	category := &entity.Category{
		Name:         req.Name,
		Abbreviation: req.Abbreviation,
		Description:  req.Description,
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		uc.log.Error("failed to create category", "name", req.Name, "error", err)
		return nil, err
	}

	uc.log.Info("category created successfully", "id", category.ID)

	return &dto.CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		Abbreviation: category.Abbreviation,
		Description:  category.Description,
		IsActive:     category.IsActive,
	}, nil
}

func (uc *categoryUsecase) GetCategory(ctx context.Context, id string) (*dto.CategoryResponse, error) {
	category, err := uc.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	return &dto.CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		Abbreviation: category.Abbreviation,
		Description:  category.Description,
		IsActive:     category.IsActive,
	}, nil
}

func (uc *categoryUsecase) UpdateCategory(ctx context.Context, id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	uc.log.Info("updating category", "id", id)

	category, err := uc.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Abbreviation != nil {
		category.Abbreviation = *req.Abbreviation
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		uc.log.Error("failed to update category", "id", id, "error", err)
		return nil, err
	}

	uc.log.Info("category updated successfully", "id", category.ID)

	return &dto.CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		Abbreviation: category.Abbreviation,
		Description:  category.Description,
		IsActive:     category.IsActive,
	}, nil
}

func (uc *categoryUsecase) DeleteCategory(ctx context.Context, id string) error {
	uc.log.Info("deleting category", "id", id)

	if err := uc.categoryRepo.Delete(ctx, id); err != nil {
		uc.log.Error("failed to delete category", "id", id, "error", err)
		return err
	}

	uc.log.Info("category deleted successfully", "id", id)
	return nil
}

func (uc *categoryUsecase) ListCategories(ctx context.Context, page, limit int) ([]*dto.CategoryResponse, int, error) {
	categories, total, err := uc.categoryRepo.List(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var result []*dto.CategoryResponse
	for _, c := range categories {
		result = append(result, &dto.CategoryResponse{
			ID:           c.ID,
			Name:         c.Name,
			Abbreviation: c.Abbreviation,
			Description:  c.Description,
			IsActive:     c.IsActive,
		})
	}

	return result, total, nil
}

func (uc *categoryUsecase) ListAllCategories(ctx context.Context) ([]*dto.CategoryResponse, error) {
	categories, err := uc.categoryRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []*dto.CategoryResponse
	for _, c := range categories {
		result = append(result, &dto.CategoryResponse{
			ID:           c.ID,
			Name:         c.Name,
			Abbreviation: c.Abbreviation,
			Description:  c.Description,
			IsActive:     c.IsActive,
		})
	}

	return result, nil
}
