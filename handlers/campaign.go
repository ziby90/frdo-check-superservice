package handlers

import (
	"errors"
	"fmt"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
)

type CampaignMain struct {
	Id                  uint      `json:"id"`   // Идентификатор
	UID                 *string   `json:"uid"`  // Идентификатор от организации
	Name                string    `json:"name"` // Наименование
	IdCampaignType      uint      `json:"id_campaign_type"`
	NameCampaignType    string    `json:"campaign_type_name"`
	IdCampaignStatus    uint      `json:"id_campaign_status"`
	NameCampaignStatus  string    `json:"campaign_status_name"`
	EducationForms      []uint    `json:"education_forms"`
	EducationFormsName  []string  `json:"education_forms_names"`
	EducationLevels     []uint    `json:"education_levels"`
	EducationLevelsName []string  `json:"education_levels_names"`
	YearStart           int64     `json:"year_start"` // Год начала компании
	YearEnd             int64     `json:"year_end"`   // Год окончания компании
	Created             time.Time `json:"created"`    // Дата создания

}
type CampaignResponse struct {
	Id                 uint      `json:"id"`   // Идентификатор
	UID                *string   `json:"uid"`  // Идентификатор от организации
	Name               string    `json:"name"` // Наименование
	IdCampaignType     uint      `json:"id_campaign_type"`
	NameCampaignType   string    `json:"campaign_type_name"`
	IdCampaignStatus   uint      `json:"id_campaign_status"`
	NameCampaignStatus string    `json:"campaign_status_name"`
	YearStart          int64     `json:"year_start"` // Год начала компании
	YearEnd            int64     `json:"year_end"`   // Год окончания компании
	Created            time.Time `json:"created"`    // Дата создания
}
type AddEndData struct {
	IdCampaign       uint      `json:"id_campaign"`
	IdEducationLevel uint      `json:"id_education_level"`
	IdEducationForm  uint      `json:"id_education_form"`
	EndDate          time.Time `json:"end_date"`
}

var CampaignSearchArray = []string{
	`name`,
}

func (result *Result) GetListCampaign() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaigns []digest.Campaign
	sortField := `year_start`
	sortOrder := `desc`
	if result.Sort.Field != `` {
		sortField = result.Sort.Field
	} else {
		result.Sort.Field = sortField
	}

	fmt.Print(result.Sort.Field, sortField)
	db := conn.Preload(`CampaignType`).Preload(`CampaignStatus`).Where(`id_organization=?`, result.User.CurrentOrganization.Id)
	for _, search := range result.Search {
		if service.SearchStringInSliceString(search[0], CampaignSearchArray) >= 0 {
			db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
		}
	}
	dbCount := db.Model(&campaigns).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Order(sortField + ` ` + sortOrder).Find(&campaigns)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Компании не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, campaign := range campaigns {

			c := map[string]interface{}{
				`id`:                   campaign.Id,
				`uid`:                  campaign.Uid,
				`name`:                 campaign.Name,
				`id_campaign_type`:     campaign.CampaignType.Id,
				`name_campaign_type`:   campaign.CampaignType.Name,
				`id_campaign_status`:   campaign.CampaignStatus.Id,
				`name_campaign_status`: campaign.CampaignStatus.Name,
				`code_campaign_status`: campaign.CampaignStatus.Code,
				`year_start`:           campaign.YearStart,
				`year_end`:             campaign.YearEnd,
				`created`:              campaign.Created,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Компании не найдены.`
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}

func (result *ResultList) GetShortListCampaign() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaigns []digest.Campaign
	sortField := `created`
	sortOrder := `asc`

	db := conn.Where(`id_organization=?`, result.User.CurrentOrganization.Id).Order(sortField + ` ` + sortOrder)

	if result.Search != `` {
		db = db.Where(`UPPER(name) like ?`, `%`+strings.ToUpper(result.Search)+`%`)
	}
	db = db.Find(&campaigns)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Компании не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, campaign := range campaigns {
			c := map[string]interface{}{
				`id`:   campaign.Id,
				`name`: campaign.Name,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Компании не найдены.`
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}

func (result *ResultInfo) GetEducationLevelCampaign(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign

	db := conn.Find(&campaign, ID)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Компания не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	var responses []interface{}
	if db.RowsAffected > 0 {
		var campEducLevels []digest.CampaignEducLevel
		db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducLevels)
		for index, _ := range campEducLevels {
			var educLevel digest.EducationLevel
			conn.First(&educLevel, campEducLevels[index].IdEducationLevel)
			responses = append(responses, map[string]interface{}{
				`id`:   educLevel.Id,
				`name`: educLevel.Name,
			})
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Компании не найдены.`
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}

func (result *ResultInfo) GetEducationFormCampaign(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign

	db := conn.Find(&campaign, ID)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Компания не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	var responses []interface{}
	if db.RowsAffected > 0 {
		var campEducForms []digest.CampaignEducForm
		db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducForms)
		for index, _ := range campEducForms {
			var educForm digest.EducationForm
			conn.First(&educForm, campEducForms[index].IdEducationForm)
			responses = append(responses, map[string]interface{}{
				`id`:   educForm.Id,
				`name`: educForm.Name,
			})
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Компании не найдены.`
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}

func (result *ResultInfo) GetInfoCampaign(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	db := conn.Find(&campaign, ID)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Компания не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		db = conn.Model(&campaign).Related(&campaign.CampaignType, `IdCampaignType`)
		db = conn.Model(&campaign).Related(&campaign.CampaignStatus, `IdCampaignStatus`)

		var campEducForms []digest.CampaignEducForm
		db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducForms)
		var campEducLevels []digest.CampaignEducLevel
		db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducLevels)

		c := CampaignMain{
			Id:                 campaign.Id,
			UID:                campaign.Uid,
			Name:               campaign.Name,
			IdCampaignType:     campaign.CampaignType.Id,
			NameCampaignType:   campaign.CampaignType.Name,
			IdCampaignStatus:   campaign.CampaignStatus.Id,
			NameCampaignStatus: campaign.CampaignStatus.Name,
			YearStart:          campaign.YearStart,
			YearEnd:            campaign.YearEnd,
			Created:            campaign.Created,
		}
		for _, campEducForm := range campEducForms {
			var educForm digest.EducationForm
			db = conn.Find(&educForm, campEducForm.IdEducationForm)
			c.EducationForms = append(c.EducationForms, educForm.Id)
			c.EducationFormsName = append(c.EducationFormsName, educForm.Name)
		}
		for _, campEducLevel := range campEducLevels {
			var educLevel digest.EducationLevel
			db = conn.Find(&educLevel, campEducLevel.IdEducationLevel)
			c.EducationLevels = append(c.EducationLevels, educLevel.Id)
			c.EducationLevelsName = append(c.EducationLevelsName, educLevel.Name)
		}
		result.Done = true
		result.Items = c
		return
	} else {
		result.Done = true
		message := `Компании не найдены.`
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}

func (result *ResultInfo) GetEndDateCampaign(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	db := conn.Find(&campaign, ID)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Компания не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	var endDate []digest.VEndApplication
	db = conn.Where(`id=? AND id_app_accept_phase IS NULL`, campaign.Id).Find(&endDate)

	if db.RowsAffected > 0 {
		var r []interface{}
		for index, _ := range endDate {
			r = append(r, map[string]interface{}{
				`id_end_application`:   endDate[index].IdEndApplication,
				`id_education_level`:   endDate[index].IdEducationLevel,
				`name_education_level`: endDate[index].EducationLevel,
				`id_education_form`:    endDate[index].IdEducationForm,
				`name_education_form`:  endDate[index].EducationForm,
				`end_date`:             endDate[index].EndDate,
				`order_end_app`:        endDate[index].OrderEndApp,
				`created`:              endDate[index].Created,
			})
		}
		result.Done = true
		result.Items = r
		return
	} else {
		message := `Не найдены даты. `
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}
func (result *ResultInfo) EditEndDateCampaign(data AddEndData) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	db := conn.Find(&campaign, data.IdCampaign)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Компания не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}

	var endDate digest.VEndApplication
	db = conn.Where(`id=? AND id_app_accept_phase IS NULL AND id_education_level=? AND id_education_form=?`, campaign.Id, data.IdEducationLevel, data.IdEducationForm).Find(&endDate)
	t := time.Now()
	var new digest.EndApplication
	if endDate.Id <= 0 {
		result.SetErrorResult(`Недопустимые значения`)
		return
	}
	if endDate.IdEndApplication == nil {
		new.Created = t
		new.IdCampaign = data.IdCampaign
		new.IdEducationForm = data.IdEducationForm
		new.IdEducationLevel = data.IdEducationLevel
		new.IdOrganization = result.User.CurrentOrganization.Id

	} else {
		var old digest.EndApplication
		db = conn.Where(`id=?`, endDate.IdEndApplication).Find(&old)
		new = old
		new.Changed = &t
	}
	new.Actual = true
	new.EndDate = data.EndDate
	db = conn.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&new)
	if db.Error != nil {
		m := db.Error.Error()
		result.Message = &m
		return
	}
	result.Done = true
	result.Items = map[string]interface{}{
		`id_end_application`: new.Id,
	}
}

func (result *ResultInfo) AddCampaign(campaignData CampaignMain, user digest.User) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var campaign digest.Campaign
	campaign.Organization.Id = user.CurrentOrganization.Id
	campaign.IdOrganization = user.CurrentOrganization.Id
	campaign.IdAuthor = user.Id
	campaign.Created = time.Now()
	campaign.Name = campaignData.Name
	if campaignData.UID != nil {
		var exist digest.Campaign
		tx.Where(`uid=? and id_organization=?`, campaignData.UID, user.CurrentOrganization.Id).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`У данной организации есть компания с данным UID`)
			tx.Commit()
			return
		}
		campaign.Uid = campaignData.UID
	}
	campaign.IdCampaignType = campaignData.IdCampaignType
	campaign.CampaignType.Id = campaignData.IdCampaignType
	db := tx.Find(&campaign.CampaignType, campaign.CampaignType.Id)
	if db.Error != nil || !campaign.CampaignType.Actual {
		result.SetErrorResult(`Не найден тип компании`)
		tx.Rollback()
		return
	}

	campaign.IdCampaignStatus = 1
	//campaign.CampaignStatus.Id = campaignData.IdCampaignStatus
	//// проверка типа
	//db = tx.Find(&campaign.CampaignStatus, campaign.CampaignStatus.Id)
	//if db.Error != nil || !campaign.CampaignStatus.Actual {
	//	result.SetErrorResult(`Статус комании не найден`)
	//	return
	//}

	campaign.YearEnd = campaignData.YearEnd
	// проверка года окончания
	if int(campaignData.YearEnd) < 1900 || int(campaignData.YearEnd) > time.Now().Year() {
		result.SetErrorResult(`Год окончания за пределами`)
		tx.Rollback()
		return
	}

	campaign.YearStart = campaignData.YearStart
	// проверка года начала
	if int(campaignData.YearStart) < 1900 || int(campaignData.YearStart) > time.Now().Year() {
		result.SetErrorResult(`Год начала за пределами`)
		tx.Rollback()
		return
	}

	if campaignData.YearStart > campaign.YearEnd {
		result.SetErrorResult(`Год начала не может быть позже года окончания`)
		tx.Rollback()
		return
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&campaign)
	if db.Error != nil {
		tx.Rollback()
		m := db.Error.Error()
		result.Message = &m
		return
	}

	for _, educLevelId := range campaignData.EducationLevels {
		var educationLevel digest.EducationLevel
		tx.Find(&educationLevel, educLevelId)
		if !educationLevel.Actual {
			result.SetErrorResult(`Уровень образования не найден`)
			tx.Rollback()
			return
		}
		row := conn.Table(`cls.edu_levels_campaign_types`).Where(`id_campaign_types=? AND id_education_level=?`, campaign.CampaignType.Id, educLevelId).Select(`id`).Row()
		var idEducLevelCampaignType uint
		err := row.Scan(&idEducLevelCampaignType)
		if err != nil || idEducLevelCampaignType <= 0 {
			result.SetErrorResult(`Данный уровень образования не соответствует типу приемной компании`)
			tx.Rollback()
			return
		}
		campaignEducLevel := digest.CampaignEducLevel{
			IdCampaign:       campaign.Id,
			IdEducationLevel: educLevelId,
			IdOrganization:   user.CurrentOrganization.Id,
		}
		db = tx.Create(&campaignEducLevel)
	}

	for _, educFormId := range campaignData.EducationForms {
		var educationForm digest.EducationForm
		tx.Find(&educationForm, educFormId)
		if !educationForm.Actual {
			result.SetErrorResult(`Форма образования не найдена`)
			tx.Rollback()
			return
		}
		campaignEducForm := digest.CampaignEducForm{
			IdCampaign:      campaign.Id,
			IdEducationForm: educFormId,
			IdOrganization:  user.CurrentOrganization.Id,
		}
		db = tx.Create(&campaignEducForm)
	}

	result.Items = campaign.Id
	result.Done = true
	tx.Commit()
}

func CheckCampaignByUser(idCampaign uint, user digest.User) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var count int
	if user.Role.Code == `administrator` {
		return nil
	}
	db := conn.Table(`cmp.campaigns`).Select(`id`).Where(`id_organization=? AND id=?`, user.CurrentOrganization.Id, idCampaign).Count(&count)
	if db.Error != nil {
		return db.Error
	}
	if count > 0 {
		return nil
	} else {
		return errors.New(`У пользователя нет доступа к данной компании `)
	}
}
