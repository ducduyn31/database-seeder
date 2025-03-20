package models

import (
	"database/sql"
	"fmt"
	"time"

	"database-test/pkg/faker"

	"github.com/pterm/pterm"
)

// Category represents a product category
type Category struct {
	ID          int
	Name        string
	Description string
	ParentID    sql.NullInt64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GenerateCategories generates n fake categories and inserts them into the database
func GenerateCategories(db *sql.DB, count int, maxDepth int) error {
	// First, create top-level categories (about 1/3 of total)
	topLevelCount := count / 3
	if topLevelCount < 1 {
		topLevelCount = 1
	}

	stmt, err := db.Prepare(`
		INSERT INTO categories (name, description, parent_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Create a progress bar for top-level categories
	topLevelBar, _ := pterm.DefaultProgressbar.
		WithTotal(topLevelCount).
		WithTitle(fmt.Sprintf("Generating %d top-level categories...", topLevelCount)).
		Start()

	// Generate top-level categories
	topLevelIDs := make([]int, 0, topLevelCount)
	for i := 0; i < topLevelCount; i++ {
		name := faker.CategoryName()
		description := faker.CategoryDescription()

		var id int
		err := stmt.QueryRow(name, description, nil).Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to insert top-level category: %w", err)
		}
		topLevelIDs = append(topLevelIDs, id)
		topLevelBar.Increment()
	}

	// Generate subcategories
	remainingCount := count - topLevelCount
	if remainingCount <= 0 {
		return nil
	}

	// Create a progress bar for subcategories
	subCategoryBar, _ := pterm.DefaultProgressbar.
		WithTotal(remainingCount).
		WithTitle(fmt.Sprintf("Generating %d subcategories...", remainingCount)).
		Start()

	// Create a map to track categories by depth
	categoriesByDepth := make(map[int][]int)
	categoriesByDepth[1] = topLevelIDs

	currentDepth := 1
	for remainingCount > 0 && currentDepth < maxDepth {
		parentIDs := categoriesByDepth[currentDepth]
		if len(parentIDs) == 0 {
			break
		}

		// Determine how many subcategories to create at this depth
		subcategoriesPerParent := remainingCount / len(parentIDs)
		if subcategoriesPerParent < 1 {
			subcategoriesPerParent = 1
		}

		nextDepthCategories := make([]int, 0)
		categoriesCreated := 0

		for _, parentID := range parentIDs {
			for j := 0; j < subcategoriesPerParent && remainingCount > 0; j++ {
				name := faker.CategoryName()
				description := faker.CategoryDescription()

				var id int
				err := stmt.QueryRow(name, description, parentID).Scan(&id)
				if err != nil {
					return fmt.Errorf("failed to insert subcategory: %w", err)
				}
				nextDepthCategories = append(nextDepthCategories, id)
				remainingCount--
				categoriesCreated++
				subCategoryBar.Increment()
			}
		}

		categoriesByDepth[currentDepth+1] = nextDepthCategories
		currentDepth++
	}

	return nil
}

// GetRandomCategoryIDs returns n random category IDs from the database
func GetRandomCategoryIDs(db *sql.DB, count int) ([]int, error) {
	rows, err := db.Query("SELECT id FROM categories ORDER BY RANDOM() LIMIT $1", count)
	if err != nil {
		return nil, fmt.Errorf("failed to get random category IDs: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan category ID: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
