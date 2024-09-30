package main

import (
	"context" // Bu importu ekleyin
	"database/sql"
	"log"
	"net/http"
	"os"

	"ecommerce-api/internal/handlers"
	"ecommerce-api/internal/repository"
	"ecommerce-api/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/cors"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get DATABASE_URL and MinIO configurations from .env file
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	minioBucket := os.Getenv("MINIO_BUCKET")

	// Connect Database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create MinIO client
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("Failed to create MinIO client:", err)
	}

	// Ensure the bucket exists
	if err := minioClient.MakeBucket(context.Background(), minioBucket, minio.MakeBucketOptions{}); err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), minioBucket)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", minioBucket)
		} else {
			log.Fatal("Failed to create bucket:", err)
		}
	}

	// Start repository and handlers
	productRepo := repository.NewProductRepository(db)
	productHandler := handlers.NewProductHandler(productRepo)

	// Create MinIO service
	minioService, err := services.NewMinIOService(minioEndpoint, minioAccessKey, minioSecretKey, minioBucket)
	if err != nil {
		log.Fatal("Failed to create MinIO service:", err)
	}

	// Update the product handler to use MinIO service
	productHandler.MinIOService = minioService

	// Create Chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)

	// Configure CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
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
