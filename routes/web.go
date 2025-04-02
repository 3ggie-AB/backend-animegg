package routes

import (
	"github.com/gofiber/fiber/v2"
)

func WebRoutes(app *fiber.App) {
	app.Static("/public", "./public")
}
