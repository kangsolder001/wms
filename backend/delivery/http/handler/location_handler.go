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

type LocationHandler struct {
	locationUC usecase.LocationUsecase
	log        logger.Logger
}

func NewLocationHandler(locationUC usecase.LocationUsecase, log logger.Logger) *LocationHandler {
	return &LocationHandler{locationUC: locationUC, log: log}
}

func (h *LocationHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	locations, total, err := h.locationUC.ListLocations(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSONWithMeta(w, http.StatusOK, locations, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

func (h *LocationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	location, err := h.locationUC.GetLocation(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, location)
}

func (h *LocationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	location, err := h.locationUC.CreateLocation(r.Context(), &req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, location)
}

func (h *LocationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	location, err := h.locationUC.UpdateLocation(r.Context(), id, &req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, location)
}

func (h *LocationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.locationUC.DeleteLocation(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "location deleted"})
}
