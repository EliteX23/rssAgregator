package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"rssAgregator/connection"
	"rssAgregator/internal/http"
	"rssAgregator/internal/logic"
	"rssAgregator/internal/postgres"
	"rssAgregator/logs"
)

var confPath = "config.json"

func init() {
	connection.InitConfig(confPath)
	logs.InitConfig(confPath)
}

func main() {
	logger := logs.CreateLogger()
	db := connection.InitDBRConnectionPG()
	articleRepo := postgres.NewArticleRepository(db)
	siteInfoRepo := postgres.NewSiteInfoRepository(db)
	siteRepo := postgres.NewSiteRepository(db)
	siteRulesRepo := postgres.NewSiteRulesRepository(db)

	articleService := logic.NewArticleService(articleRepo,logger)
	rssService := logic.NewRSSService(logger)
	siteService := logic.NewSiteService(logger,articleService, rssService, siteRepo, siteInfoRepo, siteRulesRepo)

	e := echo.New()

	e.Use(middleware.Logger())

	siteRoute := e.Group("/site")
	articleRoute := e.Group("/article")

	http.InitAndBindArticleHandler(articleRoute, articleService)
	http.InitAndBindSiteHandler(siteRoute, siteService)

	siteService.InitTaskFromDB()
	if err := e.Start(":8585"); err != nil {
		log.Panic("Произошла ошибка при старте веб-сервера!")
	}
}
