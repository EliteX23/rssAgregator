package logic

import (
	"github.com/sirupsen/logrus"
	"rssAgregator/internal/app"
)

type articleService struct {
	articleRepo app.ArticleRepository
	log *logrus.Logger
}



func NewArticleService(
	_articleRepo app.ArticleRepository,
	_logs *logrus.Logger,
) app.ArticleService {
	return &articleService{
		log:_logs,
		articleRepo: _articleRepo,
	}
}
func (a *articleService) IsExist(articleTitle string) (bool, error) {
	return a.articleRepo.IsExist(articleTitle)
}
func (a *articleService) GetByID(id int64) (*app.ArticleDTO, error) {
	article, err := a.articleRepo.GetByID(id)
	if err != nil {
		a.log.Errorf("GetByID error: %v",err)
		return nil, err
	}
	articleDTO := article.ToDTO()
	return &articleDTO, nil
}

func (a *articleService) GetList(filter app.QueryFilters) (app.PagedList, error) {
	var result app.PagedList
	totalResults, err := a.articleRepo.GetCount(filter)
	if err != nil {
		return result, err
	}
	result.Total = totalResults

	var response []app.ArticleDTO
	articleList, err := a.articleRepo.GetList(filter)
	if err != nil {
		return result, err
	}
	for _, item := range articleList {
		response = append(response, item.ToDTO())
	}
	result.Result = response
	return result, nil
}

func (a *articleService) Save(siteID int64, credit *app.ArticleDTO) (*app.ArticleDTO, error) {
	articleEntity := credit.ToEntity()
	articleEntity.SiteID = siteID
	err := a.articleRepo.Save(&articleEntity)
	if err != nil {
		return nil, err
	}
	return a.GetByID(articleEntity.ID)
}

func (a *articleService) GetLast(siteID int64) (app.ArticleDTO, error) {
	var response app.ArticleDTO
	lastArticle, err := a.articleRepo.GetLast(siteID)
	if err != nil {
		return response, err
	}
	return lastArticle.ToDTO(), nil
}

func (a *articleService) Remove(id int64) error {
	return a.articleRepo.Remove(id)
}
