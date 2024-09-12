package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
)

// Product structure
type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
}

// ProductRepository interface
type ProductRepository interface {
	CreateProduct(ctx context.Context, product Product) (uuid.UUID, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error)
	UpdateProduct(ctx context.Context, product Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

// Repo structure working in PostgreSQL
type productRepository struct {
	db *sql.DB
}

// Create new ProductRepository
func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

// Adding product
func (r *productRepository) CreateProduct(ctx context.Context, product Product) (uuid.UUID, error) {
	product.ID = uuid.New()

	query := `INSERT INTO products (id, name, description, price) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, product.ID, product.Name, product.Description, product.Price)
	if err != nil {
		log.Println("Failed to insert product:", err)
		return uuid.Nil, err
	}
	return product.ID, nil
}

// Bring Product by ID
func (r *productRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	query := `SELECT id, name, description, price FROM products WHERE id = $1`

	var product Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// Update product
func (r *productRepository) UpdateProduct(ctx context.Context, product Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.ID)
	return err
}

// Delete product
func (r *productRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
