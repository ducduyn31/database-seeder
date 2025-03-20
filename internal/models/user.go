package models

import (
	"database/sql"
	"fmt"
	"time"

	gofaker "github.com/go-faker/faker/v4"
	"github.com/pterm/pterm"
)

// Use the package-level random source

// User represents a user in the e-commerce system
type User struct {
	ID           int
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Phone        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GenerateUsers generates n fake users and inserts them into the database
func GenerateUsers(db *sql.DB, count int) error {
	stmt, err := db.Prepare(`
		INSERT INTO users (email, password_hash, first_name, last_name, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Create a progress bar
	progressBar, _ := pterm.DefaultProgressbar.
		WithTotal(count).
		WithTitle(fmt.Sprintf("Generating %d users...", count)).
		Start()

	for i := 0; i < count; i++ {
		email := gofaker.Email()
		passwordHash := gofaker.Password() // In a real app, this would be properly hashed
		firstName := gofaker.FirstName()
		lastName := gofaker.LastName()
		phone := gofaker.Phonenumber()

		var id int
		err := stmt.QueryRow(email, passwordHash, firstName, lastName, phone).Scan(&id)
		if err != nil {
			// If it's a duplicate email, try again
			if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
				i--
				continue
			}
			return fmt.Errorf("failed to insert user: %w", err)
		}

		progressBar.Increment()
	}

	return nil
}

// GetRandomUserIDs returns n random user IDs from the database
func GetRandomUserIDs(db *sql.DB, count int) ([]int, error) {
	rows, err := db.Query("SELECT id FROM users ORDER BY RANDOM() LIMIT $1", count)
	if err != nil {
		return nil, fmt.Errorf("failed to get random user IDs: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan user ID: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
