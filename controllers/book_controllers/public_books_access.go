package book_controllers

import (
	"net/http"
	"project/configs"
	"project/middlewares"
	resC "project/controllers/resource_controllers"
	bookM "project/models/book_models"
	"project/models/combined_resource_models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetBooksPublic(c echo.Context) error {
	var books []bookM.Book

	err := configs.DB.Preload("Authors").Find(&books).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	var response []struct {
		ID            uint   `json:"id"`
		Title         string `json:"title"`
		Author        string `json:"author"`
		Publisher     string `json:"publisher"`
		PublishedYear int    `json:"published_year"`
		ISBN          string `json:"isbn"`
		Summary       string `json:"summary"`
		DOI           string `json:"doi"`
	}

	for _, book := range books {
		var authorNames []string
		if book.Author != "" {
			authorNames = append(authorNames, book.Author)
		}
		for _, author := range book.Authors {
			authorNames = append(authorNames, author.GetName())
		}

		response = append(response, struct {
			ID            uint   `json:"id"`
			Title         string `json:"title"`
			Author        string `json:"author"`
			Publisher     string `json:"publisher"`
			PublishedYear int    `json:"published_year"`
			ISBN          string `json:"isbn"`
			Summary       string `json:"summary"`
			DOI           string `json:"doi"`
		}{
			ID:            book.ID,
			Title:         book.Title,
			Author:        strings.Join(authorNames, ", "),
			Publisher:     book.Publisher,
			PublishedYear: book.PublishedYear,
			ISBN:          book.ISBN,
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
		if err.Error() == "record not found" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "Book not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to retrieve book data",
			"error":   err.Error(),
		})
	}

	// Log the download activity
	userID := uint(inputID) // You'll need to implement this function
	err = middlewares.LogDownloadActivity(configs.DB, userID, uint(inputID))
	if err != nil {
		// Log the error, but don't fail the request
		println("Failed to log download activity:", err.Error())
	}

	// Increment the download counter
	err = middlewares.IncrementDownloadCounter(configs.DB)
	if err != nil {
		// Log the error, but don't fail the request
		println("Failed to increment download counter:", err.Error())
	}

	if book.DOI != "" {
		book.DOI = resC.GetDOIURL(&book.DOI)
	}
	// Combine all author names
	var authorNames []string
	if book.Author != "" {
		authorNames = append(authorNames, book.Author)
	}
	for _, author := range book.Authors {
		authorNames = append(authorNames, author.GetName())
	}
	combinedAuthors := strings.Join(authorNames, ", ")

	response := combined_resource_models.NewUnifiedItemFromBook(book, combinedAuthors)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully retrieved book data",
		"item":    response,
	})
}
