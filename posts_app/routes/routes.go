package routes

import (
	"posts_app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Create(app *fiber.App) {
	app.Post("/createPost", controllers.CreatePost)
	app.Get("/post", controllers.Post)
	app.Get("/posts", controllers.Posts)
	app.Get("/postsCount", controllers.PostsCount)
}
