package controllers

import "github.com/gofiber/fiber/v2"

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("jwt")
	return c.Redirect("login")
}
