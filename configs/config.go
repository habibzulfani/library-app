package configs

import (
	"fmt"
	"log"
	"os"

	authorM "project/models/author_models"
	bookM "project/models/book_models"
	logM "project/models/log_models"
	paperM "project/models/paper_models"
	univM "project/models/univ_models"
	userM "project/models/user_models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	i := godotenv.Load("../.env")
	if i != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	initMigrate()
	seedFakultasAndJurusan()
}

func initMigrate() {
	DB.AutoMigrate(
		&userM.User{},
		&bookM.Book{},
		&paperM.Paper{},
		&authorM.BookAuthor{},
		&authorM.PaperAuthor{},
		&logM.ActivityLog{},
		&logM.Counter{},
		&univM.Fakultas{},
		&univM.Jurusan{},
	)
	if err := MigrateBooks(DB); err != nil {
		log.Fatalf("Failed to migrate books table: %v", err)
	}

	DB.Debug().AutoMigrate(&paperM.Paper{}, &bookM.Book{})
}

func seedFakultasAndJurusan() {
	var count int64
	DB.Model(&univM.Fakultas{}).Count(&count)
	if count > 0 {
		log.Println("Fakultas data already exists, skipping seeding")
		return
	}

	fakultases := []univM.Fakultas{
		{Name: "Ekonomi", Jurusans: []univM.Jurusan{{Name: "Manajemen"}}},
		{Name: "Hukum", Jurusans: []univM.Jurusan{{Name: "Hukum"}}},
		{Name: "Ilmu Komputer", Jurusans: []univM.Jurusan{
			{Name: "Manajemen Informatika"},
			{Name: "Rekayasa Perangkat Lunak"},
			{Name: "Sistem Informasi"},
			{Name: "Teknik Informatika"},
		}},
	}

	for _, f := range fakultases {
		result := DB.Create(&f)
		if result.Error != nil {
			log.Printf("Error seeding Fakultas %s: %v", f.Name, result.Error)
		} else {
			log.Printf("Seeded Fakultas: %s", f.Name)
		}
	}
}
