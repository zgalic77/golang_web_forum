package controllers

import "github.com/gofiber/fiber/v2"

func OpenNotFound404(c *fiber.Ctx) error {
	IsUserLoggedIn(c)
	return c.Render("notFound", fiber.Map{
		"pageName": "404 page",
		"userId":   LoggedInId,
		"username": LoggedInUsername,
	})
}
