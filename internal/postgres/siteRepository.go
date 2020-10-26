package postgres

import (
	"errors"
	"github.com/gocraft/dbr"
	"rssAgregator/internal/app"
)

const siteTable = "sites"

type siteRepository struct {
	*dbr.Connection
}

func NewSiteRepository(dbConn *dbr.Connection) app.SiteRepository {
	return &siteRepository{dbConn}
}

func (s *siteRepository) GetByID(id int64) (*app.Site, error) {
	session := s.Connection.NewSession(nil)
	var site app.Site
	num, err := session.Select("*").
		From(siteTable).
		Where("id=?", id).
		Load(&site)
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, errors.New("site not found")
	}
	return &site, nil
}

func (s *siteRepository) GetList(filter app.QueryFilters) ([]app.Site, error) {
	siteList := make([]app.Site, 0, filter.Limit)
	session := s.Connection.NewSession(nil)
	b, err := filter.MakeStmt(session.Select("*").
		From(siteTable))
	if err != nil {
		return nil, err
	}
	b.Where("is_deleted = ?", false)

	_, err = b.Load(&siteList)
	if err != nil {
		return nil, err
	}
	return siteList, nil
}

func (s *siteRepository) GetAll() ([]app.Site, error) {
	var siteList []app.Site
	session := s.Connection.NewSession(nil)
	_, err := session.Select("*").
		From(siteTable).
		Where("is_deleted = ?", false).
		Load(&siteList)
	if err != nil {
		return nil, err
	}
	return siteList, nil
}

func (s *siteRepository) Save(site *app.Site) error {
	session := s.Connection.NewSession(nil)
	err := session.InsertInto(siteTable).
		Columns("link", "cron", "task_id").
		Record(site).
		Returning("id").
		Load(&site.ID)

	if err != nil {
		return err
	}
	return nil
}
func (s *siteRepository) Update(site *app.Site) error {
	session := s.Connection.NewSession(nil)
	_, err := session.Update(siteTable).
		Set("link", site.Link).
		Set("task_id", site.TaskID).
		Set("cron", site.Cron).
		Set("is_deleted", site.IsDeleted).
		Where("id=?", site.ID).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

func (s *siteRepository) Remove(id int64) error {
	session := s.Connection.NewSession(nil)
	_, err := session.Update(siteTable).
		Set("is_deleted", true).
		Where("id=?", id).
		Exec()
	if err != nil {
		return err
	}
	return nil
}
