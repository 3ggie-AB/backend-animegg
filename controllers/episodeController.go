package controllers

import (
	"fmt"
	"os"
	"strconv" // Pastikan mengimpor paket ini
	"github.com/3ggie-AB/backend-animegg/config"
	"github.com/3ggie-AB/backend-animegg/models"
	"github.com/3ggie-AB/backend-animegg/helpers"
	"github.com/gofiber/fiber/v2"
)

// Create Episode
func CreateEpisode(c *fiber.Ctx) error {
	// Ambil data episode dari request form-data
	animeIDStr := c.FormValue("anime_id") // Ambil anime_id dalam bentuk string
	episodeStr := c.FormValue("episode")  // Ambil episode dalam bentuk string
	driver := c.FormValue("driver")  // Ambil driver dalam bentuk string

	// Konversi anime_id ke uint
	animeID, err := strconv.Atoi(animeIDStr) // Gunakan Atoi untuk konversi string ke int
	if err != nil || animeID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid anime_id. Please provide a valid anime ID."})
	}

	// Konversi episode ke int
	episode, err := strconv.Atoi(episodeStr)
	if err != nil || episode <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid episode number. Please provide a valid episode number."})
	}

	// Cek apakah anime_id ada di tabel animes
	var episodeData models.Episode
	if err := config.DB.First(&episodeData, animeID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Anime not found for the given anime_id"})
	}

	// Ambil file video dari request form-data
	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Video file is required"})
	}

	// Simpan file sementara ke folder `temp/`
	uploadDir := "./temp"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create temp directory"})
		}
	}

	// Simpan sementara ke folder `temp/`
	tempFilePath := fmt.Sprintf("./temp/%s", file.Filename)
	if err := c.SaveFile(file, tempFilePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save video file"})
	}

	// Upload video ke driver yang dipilih
	var driveURL string
	if driver == "" || driver == "gdrive" {
		// Upload ke Google Drive
		var err error
		driveURL, err = helpers.UploadVideoToDrive(c, tempFilePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// Hapus file sementara setelah di-upload
	_ = os.Remove(tempFilePath)

	// Simpan URL video dari Google Drive
	newEpisode := models.Episode{
		AnimeID: uint(animeID),
		Episode: episode,
		Video:   driveURL,
	}

	// Simpan episode ke database
	if err := config.DB.Create(&newEpisode).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create episode"})
	}

	return c.JSON(fiber.Map{
		"message": "Episode berhasil dibuat!",
		"data":    newEpisode,
	})
}


// Get Episodes by Anime ID
func GetEpisodesByAnime(c *fiber.Ctx) error {
	animeID := c.Params("anime_id")
	var episodes []models.Episode
	config.DB.Where("anime_id = ?", animeID).Find(&episodes)
	return c.JSON(episodes)
}

func GetAnimeEpisode(c *fiber.Ctx) error {
	animeID := c.Params("id")
	episode := c.Params("episode")
	var episodeData models.Episode
	if err := config.DB.Where("anime_id = ? AND episode = ?", animeID, episode).First(&episodeData).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Episode not found"})
	}
	return c.JSON(episodeData)
}