package usecase

import (
	"context"
	"errors"
	"time"

	"wms/application/dto"
	"wms/domain/entity"
	"wms/domain/repository"
	"wms/infrastructure/auth"
	"wms/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)
	GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error)
	ListUsers(ctx context.Context, page, limit int) ([]*dto.UserResponse, int, error)
	UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

type authUsecase struct {
	userRepo   repository.UserRepository
	jwtService auth.JWTService
	log        logger.Logger
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtService auth.JWTService, log logger.Logger) AuthUsecase {
	return &authUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
		log:        log,
	}
}

func (uc *authUsecase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	uc.log.Info("login attempt", "username", req.Username)

	user, err := uc.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		uc.log.Error("user not found", "username", req.Username, "error", err)
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		uc.log.Error("invalid password", "username", req.Username)
		return nil, errors.New("invalid credentials")
	}

	token, err := uc.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		uc.log.Error("failed to generate token", "user_id", user.ID, "error", err)
		return nil, err
	}

	uc.log.Info("login successful", "user_id", user.ID, "username", user.Username)

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}, nil
}

func (uc *authUsecase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	uc.log.Info("registering new user", "username", req.Username, "email", req.Email)

	existing, _ := uc.userRepo.FindByUsername(ctx, req.Username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	existing, _ = uc.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.log.Error("failed to hash password", "error", err)
		return nil, err
	}

	user := &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         req.Role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.log.Error("failed to create user", "username", req.Username, "error", err)
		return nil, err
	}

	uc.log.Info("user registered successfully", "user_id", user.ID, "username", user.Username)

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	}, nil
}

func (uc *authUsecase) GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		uc.log.Error("user not found", "user_id", userID, "error", err)
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	}, nil
}

func (uc *authUsecase) ListUsers(ctx context.Context, page, limit int) ([]*dto.UserResponse, int, error) {
	users, total, err := uc.userRepo.List(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var result []*dto.UserResponse
	for _, u := range users {
		result = append(result, &dto.UserResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			FullName: u.FullName,
			Role:     u.Role,
			IsActive: u.IsActive,
		})
	}
	return result, total, nil
}

func (uc *authUsecase) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Email != "" {
		existing, _ := uc.userRepo.FindByEmail(ctx, req.Email)
		if existing != nil && existing.ID != id {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hashedPassword)
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		IsActive: user.IsActive,
	}, nil
}

func (uc *authUsecase) DeleteUser(ctx context.Context, id string) error {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	active := false
	user.IsActive = active
	user.UpdatedAt = time.Now()
	return uc.userRepo.Update(ctx, user)
}
