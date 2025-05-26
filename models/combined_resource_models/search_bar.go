package combined_resource_models

import (
	bookM "project/models/book_models"
	paparM "project/models/paper_models"
)

type SearchResult struct {
	Books  []bookM.Book   `json:"books"`
	Papers []paparM.Paper `json:"papers"`
}

type SearchResultItem struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Authors     string `json:"authors"`    // Comma-separated string of author names
	Identifier  string `json:"identifier"` // ISBN or ISSN
	Year        int    `json:"year"`
	Institution string `json:"institution"` // Publisher for books, University for papers
	DOI         string `json:"doi"`
	Type        string `json:"type"` // "Book" or "Paper"
}

type SearchResponse struct {
	Results []SearchResultItem `json:"results"`
	Total   int                `json:"total"`
}
