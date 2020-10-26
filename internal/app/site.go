package app

import (
	"github.com/gocraft/dbr"
)

type Site struct {
	ID        int64         `db:"id"`
	Link      string        `db:"link"`
	Cron      string        `db:"cron"`
	IsDeleted bool          `db:"is_deleted"`
	TaskID    dbr.NullInt64 `db:"task_id"`
	SiteRules SiteRules
	SiteInfo  SiteInfo
}

type SiteRules struct {
	ID              int64  `db:"id"`
	SiteID          int64  `db:"site_id"`
	ArticleRootName string `db:"article_root_name"`
	Title           string `db:"title"`
	URL             string `db:"url"`
	Description     string `db:"description"`
	PubDate         string `db:"pub_date"`
	IsDeleted       bool   `db:"is_deleted"`
}

type SiteRulesDTO struct {
	ArticleRootName string `json:"articleRootName"`
	Title           string `json:"title"`
	URL             string `json:"url"`
	Description     string `json:"description"`
	PubDate         string `json:"pubDate"`
}

type SiteInfo struct {
	ID          int64  `db:"id"`
	SiteID      int64  `db:"site_id"`
	Title       string `db:"title"`
	Link        string `db:"link"`
	Description string `db:"description"`
	Language    string `db:"language"`
	IsDeleted   bool   `db:"is_deleted"`
}

type SiteInfoDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Language    string `json:"language"`
}

type SiteDTO struct {
	ID        int64        `json:"id"`
	URL       string       `json:"url"`
	Cron      string       `json:"cron"`
	SiteRules SiteRulesDTO `json:"rules"`
	SiteInfo  SiteInfoDTO `json:"info,omitempty"`
}

type SiteRepository interface {
	GetByID(id int64) (*Site, error)
	GetList(filter QueryFilters) ([]Site, error)
	Save(site *Site) error
	Update(site *Site) error
	Remove(id int64) error
	GetAll() ([]Site, error)
}

type SiteInfoRepository interface {
	GetByID(id int64) (SiteInfo, error)
	GetBySiteID(siteID int64) (SiteInfo, error)
	Save(siteInfo *SiteInfo) error
	Update(siteInfo *SiteInfo) error
	Remove(id int64) error
}

type SiteRulesRepository interface {
	GetByID(id int64) (SiteRules, error)
	GetBySiteID(siteID int64) (SiteRules, error)
	Save(siteInfo *SiteRules) error
	Update(siteInfo *SiteRules) error
	Remove(id int64) error
}

type SiteService interface {
	GetByID(id int64) (*SiteDTO, error)
	GetList(filter QueryFilters) ([]SiteDTO, error)
	GetAll() ([]SiteDTO, error)
	Save(credit *SiteDTO) (*SiteDTO, error)
	Update(credit *SiteDTO) (*SiteDTO, error)
	Process(id int64) (bool, error)
	Remove(id int64) error
	InitTaskFromDB() (bool, error)
}

func (s *SiteDTO) MapToSiteEntity() Site {
	return Site{
		ID:        s.ID,
		Link:      s.URL,
		Cron:      s.Cron,
		SiteRules: SiteRules{},
		SiteInfo:  SiteInfo{},
	}
}

func (se *Site) MapToSiteDTO() SiteDTO {
	return SiteDTO{
		ID:        se.ID,
		URL:       se.Link,
		Cron:      se.Cron,
		SiteRules: SiteRulesDTO{},
		SiteInfo:  SiteInfoDTO{},
	}
}

func (si *SiteInfo) MapToDTO() SiteInfoDTO {
	return SiteInfoDTO{
		ID:          si.ID,
		Title:       si.Title,
		Link:        si.Link,
		Description: si.Description,
		Language:    si.Language,
	}
}

func (sr *SiteRules) MapToDTO() SiteRulesDTO {
	return SiteRulesDTO{
		ArticleRootName: sr.ArticleRootName,
		Title:           sr.Title,
		URL:             sr.URL,
		Description:     sr.Description,
		PubDate:         sr.PubDate,
	}
}

func (srd *SiteRulesDTO) MapToEntity() SiteRules {
	return SiteRules{
		ArticleRootName: srd.ArticleRootName,
		Title:           srd.Title,
		URL:             srd.URL,
		Description:     srd.Description,
		PubDate:         srd.PubDate,
	}
}

func (sid *SiteInfoDTO) MapToEntity() SiteInfo {
	return SiteInfo{

		Title:       sid.Title,
		Link:        sid.Link,
		Description: sid.Description,
		Language:    sid.Language,
		IsDeleted:   false,
	}
}
