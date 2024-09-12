package handlers

import (
	"context"
	"ecommerce-api/internal/repository"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Product handler structure
type ProductHandler struct {
	repo repository.ProductRepository
}

// Create a new ProductHandler
func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// Create new product
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product repository.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repo.CreateProduct(context.Background(), product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

// Fetching product by ID
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.repo.GetProductByID(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// Update product
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product repository.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateProduct(context.Background(), product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete product
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteProduct(context.Background(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
