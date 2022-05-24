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

func OpenCreatePost(c *fiber.Ctx) error {
	IsUserLoggedIn(c)
	return c.Render("createPost", fiber.Map{
		"pageName": "Create new Forum post",
		"userId":   LoggedInId,
		"username": LoggedInUsername})
}

func CreatePost(c *fiber.Ctx) error {
	IsUserLoggedIn(c)
	if LoggedInId == 0 {
		return c.Redirect("/home")
	}
	type CreatePostData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	createPostData := CreatePostData{}
	err := c.BodyParser(&createPostData)
	if err != nil {
		fmt.Println(err.Error())
		return c.Redirect("home")
	}
	title := createPostData.Title
	content := createPostData.Content

	sendData, err := json.Marshal(map[string]string{
		"title":          title,
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
	request.SetRequestURI(os.Getenv("POSTS_APP_PATH") + "/createPost")
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
		fmt.Printf("Post is not created. Error status code %v\n", resp.StatusCode())
		fmt.Println(string(resp.Body()))
		return c.Render("createPost", fiber.Map{"error": "Post is not created. Please try again later."})
	}
	return c.Redirect(("post?id=" + string(resp.Body())))
}
