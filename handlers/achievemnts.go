package handlers

import (
	"github.com/jinzhu/gorm"
	"persons/config"
	"persons/digest"
	"persons/service"
	"time"
)

type AchievementMain struct {
	Id           uint      `json:"id"`   // Идентификатор
	UID          *string   `json:"uid"`  // Идентификатор от организации
	Name         string    `json:"name"` // Наименование
	IdCampaign   uint      `json:"id_campaign"`
	IdCategory   uint      `json:"id_category"`
	NameCategory string    `json:"category_name"`
	MaxValue     int64     `json:"max_value"`
	Created      time.Time `json:"created"` // Дата создания
}

type AchievementResponse struct {
	Id           uint      `json:"id"`   // Идентификатор
	UID          *string   `json:"uid"`  // Идентификатор от организации
	Name         string    `json:"name"` // Наименование
	IdCampaign   uint      `json:"id_campaign"`
	IdCategory   uint      `json:"id_category"`
	NameCategory string    `json:"category_name"`
	MaxValue     int64     `json:"max_value"`
	Created      time.Time `json:"created"` // Дата создания
}

var sortByArray = []string{
	`uid`,
	`name`,
	`id_category`,
	`max_value`,
}
var orderArray = []string{
	`asc`,
	`desc`,
}

func (result *Result) GetListAchievement() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var achievements []digest.IndividualAchievements
	var db *gorm.DB
	if service.SearchStringInSliceString(result.Sort.Field, sortByArray) > 0 && service.SearchStringInSliceString(result.Sort.Order, orderArray) > 0 {
		db = conn.Order(result.Sort.Field + ` ` + result.Sort.Order)
	} else {
		db = conn.Order(`created asc `)
	}
	//if result.Search != `` {
	//	db = db.Where(`UPPER(name) LIKE ?`, `%`+strings.ToUpper(result.Search)+`%`)
	//}
	dbCount := db.Model(&achievements).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&achievements)
	var achievementResponse []AchievementResponse
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Достижения не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, achievement := range achievements {
			db = conn.Model(&achievement).Related(&achievement.AchievementCategory, `IdCategory`)
			c := AchievementResponse{
				Id:           achievement.Id,
				UID:          achievement.Uid,
				Name:         achievement.Name,
				IdCampaign:   achievement.IdCampaign,
				IdCategory:   achievement.AchievementCategory.Id,
				NameCategory: achievement.AchievementCategory.Name,
				MaxValue:     achievement.MaxValue,
				Created:      achievement.Created,
			}
			achievementResponse = append(achievementResponse, c)
		}
		result.Done = true
		result.Items = achievementResponse
		return
	} else {
		result.Done = true
		message := `Достижения не найдены.`
		result.Message = &message
		result.Items = []digest.IndividualAchievements{}
		return
	}
}

func (result *Result) GetListAchievementByCompanyId(campaignId uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var achievements []digest.IndividualAchievements
	var db *gorm.DB
	if service.SearchStringInSliceString(result.Sort.Field, sortByArray) >= 0 {
		order := `asc`
		if service.SearchStringInSliceString(result.Sort.Order, orderArray) >= 0 {
			order = result.Sort.Order
		}
		db = conn.Order(result.Sort.Field + ` ` + order)
	} else {
		db = conn.Order(`created asc `)
	}
	db = db.Where(`id_campaign=?`, campaignId)
	//if result.Search != `` {
	//	db = db.Where(`UPPER(name) LIKE ?`, `%`+strings.ToUpper(result.Search)+`%`)
	//}
	dbCount := db.Model(&achievements).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()

	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&achievements)
	var achievementResponse []AchievementResponse
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Достижения не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, achievement := range achievements {
			db = conn.Model(&achievement).Related(&achievement.AchievementCategory, `IdCategory`)
			c := AchievementResponse{
				Id:           achievement.Id,
				UID:          achievement.Uid,
				Name:         achievement.Name,
				IdCampaign:   achievement.IdCampaign,
				IdCategory:   achievement.AchievementCategory.Id,
				NameCategory: achievement.AchievementCategory.Name,
				MaxValue:     achievement.MaxValue,
				Created:      achievement.Created,
			}
			achievementResponse = append(achievementResponse, c)
		}
		result.Done = true
		result.Items = achievementResponse
		return
	} else {
		result.Done = true
		message := `Достижения не найдены.`
		result.Message = &message
		result.Items = []AchievementResponse{}
		return
	}
}

func (result *ResultInfo) GetInfoAchievement(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var achievement digest.IndividualAchievements
	db := conn.Find(&achievement, ID)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Достижение не найдено.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {

		db = conn.Model(&achievement).Related(&achievement.AchievementCategory, `IdCategory`)
		c := AchievementMain{
			Id:           achievement.Id,
			UID:          achievement.Uid,
			Name:         achievement.Name,
			IdCampaign:   achievement.IdCampaign,
			IdCategory:   achievement.AchievementCategory.Id,
			NameCategory: achievement.AchievementCategory.Name,
			MaxValue:     achievement.MaxValue,
			Created:      achievement.Created,
		}
		result.Done = true
		result.Items = c
		return
	} else {
		result.Done = true
		message := `Достижения не найдены.`
		result.Message = &message
		result.Items = make(map[string]string)
		return
	}
}

func (result *ResultInfo) AddAchievement(achData digest.IndividualAchievements, user digest.User) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var achievement digest.IndividualAchievements
	achievement = achData
	achievement.Organization.Id = user.CurrentOrganization.Id
	achievement.IdOrganization = user.CurrentOrganization.Id
	achievement.IdAuthor = user.Id
	achievement.Created = time.Now()

	var campaign digest.Campaign
	db := tx.Preload(`CampaignType`).Find(&campaign, achData.IdCampaign)
	if campaign.Id < 1 {
		result.SetErrorResult(`Компания не найдена`)
		tx.Rollback()
		return
	}
	achievement.IdCampaign = achData.IdCampaign

	var category digest.AchievementCategory
	db = tx.Find(&category, achData.AchievementCategory)
	if category.Id < 1 {
		result.SetErrorResult(`Категория не найдена`)
		tx.Rollback()
		return
	}
	achievement.IdCategory = achData.IdCategory
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&achievement)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}

	result.Items = achievement
	result.Done = true
	tx.Commit()
}
