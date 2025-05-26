package resource_controllers

import (
	"log"
	"net/http"
	"project/configs"
	bookM "project/models/book_models"
	resM "project/models/combined_resource_models"
	paperM "project/models/paper_models"
	"sort"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SearchBar(c echo.Context) error {
	query := c.QueryParam("q")
	log.Printf("Received search request with query: %s", query)

	query = *NormalizeDOI(&query)

	var books []bookM.Book
	var papers []paperM.Paper

	// Search in books
	bookQuery := configs.DB.Where("title LIKE ? OR publisher LIKE ? OR isbn LIKE ? OR subject LIKE ? OR language LIKE ? OR doi LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")

	// Handle year search separately
	if year, err := strconv.Atoi(query); err == nil {
		bookQuery = bookQuery.Or("published_year = ?", year)
	}

	if err := bookQuery.Or("id IN (?)",
		configs.DB.Table("book_authors").
			Where("author_name LIKE ?", "%"+query+"%").
			Select("book_id")).
		Preload("Authors").Find(&books).Error; err != nil {
		return handleDatabaseError(c, "Error searching books", err)
	}

	// Search in papers
	paperQuery := configs.DB.Where("title LIKE ? OR university LIKE ? OR department LIKE ? OR issn LIKE ? OR keywords LIKE ? OR doi LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")

	// Handle year search separately
	if year, err := strconv.Atoi(query); err == nil {
		paperQuery = paperQuery.Or("year = ?", year)
	}

	if err := paperQuery.Or("id IN (?)",
		configs.DB.Table("paper_authors").
			Where("author_name LIKE ?", "%"+query+"%").
			Select("paper_id")).
		Preload("Authors").Find(&papers).Error; err != nil {
		return handleDatabaseError(c, "Error searching papers", err)
	}

	bookResults := convertBooksToSearchItems(books)
	paperResults := convertPapersToSearchItems(papers)
	results := append(bookResults, paperResults...)

	response := resM.SearchResponse{
		Results: results,
		Total:   len(results),
	}

	return c.JSON(http.StatusOK, response)
}

func AdvancedSearch(c echo.Context) error {
	itemType := c.QueryParam("type")
	var results []resM.SearchResultItem

	switch itemType {
	case "book":
		books, err := advancedSearchBooks(c)
		if err != nil {
			return handleDatabaseError(c, "Error searching books", err)
		}
		results = convertBooksToSearchItems(books)
	case "paper":
		papers, err := advancedSearchPapers(c)
		if err != nil {
			return handleDatabaseError(c, "Error searching papers", err)
		}
		results = convertPapersToSearchItems(papers)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item type"})
	}

	response := resM.SearchResponse{
		Results: results,
		Total:   len(results),
	}

	return c.JSON(http.StatusOK, response)
}

func advancedSearchBooks(c echo.Context) ([]bookM.Book, error) {
	var books []bookM.Book
	query := configs.DB.Model(&bookM.Book{}).Preload("Authors")

	if title := c.QueryParam("title"); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if year := c.QueryParam("year"); year != "" {
        if yearInt, err := strconv.Atoi(year); err == nil {
            query = query.Where("published_year = ?", yearInt)
        }
    }
	if publisher := c.QueryParam("publisher"); publisher != "" {
		query = query.Where("publisher LIKE ?", "%"+publisher+"%")
	}
	if isbn := c.QueryParam("isbn"); isbn != "" {
		query = query.Where("isbn LIKE ?", "%"+isbn+"%")
	}
	if doi := c.QueryParam("doi"); doi != "" {
		doidata := NormalizeDOI(&doi)
		query = query.Where("doi LIKE ?", "%"+*doidata+"%")
	}
	if author := c.QueryParam("author"); author != "" {
		query = query.Where("id IN (?)",
			configs.DB.Table("book_authors").
				Where("author_name LIKE ?", "%"+author+"%").
				Select("book_id"),
		)
	}

	if err := query.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func advancedSearchPapers(c echo.Context) ([]paperM.Paper, error) {
	var papers []paperM.Paper
	query := configs.DB.Model(&paperM.Paper{}).Preload("Authors")

	if title := c.QueryParam("title"); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if year := c.QueryParam("year"); year != "" {
        if yearInt, err := strconv.Atoi(year); err == nil {
            query = query.Where("year = ?", yearInt)
        }
    }
	if issn := c.QueryParam("issn"); issn != "" {
		query = query.Where("issn LIKE ?", "%"+issn+"%")
	}
	if university := c.QueryParam("university"); university != "" {
		query = query.Where("university LIKE ?", "%"+university+"%")
	}
	if doi := c.QueryParam("doi"); doi != "" {
		doidata := NormalizeDOI(&doi)
		query = query.Where("doi LIKE ?", "%"+*doidata+"%")
	}
	if author := c.QueryParam("author"); author != "" {
		query = query.Where("id IN (?)",
			configs.DB.Table("paper_authors").
				Where("author_name LIKE ?", "%"+author+"%").
				Select("paper_id"),
		)
	}

	if err := query.Find(&papers).Error; err != nil {
		return nil, err
	}
	return papers, nil
}

func convertBooksToSearchItems(books []bookM.Book) []resM.SearchResultItem {
	var results []resM.SearchResultItem
	for _, book := range books {
		authorSet := make(map[string]struct{})
		for _, author := range book.Authors {
			authorSet[strings.TrimSpace(author.GetAuthorName())] = struct{}{}
		}
		uniqueAuthors := make([]string, 0, len(authorSet))
		for author := range authorSet {
			uniqueAuthors = append(uniqueAuthors, author)
		}
		sort.Strings(uniqueAuthors)
		results = append(results, resM.SearchResultItem{
			ID:          book.ID,
			Title:       book.Title,
			Authors:     strings.Join(uniqueAuthors, "; "),
			Identifier:  book.ISBN,
			Year:        book.PublishedYear,
			Institution: book.Publisher,
			DOI:         GetDOIURL(&book.DOI),
			Type:        "Book",
		})
	}
	return results
}

func convertPapersToSearchItems(papers []paperM.Paper) []resM.SearchResultItem {
	var results []resM.SearchResultItem
	for _, paper := range papers {
		authorSet := make(map[string]struct{})
		for _, author := range paper.Authors {
			authorSet[strings.TrimSpace(author.GetAuthorName())] = struct{}{}
		}
		uniqueAuthors := make([]string, 0, len(authorSet))
		for author := range authorSet {
			uniqueAuthors = append(uniqueAuthors, author)
		}
		sort.Strings(uniqueAuthors)
		results = append(results, resM.SearchResultItem{
			ID:          paper.ID,
			Title:       paper.Title,
			Authors:     strings.Join(uniqueAuthors, "; "),
			Identifier:  paper.ISSN,
			Year:        paper.Year,
			Institution: paper.University,
			DOI:         GetDOIURL(&paper.DOI),
			Type:        "Paper",
		})
	}
	return results
}

func handleDatabaseError(c echo.Context, message string, err error) error {
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "No results found",
		})
	}
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"message": message,
		"error":   err.Error(),
	})
}
