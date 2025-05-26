package configs

import (
	"fmt"
	bookM "project/models/book_models"

	"gorm.io/gorm"
)

func MigrateBooks(db *gorm.DB) error {
	// Check if the DOI column exists
	if db.Migrator().HasColumn(&bookM.Book{}, "doi") {
		fmt.Println("DOI column already exists")
	} else {
		// Add DOI column if it doesn't exist
		if err := db.Migrator().AddColumn(&bookM.Book{}, "doi"); err != nil {
			return fmt.Errorf("failed to add DOI column: %w", err)
		}
		fmt.Println("Added DOI column")
	}

	// Update column definitions
	columnsToAlter := []string{
		"publisher", "published_year", "summary", "subject",
		"language", "pages", "file_url", "doi",
	}

	for _, column := range columnsToAlter {
		if err := db.Migrator().AlterColumn(&bookM.Book{}, column); err != nil {
			return fmt.Errorf("failed to alter column %s: %w", column, err)
		}
		fmt.Printf("Altered column: %s\n", column)
	}

	return nil
}
