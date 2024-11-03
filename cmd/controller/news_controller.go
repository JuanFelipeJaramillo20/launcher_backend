package controller

import (
	"venecraft-back/cmd/service"
)

type NewsController struct {
	NewsService service.NewsService
}

func NewNewsController(newsService service.NewsService) *NewsController {
	return &NewsController{newsService}
}

// swagger:route POST /api/news news createNews
// Creates a new news.
//
// Responses:
//
//	201: CommonSuccess
//	400: CommonError
//	409: CommonError
//	500: CommonError
