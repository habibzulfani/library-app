package registered_papers_access

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
	paperM "project/models/paper_models"
	userM "project/models/user_models"
	"project/validators"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreatePaper(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	var paperInput paperM.CreatePaperInput
	if err := c.Bind(&paperInput); err != nil {
		log.Printf("Error binding input: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	log.Printf("Received paper input: %+v", paperInput)

	if err := validators.ValidateCreatePaper(&paperInput); err != nil {
		log.Printf("Validation error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	file, err := c.FormFile("paperFile")
	if err != nil {
		log.Printf("Error getting paper file: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Paper file is required",
			"error":   err.Error(),
		})
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	uploadDir := "./uploads/papers"
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

	paper := paperM.Paper{
		Title:      paperInput.Title,
		Advisor:    paperInput.Advisor,
		University: paperInput.University,
		Department: paperInput.Department,
		Year:       paperInput.Year,
		Abstract:   paperInput.Abstract,
		Keywords:   paperInput.Keywords,
		FileURL:    filePath,
	}

	log.Printf("Paper to be created: %+v", paper)

	err = configs.DB.Transaction(func(tx *gorm.DB) error {
		if paperInput.ISSN != "" {
			// Check if the new ISSN is unique
			var count int64
			if err := tx.Model(&paperM.Paper{}).Where("issn = ?", paperInput.ISSN).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check ISSN uniqueness: %w", err)
			}
			if count > 0 {
				return fmt.Errorf("ISSN %s is already in use", paperInput.ISSN)
			}
			// Bind the ISSN
			paper.ISSN = paperInput.ISSN
		}
		if paperInput.DOI != nil {
			paperInput.DOI = resC.NormalizeDOI(paperInput.DOI)
			if !resC.IsDOIUnique(tx, paperInput.DOI) {
				return c.JSON(http.StatusConflict, map[string]string{
					"message": "A document with this DOI already exists in the repository",
				})
			}
			paper.DOI = *paperInput.DOI
		}
		if err := tx.Create(&paper).Error; err != nil {
			log.Printf("Error creating paper: %v", err)
			if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "Duplicate entry") {
				return c.JSON(http.StatusConflict, map[string]string{
					"message": "A document with this DOI already exists in the repository",
				})
			}		
		}

		log.Printf("Paper created successfully with ID: %d", paper.ID)

		for i, authorName := range paperInput.Authors {
			log.Printf("Processing author %d: %s", i, authorName)

			var userID *uint
			if paperInput.AuthorNIMs[i] != "NON_REGISTERED" {
				var user userM.User
				if err := tx.Where("id_number = ?", paperInput.AuthorNIMs[i]).First(&user).Error; err != nil {
					log.Printf("Error finding user with NIM %s: %v", paperInput.AuthorNIMs[i], err)
					if err == gorm.ErrRecordNotFound {
						return fmt.Errorf("user with NIM/NISN %s not found", paperInput.AuthorNIMs[i])
					}
					return fmt.Errorf("error finding user: %w", err)
				}

				log.Printf("Associated registered author: %s (ID: %d)", user.Name, user.ID)
				userID = &user.ID
				authorName = user.Name
			} else {
				log.Printf("Processing non-registered author: %s", authorName)
				// userID remains nil for non-registered authors
			}

			paperAuthor := authorM.PaperAuthor{
				PaperID:    paper.ID,
				UserID:     userID,
				AuthorName: authorName,
			}

			if err := tx.Create(&paperAuthor).Error; err != nil {
				log.Printf("Error creating paper author: %v", err)
				return fmt.Errorf("failed to create paper author association: %w", err)
			}

			log.Printf("PaperAuthor created successfully for: %s", authorName)
		}

		return nil
	})

	if err != nil {
		log.Printf("Transaction error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to create paper",
			"error":   err.Error(),
		})
	}

	log.Printf("Paper created successfully: %+v", paper)
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Paper created successfully",
		"paper":   paper,
	})
}
