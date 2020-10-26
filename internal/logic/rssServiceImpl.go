package logic

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"rssAgregator/internal/app"
	"time"
)

type rssService struct {
	log *logrus.Logger
}

func NewRSSService(_log *logrus.Logger) app.RSSService {
	return &rssService{
		log: _log,
	}
}
func (r *rssService) GetRSS(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 305 {
		response, _ := ioutil.ReadAll(resp.Body)
		r.log.Errorf("bad response from client site code: %v, response %v", resp.StatusCode, string(response))
		return nil, errors.New(fmt.Sprintf("response %v", string(response)))
	}
	responseArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseArr, nil
}

func (r *rssService) GetSiteInfo(bytes []byte) (app.SiteInfoDTO, error) {
	var siteInfo app.SiteInfoDTO
	var rssFeed app.RSSFeed
	err := xml.Unmarshal(bytes, &rssFeed)
	if err != nil {
		return app.SiteInfoDTO{}, err
	}
	siteInfo.Title=rssFeed.SiteInfo.Title
	siteInfo.Description=rssFeed.SiteInfo.Description
	siteInfo.Link=rssFeed.SiteInfo.Link
	siteInfo.Language=rssFeed.SiteInfo.Language
	return siteInfo, nil
}

func (r *rssService) GetArticles(bytes []byte, rules app.SiteRulesDTO) ([]app.ArticleDTO, error) {
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(bytes)
	if err != nil {
		return nil, err
	}
	root := doc.SelectElement("rss").SelectElement("channel")
	if root==nil{
		r.log.Errorf("root doc is empty")
		return nil, errors.New("bad rss")
	}
	var articleList []app.ArticleDTO
	for _, article := range root.SelectElements(rules.ArticleRootName) {
		var articleItem app.ArticleDTO
		if title := article.SelectElement(rules.Title); title != nil {
			articleItem.Title = title.Text()
		}
		if description := article.SelectElement(rules.Description); description != nil {
			articleItem.Description = description.Text()
		}
		if pubDate := article.SelectElement(rules.PubDate); pubDate != nil {
			time, err := time.Parse(time.RFC1123Z, pubDate.Text())
			if err != nil {
				r.log.Errorf("Error while parsing date :", err)
			} else {
				articleItem.PubDate = time
			}
		}
		if url := article.SelectElement(rules.URL); url != nil {
			articleItem.Link = url.Text()
		}
		articleList = append(articleList, articleItem)
	}
	return articleList, nil
}
