package paper_models

import (
	"project/models/common"

	"gorm.io/gorm"
)

type Paper struct {
	common.BaseModel
	Title      string                   `json:"title" form:"" gorm:"not null" validate:"required,min=2,max=255"`
	Advisor    string                   `json:"advisor"`
	University string                   `json:"university"`
	Department string                   `json:"department"`
	Year       int                      `json:"year" validate:"required,gt=0,lte=2024"`
	ISSN       string                   `json:"issn" gorm:"unique;not null" validate:"required"`
	Abstract   string                   `json:"abstract" gorm:"type:text" validate:"required"`
	Keywords   string                   `json:"keywords" gorm:"type:text" validate:"required"`
	FileURL    string                   `json:"file_url"`
	Authors    []common.AuthorInterface `json:"authors" gorm:"many2many:paper_authors;"`
	DeletedAt  gorm.DeletedAt           `json:"deleted_at" gorm:"index"`
	DOI        string                   `json:"doi" form:"doi" gorm:"unique;index;type:varchar(255)"`
}

func (p *Paper) GetID() uint {
	return p.ID
}

func (p *Paper) GetTitle() string {
    return p.Title
}

func (p *Paper) GetISSN() string {
    return p.ISSN
}

func (p *Paper) GetDOI() string {
    return p.DOI
}

// Author     string                   `json:"author" gorm:"not null" validate:"required,min=2,max=255"`
