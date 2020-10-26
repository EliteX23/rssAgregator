package http

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"rssAgregator/internal/app"
	"strconv"

	"github.com/labstack/echo"
)

type httpSiteHandler struct {
	siteServ app.SiteService
}

func InitAndBindSiteHandler(
	rg *echo.Group,
	_siteServ app.SiteService,
) *echo.Group {
	h := httpSiteHandler{
		siteServ: _siteServ,
	}

	rg.GET("/", h.GetList)
	rg.GET("/:siteID/", h.GetSite)
	rg.POST("/", h.Add)
	rg.PUT("/", h.Update)
	rg.DELETE("/:siteID", h.Remove)
	rg.GET("/process/:siteID/", h.ProcessSite)
	return rg
}

func (h *httpSiteHandler) GetList(c echo.Context) error {
	var pageFilters app.QueryFilters
	if err := c.Bind(&pageFilters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "pagination settings is invalid")
	}

	response, err := h.siteServ.GetList(pageFilters)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, response)
}

func (h *httpSiteHandler) GetSite(c echo.Context) error {
	siteIDStr := c.Param("siteID")
	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "siteID is not valid")
	}
	site, err := h.siteServ.GetByID(int64(siteID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, site)
}

//ProcessSite принудительный запуск обработки для сайта
func (h *httpSiteHandler) ProcessSite(c echo.Context) error {
	siteIDStr := c.Param("siteID")
	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "siteID is not valid")
	}
	isOK, err := h.siteServ.Process(int64(siteID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, isOK)
}

func (h *httpSiteHandler) Add(c echo.Context) error {

	body := c.Request().Body
	bodyReq, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.New("bad body")
	}
	defer body.Close()

	var dto app.SiteDTO
	err = json.Unmarshal(bodyReq, &dto)
	if err != nil {
		return err
	}

	site, err := h.siteServ.Save(&dto)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, site)
}

func (h *httpSiteHandler) Update(c echo.Context) error {
	body := c.Request().Body
	bodyReq, err := ioutil.ReadAll(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad body")
	}
	defer body.Close()

	var dto app.SiteDTO
	err = json.Unmarshal(bodyReq, &dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	site, err := h.siteServ.Update(&dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, site)
}

func (h *httpSiteHandler) Remove(c echo.Context) error {
	siteIDStr := c.Param("siteID")
	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "siteID is not valid")
	}
	err = h.siteServ.Remove(int64(siteID))
	return c.JSON(http.StatusOK, err)
}
