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

type ZoneUsecase interface {
	CreateZone(ctx context.Context, req *dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetZone(ctx context.Context, id string) (*dto.ZoneResponse, error)
	UpdateZone(ctx context.Context, id string, req *dto.UpdateZoneRequest) (*dto.ZoneResponse, error)
	DeleteZone(ctx context.Context, id string) error
	ListZones(ctx context.Context, page, limit int) ([]*dto.ZoneResponse, int, error)
	ListAllZones(ctx context.Context) ([]*dto.ZoneResponse, error)
}

type zoneUsecase struct {
	zoneRepo repository.ZoneRepository
	log      logger.Logger
}

func NewZoneUsecase(zoneRepo repository.ZoneRepository, log logger.Logger) ZoneUsecase {
	return &zoneUsecase{zoneRepo: zoneRepo, log: log}
}

func (uc *zoneUsecase) CreateZone(ctx context.Context, req *dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	uc.log.Info("creating zone", "code", req.Code, "name", req.Name)

	zone := &entity.Zone{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	if err := uc.zoneRepo.Create(ctx, zone); err != nil {
		uc.log.Error("failed to create zone", "code", req.Code, "error", err)
		return nil, err
	}

	uc.log.Info("zone created successfully", "id", zone.ID)

	return &dto.ZoneResponse{
		ID:          zone.ID,
		Code:        zone.Code,
		Name:        zone.Name,
		Description: zone.Description,
		IsActive:    zone.IsActive,
	}, nil
}

func (uc *zoneUsecase) GetZone(ctx context.Context, id string) (*dto.ZoneResponse, error) {
	zone, err := uc.zoneRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("zone not found")
	}

	return &dto.ZoneResponse{
		ID:          zone.ID,
		Code:        zone.Code,
		Name:        zone.Name,
		Description: zone.Description,
		IsActive:    zone.IsActive,
	}, nil
}

func (uc *zoneUsecase) UpdateZone(ctx context.Context, id string, req *dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	uc.log.Info("updating zone", "id", id)

	zone, err := uc.zoneRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("zone not found")
	}

	if req.Code != nil {
		zone.Code = *req.Code
	}
	if req.Name != nil {
		zone.Name = *req.Name
	}
	if req.Description != nil {
		zone.Description = *req.Description
	}
	if req.IsActive != nil {
		zone.IsActive = *req.IsActive
	}

	if err := uc.zoneRepo.Update(ctx, zone); err != nil {
		uc.log.Error("failed to update zone", "id", id, "error", err)
		return nil, err
	}

	uc.log.Info("zone updated successfully", "id", zone.ID)

	return &dto.ZoneResponse{
		ID:          zone.ID,
		Code:        zone.Code,
		Name:        zone.Name,
		Description: zone.Description,
		IsActive:    zone.IsActive,
	}, nil
}

func (uc *zoneUsecase) DeleteZone(ctx context.Context, id string) error {
	uc.log.Info("deleting zone", "id", id)

	if err := uc.zoneRepo.Delete(ctx, id); err != nil {
		uc.log.Error("failed to delete zone", "id", id, "error", err)
		return err
	}

	uc.log.Info("zone deleted successfully", "id", id)
	return nil
}

func (uc *zoneUsecase) ListZones(ctx context.Context, page, limit int) ([]*dto.ZoneResponse, int, error) {
	zones, total, err := uc.zoneRepo.List(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var result []*dto.ZoneResponse
	for _, z := range zones {
		result = append(result, &dto.ZoneResponse{
			ID:          z.ID,
			Code:        z.Code,
			Name:        z.Name,
			Description: z.Description,
			IsActive:    z.IsActive,
		})
	}

	return result, total, nil
}

func (uc *zoneUsecase) ListAllZones(ctx context.Context) ([]*dto.ZoneResponse, error) {
	zones, err := uc.zoneRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []*dto.ZoneResponse
	for _, z := range zones {
		result = append(result, &dto.ZoneResponse{
			ID:          z.ID,
			Code:        z.Code,
			Name:        z.Name,
			Description: z.Description,
			IsActive:    z.IsActive,
		})
	}

	return result, nil
}
