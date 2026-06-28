package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wms/application/dto"
	"wms/application/usecase"
	"wms/delivery/http/response"
	"wms/pkg/logger"
)

type InventoryHandler struct {
	inventoryUC usecase.InventoryUsecase
	log         logger.Logger
}

func NewInventoryHandler(inventoryUC usecase.InventoryUsecase, log logger.Logger) *InventoryHandler {
	return &InventoryHandler{inventoryUC: inventoryUC, log: log}
}

func (h *InventoryHandler) ListStock(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	stocks, total, err := h.inventoryUC.GetStock(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(w, http.StatusOK, stocks, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *InventoryHandler) GetStockByItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("item_id")
	if itemID == "" {
		response.Error(w, http.StatusBadRequest, "item_id is required")
		return
	}

	stocks, err := h.inventoryUC.GetStockByItem(r.Context(), itemID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, stocks)
}

func (h *InventoryHandler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	var req dto.AdjustStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := r.Context().Value("user_id").(string)

	if err := h.inventoryUC.AdjustStock(r.Context(), &req, userID); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "stock adjusted"})
}

func (h *InventoryHandler) ListMovements(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	movements, total, err := h.inventoryUC.GetStockMovements(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(w, http.StatusOK, movements, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}
