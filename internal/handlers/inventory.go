package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/RyderDProgrammer/Inventory-Manager/internal/models"
	"github.com/RyderDProgrammer/Inventory-Manager/internal/repository"
)

type Handler struct {
	repo repository.InventoryRepository
}

func NewHandler(repo repository.InventoryRepository) *Handler {
	return &Handler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve items")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.repo.GetByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve item")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if item.Name == "" || item.SKU == "" {
		writeError(w, http.StatusBadRequest, "name and sku are required")
		return
	}
	created, err := h.repo.Create(item)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create item")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	updated, err := h.repo.Update(id, item)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update item")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.repo.Delete(id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
