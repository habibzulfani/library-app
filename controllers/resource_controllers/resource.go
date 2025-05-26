package resource_controllers

import (
	"net/http"
	"project/configs"
	"project/middlewares"
	bookM "project/models/book_models"
	"project/models/common"
	paperM "project/models/paper_models"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetResources(c echo.Context) error {
	if err := middlewares.JWTChecksRoleAdmin(c); err != nil {
		return err
	}

	var books []bookM.Book
	var papers []paperM.Paper

	// Query books table with pagination
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	limit := 50 // Adjust this value as needed
	offset := (page - 1) * limit

	if err := configs.DB.Preload("Authors").Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		return handleDatabaseError(c, "Error fetching books", err)
	}

	// Query papers table with pagination
	if err := configs.DB.Preload("Authors").Offset(offset).Limit(limit).Find(&papers).Error; err != nil {
		return handleDatabaseError(c, "Error fetching papers", err)
	}

	// Define a combined struct to hold both books and papers
	type CombinedResource struct {
		Type    string   `json:"type"`
		ID      uint     `json:"id"`
		Title   string   `json:"title"`
		Authors []string `json:"authors"`
		DOI     string   `json:"doi"`
	}

	combinedResources := make([]CombinedResource, 0, len(books)+len(papers))

	// Combine books
	for _, book := range books {
		authors := getUniqueAuthors(book.Authors)
		combinedResources = append(combinedResources, CombinedResource{
			Type:    "book",
			ID:      book.ID,
			Title:   book.Title,
			Authors: authors,
			DOI:     GetDOIURL(&book.DOI),
		})
	}

	// Combine papers
	for _, paper := range papers {
		authors := getUniqueAuthors(paper.Authors)
		combinedResources = append(combinedResources, CombinedResource{
			Type:    "paper",
			ID:      paper.ID,
			Title:   paper.Title,
			Authors: authors,
			DOI:     GetDOIURL(&paper.DOI),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    combinedResources,
		"page":    page,
		"limit":   limit,
	})
}

// Helper function to get unique authors
func getUniqueAuthors(authors []common.AuthorInterface) []string {
	authorSet := make(map[string]struct{})
	for _, author := range authors {
		authorSet[author.GetAuthorName()] = struct{}{}
	}
	uniqueAuthors := make([]string, 0, len(authorSet))
	for author := range authorSet {
		uniqueAuthors = append(uniqueAuthors, author)
	}
	sort.Strings(uniqueAuthors)
	return uniqueAuthors
}
