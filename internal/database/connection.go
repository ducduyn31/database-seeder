package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Config holds the database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DefaultConfig returns the default database configuration from docker-compose.yml
func DefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     5433, // Default to postgres1, can be overridden
		User:     "shared_user",
		Password: "shared_password",
		DBName:   "shared_db",
		SSLMode:  "disable",
	}
}

// Connect establishes a connection to the database
func Connect(config Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connected to database %s on %s:%d", config.DBName, config.Host, config.Port)
	return db, nil
}

// CreateTables creates all the necessary tables for the e-commerce database
func CreateTables(db *sql.DB) error {
	queries := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			phone VARCHAR(20),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		// Addresses table
		`CREATE TABLE IF NOT EXISTS addresses (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id),
			address_line1 VARCHAR(255) NOT NULL,
			address_line2 VARCHAR(255),
			city VARCHAR(100) NOT NULL,
			state VARCHAR(100) NOT NULL,
			postal_code VARCHAR(20) NOT NULL,
			country VARCHAR(100) NOT NULL,
			is_default BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		// Categories table
		`CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			parent_id INT REFERENCES categories(id),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		// Products table
		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			stock_quantity INT NOT NULL,
			category_id INT NOT NULL REFERENCES categories(id),
			sku VARCHAR(50) UNIQUE NOT NULL,
			weight DECIMAL(8, 2),
			dimensions VARCHAR(50),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		// Product images table
		`CREATE TABLE IF NOT EXISTS product_images (
			id SERIAL PRIMARY KEY,
			product_id INT NOT NULL REFERENCES products(id),
			image_url VARCHAR(255) NOT NULL,
			is_primary BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		// Orders table
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id),
			status VARCHAR(50) NOT NULL,
			total_amount DECIMAL(10, 2) NOT NULL,
			shipping_address_id INT NOT NULL REFERENCES addresses(id),
			billing_address_id INT NOT NULL REFERENCES addresses(id),
			payment_method VARCHAR(50) NOT NULL,
			shipping_method VARCHAR(50) NOT NULL,
			tracking_number VARCHAR(100),
			notes TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		// Order items table
		`CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INT NOT NULL REFERENCES orders(id),
			product_id INT NOT NULL REFERENCES products(id),
			quantity INT NOT NULL,
			price_per_unit DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		// Reviews table
		`CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL PRIMARY KEY,
			product_id INT NOT NULL REFERENCES products(id),
			user_id INT NOT NULL REFERENCES users(id),
			rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w\nQuery: %s", err, query)
		}
	}

	log.Println("All tables created successfully")
	return nil
}
