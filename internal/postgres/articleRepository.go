package postgres

import (
	"github.com/gocraft/dbr"
	"rssAgregator/internal/app"
)

const articleTable = "articles"

type articleRepository struct {
	*dbr.Connection
}

func NewArticleRepository(dbConn *dbr.Connection) app.ArticleRepository {
	return &articleRepository{dbConn}
}

func (a articleRepository) IsExist(articleTitle string) (bool, error) {
	var result int64
	session := a.Connection.NewSession(nil)
	_, err := session.Select("count (id)").
		From(articleTable).
		Where("title=?", articleTitle).
		Load(&result)
	if err != nil {
		return false, err
	}
	if result > 0 {
		return true, nil
	}
	return false, nil
}
func (a articleRepository) GetCount(filter app.QueryFilters) (int64, error) {
	var result int64
	session := a.Connection.NewSession(nil)
	b, err := filter.MakeStmt(session.Select("count (id)").
		From(articleTable))
	if err != nil {
		return result, err
	}
	b.Where("is_deleted = ?", false)
	if len(filter.Title) > 0 {
		b.Where("title ilike '%' || ? || '%'", filter.Title)
	}
	_, err = b.Load(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (a articleRepository) GetByID(id int64) (*app.Article, error) {
	session := a.Connection.NewSession(nil)
	var article app.Article
	_, err := session.Select("*").
		From(articleTable).
		Where("id=?", id).
		Load(&article)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (a articleRepository) GetList(filter app.QueryFilters) ([]app.Article, error) {
	articleList := make([]app.Article, 0, filter.Limit)
	session := a.Connection.NewSession(nil)
	b,err := filter.MakeStmt(session.Select("*").
		From(articleTable))

	if err != nil {
		return articleList, err
	}
		b.Where("is_deleted = ?", false)
	if len(filter.Title) > 0 {
		b.Where("title ilike '%' || ? || '%'", filter.Title)
	}
	_, err = b.Load(&articleList)
	if err != nil {
		return nil, err
	}
	return articleList, nil
}
func (a articleRepository) Save(article *app.Article) error {
	session := a.Connection.NewSession(nil)
	err := session.InsertInto(articleTable).
		Columns("site_id", "title", "link", "description", "pub_date").
		Record(article).
		Returning("id").
		Load(&article.ID)

	if err != nil {
		return err
	}
	return nil
}

func (a articleRepository) GetLast(siteID int64) (app.Article, error) {
	session := a.Connection.NewSession(nil)
	var article app.Article
	_, err := session.Select("*").
		From(articleTable).
		Where("site_id = ?", siteID).
		OrderDesc("id").
		Limit(1).
		Load(&article)
	if err != nil {
		return article, err
	}

	return article, nil
}

func (a articleRepository) Remove(id int64) error {
	session := a.Connection.NewSession(nil)
	_, err := session.Update(articleTable).
		Set("is_deleted", true).
		Where("id=?", id).
		Exec()
	if err != nil {
		return err
	}
	return nil
}
