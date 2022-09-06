package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func OpenHomepage(c *fiber.Ctx) error {

	type Post struct {
		Id             uint   `json:"id"`
		Title          string `json:"title"`
		Content        string `json:"content"`
		AuthorId       uint   `json:"authorId"`
		AuthorUsername string `json:"authorUsername"`
		CommentsCount  int64
	}

	IsUserLoggedIn(c)
	page := c.Query("page", "1")

	// Send request to posts app
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.SetRequestURI(os.Getenv("POSTS_APP_PATH") + "/posts?page=" + page)
	request.Header.SetMethod(fasthttp.MethodGet)
	request.Header.SetContentType("application/json")
	request.Header.Set("Accept-Encoding", "gzip")
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttp.Do(request, resp)
	if err != nil {
		fmt.Printf("Client error: %s\n", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		fmt.Printf("Could not get post %v\n", resp.StatusCode())
		return c.Render("post", fiber.Map{"postError": true})
	}
	posts := []Post{}
	json.Unmarshal(resp.Body(), &posts)

	// Show only first 200 characters of post content on the homepage
	for i := range posts {
		length := len(posts[i].Content)
		splitValue := math.Min(float64(length), 200)
		posts[i].Content = posts[i].Content[0:int(splitValue)]
		if splitValue == 200 {
			posts[i].Content = posts[i].Content + "..."
		}
		posts[i].CommentsCount = getCommentsCount(uint64(posts[i].Id))
	}
	return c.Render("home", fiber.Map{
		"pageName":         "Forum Homepage",
		"userId":           LoggedInId,
		"username":         LoggedInUsername,
		"page":             page,
		"posts":            posts,
		"paginationValues": getPaginationValues(getPostsCount()),
	})
}

func RedirectHomepage(c *fiber.Ctx) error {
	return c.Redirect("/home")
}

func getCommentsCount(postId uint64) int64 {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.SetRequestURI(os.Getenv("COMMENTS_APP_PATH") + "/commentsCount?post=" + strconv.FormatUint(postId, 10))
	request.Header.SetMethod(fasthttp.MethodGet)
	request.Header.SetContentType("application/json")
	request.Header.Set("Accept-Encoding", "gzip")
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttp.Do(request, resp)
	if err != nil {
		fmt.Printf("Client error: %s\n", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		fmt.Printf("Could not get comments count %v\n", resp.StatusCode())
		return 0
	}

	type CountResponse struct {
		Count int64 `json:"count"`
	}
	count := CountResponse{}
	err = json.Unmarshal(resp.Body(), &count)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return count.Count
}

func getPostsCount() int64 {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.SetRequestURI(os.Getenv("POSTS_APP_PATH") + "/postsCount")
	request.Header.SetMethod(fasthttp.MethodGet)
	request.Header.SetContentType("application/json")
	request.Header.Set("Accept-Encoding", "gzip")
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttp.Do(request, resp)
	if err != nil {
		fmt.Printf("Client error: %s\n", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		fmt.Printf("Could not get comments count %v\n", resp.StatusCode())
		return 0
	}

	type CountResponse struct {
		Count int64 `json:"count"`
	}
	count := CountResponse{}
	err = json.Unmarshal(resp.Body(), &count)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return count.Count
}

func getPaginationValues(postsCount int64) []int64 {
	var pages float64 = float64(postsCount) / 10.
	pageNumber := math.Ceil(pages)
	var paginationResult []int64
	for i := 1; i <= int(pageNumber); i++ {
		paginationResult = append(paginationResult, int64(i))
	}
	return paginationResult
}
