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

type CategoryHandler struct {
	categoryUC usecase.CategoryUsecase
	log        logger.Logger
}

func NewCategoryHandler(categoryUC usecase.CategoryUsecase, log logger.Logger) *CategoryHandler {
	return &CategoryHandler{categoryUC: categoryUC, log: log}
}

func (h *CategoryHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	items, total, err := h.categoryUC.ListCategories(c.Request.Context(), page, limit)
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

func (h *CategoryHandler) ListAll(c *gin.Context) {
	items, err := h.categoryUC.ListAllCategories(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, items)
}

func (h *CategoryHandler) Get(c *gin.Context) {
	id := c.Param("id")

	item, err := h.categoryUC.GetCategory(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, item)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.categoryUC.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, item)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.categoryUC.UpdateCategory(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, item)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.categoryUC.DeleteCategory(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "category deleted"})
}
