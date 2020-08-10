package handlers

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
)

var AppAchievementsSearchArray = []string{
	`name`,
}

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
		if db.Error.Error() == service.ErrorDbNotFound {
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
func (result *Result) GetAchievementsByApplicationId(idApplication uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var items []digest.AppAchievements

	db := conn.Where(`id_application=?`, idApplication)
	for _, search := range result.Search {
		db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
	}
	dbCount := db.Model(&items).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {
	}

	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Order(`id asc`).Preload(`IndividualAchievement`).Preload(`AchievementCategory`).Find(&items)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
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
		for index, _ := range items {
			var doc map[string]interface{}
			if items[index].IdDocument != nil {
				var document digest.VDocuments
				db = conn.Preload(`DocumentType`).Where(`id_document=?`, items[index].IdDocument).Find(&document)
				doc = map[string]interface{}{
					`id`:            document.IdDocument,
					`name_category`: document.NameSysCategories,
					`code_category`: document.NameTable,
					`name_type`:     document.DocumentType.Name,
				}
			}
			var files []digest.AchievementFiles
			db = conn.Where(`id_achievement=? AND id_application=?`, items[index].Id, idApplication).Find(&files)

			c := map[string]interface{}{
				`id`:             items[index].Id,
				`uid`:            items[index].Uid,
				`id_category`:    items[index].AchievementCategory.Id,
				`name_category`:  items[index].AchievementCategory.Name,
				`max_value`:      items[index].IndividualAchievement.MaxValue,
				`id_achievement`: items[index].IndividualAchievement.Id,
				`mark`:           items[index].Mark,
				`document`:       doc,
			}
			if len(files) > 0 {
				c[`file`] = true
				var epguDocuments []interface{}
				for _, value := range files {
					epguDocuments = append(epguDocuments, map[string]interface{}{
						`id`:        value.Id,
						`path_file`: value.PathFile,
					})
				}
				c[`epgu_documents`] = epguDocuments
			} else {
				c[`file`] = false
				c[`epgu_documents`] = []string{}
			}
			if items[index].UidEpgu != nil {
				c[`name_achievement`] = items[index].Name
			} else {
				c[`name_achievement`] = items[index].IndividualAchievement.Name
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
		result.Items = []digest.Campaign{}
		return
	}
}
func (result *Result) GetListAchievementByCompanyId(idCampaign uint) {
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
	db = db.Where(`id_campaign=?`, idCampaign)
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
		if db.Error.Error() == service.ErrorDbNotFound {
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
		if db.Error.Error() == service.ErrorDbNotFound {
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
		var campaign digest.Campaign
		db = conn.Where(`id=?`, achievement.IdCampaign).Find(&campaign)
		c := map[string]interface{}{
			`id`:            achievement.Id,
			`uid`:           achievement.Uid,
			`name`:          achievement.Name,
			`id_campaign`:   achievement.IdCampaign,
			`name_campaign`: campaign.Name,
			`id_category`:   achievement.AchievementCategory.Id,
			`name_category`: achievement.AchievementCategory.Name,
			`max_value`:     achievement.MaxValue,
			`created`:       achievement.Created,
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
	achievement.Actual = true

	var campaign digest.Campaign
	db := tx.Preload(`CampaignType`).Find(&campaign, achData.IdCampaign)
	if campaign.Id < 1 {
		result.SetErrorResult(`Компания не найдена`)
		tx.Rollback()
		return
	}
	err := CheckAddAchievements(campaign.Id)
	if err != nil {
		result.SetErrorResult(err.Error())
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
	fmt.Println(achievement.Id)
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&achievement)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}

	result.Items = achievement
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) EditAchievement(data AchievementMain) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var achievement digest.IndividualAchievements

	db := tx.Where(`id=?`, data.Id).Find(&achievement)
	if achievement.Id <= 0 {
		result.SetErrorResult(`Достижение не найдено.`)
		tx.Rollback()
		return
	}
	if achievement.IdOrganization != result.User.CurrentOrganization.Id {
		result.SetErrorResult(`Достижение принадлежит другой организации.`)
		tx.Rollback()
		return
	}
	err := CheckEditAchievements(achievement.IdCampaign)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	achievement.Name = data.Name
	if data.UID != nil && data.UID != achievement.Uid && *data.UID != `` {
		var exist digest.IndividualAchievements
		db = tx.Where(`uid=? AND id_organization=? AND id!=?`, data.UID, result.User.CurrentOrganization.Id, achievement.Id).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Достижение с данным uid уже существует у данной организации.`)
			tx.Rollback()
			return
		}
		achievement.Uid = data.UID
	}
	if achievement.IdCategory != data.IdCategory {
		var category digest.AchievementCategory
		db = tx.Find(&category, data.IdCategory)
		if category.Id < 1 {
			result.SetErrorResult(`Категория не найдена`)
			tx.Rollback()
			return
		}
		achievement.IdCategory = data.IdCategory
	}

	achievement.MaxValue = data.MaxValue
	achievement.Name = data.Name
	t := time.Now()
	achievement.Changed = &t
	achievement.IdAuthor = result.User.Id
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&achievement)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_achievement`: achievement.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) RemoveAchievement(idAchievement uint) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var old digest.IndividualAchievements
	db := tx.Find(&old, idAchievement)
	if old.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Индивидуальное достижение не найдено`)
		tx.Rollback()
		return
	}
	err := CheckEditAchievements(old.IdCampaign)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}

	db = tx.Where(`id=?`, idAchievement).Delete(&old)
	if db.Error == nil {
		result.Done = true
		tx.Commit()
		result.Items = map[string]interface{}{
			`id_achievement`: idAchievement,
		}
		return
	} else {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
}
func (result *ResultList) GetAchievementsSelectListByCompetitive(idCompetitive uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var competitive digest.CompetitiveGroup
	db := conn.Where(`id=?`, idCompetitive).Find(&competitive)
	if competitive.Id <= 0 {
		result.Done = true
		message := `Конкурсная группа не найдена.`
		result.Message = &message
		return
	}
	var items []digest.IndividualAchievements
	sortField := `created`
	sortOrder := `desc`
	db = conn.Where(`id_campaign=? AND actual`, competitive.IdCampaign).Order(sortField + ` ` + sortOrder)

	if result.Search != `` {
		db = db.Where(`UPPER(name) like ?`, `%`+strings.ToUpper(result.Search)+`%`)
	}

	db = db.Find(&items)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
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
		for _, item := range items {
			c := map[string]interface{}{
				`id`:        item.Id,
				`name`:      item.Name,
				`max_value`: item.MaxValue,
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
		result.Items = []digest.IndividualAchievements{}
		return
	}
}
func (result *ResultList) GetAchievementsSelectListByCampaign(idCampaign uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)

	var items []digest.IndividualAchievements
	sortField := `created`
	sortOrder := `desc`
	db := conn.Where(`id_campaign=? AND actual`, idCampaign).Order(sortField + ` ` + sortOrder)

	if result.Search != `` {
		db = db.Where(`UPPER(name) like ?`, `%`+strings.ToUpper(result.Search)+`%`)
	}

	db = db.Find(&items)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
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
		for _, item := range items {
			c := map[string]interface{}{
				`id`:        item.Id,
				`name`:      item.Name,
				`max_value`: item.MaxValue,
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
		result.Items = []digest.IndividualAchievements{}
		return
	}
}

func (result *ResultInfo) GetFileAppAchievement(ID uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var doc digest.AchievementFiles
	db := conn.Where(`id=?`, ID).Find(&doc)
	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			result.Done = false
			message := "Файл не найден."
			result.Message = &message
			result.Items = []interface{}{}
			return
		}
		message := "Ошибка подключения к БД."
		result.Message = &message
		return
	}

	conn.Where(`id=?`, *doc.IdAchievement).Find(&doc.Achievement)
	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			result.Done = false
			message := "Файл не найден."
			result.Message = &message
			result.Items = []interface{}{}
			return
		}
		message := "Ошибка подключения к БД."
		result.Message = &message
		return
	}

	conn.Where(`id=?`, *doc.IdApplication).Find(&doc.Application)
	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			result.Done = false
			message := "Файл не найден."
			result.Message = &message
			result.Items = []interface{}{}
			return
		}
		message := "Ошибка подключения к БД."
		result.Message = &message
		return
	}
	if doc.PathFile != nil {
		filename := *doc.PathFile
		path := getPath(doc.Application.EntrantsId, `app.achievements`, doc.Achievement.Created) + filename
		//fmt.Println("Test Test Test", doc.Achievement);
		//f, err := os.Open(path)
		//if err != nil {
		//	result.SetErrorResult(err.Error())
		//	return
		//} else {
		//	defer f.Close()
		//	reader := bufio.NewReader(f)
		//	content, _ := ioutil.ReadAll(reader)
		//	ext := mimetype.Detect(content)
		//	file := digest.FileC{
		//		Content: content,
		//		Size:    int64(len(content)),
		//		Title:   filename,
		//		Type:    ext.Extension(),
		//	}
		//	result.Items = file
		//}
		result.Items = path
	} else {
		message := "Файл не найден."
		result.Message = &message
		return
	}
	result.Done = true
	return
}
