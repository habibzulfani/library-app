package registered_books_access

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"project/configs"
	resC "project/controllers/resource_controllers"
	"project/middlewares"
	bookM "project/models/book_models"
	userM "project/models/user_models"
	"project/validators"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UpdateBook(c echo.Context) error {
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

	var book bookM.Book
	err = configs.DB.Preload("Authors").Where("id = ?", inputID).First(&book).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "Book not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to retrieve book data",
			"error":   err.Error(),
		})
	}

	var updateInput bookM.BookUpdateInput
	if err := c.Bind(&updateInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := validators.ValidateBookUpdate(&updateInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	err = configs.DB.Transaction(func(tx *gorm.DB) error {
		// Update book fields
		if updateInput.Title != nil {
			book.Title = *updateInput.Title
		}
		if updateInput.Publisher != nil {
			book.Publisher = *updateInput.Publisher
		}
		if updateInput.PublishedYear != nil {
			book.PublishedYear = *updateInput.PublishedYear
		}
		if updateInput.ISBN != nil {
			if *updateInput.ISBN != book.ISBN {
				// Check if the new ISBN is unique
				var count int64
				if err := tx.Model(&bookM.Book{}).Where("isbn = ? AND id != ?", *updateInput.ISBN, book.ID).Count(&count).Error; err != nil {
					return fmt.Errorf("failed to check ISBN uniqueness: %w", err)
				}
				if count > 0 {
					return fmt.Errorf("ISBN %s is already in use", *updateInput.ISBN)
				}
				// Update the ISBN
				book.ISBN = *updateInput.ISBN
			}
			// If the ISBN is the same as the existing one, do nothing
		}
		if updateInput.DOI != nil {
			updateInput.DOI = resC.NormalizeDOI(updateInput.DOI)
			if updateInput.DOI != &book.DOI {
				if !resC.IsDOIUnique(tx, updateInput.DOI) {
					return c.JSON(http.StatusConflict, map[string]string{
						"message": "A document with this DOI already exists in the repository",
					})
				}
				book.DOI = *updateInput.DOI
			}

		}
		if updateInput.Summary != nil {
			book.Summary = *updateInput.Summary
		}
		if updateInput.Subject != nil {
			book.Subject = *updateInput.Subject
		}
		if updateInput.Language != nil {
			book.Language = *updateInput.Language
		}
		if updateInput.Pages != nil {
			book.Pages = *updateInput.Pages
		}

		// Handle file update
		file, err := c.FormFile("bookFile")
		if err == nil {
			// New file uploaded, update it
			if book.FileURL != "" {
				// Delete the old file
				if err := os.Remove(book.FileURL); err != nil {
					log.Printf("Error deleting old file: %v", err)
					return fmt.Errorf("failed to delete old file: %v", err)
				}
			}

			// Generate new filename and save the new file
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
			uploadDir := "./uploads/books"
			if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
				log.Printf("Error creating upload directory: %v", err)
				return fmt.Errorf("failed to create upload directory: %v", err)
			}
			filepath := filepath.Join(uploadDir, filename)

			if err := configs.SaveUploadedFile(file, filepath); err != nil {
				log.Printf("Error saving uploaded file: %v", err)
				return fmt.Errorf("failed to save new file: %v", err)
			}

			book.FileURL = filepath
		}

		// Update authors if provided
		// Inside the transaction in UpdateBook function
		if len(updateInput.Authors) > 0 {
			// Fetch existing BookAuthor entries
			var existingAuthors []bookM.BookAuthor
			if err := tx.Where("book_id = ?", book.ID).Find(&existingAuthors).Error; err != nil {
				log.Printf("Error fetching existing book authors: %v", err)
				return err
			}

			// Create a map for quick lookup of existing authors
			existingAuthorsMap := make(map[uint]bookM.BookAuthor)
			for _, author := range existingAuthors {
				if author.UserID != nil {
					existingAuthorsMap[*author.UserID] = author
				}
			}

			// Process each author in the update input
			for i, authorName := range updateInput.Authors {
				var userID *uint
				if updateInput.AuthorNIMs[i] != "NON_REGISTERED" {
					var user userM.User
					if err := tx.Where("id_number = ?", updateInput.AuthorNIMs[i]).First(&user).Error; err != nil {
						if err == gorm.ErrRecordNotFound {
							return fmt.Errorf("user with NIM/NISN %s not found", updateInput.AuthorNIMs[i])
						}
						return fmt.Errorf("error finding user: %w", err)
					}
					userID = &user.ID
				}

				if userID != nil {
					// Check if this author already exists
					if existingAuthor, ok := existingAuthorsMap[*userID]; ok {
						// Update existing entry if name has changed
						if existingAuthor.AuthorName != authorName {
							existingAuthor.AuthorName = authorName
							if err := tx.Save(&existingAuthor).Error; err != nil {
								return fmt.Errorf("failed to update existing book author: %w", err)
							}
						}
						// Remove from map to mark as processed
						delete(existingAuthorsMap, *userID)
					} else {
						// Create new entry for registered author
						newAuthor := bookM.BookAuthor{
							BookID:    book.ID,
							UserID:     userID,
							AuthorName: authorName,
						}
						if err := tx.Create(&newAuthor).Error; err != nil {
							return fmt.Errorf("failed to create new book author association: %w", err)
						}
					}
				} else {
					// Handle non-registered author
					// Check if a non-registered author with the same name already exists
					var existingNonRegistered bookM.BookAuthor
					err := tx.Where("book_id = ? AND user_id IS NULL AND author_name = ?", book.ID, authorName).First(&existingNonRegistered).Error
					if err == gorm.ErrRecordNotFound {
						// Create new entry for non-registered author
						newAuthor := bookM.BookAuthor{
							BookID:    book.ID,
							UserID:     nil,
							AuthorName: authorName,
						}
						if err := tx.Create(&newAuthor).Error; err != nil {
							return fmt.Errorf("failed to create new non-registered book author: %w", err)
						}
					} else if err != nil {
						return fmt.Errorf("error checking for existing non-registered author: %w", err)
					}
					// If found, no need to update as the name is the same
				}
			}

			// Remove any remaining authors that weren't in the update input
			for _, remainingAuthor := range existingAuthorsMap {
				if err := tx.Delete(&remainingAuthor).Error; err != nil {
					return fmt.Errorf("failed to remove outdated author: %w", err)
				}
			}

			// Update the Author field in the book model
			book.Author = strings.Join(updateInput.Authors, ", ")
		}
		// Save the updated book
		if err := tx.Save(&book).Error; err != nil {
			log.Printf("Error saving updated book: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to update book data",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully updated book data",
		"book":    book,
	})
}

func joinAuthors(authors []string) string {
	return strings.Join(authors, "; ")
}
