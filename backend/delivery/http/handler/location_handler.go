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

type LocationHandler struct {
	locationUC usecase.LocationUsecase
	log        logger.Logger
}

func NewLocationHandler(locationUC usecase.LocationUsecase, log logger.Logger) *LocationHandler {
	return &LocationHandler{locationUC: locationUC, log: log}
}

func (h *LocationHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	locations, total, err := h.locationUC.ListLocations(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(c, http.StatusOK, locations, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *LocationHandler) Get(c *gin.Context) {
	id := c.Param("id")

	location, err := h.locationUC.GetLocation(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, location)
}

func (h *LocationHandler) Create(c *gin.Context) {
	var req dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	location, err := h.locationUC.CreateLocation(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, location)
}

func (h *LocationHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	location, err := h.locationUC.UpdateLocation(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, location)
}

func (h *LocationHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.locationUC.DeleteLocation(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "location deleted"})
}
