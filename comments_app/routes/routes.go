package routes

import (
	"comments_app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Create(app *fiber.App) {
	app.Get("/comments", controllers.Comments)
	app.Get("/commentsCount", controllers.CommentsCount)
	app.Post("/createComment", controllers.CreateComment)
}
