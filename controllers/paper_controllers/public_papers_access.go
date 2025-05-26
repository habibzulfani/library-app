package paper_controllers

import (
	"net/http"
	"project/configs"
	resC "project/controllers/resource_controllers"
	combRM "project/models/combined_resource_models"
	paperM "project/models/paper_models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetPapersPublic(c echo.Context) error {
	var papers []paperM.Paper

	err := configs.DB.Preload("Authors").Find(&papers).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	var response []struct {
		ID         uint   `json:"id"`
		Title      string `json:"title"`
		Authors    string `json:"authors"`
		Year       int    `json:"year"`
		ISSN       string `json:"issn"`
		University string `json:"university"`
		Abstract   string `json:"abstract"`
		DOI        string `json:"doi"`
		Type       string `json:"type"`
	}

	type AuthorInfo struct {
		Name string `json:"name"`
		NIM  string `json:"nim,omitempty"`
	}

	for _, paper := range papers {
		authorSet := make(map[string]struct{})
		for _, author := range paper.Authors {
			authorSet[strings.TrimSpace(author.GetAuthorName())] = struct{}{}
		}
		uniqueAuthors := make([]string, 0, len(authorSet))
		for author := range authorSet {
			uniqueAuthors = append(uniqueAuthors, author)
		}

		response = append(response, struct {
			ID         uint   `json:"id"`
			Title      string `json:"title"`
			Authors    string `json:"authors"`
			Year       int    `json:"year"`
			ISSN       string `json:"issn"`
			University string `json:"university"`
			Abstract   string `json:"abstract"`
			DOI        string `json:"doi"`
			Type       string `json:"type"`
		}{
			ID:         paper.ID,
			Title:      paper.Title,
			Authors:    strings.Join(uniqueAuthors, ", "),
			Year:       paper.Year,
			ISSN:       paper.ISSN,
			University: paper.University,
			Abstract:   paper.Abstract,
			DOI:        resC.GetDOIURL(&paper.DOI),
			Type:       "Paper",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
	})
}

func GetPaper(c echo.Context) error {
	inputID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid id",
			"error":   err.Error(),
		})
	}

	var paper paperM.Paper
	err = configs.DB.Preload("Authors").Where("id = ?", inputID).First(&paper).Error
	if err != nil {
		if err.Error() == "record not found" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "Paper not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to retrieve paper data",
			"error":   err.Error(),
		})
	}

	if paper.DOI != "" {
		paper.DOI = resC.GetDOIURL(&paper.DOI)
	}
	// Combine all author names
	authorSet := make(map[string]struct{})
	for _, author := range paper.Authors {
		authorSet[strings.TrimSpace(author.GetAuthorName())] = struct{}{}
	}
	uniqueAuthors := make([]string, 0, len(authorSet))
	for author := range authorSet {
		uniqueAuthors = append(uniqueAuthors, author)
	}
	combinedAuthors := strings.Join(uniqueAuthors, ", ")

	response := combRM.NewUnifiedItemFromPaper(paper, combinedAuthors)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully retrieved paper data",
		"item":    response,
	})
}
