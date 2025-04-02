package controllers

import (
	"log"
	"net/http"
	"github.com/gofiber/fiber/v2"
	"github.com/3ggie-AB/backend-animegg/config"
	"github.com/3ggie-AB/backend-animegg/models"
	"golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
    "time"
)

// Register User
func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

// Login User
func Login(c *fiber.Ctx) error {
	// Parse request body untuk mendapatkan email dan password
	var userRequest models.User
	if err := c.BodyParser(&userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Cari user berdasarkan email
	var user models.User
	if err := config.DB.Where("email = ?", userRequest.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Bandingkan password dengan hash yang ada di database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Buat session baru (token)
	tokenID := uuid.New().String()

	// Simpan session ke database
	expiresAt := time.Now().Add(365 * 24 * time.Hour)  // 1 jam dari sekarang

    session := models.Session{
        UserID:    user.ID,
        Token:     tokenID,
        ExpiresAt: expiresAt,  // Gunakan waktu yang valid
        CreatedAt: time.Now(),
    }
	if err := config.DB.Create(&session).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session",
		})
	}

	// Kirim token kepada client
	return c.JSON(fiber.Map{
		"message":   "Login successful",
		"Token": tokenID,
	})
}

// Logout User
func Logout(c *fiber.Ctx) error {
    // Ambil token dari header Authorization
    tokenID := c.Get("Authorization")

    if tokenID == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Session ID is required",
        })
    }

    // Hapus "Bearer " jika ada di depan token
    if len(tokenID) > 7 && tokenID[:7] == "Bearer " {
        tokenID = tokenID[7:]
    }

    // Hapus session berdasarkan token
    var session models.Session
    if err := config.DB.Where("token = ?", tokenID).Delete(&session).Error; err != nil {
        log.Println("Failed to delete session:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to logout",
        })
    }

    // Berhasil logout
    return c.JSON(fiber.Map{
        "message": "Logout successful",
    })
}
