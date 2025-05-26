package book_models

type CreateBookInput struct {
	Title         string   `json:"title" form:"book_title" validate:"required,min=2,max=255"`
	Authors       []string `json:"authors" form:"book_authors[]" validate:"required,min=1,dive,required"`
	AuthorNIMs    []string `json:"author_nims" form:"book_author_nims[]"`
	Publisher     string   `json:"publisher" form:"publisher" validate:"required,min=2,max=255"`
	PublishedYear int      `json:"published_year" form:"published_year" validate:"required,gt=0,lte=2024"`
	ISBN          string   `json:"isbn" form:"isbn" validate:"required,isbn"`
	DOI           *string  `json:"doi" form:"doi" validate:"omitempty,doi"`
	Summary       string   `json:"summary" form:"summary" validate:"required"`
	Subject       string   `json:"subject" form:"subject" validate:"required"`
	Language      string   `json:"language" form:"language"`
	Pages         int      `json:"pages" form:"pages"`
}
