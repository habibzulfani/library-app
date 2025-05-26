package registered_books_access

import (
	"net/http"
	"project/configs"
	resC "project/controllers/resource_controllers"
	bookM "project/models/book_models"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetBooksRegistered(c echo.Context) error {
	var books []bookM.Book

	err := configs.DB.Preload("Authors").Find(&books).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	type BookResponse struct {
		ID            uint     `json:"id"`
		Title         string   `json:"title"`
		Authors       []string `json:"author"`
		Publisher     string   `json:"publisher"`
		PublishedYear int      `json:"published_year"`
		ISBN          string   `json:"isbn"`
		Summary       string   `json:"summary"`
		DOI           string   `json:"doi"`
	}

	var response []BookResponse

	for _, book := range books {
		var authorNames []string
		if book.Author != "" {
			authorNames = append(authorNames, book.Author)
		}
		for _, author := range book.Authors {
			authorNames = append(authorNames, author.GetName())
		}

		response = append(response, BookResponse{
			ID:            book.ID,
			Title:         book.Title,
			Publisher:     book.Publisher,
			PublishedYear: book.PublishedYear,
			ISBN:          book.ISBN,
			Authors:       authorNames,
			Summary:       book.Summary,
			DOI:           resC.GetDOIURL(&book.DOI),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
	})

}

func GetBook(c echo.Context) error {
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

	// Create a response structure that includes author details
	type AuthorDetail struct {
		Name string `json:"name"`
		NIM  string `json:"nim,omitempty"`
	}

	type BookResponse struct {
		ID            uint           `json:"id"`
		Title         string         `json:"title"`
		Publisher     string         `json:"publisher"`
		PublishedYear int            `json:"published_year"`
		ISBN          string         `json:"isbn"`
		Summary       string         `json:"summary"`
		Subject       string         `json:"subject"`
		Language      string         `json:"language"`
		Pages         int            `json:"pages"`
		DOI           string         `json:"doi"`
		FileURL       string         `json:"file_url"`
		AuthorDetails []AuthorDetail `json:"author_details"`
	}

	response := BookResponse{
		ID:            book.ID,
		Title:         book.Title,
		Publisher:     book.Publisher,
		PublishedYear: book.PublishedYear,
		ISBN:          book.ISBN,
		Summary:       book.Summary,
		Subject:       book.Subject,
		Language:      book.Language,
		Pages:         book.Pages,
		DOI:           resC.GetDOIURL(&book.DOI),
		FileURL:       book.FileURL,
	}

	// Handle registered authors
	for _, author := range book.Authors {
		response.AuthorDetails = append(response.AuthorDetails, AuthorDetail{
			Name: author.GetName(),
			NIM:  author.GetNIM(),
		})
	}

	// Handle unregistered author if no registered authors are found
	if len(response.AuthorDetails) == 0 && book.Author != "" {
		response.AuthorDetails = append(response.AuthorDetails, AuthorDetail{
			Name: book.Author,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully retrieved book data",
		"book":    response,
	})
}
