package controllers

import (
	"log"
	"net/http"
	"github.com/gofiber/fiber/v2"
	"github.com/3ggie-AB/backend-animegg/config"
	"github.com/3ggie-AB/backend-animegg/models"
	"golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
	"gorm.io/gorm"
    "time"
)

// Register User
func Register(c *fiber.Ctx) error {
	blue := "\033[34m"
	// red := "\033[31m"
	reset := "\033[0m"
	var data models.User
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Input tidak valid",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Hash Password Gagal",
		})
	}

	data.Password = string(hashedPassword)
	data.Role = "user"

	// Create user
	if err := config.DB.Create(&data).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Maaf Tidak dapat Membuat User",
		})
	}

	log.Println(blue + "[Info REGISTER]" + reset + " : " + blue + "Request name: \"" + data.Name + "\"" + " Request email: \"" + data.Email + "\"" +reset)
	return c.JSON(fiber.Map{
		"success": true,
        "message": "Berhasil Membuat User Silahkan Login",
	})
}

// Login User
func Login(c *fiber.Ctx) error {
	blue := "\033[34m"
	// red := "\033[31m"
	reset := "\033[0m"
	// Parse request body untuk mendapatkan email dan password
	var userRequest models.User
	if err := c.BodyParser(&userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success" : false,
			"error": "Membutuhkan Data Email dan Password",
		})
	}
	
	// Cari user berdasarkan email
	var user models.User
	if err := config.DB.Where("email = ?", userRequest.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
            "message": "Email atau Password Salah",
		})
	}
	
	// Bandingkan password dengan hash yang ada di database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
            "message": "Email atau Password Salah",
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
			"success": false,
            "message": "Gagal membuat Session",
		})
	}
	
	log.Println(blue + "[Info LOGIN]" + reset + " : " + blue + " Request email: \"" + userRequest.Email + "\"" + reset)
	// Kirim token kepada client
	return c.JSON(fiber.Map{
		"success": true,
        "message": "Login Berhasil",
        "token": tokenID,
	})
}

// Logout User
func Logout(c *fiber.Ctx) error {
    // Ambil token dari header Authorization
    tokenID := c.Get("Authorization")

    if tokenID == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
            "message": "Token Session Harus di Isi",
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
			"success": false,
			"message": "Gagal menghapus session",
        })
    }

    // Berhasil logout
    return c.JSON(fiber.Map{
		"success": true,
        "message": "Logout successful",
    })

}

// Controller Cek Token
func CekToken(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	tokenID := c.Get("Authorization")
	if tokenID == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
            "message": "Token Session Harus di Isi",
        })
    }

	// Hapus "Bearer " jika ada di depan token
	if len(tokenID) > 7 && tokenID[:7] == "Bearer " {
		tokenID = tokenID[7:]
	}

	// Cek session di database
	var session models.Session
	if err := config.DB.Where("token = ?", tokenID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "success": false,
                "message":   "Token Session Tidak Valid",
            })
        }
		// Jika error selain record not found, kembalikan error internal server
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
	}

	// Jika session ditemukan tapi sudah expired
	if session.ExpiresAt.Before(time.Now().UTC()) {
		config.DB.Delete(&session)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message":   "Token Session Kadaluarsa",
		})
	}

	// Cek apakah user masih ada di database
	var user models.User
	if err := config.DB.First(&user, session.UserID).Error; err != nil {
		// Jika user tidak ditemukan, hapus session dan kirim respon
		config.DB.Delete(&session)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message":   "Token Session Tidak Valid",
		})
	}

	// Jika semua sudah benar, kirim respon kepada client
    return c.JSON(fiber.Map{
        "success": true,
        "role":    user.Role,
        "message": "Token Session Valid",
    })
}
