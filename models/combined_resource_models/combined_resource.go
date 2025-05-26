package combined_resource_models

import (
	"project/models/book_models"
	"project/models/paper_models"
)

var CombinedResources []struct {
	Type   string `json:"type"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type UnifiedItem struct {
	Type            string `json:"type"`
	ID              uint   `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	CombinedAuthors string `json:"combined_authors"`
	Year            int    `json:"year"`
	Abstract        string `json:"abstract,omitempty"`
	Summary         string `json:"summary,omitempty"`
	DOI             string `json:"doi,omitempty"`
	FileURL         string `json:"file_url"`

	// Paper-specific fields
	Advisor    string `json:"advisor,omitempty"`
	University string `json:"university,omitempty"`
	Department string `json:"department,omitempty"`
	ISSN       string `json:"issn,omitempty"`
	Keywords   string `json:"keywords,omitempty"`

	// Book-specific fields
	Publisher string `json:"publisher,omitempty"`
	ISBN      string `json:"isbn,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Language  string `json:"language,omitempty"`
	Pages     int    `json:"pages,omitempty"`
}

func NewUnifiedItemFromPaper(paper paper_models.Paper, combinedAuthors string) UnifiedItem {
	return UnifiedItem{
		Type:            "paper",
		ID:              paper.ID,
		Title:           paper.Title,
		CombinedAuthors: combinedAuthors,
		Year:            paper.Year,
		Abstract:        paper.Abstract,
		FileURL:         paper.FileURL,
		Advisor:         paper.Advisor,
		University:      paper.University,
		Department:      paper.Department,
		ISSN:            paper.ISSN,
		DOI:             paper.DOI,
		Keywords:        paper.Keywords,
	}
}

func NewUnifiedItemFromBook(book book_models.Book, combinedAuthors string) UnifiedItem {
	return UnifiedItem{
		Type:            "book",
		ID:              book.ID,
		Title:           book.Title,
		CombinedAuthors: combinedAuthors,
		Year:            book.PublishedYear,
		Summary:         book.Summary,
		FileURL:         book.FileURL,
		Publisher:       book.Publisher,
		ISBN:            book.ISBN,
		DOI:             book.DOI,
		Subject:         book.Subject,
		Language:        book.Language,
		Pages:           book.Pages,
	}
}
