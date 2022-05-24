package routes

import (
	"main_app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Create(app *fiber.App) {

	// Homepage routes
	app.Get("/", controllers.RedirectHomepage)
	app.Get("/home", controllers.OpenHomepage)

	// User authentication routes
	app.Get("/login", controllers.OpenLogin)
	app.Post("/login", controllers.Login)
	app.Get("/logout", controllers.Logout)
	app.Get("/register", controllers.OpenRegister)
	app.Post("/register", controllers.Register)

	// Get post route
	app.Get("/post", controllers.OpenPost)

	// Add comment to the post route
	app.Post("/post", controllers.CreateComment)

	// Create post routes
	app.Get("/createPost", controllers.OpenCreatePost)
	app.Post("/createPost", controllers.CreatePost)

	// Get 404 for non-implemented pages routes
	app.Get("*", controllers.OpenNotFound404)
}
