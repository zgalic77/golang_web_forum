package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

func OpenRegister(c *fiber.Ctx) error {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	IsUserLoggedIn(c)
	if LoggedInId != 0 {
		c.Redirect("/home")
	}
	return c.Render("register", fiber.Map{})
}

func Register(c *fiber.Ctx) error {
	// Get data from form
	type RegisterData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	registerData := RegisterData{}
	err := c.BodyParser(&registerData)
	if err != nil {
		fmt.Println(err.Error())
		return c.Redirect("home")
	}
	username := registerData.Username
	password := registerData.Password
	email := registerData.Email
	sendData, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
		"email":    email,
	})
	if err != nil {
		fmt.Println(err.Error())
		return c.Redirect("home")
	}

	// Send request to users app
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.SetRequestURI(os.Getenv("USERS_APP_PATH") + "/register")
	request.Header.SetMethod(fasthttp.MethodPost)
	request.Header.SetContentType("application/json")
	request.SetBody(sendData)
	request.Header.Set("Accept-Encoding", "gzip")
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err = fasthttp.Do(request, resp)
	if err != nil {
		fmt.Printf("Client error: %s\n", err)
	}

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		fmt.Printf("Unsuccessful register with status code %v\n", resp.StatusCode())
		if strings.Contains(string(resp.Body()), "idx_users_username") {
			return c.Render("register", fiber.Map{"usernameError": true})
		}
		if strings.Contains(string(resp.Body()), "idx_users_email") {
			return c.Render("register", fiber.Map{"emailError": true})
		}
		return c.Render("register", fiber.Map{})
	}
	return c.Redirect("/login?id=" + string(c.Body()))
}
