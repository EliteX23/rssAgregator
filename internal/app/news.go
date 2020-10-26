package app

import (
	"time"
)

type Article struct {
	ID          int64     `db:"id"`
	SiteID      int64     `db:"site_id"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	Description string    `db:"description"`
	PubDate     time.Time `db:"pub_date"`
	IsDeleted   bool      `db:"is_deleted"`
}

type ArticleDTO struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"pubDate"`
}

type ArticleRepository interface {
	GetByID(id int64) (*Article, error)
	GetList(filter QueryFilters) ([]Article, error)
	GetCount(filter QueryFilters) (int64, error)
	Save(article *Article) error
	GetLast(siteID int64) (Article, error)
	IsExist(articleTitle string) (bool, error)
	Remove(id int64) error
}

type ArticleService interface {
	GetByID(id int64) (*ArticleDTO, error)
	IsExist(articleTitle string) (bool, error)
	GetList(filter QueryFilters) (PagedList, error)
	Save(siteID int64, credit *ArticleDTO) (*ArticleDTO, error)
	GetLast(siteID int64) (ArticleDTO, error)
	Remove(id int64) error
}

func (ad *ArticleDTO) ToEntity() Article {
	return Article{
		ID:          ad.ID,
		Title:       ad.Title,
		Link:        ad.Link,
		Description: ad.Description,
		PubDate:     ad.PubDate,
	}
}

func (a *Article) ToDTO() ArticleDTO {
	return ArticleDTO{
		ID:          a.ID,
		Title:       a.Title,
		Link:        a.Link,
		Description: a.Description,
		PubDate:     a.PubDate,
	}
}
