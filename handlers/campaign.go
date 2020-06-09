package handlers

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
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
	NameCampaignType    string    `json:"name_campaign_type"`
	IdCampaignStatus    uint      `json:"id_campaign_status"`
	NameCampaignStatus  string    `json:"name_campaign_status"`
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
type ChangeStatusCampaign struct {
	Campaign         digest.Campaign       `gorm:"foreignkey:IdCampaign"`
	IdCampaign       uint                  `json:"id_campaign"`
	CampaignStatus   digest.CampaignStatus `gorm:"foreignkey:IdCampaignStatus"`
	IdCampaignStatus *uint                 `json:"id_campaign_status"`
	CodeStatus       string                `json:"code"`
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
		var educationForms []uint
		var educationFormsName []string
		for _, campEducForm := range campEducForms {
			var educForm digest.EducationForm
			db = conn.Find(&educForm, campEducForm.IdEducationForm)
			educationForms = append(educationForms, educForm.Id)
			educationFormsName = append(educationFormsName, educForm.Name)
		}
		if len(educationForms) > 0 {
			c[`education_forms`] = educationForms
			c[`education_forms_names`] = educationFormsName
		}

		var educationLevels []uint
		var educationLevelsName []string
		for _, campEducLevel := range campEducLevels {
			var educLevel digest.EducationLevel
			db = conn.Find(&educLevel, campEducLevel.IdEducationLevel)
			educationLevels = append(educationLevels, educLevel.Id)
			educationLevelsName = append(educationLevelsName, educLevel.Name)
		}
		if len(educationLevels) > 0 {
			c[`education_levels`] = educationLevels
			c[`education_levels_names`] = educationLevelsName
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
	err := CheckEditEndDate(data.IdCampaign)
	if err != nil {
		result.SetErrorResult(err.Error())
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
	year := data.EndDate.Year()
	if year > int(campaign.YearEnd) || year < int(campaign.YearStart) {
		result.SetErrorResult(`Контрольная дата должна назодиться в диапозоне проведения приемной компании `)
		return
	}

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
func (result *ResultInfo) EditCampaign(data CampaignMain) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var campaign digest.Campaign
	err := CheckEditCampaign(data.Id)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	db := tx.Where(`id=? AND actual`, data.Id).Find(&campaign)
	if campaign.Id <= 0 {
		result.SetErrorResult(`Приемная компания не найдена`)
		tx.Rollback()
		return
	}

	if campaign.IdOrganization != result.User.CurrentOrganization.Id {
		result.SetErrorResult(`Приемная компания принадлежит другой организации.`)
		tx.Rollback()
		return
	}

	campaign.Name = strings.TrimSpace(data.Name)
	if data.UID != nil && data.UID != campaign.Uid && *data.UID != `` {
		var exist digest.Campaign
		db = tx.Where(`uid=? AND id_organization=? AND id!=?`, data.UID, result.User.CurrentOrganization.Id, campaign.Id).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Приемная компания с данным uid уже существует у данной организации.`)
			tx.Rollback()
			return
		}
		campaign.Uid = data.UID
	}

	if campaign.IdCampaignType != data.IdCampaignType {
		var category digest.CampaignType
		db = tx.Find(&category, data.IdCampaignType)
		if category.Id < 1 {
			result.SetErrorResult(`Тип приемной компании не найден`)
			tx.Rollback()
			return
		}
		campaign.IdCampaignType = data.IdCampaignType
	}
	t := time.Now()
	campaign.Changed = &t

	// проверка года окончания
	if int(data.YearEnd) < 1900 || int(data.YearEnd) > time.Now().Year() {
		result.SetErrorResult(`Год окончания за пределами`)
		tx.Rollback()
		return
	}
	campaign.YearEnd = data.YearEnd
	// проверка года начала
	if int(data.YearStart) < 1900 || int(data.YearStart) > time.Now().Year() {
		result.SetErrorResult(`Год начала за пределами`)
		tx.Rollback()
		return
	}
	campaign.YearStart = data.YearStart
	if data.YearStart > campaign.YearEnd {
		result.SetErrorResult(`Год начала не может быть позже года окончания`)
		tx.Rollback()
		return
	}
	if len(data.EducationLevels) > 0 {
		db = tx.Where(`id_campaign=?`, campaign.Id).Delete(digest.CampaignEducLevel{})
		for _, educLevelId := range data.EducationLevels {
			var educationLevel digest.EducationLevel
			tx.Find(&educationLevel, educLevelId)
			if !educationLevel.Actual {
				result.SetErrorResult(`Уровень образования не найден`)
				tx.Rollback()
				return
			}
			row := conn.Table(`cls.edu_levels_campaign_types`).Where(`id_campaign_types=? AND id_education_level=?`, campaign.IdCampaignType, educLevelId).Select(`id`).Row()
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
				IdOrganization:   result.User.CurrentOrganization.Id,
			}
			db = tx.Create(&campaignEducLevel)
		}
	}

	if len(data.EducationForms) > 0 {
		db = tx.Where(`id_campaign=?`, campaign.Id).Delete(digest.CampaignEducForm{})
		for _, educFormId := range data.EducationForms {
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
				IdOrganization:  result.User.CurrentOrganization.Id,
			}
			db = tx.Create(&campaignEducForm)
		}
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&campaign)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
	result.Done = true
	result.Items = map[string]interface{}{
		`id_campaign`: campaign.Id,
	}
	tx.Commit()
	return

}
func (result *ResultInfo) SetStatusCampaign(data ChangeStatusCampaign) {
	conn := &config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	if data.CodeStatus == `` {
		result.SetErrorResult(`Пустой статус`)
		tx.Rollback()
		return
	}
	status, err := GetCampaignStatusByCode(data.CodeStatus)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	campaign, err := digest.GetCampaign(data.IdCampaign)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}

	if campaign.IdCampaignStatus == status.Id {
		result.SetErrorResult(`Приемная компания уже в этом статусе`)
		tx.Rollback()
		return
	}
	campaign.IdCampaignStatus = status.Id
	db := tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&campaign)
	if db.Error != nil {
		result.SetErrorResult(`Ошибка при изменении статуса применой компании ` + db.Error.Error())
		tx.Rollback()
		return
	}
	//if campaign.campaign != nil {
	//	sendToEpgu.InitConnect(config.Db.ConnGORM, config.Db.ConnSmevGorm)
	//	err = sendToEpgu.PrepareSendStatementResponse(*campaign.UidEpgu, sendToEpgu.NewApplication)
	//	fmt.Println(err)
	//}
	result.Items = map[string]interface{}{
		`id_campaign`: campaign.Id,
		`new_status`:  campaign.IdCampaignStatus,
	}
	result.Done = true
	tx.Commit()
	return
}
func (result *ResultInfo) GetCampaignStatuses(keys map[string][]string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var db *gorm.DB
	var statuses []digest.CampaignStatus
	//db = conn.Select(`id, name`).Table(`cls.v_direction_specialty`)
	db = conn.Where(`actual`)
	if len(keys[`search`]) > 0 {
		db = db.Where(`UPPER(name) LIKE ?`, `%`+strings.ToUpper(keys[`search`][0])+`%`)
	}

	db = db.Find(&statuses)
	if db.Error != nil {
		message := db.Error.Error()
		result.Message = &message
		return
	}
	var response []interface{}
	for index, _ := range statuses {
		response = append(response, map[string]interface{}{
			"id":   statuses[index].Id,
			"name": statuses[index].Name,
			"code": statuses[index].Code,
		})
	}
	result.Done = true
	result.Items = response
	return
}
func GetCampaignStatusByCode(code string) (*digest.CampaignStatus, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item digest.CampaignStatus
	db := conn.Where(`code=?`, code).Find(&item)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Статус не найден. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if item.Id <= 0 {
		return nil, errors.New(`Статус не найден. `)
	}
	return &item, nil
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

	var exist []digest.Campaign
	tx.Where(`id_campaign_type=? AND year_start=? AND id_organization=?`, campaignData.IdCampaignType, campaignData.YearStart, user.CurrentOrganization.Id).Find(&exist)
	if len(exist) > 0 {
		result.SetErrorResult(`У данной организации уже есть приемная компания с заданным типом и годом начала`)
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
