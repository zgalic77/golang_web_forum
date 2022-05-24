package controllers

import (
	"fmt"
	"posts_app/database"
	"posts_app/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Post(c *fiber.Ctx) error {
	postId, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"id": "0",
		})
	}
	post := models.Post{
		Id: uint(postId),
	}

	result := database.DB.Model(models.Post{}).First(&post)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"id": "0",
		})
	}
	return c.JSON(post)
}

func Posts(c *fiber.Ctx) error {
	page, err := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
		})
	}
	offset := (int(page) - 1) * 10
	type Post models.Post
	posts := []Post{}
	result := database.DB.Order("id DESC").Offset(int(offset)).Limit(10).Find(&posts)

	if len(posts) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
		})
	}

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
		})
	}
	return c.JSON(posts)
}

func PostsCount(c *fiber.Ctx) error {
	posts := models.Post{}
	var postsNum int64
	result := database.DB.Find(&posts).Count(&postsNum)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": postsNum,
	})
}
