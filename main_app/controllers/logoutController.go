package controllers

import "github.com/gofiber/fiber/v2"

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("jwt")
	LoggedInUsername = ""
	LoggedInId = 0
	return c.Redirect("login")
}
