package seed

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"database-test/internal/database"
	"database-test/internal/models"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	// Flags for data generation
	userCount        int
	addressesPerUser int
	categoryCount    int
	maxCategoryDepth int
	productCount     int
	imagesPerProduct int
	orderCount       int
	maxItemsPerOrder int
	reviewCount      int
	allFlag          bool
)

// Command represents the seed command
var Command = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with fake data",
	Long: `Seed command generates and inserts fake data into your database.
You can specify which types of data to generate and how many records to create.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get database configuration from flags
		portStr := cmd.Flag("port").Value.String()
		var port int
		fmt.Sscanf(portStr, "%d", &port)

		config := database.Config{
			Host:     cmd.Flag("host").Value.String(),
			Port:     port,
			User:     cmd.Flag("user").Value.String(),
			Password: cmd.Flag("password").Value.String(),
			DBName:   cmd.Flag("dbname").Value.String(),
			SSLMode:  cmd.Flag("sslmode").Value.String(),
		}

		// Connect to database
		db, err := database.Connect(config)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		// Create tables if they don't exist
		if err := database.CreateTables(db); err != nil {
			log.Fatalf("Failed to create tables: %v", err)
		}

		// Start timing
		startTime := time.Now()

		// Show header
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Println("E-Commerce Database Seeder")
		pterm.Println() // Empty line

		// Seed data based on flags
		if allFlag || userCount > 0 {
			if err := seedUsers(db, userCount); err != nil {
				pterm.Error.Println("Failed to seed users:", err)
				return
			}
		}

		if allFlag || (userCount > 0 && addressesPerUser > 0) {
			if err := seedAddresses(db, userCount, addressesPerUser); err != nil {
				pterm.Error.Println("Failed to seed addresses:", err)
				return
			}
		}

		if allFlag || categoryCount > 0 {
			if err := seedCategories(db, categoryCount, maxCategoryDepth); err != nil {
				pterm.Error.Println("Failed to seed categories:", err)
				return
			}
		}

		if allFlag || productCount > 0 {
			if err := seedProducts(db, productCount, imagesPerProduct); err != nil {
				pterm.Error.Println("Failed to seed products:", err)
				return
			}
		}

		if allFlag || orderCount > 0 {
			if err := seedOrders(db, orderCount, maxItemsPerOrder); err != nil {
				pterm.Error.Println("Failed to seed orders:", err)
				return
			}
		}

		if allFlag || reviewCount > 0 {
			if err := seedReviews(db, reviewCount); err != nil {
				pterm.Error.Println("Failed to seed reviews:", err)
				return
			}
		}

		// Print summary
		duration := time.Since(startTime)
		pterm.Println() // Empty line
		pterm.Success.Println("Seeding completed successfully!")
		pterm.Info.Printf("Total time: %s\n", duration)
	},
}

func init() {
	// Add flags for data generation
	Command.Flags().IntVar(&userCount, "users", 100, "Number of users to generate")
	Command.Flags().IntVar(&addressesPerUser, "addresses-per-user", 2, "Number of addresses per user")
	Command.Flags().IntVar(&categoryCount, "categories", 30, "Number of categories to generate")
	Command.Flags().IntVar(&maxCategoryDepth, "category-depth", 3, "Maximum depth of category hierarchy")
	Command.Flags().IntVar(&productCount, "products", 1000, "Number of products to generate")
	Command.Flags().IntVar(&imagesPerProduct, "images-per-product", 3, "Number of images per product")
	Command.Flags().IntVar(&orderCount, "orders", 500, "Number of orders to generate")
	Command.Flags().IntVar(&maxItemsPerOrder, "max-items-per-order", 5, "Maximum number of items per order")
	Command.Flags().IntVar(&reviewCount, "reviews", 300, "Number of reviews to generate")
	Command.Flags().BoolVar(&allFlag, "all", false, "Generate all types of data")
}

// Helper functions to seed different types of data

func seedUsers(db *sql.DB, count int) error {
	pterm.DefaultSection.Println("Seeding Users")
	spinner, _ := pterm.DefaultSpinner.
		WithShowTimer(true).
		WithText("Generating users...").
		Start()

	err := models.GenerateUsers(db, count)

	if err != nil {
		spinner.Fail("Failed to generate users")
		return err
	}

	spinner.Success("Successfully generated " + pterm.Green(fmt.Sprintf("%d", count)) + " users")
	return nil
}

func seedAddresses(db *sql.DB, userCount, addressesPerUser int) error {
	pterm.DefaultSection.Println("Seeding Addresses")
	spinner, _ := pterm.DefaultSpinner.
		WithShowTimer(true).
		WithText("Generating addresses...").
		Start()

	err := models.GenerateAddresses(db, userCount, addressesPerUser)

	if err != nil {
		spinner.Fail("Failed to generate addresses")
		return err
	}

	totalAddresses := userCount * addressesPerUser
	spinner.Success("Successfully generated " + pterm.Green(fmt.Sprintf("%d", totalAddresses)) + " addresses")
	return nil
}

func seedCategories(db *sql.DB, count, maxDepth int) error {
	pterm.DefaultSection.Println("Seeding Categories")
	spinner, _ := pterm.DefaultSpinner.
		WithShowTimer(true).
		WithText("Generating categories...").
		Start()

	err := models.GenerateCategories(db, count, maxDepth)

	if err != nil {
		spinner.Fail("Failed to generate categories")
		return err
	}

	spinner.Success("Successfully generated " + pterm.Green(fmt.Sprintf("%d", count)) + " categories")
	return nil
}

func seedProducts(db *sql.DB, count, imagesPerProduct int) error {
	pterm.DefaultSection.Println("Seeding Products")
	spinner, _ := pterm.DefaultSpinner.
		WithShowTimer(true).
		WithText("Generating products...").
		Start()

	err := models.GenerateProducts(db, count, imagesPerProduct)

	if err != nil {
		spinner.Fail("Failed to generate products")
		return err
	}

	spinner.Success("Successfully generated " + pterm.Green(fmt.Sprintf("%d", count)) + " products with " +
		pterm.Green(fmt.Sprintf("%d", count*imagesPerProduct)) + " images")
	return nil
}

func seedOrders(db *sql.DB, count, maxItemsPerOrder int) error {
	pterm.DefaultSection.Println("Seeding Orders")
	spinner, _ := pterm.DefaultSpinner.
		WithShowTimer(true).
		WithText("Generating orders...").
		Start()

	err := models.GenerateOrders(db, count, maxItemsPerOrder)

	if err != nil {
		spinner.Fail("Failed to generate orders")
		return err
	}

	spinner.Success("Successfully generated " + pterm.Green(fmt.Sprintf("%d", count)) + " orders")
	return nil
}

func seedReviews(db *sql.DB, count int) error {
	pterm.DefaultSection.Println("Seeding Reviews")
	spinner, _ := pterm.DefaultSpinner.
		WithShowTimer(true).
		WithText("Generating reviews...").
		Start()

	err := models.GenerateReviews(db, count)

	if err != nil {
		spinner.Fail("Failed to generate reviews")
		return err
	}

	spinner.Success("Successfully generated " + pterm.Green(fmt.Sprintf("%d", count)) + " reviews")
	return nil
}
