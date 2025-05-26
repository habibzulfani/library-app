package resource_controllers

import (
	bookM "project/models/book_models"
	paperM "project/models/paper_models"
	"strings"

	"gorm.io/gorm"
)

// New function to normalize DOI input
func NormalizeDOI(doi *string) *string {
	if doi == nil {
		return nil
	}
	normalized := strings.TrimPrefix(*doi, "https://doi.org/")
	return &normalized
}

// New function to get DOI URL
func GetDOIURL(doi *string) string {
	if doi == nil {
		return ""
	}
	return "https://doi.org/" + *doi
}

func IsDOIUnique(db *gorm.DB, doi *string) bool {
	if doi == nil {
		return true // Nil DOIs are considered unique (as they're optional)
	}

	var count int64
	db.Model(&paperM.CreatePaperInput{}).Where("doi = ?", doi).Count(&count)
	if count > 0 {
		return false
	}

	db.Model(&bookM.CreateBookInput{}).Where("doi = ?", doi).Count(&count)
	return count == 0
}

