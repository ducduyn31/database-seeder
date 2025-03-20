package models

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"database-test/pkg/faker"

	"github.com/pterm/pterm"
)

// Create a package-level random source
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// Address represents a shipping or billing address
type Address struct {
	ID           int
	UserID       int
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   string
	Country      string
	IsDefault    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GenerateAddresses generates n fake addresses for each user and inserts them into the database
func GenerateAddresses(db *sql.DB, usersCount, addressesPerUser int) error {
	userIDs, err := GetRandomUserIDs(db, usersCount)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`
		INSERT INTO addresses (user_id, address_line1, address_line2, city, state, postal_code, country, is_default)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	totalAddresses := usersCount * addressesPerUser

	// Create a progress bar
	progressBar, _ := pterm.DefaultProgressbar.
		WithTotal(totalAddresses).
		WithTitle(fmt.Sprintf("Generating %d addresses for %d users (%d per user)...", totalAddresses, usersCount, addressesPerUser)).
		Start()

	addressesGenerated := 0

	for _, userID := range userIDs {
		for j := 0; j < addressesPerUser; j++ {
			addressLine1 := fmt.Sprintf("%d %s St", random.Intn(1000)+1, faker.City())
			var addressLine2 sql.NullString
			if faker.Boolean() {
				addressLine2 = sql.NullString{String: fmt.Sprintf("Apt %d", random.Intn(100)+1), Valid: true}
			}
			city := faker.City()
			state := faker.State()
			postalCode := faker.Zip()
			country := faker.CountryAbbr()
			isDefault := j == 0 // First address is default

			var id int
			err := stmt.QueryRow(
				userID, addressLine1, addressLine2, city, state, postalCode, country, isDefault,
			).Scan(&id)
			if err != nil {
				return fmt.Errorf("failed to insert address: %w", err)
			}

			addressesGenerated++
			progressBar.Increment()
		}
	}

	return nil
}

// GetRandomAddressIDs returns n random address IDs from the database
func GetRandomAddressIDs(db *sql.DB, count int) ([]int, error) {
	rows, err := db.Query("SELECT id FROM addresses ORDER BY RANDOM() LIMIT $1", count)
	if err != nil {
		return nil, fmt.Errorf("failed to get random address IDs: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan address ID: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// GetRandomAddressIDsByUser returns a random address ID for each user ID
func GetRandomAddressIDsByUser(db *sql.DB, userIDs []int) (map[int]int, error) {
	if len(userIDs) == 0 {
		return make(map[int]int), nil
	}

	// Prepare the query with a dynamic number of placeholders
	query := "SELECT id, user_id FROM addresses WHERE user_id IN ("
	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query += ") ORDER BY RANDOM()"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get random address IDs by user: %w", err)
	}
	defer rows.Close()

	result := make(map[int]int)
	for rows.Next() {
		var addressID, userID int
		if err := rows.Scan(&addressID, &userID); err != nil {
			return nil, fmt.Errorf("failed to scan address and user ID: %w", err)
		}
		// Only add the first address for each user (which will be random due to ORDER BY RANDOM())
		if _, exists := result[userID]; !exists {
			result[userID] = addressID
		}
	}

	return result, nil
}
