package middlewares

import (
    "github.com/gofiber/fiber/v2"
    "github.com/3ggie-AB/backend-animegg/models"
    "github.com/3ggie-AB/backend-animegg/config"
    "gorm.io/gorm"
)

func CheckToken(c *fiber.Ctx) error {
    // Ambil session ID dari header Authorization (misalnya Bearer <session_id>)
    tokenID := c.Get("Authorization")

    if tokenID == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Session ID is required",
        })
    }

    // Cek apakah format header valid (misalnya Bearer <session_id>)
    if len(tokenID) < 7 || tokenID[:7] != "Bearer " {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid session format",
        })
    }

    // Ambil session ID setelah "Bearer "
    tokenID = tokenID[7:]

    // Cek session di database
    var session models.Session
    if err := config.DB.Where("token = ?", tokenID).First(&session).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid session",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to validate session",
        })
    }

    // Simpan user_id di locals untuk diakses oleh handler berikutnya
    c.Locals("user_id", session.UserID)

    return c.Next()
}
