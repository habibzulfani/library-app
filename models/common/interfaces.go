package common

import "time"

type BaseModel struct {
    ID        uint       `json:"id" gorm:"primaryKey;autoIncrement"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type PaperInterface interface {
    GetID() uint
    GetTitle() string
    GetISSN() string
    GetDOI() string
    // Add other necessary methods
}
type BookInterface interface {
    GetID() uint
    GetTitle() string
    GetISBN() string
    GetDOI() string
    // Add other necessary methods
}

type UserInterface interface {
    GetID() uint
    GetName() string
    // Add other necessary methods
}

type AuthorInterface interface {
    GetPaperID() uint
    GetBookID() uint
    GetUserID() *uint
    GetAuthorName() string
}