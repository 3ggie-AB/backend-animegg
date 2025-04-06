package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
    "github.com/3ggie-AB/backend-animegg/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Ambil konfigurasi dari .env
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

	// Koneksi ke database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	fmt.Println("Database berhasil terkoneksi!")
	DB = db
	DB.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Notification{},
		&models.Process{},
		&models.Anime{},
		&models.Episode{},
		&models.Video{},
		&models.Genre{},
		&models.Studio{},
		&models.Season{},
		&models.AnimeGenre{},
	)
	
}
