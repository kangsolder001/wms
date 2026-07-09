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

type ZoneHandler struct {
	zoneUC usecase.ZoneUsecase
	log    logger.Logger
}

func NewZoneHandler(zoneUC usecase.ZoneUsecase, log logger.Logger) *ZoneHandler {
	return &ZoneHandler{zoneUC: zoneUC, log: log}
}

func (h *ZoneHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	items, total, err := h.zoneUC.ListZones(c.Request.Context(), page, limit)
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

func (h *ZoneHandler) ListAll(c *gin.Context) {
	items, err := h.zoneUC.ListAllZones(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, items)
}

func (h *ZoneHandler) Get(c *gin.Context) {
	id := c.Param("id")

	item, err := h.zoneUC.GetZone(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, item)
}

func (h *ZoneHandler) Create(c *gin.Context) {
	var req dto.CreateZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.zoneUC.CreateZone(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, item)
}

func (h *ZoneHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.zoneUC.UpdateZone(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, item)
}

func (h *ZoneHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.zoneUC.DeleteZone(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "zone deleted"})
}
