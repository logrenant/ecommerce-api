package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"ecommerce-api/internal/handlers"
	"ecommerce-api/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// get DATABASE_URL from .env file
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	// Connect Database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Start repository and handlers
	productRepo := repository.NewProductRepository(db)
	productHandler := handlers.NewProductHandler(productRepo)

	// Create Chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)

	// Configure CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Veya "*" tüm kökenlere izin verir
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Apply CORS middleware
	r.Use(corsHandler.Handler)

	// Basic GET endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the E-commerce API!"))
	})

	// Routes for CRUD functions
	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products/{id}", productHandler.GetProductByID)
	r.Get("/products", productHandler.GetAllProducts)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)

	// Start server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
