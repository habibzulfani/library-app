package user_models

import (
	"project/models/common"
	univM "project/models/univ_models"

	"gorm.io/gorm"
)

type User struct {
	common.BaseModel
	Email      string                   `json:"email" form:"email" gorm:"column:email;unique;not null" validate:"required,email"`
	Name       string                   `json:"name" form:"name" gorm:"column:name;type:longtext;not null" validate:"required,min=2,max=100"`
	Password   string                   `json:"-" form:"password" gorm:"column:password;type:longtext;not null" validate:"required,min=8"`
	Role       string                   `json:"role" form:"role" gorm:"column:role;type:longtext;not null" validate:"required,oneof=user admin"`
	Jurusan    string                   `json:"jurusan" form:"jurusan" gorm:"column:jurusan;type:longtext;not null"`
	Address    string                   `json:"address" form:"address" gorm:"column:address;type:text"`
	UserType   string                   `json:"user_type" form:"user_type" gorm:"column:user_type;type:longtext;not null" validate:"required,oneof=student teacher"`
	IDNumber   string                   `json:"id_number" form:"id_number" gorm:"column:id_number;unique;not null" validate:"required"`
	JurusanID  uint                     `json:"jurusan_id" form:"jurusan_id" gorm:"column:jurusan_id;not null"`
	JurusanRef univM.Jurusan            `json:"jurusan_ref" gorm:"foreignKey:JurusanID"`
	Papers     []common.AuthorInterface `json:"papers" gorm:"many2many:paper_authors;"`
	Books      []common.AuthorInterface `json:"books" gorm:"many2many:book_authors;"`
	DeletedAt  gorm.DeletedAt           `json:"deleted_at" gorm:"index"`
}

func (u *User) GetID() uint {
	return u.ID
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetNIM() string {
	return u.IDNumber
}
