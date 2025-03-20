package models

import (
	"database/sql"
	"fmt"
	"time"

	"database-test/pkg/faker"

	"github.com/pterm/pterm"
)

// Use the package-level random source

// Product represents a product in the e-commerce system
type Product struct {
	ID            int
	Name          string
	Description   string
	Price         float64
	StockQuantity int
	CategoryID    int
	SKU           string
	Weight        sql.NullFloat64
	Dimensions    sql.NullString
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ProductImage represents an image associated with a product
type ProductImage struct {
	ID        int
	ProductID int
	ImageURL  string
	IsPrimary bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GenerateProducts generates n fake products and inserts them into the database
func GenerateProducts(db *sql.DB, count int, imagesPerProduct int) error {
	// Get random category IDs
	categoryIDs, err := GetRandomCategoryIDs(db, count)
	if err != nil {
		return err
	}

	// If we don't have enough categories, reuse them
	for len(categoryIDs) < count {
		categoryIDs = append(categoryIDs, categoryIDs...)
	}
	categoryIDs = categoryIDs[:count]

	// Prepare product statement
	productStmt, err := db.Prepare(`
		INSERT INTO products (
			name, description, price, stock_quantity, category_id, 
			sku, weight, dimensions, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare product statement: %w", err)
	}
	defer productStmt.Close()

	// Prepare image statement
	imageStmt, err := db.Prepare(`
		INSERT INTO product_images (
			product_id, image_url, is_primary, created_at, updated_at
		)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare image statement: %w", err)
	}
	defer imageStmt.Close()

	// Create a progress bar
	progressBar, _ := pterm.DefaultProgressbar.
		WithTotal(count).
		WithTitle(fmt.Sprintf("Generating %d products...", count)).
		Start()

	for i := 0; i < count; i++ {
		// Generate product data
		name := faker.ProductName()
		description := faker.ProductDescription()
		price := faker.Price(9.99, 999.99)
		stockQuantity := random.Intn(1000) + 1
		categoryID := categoryIDs[i]
		sku := faker.SKU()

		// 80% chance of having weight
		var weight sql.NullFloat64
		if random.Float64() < 0.8 {
			weight = sql.NullFloat64{
				Float64: faker.Weight(0.1, 20.0),
				Valid:   true,
			}
		}

		// 70% chance of having dimensions
		var dimensions sql.NullString
		if random.Float64() < 0.7 {
			dimensions = sql.NullString{
				String: faker.Dimensions(),
				Valid:  true,
			}
		}

		// Insert product
		var productID int
		err := productStmt.QueryRow(
			name, description, price, stockQuantity, categoryID,
			sku, weight, dimensions,
		).Scan(&productID)

		if err != nil {
			return fmt.Errorf("failed to insert product: %w", err)
		}

		// Generate images for this product
		for j := 0; j < imagesPerProduct; j++ {
			imageURL := faker.ImageURL(productID)
			isPrimary := j == 0 // First image is primary

			_, err := imageStmt.Exec(productID, imageURL, isPrimary)
			if err != nil {
				return fmt.Errorf("failed to insert product image: %w", err)
			}
		}

		progressBar.Increment()
	}

	return nil
}

// GetRandomProductIDs returns n random product IDs from the database
func GetRandomProductIDs(db *sql.DB, count int) ([]int, error) {
	rows, err := db.Query("SELECT id FROM products ORDER BY RANDOM() LIMIT $1", count)
	if err != nil {
		return nil, fmt.Errorf("failed to get random product IDs: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan product ID: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
