package routes

import (
	"users_app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Create(app *fiber.App) {

	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Get("/user", controllers.User)
	app.Post("/logout", controllers.Logout)

}
