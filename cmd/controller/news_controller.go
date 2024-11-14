package controller

import (
	"log"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/middlewares"
	"venecraft-back/cmd/service"
	"venecraft-back/cmd/utils"

	"github.com/gin-gonic/gin"
)

// Request model for news creation
// swagger:model CreateNewsRequest

type CreateNewsRequest struct {
	// Title for the news article
	// required: true
	Title string `json:"title"`

	// Main content of the news
	// required: true
	Content string `json:"content"`

	// User ID
	// required: true
	CreatedBy uint64 `json:"created_by"`

	// Image for the news article
	// required: true
	Image *multipart.FileHeader `form:"image"`
}

// Parameters for creating a news
// swagger:parameters createNews
type CreateNewsParams struct {
	// news details for a news article creation
	// in: body
	// required: true
	Body CreateNewsRequest
}

// Parameters for retrieving, updating, or deleting a news by ID
// swagger:parameters getNewsByID updateNews deleteNews
type UserNewsParams struct {
	// ID of the news
	// in: path
	// required: true
	ID uint64 `json:"id"`
}

type NewsController struct {
	NewsService service.NewsService
}

func NewNewsController(newsService service.NewsService) *NewsController {
	return &NewsController{NewsService: newsService}
}

// swagger:route POST /api/news news createNews
// Creates a new news article.
//
// Security:
//   - BearerAuth: []
//
// Responses:
//
//	201: CommonSuccess
//	400: CommonError
//	403: CommonError
//	500: CommonError
func (nc *NewsController) CreateNews(c *gin.Context) {
	_, roles, authenticated := middlewares.GetLoggedInUser(c)
	if !authenticated || !slices.Contains(roles, "ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Retrieve other form data manually
	title := c.PostForm("title")
	content := c.PostForm("content")
	createdByStr := c.PostForm("created_by")

	// Parse `created_by` as an integer
	createdBy, err := strconv.ParseUint(createdByStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid created_by field"})
		return
	}

	// Retrieve the image file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	// Open the file to upload it to S3
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to open image file"})
		return
	}
	defer func(fileContent multipart.File) {
		err := fileContent.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to close image file"})
		}
	}(fileContent)

	// Upload the image to S3
	imageURL, err := utils.UploadFileToS3(fileContent, file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	// Create the news entity with the image URL
	news := entity.News{
		Title:     title,
		Content:   content,
		CreatedBy: createdBy,
		ImageURL:  imageURL, // Assuming `ImageURL` is added to the News struct
	}

	// Attempt to create the news entry in the database
	if err := nc.NewsService.CreateNews(&news); err != nil {
		// If there is an error, delete the uploaded file from S3
		deleteErr := utils.DeleteFileFromS3(file.Filename)
		if deleteErr != nil {
			log.Printf("Failed to delete file from S3 after error: %v", deleteErr)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, news)
}

// swagger:route GET /api/news news getAllNews
// Returns all news articles.
//
// Security:
//   - BearerAuth: []
//
// Responses:
//
//	200: []News
//	500: CommonError
func (nc *NewsController) GetAllNews(c *gin.Context) {
	userID, _, authenticated := middlewares.GetLoggedInUser(c)
	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	newsList, err := nc.NewsService.GetAllNews(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newsList)
}

// swagger:route GET /api/news/latest news getLatestNews
// Returns the latest news articles.
//
// Security:
//   - BearerAuth: []
//
// Responses:
//
//	200: []News
//	500: CommonError
func (nc *NewsController) GetLatestNews(c *gin.Context) {
	userID, _, authenticated := middlewares.GetLoggedInUser(c)
	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	newsList, err := nc.NewsService.GetLatestNews(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newsList)
}

// swagger:route GET /api/news/{id} news getNewsByID
// Returns a news article by its ID.
//
// Security:
//   - BearerAuth: []
//
// Responses:
//
//	200: News
//	400: CommonError
//	404: CommonError
//	500: CommonError
func (nc *NewsController) GetNewsByID(c *gin.Context) {
	userID, _, authenticated := middlewares.GetLoggedInUser(c)
	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	news, err := nc.NewsService.GetNewsByID(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if news == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	c.JSON(http.StatusOK, news)
}

// swagger:route PUT /api/news news updateNews
// Updates an existing news article.
//
// Security:
//   - BearerAuth: []
//
// Responses:
//
//	200: CommonSuccess
//	400: CommonError
//	403: CommonError
//	500: CommonError
func (nc *NewsController) UpdateNews(c *gin.Context) {
	var news entity.News
	_, roles, authenticated := middlewares.GetLoggedInUser(c)
	if !authenticated || !slices.Contains(roles, "ADMIN") && !slices.Contains(roles, "MODERATOR") {
		c.JSON(http.StatusForbidden, gin.H{"error": " ADMIN ro MODERATOR access required"})
		return
	}

	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := nc.NewsService.UpdateNews(&news)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News updated successfully"})
}

// swagger:route DELETE /api/news/{id} news deleteNews
// Deletes a news article by its ID.
//
// Security:
//   - BearerAuth: []
//
// Responses:
//
//	200: CommonSuccess
//	400: CommonError
//	500: CommonError
func (nc *NewsController) DeleteNews(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	if err := nc.NewsService.DeleteNews(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully"})
}

// swagger:route POST /api/news/{id}/like news likeNews
// Likes or unlikes a news post by ID.
//
// Security:
//   - BearerAuth: []
//
// Responses:
//
//	200: CommonSuccess
//	400: CommonError
//	500: CommonError
func (nc *NewsController) ToggleLikeNews(c *gin.Context) {
	userID, roles, authenticated := middlewares.GetLoggedInUser(c)
	if !authenticated || (!slices.Contains(roles, "ADMIN") && !slices.Contains(roles, "PLAYER")) {
		c.JSON(http.StatusForbidden, gin.H{"error": "User access required"})
		return
	}

	newsIDStr := c.Param("id")
	newsID, err := strconv.ParseUint(newsIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	liked, err := nc.NewsService.ToggleLikeNews(userID, newsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := "News unliked successfully"
	if liked {
		message = "News liked successfully"
	}
	c.JSON(http.StatusOK, gin.H{"message": message, "liked": liked})
}
