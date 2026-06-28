package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wms/application/dto"
	"wms/application/usecase"
	"wms/delivery/http/response"
	"wms/pkg/logger"

	"github.com/go-chi/chi/v5"
)

type OutboundHandler struct {
	outboundUC usecase.OutboundUsecase
	log        logger.Logger
}

func NewOutboundHandler(outboundUC usecase.OutboundUsecase, log logger.Logger) *OutboundHandler {
	return &OutboundHandler{outboundUC: outboundUC, log: log}
}

func (h *OutboundHandler) ListSalesOrders(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	sos, total, err := h.outboundUC.ListSalesOrders(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(w, http.StatusOK, sos, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *OutboundHandler) GetSalesOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	so, err := h.outboundUC.GetSalesOrder(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, so)
}

func (h *OutboundHandler) CreateSalesOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSalesOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := r.Context().Value("user_id").(string)

	so, err := h.outboundUC.CreateSalesOrder(r.Context(), &req, userID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, so)
}

func (h *OutboundHandler) PickOrder(w http.ResponseWriter, r *http.Request) {
	soID := chi.URLParam(r, "id")

	var req dto.PickListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := r.Context().Value("user_id").(string)

	if err := h.outboundUC.PickOrder(r.Context(), soID, &req, userID); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "order picked"})
}

func (h *OutboundHandler) ShipOrder(w http.ResponseWriter, r *http.Request) {
	soID := chi.URLParam(r, "id")

	var req dto.ShipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := r.Context().Value("user_id").(string)

	shipment, err := h.outboundUC.ShipOrder(r.Context(), soID, &req, userID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, shipment)
}
