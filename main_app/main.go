package main

import (
	"log"
	"main_app/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		},
	})
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Static("/", os.Getenv("STATIC_FILES_PATH"))
	routes.Create(app)
	app.Listen(os.Getenv("LOCALHOST"))
}
