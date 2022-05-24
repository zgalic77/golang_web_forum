package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"users_app/database"
	"users_app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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
		Id             string `json:"id"`
		Content        string `json:"content"`
		PostId         string `json:"postId"`
		AuthorId       string `json:"authorId"`
		AuthorUsername string `json:"authorUsername"`
	}
	var body Request

	// Get password cost value from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file.")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
		})
	}
	err = c.BodyParser(&body)
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
