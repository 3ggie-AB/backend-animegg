package controllers

import (
	"github.com/3ggie-AB/backend-animegg/config"
	"github.com/3ggie-AB/backend-animegg/helpers"
	"github.com/3ggie-AB/backend-animegg/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// Create Anime
func CreateAnime(c *fiber.Ctx) error {
	// Ambil input dari form-data
	title := c.FormValue("title")
	seasonIDStr := c.FormValue("season_id")
	studioIDStr := c.FormValue("studio_id")
	enTitle := c.FormValue("en_title")
	status := c.FormValue("status")
	description := c.FormValue("description")
	isHidden := c.FormValue("is_hidden") // false or true, depends on form data

	// Validasi input manual
	if title == "" || seasonIDStr == "" || studioIDStr == "" || description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Semua field harus diisi"})
	}

	// Convert season_id and studio_id to uint
	seasonID, err := strconv.ParseUint(seasonIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid season_id"})
	}
	studioID, err := strconv.ParseUint(studioIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid studio_id"})
	}

	// Upload photo
	photoURL, err := helpers.UploadFoto(c, "anime", "photo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Pemetaan data yang sudah divalidasi ke model Anime
	anime := models.Anime{
		SeasonID:    uint(seasonID),
		StudioID:    uint(studioID),
		Title:       title,
		EnTitle:     enTitle,
		Status:      status,
		IsHidden:    isHidden == "true", // Convert "true"/"false" to boolean
		Description: description,
		Photo:       photoURL,
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
	return c.JSON(fiber.Map{"data": animes, "message": "Berhasil Ambil data anime"})
}

// Get Single Anime by ID
func GetAnime(c *fiber.Ctx) error {
	id := c.Params("id")
	var anime models.Anime

	if err := config.DB.Preload("Episodes").First(&anime, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Anime not found"})
	}

	return c.JSON(
		fiber.Map{"data": anime, "message": "Berhasil Ambil data anime"},
	)
}
