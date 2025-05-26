package registered_papers_access

import (
	"net/http"
	"project/configs"
	resC "project/controllers/resource_controllers"
	paperM "project/models/paper_models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetPapersRegistered(c echo.Context) error {
	var papers []paperM.Paper

	err := configs.DB.Preload("Authors").Find(&papers).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	type paperResponse struct {
		ID         uint     `json:"id"`
		Title      string   `json:"title"`
		Authors    []string `json:"authors"`
		Year       int      `json:"year"`
		ISSN       string   `json:"issn"`
		University string   `json:"university"`
		Abstract   string   `json:"abstract"`
		DOI        string   `json:"doi"`
	}

	var response []paperResponse

	for _, paper := range papers {
		authorSet := make(map[string]struct{})
		for _, author := range paper.Authors {
			authorSet[strings.TrimSpace(author.GetAuthorName())] = struct{}{}
		}
		uniqueAuthors := make([]string, 0, len(authorSet))
		for author := range authorSet {
			uniqueAuthors = append(uniqueAuthors, author)
		}

		response = append(response, paperResponse{
			ID:         paper.ID,
			Title:      paper.Title,
			ISSN:       paper.ISSN,
			Year:       paper.Year,
			University: paper.University,
			Abstract:   paper.Abstract,
			Authors:    uniqueAuthors,
			DOI:        resC.GetDOIURL(&paper.DOI),
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
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "Paper not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to retrieve paper data",
			"error":   err.Error(),
		})
	}

	// Create a response structure that includes author details
	type AuthorDetail struct {
		Name string `json:"name"`
		NIM  string `json:"nim"`
	}

	type PaperResponse struct {
		ID            uint           `json:"id"`
		Title         string         `json:"title"`
		Advisor       string         `json:"advisor"`
		University    string         `json:"university"`
		Department    string         `json:"department"`
		Year          int            `json:"year"`
		ISSN          string         `json:"issn"`
		DOI           string         `json:"doi"`
		Abstract      string         `json:"abstract"`
		Keywords      string         `json:"keywords"`
		FileURL       string         `json:"file_url"`
		AuthorDetails []AuthorDetail `json:"author_details"`
	}

	response := PaperResponse{
		ID:         paper.ID,
		Title:      paper.Title,
		Advisor:    paper.Advisor,
		University: paper.University,
		Department: paper.Department,
		Year:       paper.Year,
		ISSN:       paper.ISSN,
		DOI:        resC.GetDOIURL(&paper.DOI),	
		Abstract:   paper.Abstract,
		Keywords:   paper.Keywords,
		FileURL:    paper.FileURL,
	}

	// Handle registered authors
	for _, author := range paper.Authors {
		response.AuthorDetails = append(response.AuthorDetails, AuthorDetail{
			Name: author.GetAuthorName(),
		})
	}

		// Handle unregistered author if no registered authors are found
		if len(response.AuthorDetails) == 0 && paper.Author != "" {
			// Split the Author string in case it contains multiple names
			authorNames := strings.Split(paper.Author, ",")
			for _, name := range authorNames {
				response.AuthorDetails = append(response.AuthorDetails, AuthorDetail{
					Name: strings.TrimSpace(name),
					NIM:  "", // Unregistered authors don't have NIMs
				})
			}
		}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully retrieved paper data",
		"paper":   response,
	})
}
