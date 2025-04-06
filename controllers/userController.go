package controllers

import (
	"log"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/3ggie-AB/backend-animegg/models"
	"github.com/3ggie-AB/backend-animegg/config"
	"github.com/3ggie-AB/backend-animegg/helpers"
	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}
var validate = validator.New()

// Ambil semua user
func GetUsers(c *fiber.Ctx) error {
	blue := "\033[34m"
	reset := "\033[0m"
	var users []models.User
	err := config.DB.Find(&users).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal Mendapatkan Data User",
		})
	}

	log.Println(blue + "[Info GET USERS]" + reset + " : " + blue + "Berhasil Mendapatkan Data User" + reset)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil Mendapatkan Data User",
		"data":    users,
	})
}

// Buat user baru
func CreateUser(c *fiber.Ctx) error {
	var userRequest UserRequest
	err := c.BodyParser(&userRequest)
	if err != nil {
		helpers.Log(c,"warning", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Gagal Membaca Data User",
		})
	}

	err = validate.Struct(userRequest)
	if err != nil {
		helpers.Log(c,"warning", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.Log(c,"warning", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal Membuat User",
		})
	}
	
	user := models.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: string(hashedPassword),
		Role:     userRequest.Role,
	}
	err = config.DB.Create(&user).Error
	if err != nil {
		helpers.Log(c,"warning", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal Membuat User",
		})
	}
	
	helpers.Log(c,"info", "Berhasil Membuat User "+user.Name)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil Membuat User",
		"data":    user,
	})
}

// Update user
func UpdateUser(c *fiber.Ctx) error {
	blue := "\033[34m"
	reset := "\033[0m"
	var userRequest UserRequest
	err := c.BodyParser(&userRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Gagal Membaca Data User",
        })
    }

	err = validate.Struct(userRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": err.Error(),
        })
    }
	
    id := c.Params("id")
	var user models.User
	err = config.DB.First(&user, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User tidak ditemukan",
		})
	}

	user.Name = userRequest.Name
	user.Email = userRequest.Email
	user.Password = userRequest.Password
	user.Role = userRequest.Role
	err = config.DB.Save(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal Mengupdate User",
			"error": err.Error(),
		})
	}

	log.Println(blue + "[Info UPDATE USER]" + reset + " : " + blue + "Berhasil Mengupdate User" + reset)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil Mengupdate User",
		"data":    user,
    })
}

// Hapus user
func DeleteUser(c *fiber.Ctx) error {
	blue := "\033[34m"
	reset := "\033[0m"
	id := c.Params("id")
	var user models.User
	err := config.DB.First(&user, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User tidak ditemukan",
		})
	}

	err = config.DB.Delete(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal Menghapus User",
			"error": err.Error(),
		})
	}

	log.Println(blue + "[Info DELETE USER]" + reset + " : " + blue + "Berhasil Menghapus User" + reset)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil Menghapus User",
	})
}

// Ambil user berdasarkan ID
func GetUser(c *fiber.Ctx) error {
	blue := "\033[34m"
	reset := "\033[0m"
	id := c.Params("id")
	var user models.User
	err := config.DB.First(&user, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User tidak ditemukan",
		})
	}

	log.Println(blue + "[Info GET USER BY ID]" + reset + " : " + blue + "Berhasil Mendapatkan Data User" + reset)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil Mendapatkan Data User",
		"data":    user,
	})
}