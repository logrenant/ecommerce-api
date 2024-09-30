package handlers

import (
	"context"
	"ecommerce-api/internal/repository"
	"ecommerce-api/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Product handler structure
type ProductHandler struct {
	repo         repository.ProductRepository
	MinIOService *services.MinIOService
}

// Create a new ProductHandler
func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// Create new product
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var product repository.Product
	product.Name = r.FormValue("name")
	product.Description = r.FormValue("description")
	price := r.FormValue("price")
	product.Price, err = strconv.ParseFloat(price, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	imageName := fmt.Sprintf("%s.jpg", uuid.New().String())
	imageURL, err := h.MinIOService.UploadFile(context.Background(), file, imageName)
	if err != nil {
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	fmt.Println("image", imageURL)

	product.ImageURL = imageURL

	id, err := h.repo.CreateProduct(context.Background(), product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        id,
		"image_url": product.ImageURL,
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

// Get all products
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	// Fetch products from the repository
	products, err := h.repo.GetAllProducts(context.Background())
	if err != nil {
		log.Println("Error fetching products:", err)
		http.Error(w, "Unable to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Println("Error encoding products:", err)
		http.Error(w, "Unable to encode products to JSON", http.StatusInternalServerError)
	}
}

// Update product
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var product repository.Product
	if err := r.ParseMultipartForm(10 << 20); err != nil { // FormData için
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	// Gelen verileri al
	product.Name = r.FormValue("name")
	product.Description = r.FormValue("description")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		http.Error(w, "invalid price", http.StatusBadRequest)
		return
	}
	product.Price = price
	product.ID = id

	// Resim dosyası varsa güncelle
	if file, _, err := r.FormFile("image"); err == nil {
		// MinIO'da yükleme işlemini yap
		imageURL, err := h.MinIOService.UploadFile(r.Context(), file, product.ID.String())
		if err != nil {
			http.Error(w, "failed to upload image", http.StatusInternalServerError)
			return
		}
		product.ImageURL = imageURL // Resim URL'sini güncelle
	}

	// Veritabanında güncelle
	if err := h.repo.UpdateProduct(r.Context(), product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product) // Güncellenmiş ürünü döndür
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
