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

type OutboundHandler struct {
	outboundUC usecase.OutboundUsecase
	log        logger.Logger
}

func NewOutboundHandler(outboundUC usecase.OutboundUsecase, log logger.Logger) *OutboundHandler {
	return &OutboundHandler{outboundUC: outboundUC, log: log}
}

func (h *OutboundHandler) ListSalesOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	sos, total, err := h.outboundUC.ListSalesOrders(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(c, http.StatusOK, sos, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *OutboundHandler) GetSalesOrder(c *gin.Context) {
	id := c.Param("id")

	so, err := h.outboundUC.GetSalesOrder(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, so)
}

func (h *OutboundHandler) CreateSalesOrder(c *gin.Context) {
	var req dto.CreateSalesOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := c.GetString("user_id")

	so, err := h.outboundUC.CreateSalesOrder(c.Request.Context(), &req, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, so)
}

func (h *OutboundHandler) PickOrder(c *gin.Context) {
	soID := c.Param("id")

	var req dto.PickListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := c.GetString("user_id")

	if err := h.outboundUC.PickOrder(c.Request.Context(), soID, &req, userID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "order picked"})
}

func (h *OutboundHandler) ShipOrder(c *gin.Context) {
	soID := c.Param("id")

	var req dto.ShipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := c.GetString("user_id")

	shipment, err := h.outboundUC.ShipOrder(c.Request.Context(), soID, &req, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, shipment)
}
