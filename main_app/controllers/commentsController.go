package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func CreateComment(c *fiber.Ctx) error {

	// If user is not logged in redirect them to login page
	IsUserLoggedIn(c)
	if LoggedInId == 0 {
		return c.Redirect("login")
	}

	postId := c.Query("id", "0")
	if postId == "0" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
		})
	}

	type CreateCommentData struct {
		Content string `json:"content"`
	}
	createCommentData := CreateCommentData{}
	err := c.BodyParser(&createCommentData)
	if err != nil {
		fmt.Println(err.Error())
		return c.Redirect("home")
	}
	content := createCommentData.Content

	sendData, err := json.Marshal(map[string]string{
		"postId":         postId,
		"content":        content,
		"authorId":       strconv.FormatUint(uint64(LoggedInId), 10),
		"authorUsername": LoggedInUsername,
	})
	if err != nil {
		fmt.Println(err.Error())
		return c.Redirect("home")
	}

	// Send request to posts app
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.SetRequestURI(os.Getenv("COMMENTS_APP_PATH") + "/createComment")
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
		fmt.Printf("Comment is not created. Error status code %v\n", resp.StatusCode())
		fmt.Println(string(resp.Body()))
		return c.Render("post?id="+postId, fiber.Map{
			"commentError": true,
		})
	}
	return c.Redirect(("post?id=" + postId))
}
