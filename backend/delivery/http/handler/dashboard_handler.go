package handler

import (
	"net/http"

	"wms/application/usecase"
	"wms/delivery/http/response"
	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardUC usecase.DashboardUsecase
	log         logger.Logger
}

func NewDashboardHandler(dashboardUC usecase.DashboardUsecase, log logger.Logger) *DashboardHandler {
	return &DashboardHandler{dashboardUC: dashboardUC, log: log}
}

func (h *DashboardHandler) GetSummary(c *gin.Context) {
	summary, err := h.dashboardUC.GetSummary(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, summary)
}
