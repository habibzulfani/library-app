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
	authorM "project/models/author_models"
	bookM "project/models/book_models"
	userM "project/models/user_models"
	"project/validators"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateBook(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	var bookInput bookM.CreateBookInput
	if err := c.Bind(&bookInput); err != nil {
		log.Printf("Error binding input: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	log.Printf("Received book input: %+v", bookInput)

	if err := validators.ValidateCreateBook(&bookInput); err != nil {
		log.Printf("Validation error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	if len(bookInput.Authors) != len(bookInput.AuthorNIMs) {
		log.Printf("Mismatch in Authors and AuthorNIMs length: Authors=%d, AuthorNIMs=%d", len(bookInput.Authors), len(bookInput.AuthorNIMs))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Authors and AuthorNIMs must have the same length",
		})
	}

	file, err := c.FormFile("bookFile")
	if err != nil {
		log.Printf("Error getting book file: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Book file is required",
			"error":   err.Error(),
		})
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	uploadDir := "./uploads/books"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Error creating upload directory: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
	}
	filePath := filepath.Join(uploadDir, filename)

	if err := configs.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Error saving uploaded file: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to save file",
			"error":   err.Error(),
		})
	}

	book := bookM.Book{
		Title:         bookInput.Title,
		Author:        strings.Join(bookInput.Authors, ", "),
		Publisher:     bookInput.Publisher,
		PublishedYear: bookInput.PublishedYear,
		Summary:       bookInput.Summary,
		Subject:       bookInput.Subject,
		Language:      bookInput.Language,
		Pages:         bookInput.Pages,
		FileURL:       filePath,
	}

	log.Printf("Book to be created: %+v", book)

	err = configs.DB.Transaction(func(tx *gorm.DB) error {
		if bookInput.ISBN != "" {
			// Check if the new ISBN is unique
			var count int64
			if err := tx.Model(&bookM.Book{}).Where("isbn = ?", bookInput.ISBN).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check ISBN uniqueness: %w", err)
			}
			if count > 0 {
				return fmt.Errorf("ISBN %s is already in use", bookInput.ISBN)
			}
			// Bind the ISBN
			book.ISBN = bookInput.ISBN
		}
		if bookInput.DOI != nil {
			bookInput.DOI = resC.NormalizeDOI(bookInput.DOI)
			if !resC.IsDOIUnique(tx, bookInput.DOI) {
				return c.JSON(http.StatusConflict, map[string]string{
					"message": "A document with this DOI already exists in the repository",
				})
			}
			book.DOI = *bookInput.DOI
		}

		if err := tx.Create(&book).Error; err != nil {
			log.Printf("Error creating book: %v", err)
			if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "Duplicate entry") {
				return c.JSON(http.StatusConflict, map[string]string{
					"message": "A document with this DOI already exists in the repository",
				})
			}
		}

		log.Printf("Book created successfully with ID: %d", book.ID)

		for i, authorName := range bookInput.Authors {
			log.Printf("Processing author %d: %s", i, authorName)

			var userID *uint
			if bookInput.AuthorNIMs[i] != "NON_REGISTERED" {
				var user userM.User
				if err := tx.Where("id_number = ?", bookInput.AuthorNIMs[i]).First(&user).Error; err != nil {
					log.Printf("Error finding user with NIM %s: %v", bookInput.AuthorNIMs[i], err)
					if err == gorm.ErrRecordNotFound {
						return fmt.Errorf("user with NIM/NISN %s not found", bookInput.AuthorNIMs[i])
					}
					return fmt.Errorf("error finding user: %w", err)
				}

				log.Printf("Associated registered author: %s (ID: %d)", user.Name, user.ID)
				userID = &user.ID
			} else {
				log.Printf("Processing non-registered author: %s", authorName)
				// userID remains nil for non-registered authors
			}

			bookAuthor := authorM.BookAuthor{
				BookID:     book.ID,
				UserID:     userID,
				AuthorName: authorName,
			}

			if err := tx.Create(&bookAuthor).Error; err != nil {
				log.Printf("Error creating book author: %v", err)
				return fmt.Errorf("failed to create book author association: %w", err)
			}

			log.Printf("BookAuthor created successfully for: %s", authorName)
		}

		return nil
	})

	if err != nil {
		log.Printf("Transaction error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to create book",
			"error":   err.Error(),
		})
	}

	log.Printf("Book created successfully: %+v", book)
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Book created successfully",
		"book":    book,
	})
}
