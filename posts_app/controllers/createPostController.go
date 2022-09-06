package controllers

import (
	"fmt"
	"posts_app/database"
	"posts_app/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreatePost(c *fiber.Ctx) error {
	requestURI := (string(c.Request().RequestURI()))
	if requestURI != "/createPost" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "RequestURI not allowed",
		})
	}
	type Request struct {
		Id             string `json:"id"`
		Title          string `json:"title"`
		Content        string `json:"content"`
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
	post := models.Post{
		Title:          body.Title,
		Content:        body.Content,
		AuthorId:       authorId,
		AuthorUsername: body.AuthorUsername,
		Timestamp:      time.Now(),
	}
	if dbc := database.DB.Create(&post); dbc.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": dbc.Error.Error(),
		})
	}
	return c.JSON(post.Id)
}
