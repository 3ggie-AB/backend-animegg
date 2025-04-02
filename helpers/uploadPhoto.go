package helpers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"
)

// Fungsi untuk mengupload foto
func UploadFoto(c *fiber.Ctx, folder, input string) (string, error) {
	// Mengatur folder default jika tidak diberikan
	if folder == "" {
		folder = "uploads"
	}

	if input == "" {
		input = "photo"
	}

	// Mengambil file yang diupload
	file, err := c.FormFile(input)
	if err != nil {
		return "", fmt.Errorf("No file uploaded: %v", err)
	}

	// Validasi ekstensi file (hanya menerima gambar)
	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	extension := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("Invalid file type: %s. Only JPG, JPEG, PNG, and GIF are allowed", extension)
	}

	// Membuat folder jika belum ada
	uploadDir := filepath.Clean("./public/" + folder)
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("Failed to create upload directory: %v", err)
		}
	}

	// Menyimpan file dengan nama unik (UUID)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), extension)
	savePath := filepath.Join(uploadDir, newFileName)

	// Menyimpan file ke folder
	if err := c.SaveFile(file, savePath); err != nil {
		return "", fmt.Errorf("Failed to save file: %v", err)
	}

	// 🔥 Ambil Base URL dari Request
	baseURL := c.BaseURL()

	// 🔥 Kembalikan URL lengkap dengan protocol & domain
	fullURL := fmt.Sprintf("%s/public/%s/%s", baseURL, folder, newFileName)

	return fullURL, nil
}
