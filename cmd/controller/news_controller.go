package controller

import (
	"net/http"
	"slices"
	"strconv"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/middlewares"
	"venecraft-back/cmd/service"

	"github.com/gin-gonic/gin"
)

type NewsController struct {
	NewsService service.NewsService
}

func NewNewsController(newsService service.NewsService) *NewsController {
	return &NewsController{NewsService: newsService}
}

// swagger:route POST /api/news news createNews
// Creates a new news article.
//
// Responses:
//
//	201: CommonSuccess
//	400: CommonError
//	403: CommonError
//	500: CommonError
func (nc *NewsController) CreateNews(c *gin.Context) {
	userID, roles, authenticated := middlewares.GetLoggedInUser(c)
	if !authenticated || !slices.Contains(roles, "ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var news entity.News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	news.CreatedBy = userID
	if err := nc.NewsService.CreateNews(&news); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, news)
}

// swagger:route GET /api/news news getAllNews
// Returns all news articles.
//
// Responses:
//
//	200: NewsListResponse
//	500: CommonError
func (nc *NewsController) GetAllNews(c *gin.Context) {
	news, err := nc.NewsService.GetAllNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}

// swagger:route GET /api/news/latest news getLatestNews
// Returns the latest news articles.
//
// Responses:
//
//   200: NewsListResponse
//   500: CommonError

func (nc *NewsController) GetLatestNews(c *gin.Context) {
	news, err := nc.NewsService.GetLatestNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}

// swagger:route GET /api/news/{id} news getNewsByID
// Returns a news article by its ID.
//
// Responses:
//
//	200: NewsResponse
//	400: CommonError
//	404: CommonError
//	500: CommonError
func (nc *NewsController) GetNewsByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	news, err := nc.NewsService.GetNewsByID(id)
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
// Responses:
//
//	200: NewsResponse
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
