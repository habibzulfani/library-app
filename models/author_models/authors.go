package author_models

import (
	"project/models/common"

	"gorm.io/gorm"
)

type PaperAuthor struct {
	common.BaseModel
	PaperID    uint                  `json:"paper_id" gorm:"not null;index:idx_paper_user,unique:false"`
	UserID     *uint                 `json:"user_id" gorm:"index:idx_paper_user,unique:false"`
	AuthorName string                `json:"author_name" gorm:"not null"`
	Paper      common.PaperInterface `json:"paper" gorm:"foreignKey:PaperID"`
	User       common.UserInterface  `json:"user" gorm:"foreignKey:UserID"`
	DeletedAt  gorm.DeletedAt        `json:"deleted_at" gorm:"index"`
}

func (pa *PaperAuthor) GetPaperID() uint {
	return pa.PaperID
}
func (pa *PaperAuthor) GetUserID() *uint {
	return pa.UserID
}
func (pa *PaperAuthor) GetAuthorName() string {
	return pa.AuthorName
}

type BookAuthor struct {
	common.BaseModel
	BookID     uint                 `json:"book_id" gorm:"not null;index:idx_paper_user,unique:false"`
	UserID     *uint                `json:"user_id" gorm:"index:idx_paper_user,unique:false"`
	AuthorName string               `json:"author_name" gorm:"not null"`
	Book       common.BookInterface `json:"book" gorm:"foreignKey:BookID"`
	User       common.UserInterface `json:"user" gorm:"foreignKey:UserID"`
	DeletedAt  gorm.DeletedAt       `json:"deleted_at" gorm:"index"`
}

func (ba *BookAuthor) GetPaperID() uint {
	return ba.BookID
}
func (ba *BookAuthor) GetUserID() *uint {
	return ba.UserID
}
func (ba *BookAuthor) GetAuthorName() string {
	return ba.AuthorName
}

// import (
// 	"project/models/common"
// 	"project/models/paper_models"
// 	"project/models/user_models"
// )

// type PaperAuthor struct {
// 	common.BaseModel
// 	PaperID    uint               `json:"paper_id" gorm:"not null"`
// 	UserID     *uint              `json:"user_id"`
// 	AuthorName string             `json:"author_name" gorm:"not null"`
// 	Paper      paper_models.Paper `json:"paper" gorm:"foreignKey:PaperID"`
// 	User       *user_models.User  `json:"user" gorm:"foreignKey:UserID"`
// }
