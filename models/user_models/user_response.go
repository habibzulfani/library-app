package user_models

import (
	"time"

	"gorm.io/gorm"
)

type UserLoginResponse struct {
	Id        int    `json:"id" form:"id"`
	Email     string `json:"email" form:"email"`
	Name      string `json:"name" form:"name"`
	UserType  string `json:"user_type" form:"user_type"`
	IDNumber  string `json:"id_number" form:"id_number"`
	Fakultas  string `json:"fakultas" form:"fakultas"`
	Jurusan   string `json:"jurusan" form:"jurusan"`
	Address   string `json:"address" form:"address"`
	Role      string `json:"role" form:"role"`
	Token     string `json:"token" form:"token"`
}

type UserRequestResponse struct {
	ID        int            `json:"id" form:"id"`
	Email     string         `json:"email" form:"email"`
	Name      string         `json:"name" form:"name"`
	UserType  string         `json:"user_type" form:"user_type"`
	IDNumber  string         `json:"id_number" form:"id_number"`
	Fakultas  string         `json:"fakultas" form:"fakultas"`
	Jurusan   string         `json:"jurusan" form:"jurusan"`
	Address   string         `json:"address" form:"address"`
	Role      string         `json:"role" form:"role"`
	CreatedAt time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type UserUpdateResponse struct {
	ID        int            `json:"id" form:"id"`
	Email     string         `json:"email" form:"email"`
	Name      string         `json:"name" form:"name"`
	UserType  string         `json:"user_type" form:"user_type"`
	IDNumber  string         `json:"id_number" form:"id_number"`
	Fakultas  string         `json:"fakultas" form:"fakultas"`
	Jurusan   string         `json:"jurusan" form:"jurusan"`
	Address   string         `json:"address" form:"address"`
	CreatedAt time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
