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

type TransferHandler struct {
	transferUC usecase.TransferUsecase
	log        logger.Logger
}

func NewTransferHandler(transferUC usecase.TransferUsecase, log logger.Logger) *TransferHandler {
	return &TransferHandler{transferUC: transferUC, log: log}
}

func (h *TransferHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	transfers, total, err := h.transferUC.ListTransfers(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(w, http.StatusOK, transfers, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *TransferHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := r.Context().Value("user_id").(string)

	transfer, err := h.transferUC.CreateTransfer(r.Context(), &req, userID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, transfer)
}

func (h *TransferHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.transferUC.CompleteTransfer(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "transfer completed"})
}
