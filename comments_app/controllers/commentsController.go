package controllers

import (
	"strconv"
	"users_app/database"
	"users_app/models"

	"github.com/gofiber/fiber/v2"
)

func Comments(c *fiber.Ctx) error {
	page, err := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
		})
	}
	postId, err := strconv.ParseInt(c.Query("post", "0"), 10, 64)
	if err != nil || postId < 1 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
		})
	}
	offset := (int(page) - 1) * 10
	type Comment models.Comment
	comments := []Comment{}
	result := database.DB.Order("id ASC").Offset(int(offset)).Limit(10).Find(&comments, "post_id = ?", postId)

	if len(comments) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
		})
	}

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
		})
	}
	return c.JSON(comments)
}

func CommentsCount(c *fiber.Ctx) error {
	postId, err := strconv.ParseInt(c.Query("post", "0"), 10, 64)
	if err != nil || postId < 1 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
		})
	}
	type Comment models.Comment
	comment := Comment{}
	var count int64
	result := database.DB.Model(&comment).Where("post_id = ?", postId).Count(&count)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": count,
	})
}
