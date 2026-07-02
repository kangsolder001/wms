package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"wms/application/dto"
	"wms/application/usecase"
	"wms/delivery/http/response"
	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

type InboundHandler struct {
	inboundUC usecase.InboundUsecase
	log       logger.Logger
}

func NewInboundHandler(inboundUC usecase.InboundUsecase, log logger.Logger) *InboundHandler {
	return &InboundHandler{inboundUC: inboundUC, log: log}
}

func (h *InboundHandler) ListPurchaseOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	pos, total, err := h.inboundUC.ListPurchaseOrders(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(c, http.StatusOK, pos, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *InboundHandler) GetPurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	po, err := h.inboundUC.GetPurchaseOrder(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, po)
}

func (h *InboundHandler) CreatePurchaseOrder(c *gin.Context) {
	var req dto.CreatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := c.GetString("user_id")

	po, err := h.inboundUC.CreatePurchaseOrder(c.Request.Context(), &req, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, po)
}

func (h *InboundHandler) ReceiveGoods(c *gin.Context) {
	poID := c.Param("id")

	var req struct {
		Notes string `json:"notes"`
		Items []struct {
			ItemID     string  `json:"item_id"`
			Quantity   float64 `json:"quantity"`
			LocationID string  `json:"location_id"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Items) == 0 {
		response.Error(c, http.StatusBadRequest, "items is required")
		return
	}

	userID := c.GetString("user_id")

	grnNumber := fmt.Sprintf("GRN-%s", time.Now().Format("20060102150405"))

	receiveReq := &dto.ReceiveGoodsRequest{
		GRNNumber: grnNumber,
		Notes:     req.Notes,
	}
	for _, item := range req.Items {
		receiveReq.Items = append(receiveReq.Items, dto.ReceiveItemRequest{
			ItemID:     item.ItemID,
			Quantity:   item.Quantity,
			LocationID: item.LocationID,
		})
	}

	grn, err := h.inboundUC.ReceiveGoods(c.Request.Context(), poID, receiveReq, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, grn)
}

func (h *InboundHandler) ApprovePurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.inboundUC.ApprovePurchaseOrder(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "purchase order approved"})
}
