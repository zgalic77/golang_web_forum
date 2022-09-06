package controllers

import (
	"fmt"
	"strconv"
	"time"
	"users_app/database"
	"users_app/models"

	"github.com/gofiber/fiber/v2"
)

func CreateComment(c *fiber.Ctx) error {
	requestURI := (string(c.Request().RequestURI()))
	if requestURI != "/createComment" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "RequestURI not allowed",
		})
	}
	type Request struct {
		Content        string `json:"content"`
		PostId         string `json:"postId"`
		AuthorId       string `json:"authorId"`
		AuthorUsername string `json:"authorUsername"`
	}
	var body Request

	err := c.BodyParser(&body)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}
	authorId, err := strconv.ParseUint(body.AuthorId, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
		})
	}
	postId, err := strconv.ParseUint(body.PostId, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
		})
	}
	comment := models.Comment{
		Content:        body.Content,
		PostId:         postId,
		Timestamp:      time.Now(),
		AuthorId:       authorId,
		AuthorUsername: body.AuthorUsername,
	}
	if dbc := database.DB.Create(&comment); dbc.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": dbc.Error.Error(),
		})
	}
	return c.JSON(comment.Id)
}
