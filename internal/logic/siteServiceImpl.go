package logic

import (
	"github.com/gocraft/dbr"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"rssAgregator/internal/app"
)

type siteService struct {
	log *logrus.Logger
	rssService     app.RSSService
	scheduler      *cron.Cron
	articleService app.ArticleService
	siteRepo       app.SiteRepository
	siteInfoRepo   app.SiteInfoRepository
	siteRulesRepo  app.SiteRulesRepository
}

func NewSiteService(
	_log *logrus.Logger,
	_articleServ app.ArticleService,
	_rssService app.RSSService,
	_siteRepo app.SiteRepository,
	_siteInfoRepo app.SiteInfoRepository,
	_siteRulesRepo app.SiteRulesRepository,
) app.SiteService {
	service := &siteService{
		log: _log,
		rssService:     _rssService,
		scheduler:      cron.New(),
		articleService: _articleServ,
		siteRepo:       _siteRepo,
		siteInfoRepo:   _siteInfoRepo,
		siteRulesRepo:  _siteRulesRepo,
	}
	service.scheduler.Start()
	return service
}

func (s *siteService) GetByID(id int64) (*app.SiteDTO, error) {
	site, err := s.siteRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	siteDTO := site.MapToSiteDTO()
	siteInfo, err := s.siteInfoRepo.GetBySiteID(site.ID)
	if err != nil {
		return nil, err
	}
	if siteInfo.ID > 0 {
		siteInfoDto := siteInfo.MapToDTO()
		siteDTO.SiteInfo = siteInfoDto
	}
	siteRules, err := s.siteRulesRepo.GetBySiteID(site.ID)
	if err != nil {
		return nil, err
	}
	if siteRules.ID > 0 {
		siteRulesDto := siteRules.MapToDTO()
		siteDTO.SiteRules = siteRulesDto
	}
	return &siteDTO, err
}

func (s *siteService) GetList(filter app.QueryFilters) ([]app.SiteDTO, error) {
	var response []app.SiteDTO
	siteList, err := s.siteRepo.GetList(filter)
	if err != nil {
		return nil, err
	}
	for _, item := range siteList {
		response = append(response, item.MapToSiteDTO())
	}
	return response, nil
}

func (s *siteService) Save(credit *app.SiteDTO) (*app.SiteDTO, error) {
	siteEntity := credit.MapToSiteEntity()

	err := s.siteRepo.Save(&siteEntity)
	if err != nil {
		return nil, err
	}
	credit.ID = siteEntity.ID
	taskID, err := s.addCronTask(credit)
	if err != nil {
		return nil, err
	}
	siteEntity.TaskID = dbr.NewNullInt64(taskID)
	err = s.siteRepo.Update(&siteEntity)
	if err != nil {
		s.removeCronTask(taskID)
		return nil, err
	}

	siteRulesEntity := credit.SiteRules.MapToEntity()
	siteRulesEntity.SiteID = siteEntity.ID
	err = s.siteRulesRepo.Save(&siteRulesEntity)
	if err != nil {
		s.removeCronTask(taskID)
		return nil, err
	}
	return s.GetByID(siteEntity.ID)
}

func (s *siteService) Update(credit *app.SiteDTO) (*app.SiteDTO, error) {
	site, err := s.siteRepo.GetByID(credit.ID)
	if err != nil {
		return nil, err
	}

	siteEntity := credit.MapToSiteEntity()
	if site.Cron != credit.Cron {
		s.removeCronTask(site.TaskID.Int64)

		taskID, err := s.addCronTask(credit)
		if err != nil {
			return nil, err
		}
		siteEntity.TaskID = dbr.NewNullInt64(taskID)
	}

	err = s.siteRepo.Update(&siteEntity)
	if err != nil {
		return nil, err
	}

	siteRulesEntity := credit.SiteRules.MapToEntity()
	siteRulesEntity.SiteID = siteEntity.ID
	err = s.siteRulesRepo.Update(&siteRulesEntity)
	if err != nil {
		return nil, err
	}
	if credit.SiteInfo.ID > 0 && len(credit.SiteInfo.Title)>0 {
		siteInfo, err := s.siteInfoRepo.GetBySiteID(credit.ID)
		if err != nil {
			return nil, err
		}

		siteInfoEntity := credit.SiteInfo.MapToEntity()
		siteInfoEntity.SiteID = credit.ID
		if siteInfo.ID > 0 {
			err = s.siteInfoRepo.Update(&siteInfoEntity)
			if err != nil {
				return nil, err
			}
		}
	}

	if credit.SiteInfo.ID == 0 && len(credit.SiteInfo.Title)>0{
		siteInfoEntity := credit.SiteInfo.MapToEntity()
		siteInfoEntity.SiteID = credit.ID
		err = s.siteInfoRepo.Save(&siteInfoEntity)
		if err != nil {
			return nil, err
		}
	}
	return s.GetByID(siteEntity.ID)
}

func (s *siteService) Remove(id int64) error {
	site, err := s.siteRepo.GetByID(id)
	if err != nil {
		return err
	}
	if site.ID == 0 {
		return nil
	}

	info, err := s.siteInfoRepo.GetBySiteID(id)
	if err != nil {
		return err
	}
	if info.ID > 0 {
		err = s.siteInfoRepo.Remove(info.ID)
		if err != nil {
			return err
		}
	}

	rules, err := s.siteRulesRepo.GetBySiteID(id)
	if err != nil {
		return err
	}
	if rules.ID > 0 {
		err = s.siteRulesRepo.Remove(rules.ID)
		if err != nil {
			return err
		}
	}
	s.removeCronTask(site.TaskID.Int64)
	return s.siteRepo.Remove(site.ID)
}

func (s siteService) GetAll() ([]app.SiteDTO, error) {
	var response []app.SiteDTO
	siteList, err := s.siteRepo.GetAll()
	if err != nil {
		return nil, err
	}
	for _, item := range siteList {
		currentItem, err := s.GetByID(item.ID)
		if err != nil {
			return nil, err
		}
		response = append(response, *currentItem)
	}
	return response, nil
}

func (s *siteService) Process(siteID int64) (bool, error) {
	s.log.Infof("begin process siteID %v", siteID)
	site, err := s.GetByID(siteID)
	if err != nil {
		s.log.Error(err)
		return false, err
	}

	rssArr, err := s.rssService.GetRSS(site.URL)
	if err != nil {
		s.log.Errorf("can`t get rss %v %v",site.URL,err)
		return false, err
	}

	if site.SiteInfo.ID == 0 {
		siteInfo, err := s.rssService.GetSiteInfo(rssArr)
		if err != nil {
			s.log.Errorf("can`t get siteInfo",err)
			return false, err
		}
		site.SiteInfo = siteInfo
		_, err = s.Update(site)
		if err != nil {
			s.log.Errorf("can`t update siteEntity",err)
			return false, err
		}
	}

	articles, err := s.rssService.GetArticles(rssArr, site.SiteRules)
	if err != nil {
		return false, err
	}

	for _, item := range articles {

		// не самая экономная проверка
		isExist, err := s.articleService.IsExist(item.Title)
		if err != nil {
			s.log.Errorf("can`t check article in db %v %v", item.Title,err)
			return false, err
		}

		if isExist {
			continue
		}
		_, err = s.articleService.Save(site.ID, &item)
		if err != nil {
			return false, err
		}
	}
	s.log.Infof("process siteID %v finish ", siteID)
	return true, nil
}

func (s *siteService) InitTaskFromDB() (bool, error) {
	s.log.Info("InitTaskFromDB start")
	siteList, err := s.GetAll()
	if err != nil {
		return false, err
	}
	for i, site := range siteList {
		taskID, err := s.addCronTask(&siteList[i])
		if err != nil {
			s.log.Errorf("can`t add task %v", err)
		}
		siteEnt := site.MapToSiteEntity()
		siteEnt.TaskID = dbr.NewNullInt64(taskID)
		err = s.siteRepo.Update(&siteEnt)
		if err != nil {
			s.log.Errorf("can`t update entity %v", err)
		}
	}
	s.log.Info("InitTaskFromDB finish")
	return true, nil
}

func (s *siteService) addCronTask(site *app.SiteDTO) (int64, error) {
	entryID, err := s.scheduler.AddFunc(site.Cron, func() { s.Process(site.ID) })
	return int64(entryID), err
}

func (s *siteService) removeCronTask(taskID int64) {
	s.scheduler.Remove(cron.EntryID(taskID))
}
