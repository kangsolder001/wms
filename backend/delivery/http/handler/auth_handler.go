package handler

import (
	"encoding/json"
	"net/http"

	"wms/application/dto"
	"wms/application/usecase"
	"wms/delivery/http/response"
	"wms/pkg/logger"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	authUC usecase.AuthUsecase
	log    logger.Logger
}

func NewAuthHandler(authUC usecase.AuthUsecase, log logger.Logger) *AuthHandler {
	return &AuthHandler{authUC: authUC, log: log}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authUC.Login(r.Context(), &dto.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authUC.Register(r.Context(), &dto.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     req.Role,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	result, err := h.authUC.GetProfile(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 50

	result, total, err := h.authUC.ListUsers(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(w, http.StatusOK, result, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authUC.UpdateUser(r.Context(), id, &req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.authUC.DeleteUser(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "user deactivated"})
}
