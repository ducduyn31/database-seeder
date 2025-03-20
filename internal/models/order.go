package models

import (
	"database/sql"
	"fmt"
	"time"

	"database-test/pkg/faker"

	"github.com/pterm/pterm"
)

// Order represents an order in the e-commerce system
type Order struct {
	ID                int
	UserID            int
	Status            string
	TotalAmount       float64
	ShippingAddressID int
	BillingAddressID  int
	PaymentMethod     string
	ShippingMethod    string
	TrackingNumber    sql.NullString
	Notes             sql.NullString
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID           int
	OrderID      int
	ProductID    int
	Quantity     int
	PricePerUnit float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GenerateOrders generates n fake orders and inserts them into the database
func GenerateOrders(db *sql.DB, count int, maxItemsPerOrder int) error {
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
	productIDs, err := GetRandomProductIDs(db, 100) // Get a pool of products to choose from
	if err != nil {
		return err
	}

	// Get addresses for each user
	userAddresses, err := GetRandomAddressIDsByUser(db, userIDs)
	if err != nil {
		return err
	}

	// Prepare order statement
	orderStmt, err := db.Prepare(`
		INSERT INTO orders (
			user_id, status, total_amount, shipping_address_id, billing_address_id,
			payment_method, shipping_method, tracking_number, notes, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare order statement: %w", err)
	}
	defer orderStmt.Close()

	// Prepare order item statement
	itemStmt, err := db.Prepare(`
		INSERT INTO order_items (
			order_id, product_id, quantity, price_per_unit, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare order item statement: %w", err)
	}
	defer itemStmt.Close()

	// Get product prices
	productPrices := make(map[int]float64)
	for _, productID := range productIDs {
		var price float64
		err := db.QueryRow("SELECT price FROM products WHERE id = $1", productID).Scan(&price)
		if err != nil {
			return fmt.Errorf("failed to get product price: %w", err)
		}
		productPrices[productID] = price
	}

	// Create a progress bar
	progressBar, _ := pterm.DefaultProgressbar.
		WithTotal(count).
		WithTitle(fmt.Sprintf("Generating %d orders...", count)).
		Start()

	for i := 0; i < count; i++ {
		userID := userIDs[i]

		// Get address IDs for this user
		addressID, ok := userAddresses[userID]
		if !ok {
			// If no address found, skip this order
			pterm.Warning.Printf("No address found for user %d, skipping order", userID)
			continue
		}

		// Use the same address for shipping and billing (could be randomized)
		shippingAddressID := addressID
		billingAddressID := addressID

		// Generate order data
		status := faker.OrderStatus()
		paymentMethod := faker.PaymentMethod()
		shippingMethod := faker.ShippingMethod()

		// 70% chance of having a tracking number if status is not "Pending"
		var trackingNumber sql.NullString
		if status != "Pending" && random.Float64() < 0.7 {
			trackingNumber = sql.NullString{
				String: faker.TrackingNumber(),
				Valid:  true,
			}
		}

		// 30% chance of having notes
		var notes sql.NullString
		if random.Float64() < 0.3 {
			notes = sql.NullString{
				String: "Please deliver to the front door.",
				Valid:  true,
			}
		}

		// Generate order items
		numItems := random.Intn(maxItemsPerOrder) + 1
		orderItems := make([]struct {
			ProductID    int
			Quantity     int
			PricePerUnit float64
		}, numItems)

		totalAmount := 0.0

		// Select random products for this order
		for j := 0; j < numItems; j++ {
			productID := productIDs[random.Intn(len(productIDs))]
			quantity := random.Intn(5) + 1
			price := productPrices[productID]

			orderItems[j] = struct {
				ProductID    int
				Quantity     int
				PricePerUnit float64
			}{
				ProductID:    productID,
				Quantity:     quantity,
				PricePerUnit: price,
			}

			totalAmount += float64(quantity) * price
		}

		// Insert order
		var orderID int
		err := orderStmt.QueryRow(
			userID, status, totalAmount, shippingAddressID, billingAddressID,
			paymentMethod, shippingMethod, trackingNumber, notes,
		).Scan(&orderID)

		if err != nil {
			return fmt.Errorf("failed to insert order: %w", err)
		}

		// Insert order items
		for _, item := range orderItems {
			_, err := itemStmt.Exec(
				orderID, item.ProductID, item.Quantity, item.PricePerUnit,
			)
			if err != nil {
				return fmt.Errorf("failed to insert order item: %w", err)
			}
		}

		progressBar.Increment()
	}

	return nil
}
