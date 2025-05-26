package book_models

import (
	"project/models/common"

	"gorm.io/gorm"
)

type Book struct {
	common.BaseModel
	Title         string                   `json:"title" gorm:"not null" validate:"required,min=2,max=255"`
	Publisher     string                   `json:"publisher" gorm:"null" validate:"required,min=2,max=255"`
	PublishedYear int                      `json:"published_year" gorm:"null" validate:"required,gt=0,lte=2024"`
	ISBN          string                   `json:"isbn" gorm:"unique;not null" validate:"required"`
	Summary       string                   `json:"summary" gorm:"null" gorm:"type:text" validate:"required"`
	Subject       string                   `json:"subject" gorm:"null"`
	Language      string                   `json:"language" gorm:"null"`
	Pages         int                      `json:"pages" gorm:"null"`
	FileURL       string                   `json:"file_url"`
	Authors       []common.AuthorInterface `json:"authors" gorm:"many2many:book_authors;"`
	DeletedAt     gorm.DeletedAt           `json:"deleted_at" gorm:"index"`
	DOI           string                   `json:"doi" form:"doi" gorm:"unique;index;type:varchar(255)"`
}

func (b *Book) GetID() uint {
	return b.ID
}

func (b *Book) GetTitle() string {
	return b.Title
}

func (b *Book) GetISBN() string {
	return b.ISBN
}

func (b *Book) GetDOI() string {
	return b.DOI
}

// Author        string                   `json:"author" gorm:"not null" validate:"required,min=2,max=255"`
