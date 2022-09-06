package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

var LoggedInUsername string = ""
var LoggedInId uint = 0

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserData struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

func OpenLogin(c *fiber.Ctx) error {
	IsUserLoggedIn(c)
	if LoggedInId != 0 {
		c.Redirect("/home")
	}

	return c.Render("login", fiber.Map{})
}

func Login(c *fiber.Ctx) error {
	// Get data from form
	loginData := LoginData{}
	err := c.BodyParser(&loginData)
	if err != nil {
		fmt.Println(err.Error())
		return c.Redirect("home")
	}
	username := loginData.Username
	password := loginData.Password
	sendData, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		fmt.Println(err.Error())
		return c.Redirect("home")
	}

	// Send request to users app
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.SetRequestURI(os.Getenv("USERS_APP_PATH") + "/login")
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

	if resp.StatusCode() != fasthttp.StatusOK {
		fmt.Printf("Unsuccessful login with status code %v\n", resp.StatusCode())
		return c.Render("login", fiber.Map{"loginError": true})
	}
	setCookieValue := string(resp.Header.Peek("Set-Cookie"))
	cookieInfo := strings.Split(setCookieValue, "; ")
	cookieName := strings.Split(cookieInfo[0], "=")[0]
	cookieValue := strings.Split(cookieInfo[0], "=")[1]
	if err != nil {
		log.Fatal(err)
	}
	cookie := fiber.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Expires:  time.Now().Add(time.Hour * 48),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.Redirect("/home")
}

func IsUserLoggedIn(c *fiber.Ctx) bool {
	// Checks if user is already logged in
	// If user is logged in then set cookie expiration date 48 hours in the future
	cookieValue := c.Cookies("jwt")
	if cookieValue == "" {
		LoggedInUsername = ""
		LoggedInId = 0
		return false
	}
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.SetRequestURI(os.Getenv("USERS_APP_PATH") + "/user")
	request.Header.SetContentType("application/json")
	request.Header.SetCookie("jwt", cookieValue)
	request.Header.Set("Accept-Encoding", "gzip")
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := fasthttp.Do(request, resp)
	if err != nil {
		fmt.Printf("Client error: %s\n", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		LoggedInUsername = ""
		LoggedInId = 0
		return false
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    c.Cookies("jwt"),
		Expires:  time.Now().Add(time.Hour * 48),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	userData := UserData{}
	err = json.Unmarshal(resp.Body(), &userData)
	if err != nil {
		fmt.Println(err.Error())
		LoggedInUsername = ""
		LoggedInId = 0
		return false
	}
	LoggedInUsername = userData.Username
	LoggedInId = userData.Id
	if err != nil {
		fmt.Println(err.Error())
		LoggedInUsername = ""
		LoggedInId = 0
		return false
	}
	return true
}
