package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

func OpenPost(c *fiber.Ctx) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	type Post struct {
		Id                 string    `json:"id"`
		Title              string    `json:"title"`
		Content            string    `json:"content"`
		Timestamp          time.Time `json:"timestamp"`
		AuthorId           string    `json:"authorId"`
		AuthorUsername     string    `json:"authorUsername"`
		FormattedTimestamp string
	}

	type Comment struct {
		Id                 uint      `json:"id"`
		Content            string    `json:"content"`
		Timestamp          time.Time `json:"timestamp"`
		AuthorId           uint      `json:"authorId"`
		AuthorUsername     string    `json:"authorUsername"`
		FormattedTimestamp string
	}

	IsUserLoggedIn(c)
	postId := c.Query("id")
	page := c.Query("page")

	// Send request to posts app
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.SetRequestURI(os.Getenv("POSTS_APP_PATH") + "/post?id=" + postId)
	request.Header.SetMethod(fasthttp.MethodGet)
	request.Header.SetContentType("application/json")
	request.Header.Set("Accept-Encoding", "gzip")
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err = fasthttp.Do(request, resp)
	if err != nil {
		fmt.Printf("Client error: %s\n", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		fmt.Printf("Could not get post %v\n", resp.StatusCode())
		return c.Render("post", fiber.Map{"postError": true})
	}
	post := Post{}
	json.Unmarshal(resp.Body(), &post)
	post.FormattedTimestamp = post.Timestamp.Format("02.01.2006 15:04")

	// Get comments
	requestComments := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(requestComments)
	requestComments.SetRequestURI(os.Getenv("COMMENTS_APP_PATH") + "/comments?post=" + postId + "&page=" + page)
	requestComments.Header.SetMethod(fasthttp.MethodGet)
	requestComments.Header.SetContentType("application/json")
	requestComments.Header.Set("Accept-Encoding", "gzip")
	respComments := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(respComments)
	err = fasthttp.Do(requestComments, respComments)
	if err != nil {
		fmt.Printf("Client error: %s\n", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		fmt.Printf("Could not get post %v\n", resp.StatusCode())
		return c.Render("post", fiber.Map{"postError": true})
	}
	comments := []Comment{}
	hasComments := false
	if respComments.StatusCode() == fiber.StatusOK {
		json.Unmarshal(respComments.Body(), &comments)
		hasComments = true
	}

	// Format comment timestamps as dd.mm.yyyy hh:mm
	for i := range comments {
		comments[i].FormattedTimestamp = comments[i].Timestamp.Format("02.01.2006 15:04")
	}
	// Convert postId to int64 to call getCommentsCount method

	convertedPostId, err := strconv.ParseUint(postId, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		convertedPostId = 0
	}
	paginationValues := getPaginationValues(getCommentsCount(convertedPostId))
	return c.Render("post", fiber.Map{
		"pageName":         "Post",
		"userId":           LoggedInId,
		"username":         LoggedInUsername,
		"post":             post,
		"hasComments":      hasComments,
		"comments":         comments,
		"paginationValues": paginationValues,
		"postId":           postId,
	})
}
