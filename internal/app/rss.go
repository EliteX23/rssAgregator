package app

type RSSService interface {
	GetRSS(url string) ([]byte, error)
	GetSiteInfo([]byte) (SiteInfoDTO, error)
	GetArticles([]byte, SiteRulesDTO) ([]ArticleDTO, error)
}

type RSSFeed struct {
	SiteInfo RSSSite `xml:"channel"`
}
type RSSSite struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
}

func (r *RSSSite) MapToSiteInfo() SiteInfoDTO {
	return SiteInfoDTO{
		Title:       r.Title,
		Link:        r.Link,
		Description: r.Description,
		Language:    r.Language,
	}
}
