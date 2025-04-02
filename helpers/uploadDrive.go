package helpers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Fungsi untuk mengupload video ke Google Drive// Fungsi untuk mengupload video ke Google Drive
func UploadVideoToDrive(c *fiber.Ctx, tempFilePath string) (string, error) {
    // Validasi ekstensi file (hanya menerima video)
    allowedExtensions := map[string]bool{".mp4": true, ".avi": true, ".mov": true, ".mkv": true}
    extension := strings.ToLower(filepath.Ext(tempFilePath))
    if !allowedExtensions[extension] {
        return "", fmt.Errorf("Invalid file type: %s. Only MP4, AVI, MOV, MKV are allowed", extension)
    }

    // 🔥 Upload ke Google Drive
    driveURL, err := uploadToGoogleDrive(tempFilePath, filepath.Base(tempFilePath))
    if err != nil {
        return "", err
    }

    return driveURL, nil
}


// Fungsi untuk mengupload file ke Google Drive
func uploadToGoogleDrive(filePath, originalFileName string) (string, error) {
    // Load kredensial Google Drive
    credFile := "animegg-drive.json"
    b, err := os.ReadFile(credFile)
    if err != nil {
        return "", fmt.Errorf("error reading service account file: %v", err)
    }

    // Autentikasi ke Google Drive API
    config, err := google.JWTConfigFromJSON(b, drive.DriveFileScope)
    if err != nil {
        return "", fmt.Errorf("error creating JWT config: %v", err)
    }

    client := config.Client(context.Background())
    srv, err := drive.New(client)
    if err != nil {
        return "", fmt.Errorf("error creating Drive client: %v", err)
    }

    // Buka file
    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("error opening file: %v", err)
    }
    defer file.Close()

    // Set metadata file di Google Drive
    driveFile := &drive.File{
        Name:    originalFileName,     // Pakai nama asli file
        Parents: []string{"1AW9FxYl4qBKD2jXRY3P9fiLlxYltMdus"}, // Ganti dengan ID folder Google Drive
    }

    // Upload file ke Google Drive
    uploadedFile, err := srv.Files.Create(driveFile).Media(file).Do()
    if err != nil {
        return "", fmt.Errorf("error uploading file: %v", err)
    }

    // Kembalikan URL Google Drive
    return fmt.Sprintf("https://drive.google.com/file/d/%s/view", uploadedFile.Id), nil
}
