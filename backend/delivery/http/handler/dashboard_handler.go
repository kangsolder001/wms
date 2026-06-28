package handler

import (
	"net/http"

	"wms/application/usecase"
	"wms/delivery/http/response"
	"wms/pkg/logger"
)

type DashboardHandler struct {
	dashboardUC usecase.DashboardUsecase
	log         logger.Logger
}

func NewDashboardHandler(dashboardUC usecase.DashboardUsecase, log logger.Logger) *DashboardHandler {
	return &DashboardHandler{dashboardUC: dashboardUC, log: log}
}

func (h *DashboardHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.dashboardUC.GetSummary(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, summary)
}
