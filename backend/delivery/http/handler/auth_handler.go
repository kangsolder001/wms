package handler

import (
	"net/http"

	"wms/application/dto"
	"wms/application/usecase"
	"wms/delivery/http/response"
	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUC usecase.AuthUsecase
	log    logger.Logger
}

func NewAuthHandler(authUC usecase.AuthUsecase, log logger.Logger) *AuthHandler {
	return &AuthHandler{authUC: authUC, log: log}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authUC.Login(c.Request.Context(), &dto.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, result)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authUC.Register(c.Request.Context(), &dto.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     req.Role,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, result)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.authUC.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, result)
}

func (h *AuthHandler) ListUsers(c *gin.Context) {
	page := 1
	limit := 50

	result, total, err := h.authUC.ListUsers(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(c, http.StatusOK, result, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authUC.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, result)
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.authUC.DeleteUser(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "user deactivated"})
}
