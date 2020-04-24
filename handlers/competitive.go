package handlers

import (
	"github.com/jinzhu/gorm"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
)

type AddCompetitiveGroup struct {
	Id                uint                   `gorm:"primary_key" json:"id"` // Идентификатор
	UID               *string                `xml:"UID" json:"uid" `        // Идентификатор от организации
	Direction         digest.Direction       `gorm:"foreignkey:IdDirection" json:"-"`
	IdDirection       uint                   `json:"id_direction"` // Идентификатор направления
	Name              string                 `json:"name"`
	EducationForm     digest.EducationForm   `gorm:"foreignkey:IdEducationForm" json:"-"`
	IdEducationForm   uint                   `json:"id_education_form"`
	EducationLevel    digest.EducationLevel  `gorm:"foreignkey:IdEducationLevel" json:"-"`
	IdEducationLevel  uint                   `json:"id_education_level"`
	EducationSource   digest.EducationSource `gorm:"foreignkey:IdEducationSource" json:"-"`
	IdEducationSource uint                   `json:"id_education_source"`
	LevelBudget       digest.LevelBudget     `gorm:"foreignkey:IdLevelBudget" json:"-"`
	IdLevelBudget     *uint                  `json:"id_level_budget"`
	Campaign          digest.Campaign        `gorm:"foreignkey:IdCampaign" json:"-"`
	IdCampaign        uint                   `json:"id_campaign"`
	Number            int64                  `json:"number"`
}

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

func (result *ResultInfo) AddCompetitive(data AddCompetitiveGroup) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var competitive digest.CompetitiveGroup
	competitive.Organization.Id = result.User.CurrentOrganization.Id
	competitive.IdOrganization = result.User.CurrentOrganization.Id
	competitive.IdAuthor = result.User.Id
	competitive.Created = time.Now()
	competitive.IdEducationLevel = data.IdEducationLevel
	competitive.Name = strings.TrimSpace(data.Name)

	var campaign digest.Campaign
	db := tx.Preload(`CampaignType`).Find(&campaign, data.IdCampaign)
	if campaign.Id < 1 {
		result.SetErrorResult(`Компания не найдена`)
		tx.Rollback()
		return
	}
	row := conn.Table(`cls.edu_levels_campaign_types`).Where(`id_campaign_types=? AND id_education_level=?`, campaign.CampaignType.Id, data.IdEducationLevel).Select(`id`).Row()
	var idEducLevelCampaignType uint
	err := row.Scan(&idEducLevelCampaignType)
	if err != nil || idEducLevelCampaignType <= 0 {
		result.SetErrorResult(`Данный уровень образования не соответствует типу приемной компании`)
		tx.Rollback()
		return
	}
	competitive.IdCampaign = data.IdCampaign
	var direction digest.Direction
	db = tx.Where(`id_education_level=?`, data.IdEducationLevel).Find(&direction, data.IdDirection)
	if direction.Id < 1 {
		result.SetErrorResult(`Специальность не найдена`)
		tx.Rollback()
		return
	}
	competitive.IdDirection = data.IdDirection

	var educForm digest.EducationForm
	db = tx.Find(&educForm, data.IdEducationForm)
	if direction.Id < 1 {
		result.SetErrorResult(`Форма образования не найдена`)
		tx.Rollback()
		return
	}
	competitive.IdEducationForm = data.IdEducationForm

	var educSource digest.EducationSource
	db = tx.Find(&educSource, data.IdEducationSource)
	if direction.Id < 1 {
		result.SetErrorResult(`Форма образования не найдена`)
		tx.Rollback()
		return
	}
	competitive.IdEducationSource = data.IdEducationSource

	if data.IdLevelBudget != nil {
		var levelBudget digest.LevelBudget
		db = tx.Find(&levelBudget, *data.IdLevelBudget)
		if direction.Id < 1 {
			result.SetErrorResult(`Уровень бюджета не найден`)
			tx.Rollback()
			return
		}
		competitive.IdLevelBudget = data.IdLevelBudget
	}

	var exist digest.CompetitiveGroup
	tx.Where(`id_direction=? AND id_campaign=? AND id_education_level=? AND id_education_source=? AND id_education_form=?`, data.IdDirection, data.IdCampaign, data.IdEducationLevel, data.IdEducationSource, data.IdEducationForm).Find(&exist)
	if exist.Id > 0 {
		result.SetErrorResult(`Подобная конкурсная группа уже существует.`)
		tx.Rollback()
		return
	}
	if data.UID != nil {
		var exist digest.CompetitiveGroup
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *data.UID).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`КОнкурсная группа с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		competitive.UID = data.UID
	}
	competitive.Actual = true

	switch competitive.IdEducationSource {
	case 1: // Бюджетные места
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			competitive.BudgetO = data.Number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.BudgetOz = data.Number
			break
		case 3: // Заочная
			competitive.BudgetZ = data.Number
			break
		default:
			result.SetErrorResult(`Ошибка`)
			tx.Rollback()
			return
		}
		break
	case 2: // Квота приема лиц, имеющих особое право
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			competitive.QuotaO = data.Number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.QuotaOz = data.Number
			break
		case 3: // Заочная
			competitive.QuotaZ = data.Number
			break
		default:
			result.SetErrorResult(`Ошибка`)
			tx.Rollback()
			return
		}
		break
	case 3: // С оплатой обучения
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			competitive.PaidO = data.Number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.PaidOz = data.Number
			break
		case 3: // Заочная
			competitive.PaidZ = data.Number
			break
		default:
			result.SetErrorResult(`Ошибка`)
			tx.Rollback()
			return
		}
		break
	case 4: // Целевой прием
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			competitive.TargetO = data.Number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.TargetOz = data.Number
			break
		case 3: // Заочная
			competitive.TargetZ = data.Number
			break
		default:
			result.SetErrorResult(`Ошибка`)
			tx.Rollback()
			return
		}
	default:
		result.SetErrorResult(`Ошибка`)
		tx.Rollback()
		return
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&competitive)
	if db.Error == nil {
		result.Items = competitive.Id
		result.Done = true
		tx.Commit()
	} else {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
}
