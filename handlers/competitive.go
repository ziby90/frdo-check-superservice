package handlers

import (
	"github.com/jinzhu/gorm"
	"persons/config"
	"persons/digest"
	"persons/service"
)

func (result *Result) GetListCompetitiveGroupsByCompanyId(campaignId uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var competitiveGroups []digest.CompetitiveGroup
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
	dbCount := db.Model(&competitiveGroups).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()

	db = db.Preload(`EducationLevel`).Preload(`ReceptionCampaign`).Preload(`LevelBudget`).Preload(`EducationSource`).Preload(`EducationForm`).Preload(`Direction`).Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&competitiveGroups)
	var responses []interface{}
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
		for _, competitveGroup := range competitiveGroups {
			c := map[string]interface{}{
				`id`:                    competitveGroup.Id,
				`name`:                  competitveGroup.Name,
				`id_campaign`:           competitveGroup.IdCampaign,
				`name_campaign`:         competitveGroup.ReceptionCampaign.Name,
				`year_start_campaign`:   competitveGroup.ReceptionCampaign.YearStart,
				`year_end_campaign`:     competitveGroup.ReceptionCampaign.YearEnd,
				`id_education_form`:     competitveGroup.EducationForm.Id,
				`name_education_form`:   competitveGroup.EducationForm.Name,
				`id_education_source`:   competitveGroup.EducationSource.Id,
				`name_education_source`: competitveGroup.EducationSource.Name,
				`id_education_level`:    competitveGroup.EducationLevel.Id,
				`name_education_level`:  competitveGroup.EducationLevel.Name,
				`id_direction`:          competitveGroup.Direction.Id,
				`name_direction`:        competitveGroup.Direction.Name,
				`id_level_budget`:       competitveGroup.LevelBudget.Id,
				`name_level_budget`:     competitveGroup.LevelBudget.Name,
				`id_author`:             competitveGroup.Id,
				`actual`:                competitveGroup.Id,
				`uid`:                   competitveGroup.Id,
				`id_organization`:       competitveGroup.Id,
				`created`:               competitveGroup.Id,
				`changed`:               competitveGroup.Id,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Достижения не найдены.`
		result.Message = &message
		result.Items = []digest.CompetitiveGroup{}
		return
	}
}

//func (result *ResultInfo) GetInfoAchievement(ID uint) {
//	result.Done = false
//	conn := config.Db.ConnGORM
//	conn.LogMode(config.Conf.Dblog)
//	var achievement digest.IndividualAchievements
//	db := conn.Find(&achievement, ID)
//	if db.Error != nil {
//		if db.Error.Error() == `record not found` {
//			result.Done = true
//			message := `Достижение не найдено.`
//			result.Message = &message
//			return
//		}
//		message := `Ошибка подключения к БД. `
//		result.Message = &message
//		return
//	}
//	if db.RowsAffected > 0 {
//
//		db = conn.Model(&achievement).Related(&achievement.AchievementCategory, `IdCategory`)
//		c := AchievementMain{
//			Id:           achievement.Id,
//			UID:          achievement.Uid,
//			Name:         achievement.Name,
//			IdCampaign:   achievement.IdCampaign,
//			IdCategory:   achievement.AchievementCategory.Id,
//			NameCategory: achievement.AchievementCategory.Name,
//			MaxValue:     achievement.MaxValue,
//			Created:      achievement.Created,
//		}
//		result.Done = true
//		result.Items = c
//		return
//	} else {
//		result.Done = true
//		message := `Достижения не найдены.`
//		result.Message = &message
//		result.Items = make(map[string]string)
//		return
//	}
//}
