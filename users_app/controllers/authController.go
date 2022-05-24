package controllers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"users_app/database"
	"users_app/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	type Request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Request

	// Get password cost value from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file.")
	}
	password_cost, err := strconv.Atoi(os.Getenv("PASSWORD_COST"))
	if err != nil {
		return err
	}

	err = c.BodyParser(&body)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}

	password, err := bcrypt.GenerateFromPassword([]byte(body.Password), password_cost)
	if err != nil {
		panic(fmt.Sprintf("Error encrypting password!\n%s", err))
	}

	user := models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: password,
	}

	if dbc := database.DB.Create(&user); dbc.Error != nil {
		c.Context().SetStatusCode(fiber.ErrBadRequest.Code)
		return c.JSON(fiber.Map{
			"message": dbc.Error.Error(),
		})
	}

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body Request

	// Get JWT secret key from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file.")
	}
	secret_key := os.Getenv("JWT_KEY")
	if err != nil {
		panic(err)
	}

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	// Search for user in the database
	var user models.User
	database.DB.Where("username = ?", body.Username).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User is not found.",
		})
	}

	// Check if password is correct
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(body.Password))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Password is incorrect.",
		})
	}

	// Create JWT token that is valid for 48 hours
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	})

	token, err := claims.SignedString([]byte(secret_key))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not log in.",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "JWT token has been created successfully.",
	})
}

func User(c *fiber.Ctx) error {

	cookie := c.Cookies("jwt")

	// Get JWT secret key from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file.")
	}
	secret_key := os.Getenv("JWT_KEY")
	if err != nil {
		panic(err)
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "User is not logged in.",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)
	var user models.User
	database.DB.Where("id = ?", claims.Issuer).First(&user)
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {

	// On logout set JWT expiration time in the past
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Minute),
		HTTPOnly: true,
	}

	// Set new cookie values and return message
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "User has been logged out.",
	})
}
