package paper_models

type PaperUpdateInput struct {
	Title      *string  `json:"title" form:"paper_title" validate:"omitempty,min=2,max=255"`
	Authors    []string `json:"authors" form:"paper_authors[]" validate:"omitempty,min=1,dive,required"`
	AuthorNIMs []string `json:"author_nims" form:"paper_author_nims[]"`
	Advisor    *string  `json:"advisor" form:"advisor"`
	University *string  `json:"university" form:"university"`
	Department *string  `json:"department" form:"department"`
	Year       *int     `json:"year" form:"year" validate:"omitempty,gt=0,lte=2024"`
	ISSN       *string  `json:"issn" form:"issn" validate:"required,issn"`
	DOI        *string   `json:"doi" form:"doi" validate:"omitempty,doi"`
	Abstract   *string  `json:"abstract" form:"abstract" validate:"omitempty"`
	Keywords   *string  `json:"keywords" form:"keywords" validate:"omitempty"`
}
