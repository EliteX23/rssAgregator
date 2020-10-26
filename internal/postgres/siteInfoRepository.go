package postgres

import (
	"github.com/gocraft/dbr"
	"rssAgregator/internal/app"
)

const siteInfoTable = "site_info"

type siteInfoRepository struct {
	*dbr.Connection
}

func NewSiteInfoRepository(dbConn *dbr.Connection) app.SiteInfoRepository {
	return &siteInfoRepository{dbConn}
}

func (s *siteInfoRepository) GetByID(id int64) (app.SiteInfo, error) {
	session := s.Connection.NewSession(nil)
	var site app.SiteInfo
	_, err := session.Select("*").
		From(siteInfoTable).
		Where("id=?", id).
		Load(&site)
	if err != nil {
		return site, err
	}
	return site, nil
}
func (s *siteInfoRepository) GetBySiteID(siteID int64) (app.SiteInfo, error) {
	session := s.Connection.NewSession(nil)
	var site app.SiteInfo
	_, err := session.Select("*").
		From(siteInfoTable).
		Where("site_id=?", siteID).
		Load(&site)
	if err != nil {
		return site, err
	}
	return site, nil
}

func (s *siteInfoRepository) Save(siteInfo *app.SiteInfo) error {
	session := s.Connection.NewSession(nil)
	err := session.InsertInto(siteInfoTable).
		Columns("site_id", "title", "link", "description", "language").
		Record(siteInfo).
		Returning("id").
		Load(&siteInfo.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *siteInfoRepository) Update(siteInfo *app.SiteInfo) error {
	session := s.Connection.NewSession(nil)
	_, err := session.Update(siteInfoTable).
		Set("language", siteInfo.Language).
		Set("description", siteInfo.Description).
		Set("title", siteInfo.Title).
		Set("link", siteInfo.Link).
		Set("is_deleted", siteInfo.IsDeleted).
		Where("id=?", siteInfo.ID).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

func (s *siteInfoRepository) Remove(id int64) error {
	session := s.Connection.NewSession(nil)
	_, err := session.Update(siteInfoTable).
		Set("is_deleted", true).
		Where("id=?", id).
		Exec()
	if err != nil {
		return err
	}
	return nil
}
