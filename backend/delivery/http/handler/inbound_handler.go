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

type InboundHandler struct {
	inboundUC usecase.InboundUsecase
	log       logger.Logger
}

func NewInboundHandler(inboundUC usecase.InboundUsecase, log logger.Logger) *InboundHandler {
	return &InboundHandler{inboundUC: inboundUC, log: log}
}

func (h *InboundHandler) ListPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	pos, total, err := h.inboundUC.ListPurchaseOrders(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(w, http.StatusOK, pos, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *InboundHandler) GetPurchaseOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	po, err := h.inboundUC.GetPurchaseOrder(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, po)
}

func (h *InboundHandler) CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePurchaseOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := r.Context().Value("user_id").(string)

	po, err := h.inboundUC.CreatePurchaseOrder(r.Context(), &req, userID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, po)
}

func (h *InboundHandler) ReceiveGoods(w http.ResponseWriter, r *http.Request) {
	poID := chi.URLParam(r, "id")

	var req dto.ReceiveGoodsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := r.Context().Value("user_id").(string)

	grn, err := h.inboundUC.ReceiveGoods(r.Context(), poID, &req, userID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, grn)
}
