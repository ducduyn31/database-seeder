package models

import (
	"database/sql"
	"fmt"
	"time"

	"database-test/pkg/faker"

	"github.com/pterm/pterm"
)

// Review represents a product review
type Review struct {
	ID        int
	ProductID int
	UserID    int
	Rating    int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GenerateReviews generates reviews for products
func GenerateReviews(db *sql.DB, count int) error {
	// Get random user IDs
	userIDs, err := GetRandomUserIDs(db, count)
	if err != nil {
		return err
	}

	// If we don't have enough users, reuse them
	for len(userIDs) < count {
		userIDs = append(userIDs, userIDs...)
	}
	userIDs = userIDs[:count]

	// Get random product IDs
	productIDs, err := GetRandomProductIDs(db, count)
	if err != nil {
		return err
	}

	// If we don't have enough products, reuse them
	for len(productIDs) < count {
		productIDs = append(productIDs, productIDs...)
	}
	productIDs = productIDs[:count]

	// Prepare review statement
	stmt, err := db.Prepare(`
		INSERT INTO reviews (
			product_id, user_id, rating, title, content, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Create a progress bar
	progressBar, _ := pterm.DefaultProgressbar.
		WithTotal(count).
		WithTitle(fmt.Sprintf("Generating %d reviews...", count)).
		Start()

	reviewsGenerated := 0
	for i := 0; i < count; i++ {
		productID := productIDs[i]
		userID := userIDs[i]

		// Generate review data
		rating := random.Intn(5) + 1 // 1-5 stars
		title := faker.ReviewTitle()
		content := faker.ReviewContent()

		// Insert review
		var id int
		err := stmt.QueryRow(
			productID, userID, rating, title, content,
		).Scan(&id)

		if err != nil {
			// Check if this user already reviewed this product
			if err.Error() == "pq: duplicate key value violates unique constraint \"reviews_product_id_user_id_key\"" {
				// Skip this review
				continue
			}
			return fmt.Errorf("failed to insert review: %w", err)
		}

		reviewsGenerated++
		progressBar.Increment()
	}

	return nil
}
