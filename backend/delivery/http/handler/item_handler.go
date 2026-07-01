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

type ItemHandler struct {
	itemUC usecase.ItemUsecase
	log    logger.Logger
}

func NewItemHandler(itemUC usecase.ItemUsecase, log logger.Logger) *ItemHandler {
	return &ItemHandler{itemUC: itemUC, log: log}
}

func (h *ItemHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	items, total, err := h.itemUC.ListItems(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(c, http.StatusOK, items, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *ItemHandler) Get(c *gin.Context) {
	id := c.Param("id")

	item, err := h.itemUC.GetItem(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, item)
}

func (h *ItemHandler) Create(c *gin.Context) {
	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.itemUC.CreateItem(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, item)
}

func (h *ItemHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.itemUC.UpdateItem(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, item)
}

func (h *ItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.itemUC.DeleteItem(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "item deleted"})
}
