package resource_controllers

import (
	"fmt"
	"net/http"
	"os"
	"project/configs"
	"project/middlewares"
	authorM "project/models/author_models"
	bookM "project/models/book_models"
	paperM "project/models/paper_models"

	"github.com/labstack/echo/v4"
)

type DeleteItem struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

func DeleteBibliographies(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	var req struct {
		Items []DeleteItem `json:"items"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	type DeletionResult struct {
		ID      int    `json:"id"`
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	results := []DeletionResult{}

	for _, item := range req.Items {
		result := DeletionResult{ID: item.ID, Success: true}

		var err error
		switch item.Type {
		case "book":
			err = deleteBookById(item.ID)
		case "paper":
			err = deletePaperById(item.ID)
		default:
			// Try to determine the type from the database
			if exists, _ := bookExists(item.ID); exists {
				err = deleteBookById(item.ID)
			} else if exists, _ := paperExists(item.ID); exists {
				err = deletePaperById(item.ID)
			} else {
				result.Success = false
				result.Message = "Item not found or invalid type"
			}
		}

		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("Failed to delete: %v", err)
		}

		results = append(results, result)
	}

	failedDeletions := []DeletionResult{}
	for _, result := range results {
		if !result.Success {
			failedDeletions = append(failedDeletions, result)
		}
	}

	if len(failedDeletions) > 0 {
		return c.JSON(http.StatusPartialContent, map[string]interface{}{
			"success": false,
			"message": "Some items failed to delete",
			"results": results,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "All items deleted successfully",
		"results": results,
	})
}

func bookExists(id int) (bool, error) {
	var count int64
	result := configs.DB.Model(&bookM.Book{}).Where("id = ?", id).Count(&count)
	return count > 0, result.Error
}

func paperExists(id int) (bool, error) {
	var count int64
	result := configs.DB.Model(&paperM.Paper{}).Where("id = ?", id).Count(&count)
	return count > 0, result.Error
}

func deleteBookById(id int) error {
	tx := configs.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// Fetch the book to get the file path
	var book bookM.Book
	if err := tx.First(&book, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete associated book_authors
	if err := tx.Where("book_id = ?", id).Delete(&authorM.BookAuthor{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the book from the database
	if err := tx.Delete(&book).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the file
	if err := os.Remove(book.FileURL); err != nil {
		// Log the error but don't fail the transaction
		fmt.Printf("Failed to delete file for book %d: %v\n", id, err)
	}

	return tx.Commit().Error
}

func deletePaperById(id int) error {
    tx := configs.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Error; err != nil {
        return err
    }

    // Fetch the paper to get the file path
    var paper paperM.Paper
    if err := tx.First(&paper, id).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Delete associated paper_authors
    if err := tx.Where("paper_id = ?", id).Delete(&authorM.PaperAuthor{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Delete the paper from the database
    if err := tx.Delete(&paper).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Delete the file
    if err := os.Remove(paper.FileURL); err != nil {
        // Log the error but don't fail the transaction
        fmt.Printf("Failed to delete file for paper %d: %v\n", id, err)
    }

    return tx.Commit().Error
}
