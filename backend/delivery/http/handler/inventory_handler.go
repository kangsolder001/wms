package handler

import (
	"net/http"
	"strconv"

	"wms/application/dto"
	"wms/application/usecase"
	"wms/delivery/http/response"
	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryUC usecase.InventoryUsecase
	log         logger.Logger
}

func NewInventoryHandler(inventoryUC usecase.InventoryUsecase, log logger.Logger) *InventoryHandler {
	return &InventoryHandler{inventoryUC: inventoryUC, log: log}
}

func (h *InventoryHandler) ListStock(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	itemID := c.Query("item_id")
	locationID := c.Query("location_id")
	search := c.Query("search")

	stocks, total, err := h.inventoryUC.GetStock(c.Request.Context(), page, limit, itemID, locationID, search)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(c, http.StatusOK, stocks, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *InventoryHandler) GetStockByItem(c *gin.Context) {
	itemID := c.Query("item_id")
	if itemID == "" {
		response.Error(c, http.StatusBadRequest, "item_id is required")
		return
	}

	stocks, err := h.inventoryUC.GetStockByItem(c.Request.Context(), itemID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, stocks)
}

func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	var req dto.AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := c.GetString("user_id")

	if err := h.inventoryUC.AdjustStock(c.Request.Context(), &req, userID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "stock adjusted"})
}

func (h *InventoryHandler) ListMovements(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	movements, total, err := h.inventoryUC.GetStockMovements(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(c, http.StatusOK, movements, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}
