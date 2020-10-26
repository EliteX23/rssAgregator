package http

import (
	"net/http"
	"rssAgregator/internal/app"
	"strconv"

	"github.com/labstack/echo"
)

type httpArticleHandler struct {
	articleServ app.ArticleService
}

func InitAndBindArticleHandler(
	rg *echo.Group,
	_articleServ app.ArticleService,
) *echo.Group {
	h := httpArticleHandler{
		articleServ: _articleServ,
	}

	rg.GET("/", h.GetList)
	rg.GET("/:articleID/", h.Get)
	return rg
}

func (h *httpArticleHandler) GetList(c echo.Context) error {
	var pageFilters app.QueryFilters
	if err := c.Bind(&pageFilters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "pagination settings is invalid")
	}
	response, err := h.articleServ.GetList(pageFilters)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, response)
}

func (h *httpArticleHandler) Get(c echo.Context) error {
	articleIDStr := c.Param("articleID")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "articleID is not valid")
	}
	article, err := h.articleServ.GetByID(int64(articleID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, article)
}
