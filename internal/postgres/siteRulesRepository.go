package postgres

import (
	"github.com/gocraft/dbr"
	"rssAgregator/internal/app"
)

const siteRulesTable = "site_rules"

type siteRulesRepository struct {
	*dbr.Connection
}

func NewSiteRulesRepository(dbConn *dbr.Connection) app.SiteRulesRepository {
	return &siteRulesRepository{dbConn}
}
func (s siteRulesRepository) GetBySiteID(siteID int64) (app.SiteRules, error) {
	session := s.Connection.NewSession(nil)
	var siteRules app.SiteRules
	_, err := session.Select("*").
		From(siteRulesTable).
		Where("site_id=?", siteID).
		Load(&siteRules)
	if err != nil {
		return siteRules, err
	}
	return siteRules, nil
}

func (s siteRulesRepository) GetByID(id int64) (app.SiteRules, error) {
	session := s.Connection.NewSession(nil)
	var siteRules app.SiteRules
	_, err := session.Select("*").
		From(siteRulesTable).
		Where("id=?", id).
		Load(&siteRules)
	if err != nil {
		return siteRules, err
	}
	return siteRules, nil
}
func (s siteRulesRepository) Save(siteR *app.SiteRules) error {
	session := s.Connection.NewSession(nil)
	err := session.InsertInto(siteRulesTable).
		Columns("site_id", "article_root_name", "title", "description", "url", "pub_date").
		Record(siteR).
		Returning("id").
		Load(&siteR.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s siteRulesRepository) Update(siteR *app.SiteRules) error {
	session := s.Connection.NewSession(nil)
	_, err := session.Update(siteRulesTable).
		Set("pub_date", siteR.PubDate).
		Set("description", siteR.Description).
		Set("title", siteR.Title).
		Set("article_root_name", siteR.ArticleRootName).
		Set("is_deleted", siteR.IsDeleted).
		Where("id=?", siteR.ID).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

func (s siteRulesRepository) Remove(id int64) error {
	session := s.Connection.NewSession(nil)
	_, err := session.Update(siteRulesTable).
		Set("is_deleted", true).
		Where("id=?", id).
		Exec()
	if err != nil {
		return err
	}
	return nil
}
