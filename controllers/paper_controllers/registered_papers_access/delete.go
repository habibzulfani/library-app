package registered_papers_access

import (
	"net/http"
	"project/configs"
	"project/middlewares"
	paperM "project/models/paper_models"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// func DeletePaper(c echo.Context) error {
// 	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
// 		return err
// 	}

// 	inputID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]interface{}{
// 			"message": "Invalid id",
// 			"error":   err.Error(),
// 		})
// 	}

// 	var paper paperM.Paper
// 	result := configs.DB.Delete(&paper, inputID)

// 	if result.RowsAffected == 0 {
// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 			"message": "Paper not found",
// 		})
// 	}

// 	if result.Error != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 			"message": "Failed to delete paper",
// 			"error":   result.Error.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"search": "Successfully deleted paper data",
// 	})
// }

func DeletePaper(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	inputID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid id",
			"error":   err.Error(),
		})
	}

	err = configs.DB.Transaction(func(tx *gorm.DB) error {
		// Soft delete the book
		result := tx.Model(&paperM.Paper{}).Where("id = ?", inputID).Update("deleted_at", time.Now())
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "Book not found")
		}

		// Remove associations in the book_authors table
		if err := tx.Exec("DELETE FROM book_authors WHERE book_id = ?", inputID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			return c.JSON(he.Code, map[string]interface{}{
				"message": he.Message,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to delete book",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully soft-deleted book data",
	})
}

// If you need to retrieve soft-deleted books (for admin purposes)
func GetAllBooksIncludingDeleted(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	var books []paperM.Paper
	if err := configs.DB.Unscoped().Find(&books).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to retrieve books",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, books)
}

// If you need to permanently delete a book
func PermanentlyDeleteBook(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	inputID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid id",
			"error":   err.Error(),
		})
	}

	err = configs.DB.Transaction(func(tx *gorm.DB) error {
		// Remove associations in the book_authors table
		if err := tx.Exec("DELETE FROM book_authors WHERE book_id = ?", inputID).Error; err != nil {
			return err
		}

		// Permanently delete the book
		if err := tx.Unscoped().Delete(&paperM.Paper{}, inputID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to permanently delete book",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully permanently deleted book data",
	})
}

type DeleteBooksRequest struct {
	BookIDs []uint `json:"book_ids" validate:"required,min=1"`
}

func DeleteMultipleBooks(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	var req DeleteBooksRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if len(req.BookIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "No book IDs provided",
		})
	}

	var deletedCount int64
	err := configs.DB.Transaction(func(tx *gorm.DB) error {
		// Soft delete the books
		result := tx.Model(&paperM.Paper{}).Where("id IN ?", req.BookIDs).Update("deleted_at", time.Now())
		if result.Error != nil {
			return result.Error
		}
		deletedCount = result.RowsAffected

		// Remove associations in the book_authors table
		if err := tx.Exec("DELETE FROM book_authors WHERE book_id IN ?", req.BookIDs).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to delete books",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Successfully soft-deleted books",
		"deleted_count": deletedCount,
	})
}

// Function to restore soft-deleted books
func RestoreMultipleBooks(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	var req DeleteBooksRequest // We can reuse the same struct for restoration
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if len(req.BookIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "No book IDs provided",
		})
	}

	result := configs.DB.Model(&paperM.Paper{}).Where("id IN ?", req.BookIDs).Update("deleted_at", nil)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to restore books",
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":        "Successfully restored books",
		"restored_count": result.RowsAffected,
	})
}

// Function to permanently delete multiple books
func PermanentlyDeleteMultipleBooks(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	var req DeleteBooksRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if len(req.BookIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "No book IDs provided",
		})
	}

	err := configs.DB.Transaction(func(tx *gorm.DB) error {
		// Remove associations in the book_authors table
		if err := tx.Exec("DELETE FROM book_authors WHERE book_id IN ?", req.BookIDs).Error; err != nil {
			return err
		}

		// Permanently delete the books
		result := tx.Unscoped().Delete(&paperM.Paper{}, req.BookIDs)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to permanently delete books",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully permanently deleted books",
	})
}
