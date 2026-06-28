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

type LocationUsecase interface {
	CreateLocation(ctx context.Context, req *dto.CreateLocationRequest) (*dto.LocationResponse, error)
	GetLocation(ctx context.Context, id string) (*dto.LocationResponse, error)
	UpdateLocation(ctx context.Context, id string, req *dto.UpdateLocationRequest) (*dto.LocationResponse, error)
	DeleteLocation(ctx context.Context, id string) error
	ListLocations(ctx context.Context, page, limit int) ([]*dto.LocationResponse, int, error)
}

type locationUsecase struct {
	locationRepo repository.LocationRepository
	log          logger.Logger
}

func NewLocationUsecase(locationRepo repository.LocationRepository, log logger.Logger) LocationUsecase {
	return &locationUsecase{locationRepo: locationRepo, log: log}
}

func (uc *locationUsecase) CreateLocation(ctx context.Context, req *dto.CreateLocationRequest) (*dto.LocationResponse, error) {
	uc.log.Info("creating location", "code", req.Code, "name", req.Name)

	existing, _ := uc.locationRepo.FindByCode(ctx, req.Code)
	if existing != nil {
		return nil, errors.New("location code already exists")
	}

	location := &entity.Location{
		Code:      req.Code,
		Name:      req.Name,
		Zone:      req.Zone,
		Aisle:     req.Aisle,
		Rack:      req.Rack,
		Level:     req.Level,
		Bin:       req.Bin,
		Type:      req.Type,
		Capacity:  req.Capacity,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := uc.locationRepo.Create(ctx, location); err != nil {
		uc.log.Error("failed to create location", "code", req.Code, "error", err)
		return nil, err
	}

	uc.log.Info("location created successfully", "id", location.ID, "code", location.Code)

	return &dto.LocationResponse{
		ID:       location.ID,
		Code:     location.Code,
		Name:     location.Name,
		Zone:     location.Zone,
		Aisle:    location.Aisle,
		Rack:     location.Rack,
		Level:    location.Level,
		Bin:      location.Bin,
		Type:     location.Type,
		Capacity: location.Capacity,
		IsActive: location.IsActive,
	}, nil
}

func (uc *locationUsecase) GetLocation(ctx context.Context, id string) (*dto.LocationResponse, error) {
	location, err := uc.locationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("location not found")
	}

	return &dto.LocationResponse{
		ID:       location.ID,
		Code:     location.Code,
		Name:     location.Name,
		Zone:     location.Zone,
		Aisle:    location.Aisle,
		Rack:     location.Rack,
		Level:    location.Level,
		Bin:      location.Bin,
		Type:     location.Type,
		Capacity: location.Capacity,
		IsActive: location.IsActive,
	}, nil
}

func (uc *locationUsecase) UpdateLocation(ctx context.Context, id string, req *dto.UpdateLocationRequest) (*dto.LocationResponse, error) {
	uc.log.Info("updating location", "id", id)

	location, err := uc.locationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("location not found")
	}

	if req.Name != nil {
		location.Name = *req.Name
	}
	if req.Zone != nil {
		location.Zone = *req.Zone
	}
	if req.Aisle != nil {
		location.Aisle = *req.Aisle
	}
	if req.Rack != nil {
		location.Rack = *req.Rack
	}
	if req.Level != nil {
		location.Level = *req.Level
	}
	if req.Bin != nil {
		location.Bin = *req.Bin
	}
	if req.Type != nil {
		location.Type = *req.Type
	}
	if req.Capacity != nil {
		location.Capacity = *req.Capacity
	}
	if req.IsActive != nil {
		location.IsActive = *req.IsActive
	}

	if err := uc.locationRepo.Update(ctx, location); err != nil {
		uc.log.Error("failed to update location", "id", id, "error", err)
		return nil, err
	}

	uc.log.Info("location updated successfully", "id", location.ID)

	return &dto.LocationResponse{
		ID:       location.ID,
		Code:     location.Code,
		Name:     location.Name,
		Zone:     location.Zone,
		Aisle:    location.Aisle,
		Rack:     location.Rack,
		Level:    location.Level,
		Bin:      location.Bin,
		Type:     location.Type,
		Capacity: location.Capacity,
		IsActive: location.IsActive,
	}, nil
}

func (uc *locationUsecase) DeleteLocation(ctx context.Context, id string) error {
	uc.log.Info("deleting location", "id", id)

	if err := uc.locationRepo.Delete(ctx, id); err != nil {
		uc.log.Error("failed to delete location", "id", id, "error", err)
		return err
	}

	uc.log.Info("location deleted successfully", "id", id)
	return nil
}

func (uc *locationUsecase) ListLocations(ctx context.Context, page, limit int) ([]*dto.LocationResponse, int, error) {
	locations, total, err := uc.locationRepo.List(ctx, page, limit)
	if err != nil {
		uc.log.Error("failed to list locations", "error", err)
		return nil, 0, err
	}

	var result []*dto.LocationResponse
	for _, loc := range locations {
		result = append(result, &dto.LocationResponse{
			ID:       loc.ID,
			Code:     loc.Code,
			Name:     loc.Name,
			Zone:     loc.Zone,
			Aisle:    loc.Aisle,
			Rack:     loc.Rack,
			Level:    loc.Level,
			Bin:      loc.Bin,
			Type:     loc.Type,
			Capacity: loc.Capacity,
			IsActive: loc.IsActive,
		})
	}

	return result, total, nil
}
