package controllers

import (
	"github.com/3ggie-AB/backend-animegg/config"
	"github.com/3ggie-AB/backend-animegg/helpers"
	"github.com/3ggie-AB/backend-animegg/models"
	"github.com/gofiber/fiber/v2"
)

// Create Anime
func CreateAnime(c *fiber.Ctx) error {
	// Ambil input dari form-data
	title := c.FormValue("title")
	tags := c.FormValue("tags")
	description := c.FormValue("description")

	// Validasi input manual
	if title == "" || tags == "" || description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Semua field harus diisi"})
	}

	// Upload foto
	photoURL, err := helpers.UploadFoto(c, "anime", "photo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Pemetaan data yang sudah divalidasi ke model Anime
	anime := models.Anime{
		Title:       title,
		Photo:       photoURL,
		Tags:        tags,
		Description: description,
	}

	// Simpan ke database
	if err := config.DB.Create(&anime).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create anime"})
	}

	return c.JSON(fiber.Map{
		"message": "Anime Berhasil dibuat!",
		"data":   anime,
	})
}

// Get All Anime
func GetAnimes(c *fiber.Ctx) error {
	var animes []models.Anime
	config.DB.Preload("Episodes").Find(&animes)
	return c.JSON(animes)
}

// Get Single Anime by ID
func GetAnime(c *fiber.Ctx) error {
	id := c.Params("id")
	var anime models.Anime

	if err := config.DB.Preload("Episodes").First(&anime, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Anime not found"})
	}

	return c.JSON(anime)
}
