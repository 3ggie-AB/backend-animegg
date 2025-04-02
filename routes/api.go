package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/3ggie-AB/backend-animegg/controllers"
	"github.com/3ggie-AB/backend-animegg/middlewares"
)

func ApiRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Routes untuk Authentication
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Post("/logout", middlewares.CheckToken, controllers.Logout)

	// Routes untuk Anime
	api.Post("/anime", middlewares.CheckToken, controllers.CreateAnime)
	api.Get("/anime", controllers.GetAnimes)
	api.Get("/anime/:id", controllers.GetAnime)
	
	// Routes untuk Episode
	api.Get("/anime/:id/episode/:episode", controllers.GetAnimeEpisode)
	api.Post("/episode", middlewares.CheckToken, controllers.CreateEpisode)
	api.Get("/episode/:anime_id", controllers.GetEpisodesByAnime)
}
