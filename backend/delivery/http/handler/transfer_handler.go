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

type TransferHandler struct {
	transferUC usecase.TransferUsecase
	log        logger.Logger
}

func NewTransferHandler(transferUC usecase.TransferUsecase, log logger.Logger) *TransferHandler {
	return &TransferHandler{transferUC: transferUC, log: log}
}

func (h *TransferHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	transfers, total, err := h.transferUC.ListTransfers(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(c, http.StatusOK, transfers, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *TransferHandler) Create(c *gin.Context) {
	var req dto.CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := c.GetString("user_id")

	transfer, err := h.transferUC.CreateTransfer(c.Request.Context(), &req, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, transfer)
}

func (h *TransferHandler) Complete(c *gin.Context) {
	id := c.Param("id")

	if err := h.transferUC.CompleteTransfer(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "transfer completed"})
}
