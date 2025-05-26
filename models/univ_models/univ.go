package univ_models

import (
	"gorm.io/gorm"
)

type Fakultas struct {
	gorm.Model
	Name     string    `json:"name" gorm:"unique;not null"`
	Jurusans []Jurusan `json:"jurusans" gorm:"foreignKey:FakultasID"`
}

type Jurusan struct {
	gorm.Model
	Name       string   `json:"name" gorm:"not null"`
	FakultasID uint     `json:"fakultas_id" gorm:"not null"`
	Fakultas   Fakultas `json:"fakultas" gorm:"foreignKey:FakultasID"`
}
