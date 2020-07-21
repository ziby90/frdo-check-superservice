package handlers

import (
	"github.com/jinzhu/gorm"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strconv"
	"strings"
	"time"
)

var CompetitiveSearchArray = []string{
	`name`,
}

type CompetitiveGroup struct {
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
	Number            *int64                 `json:"number"`
	Comment           *string                `json:"comment"`
}

type DirectionCompetitiveGroups struct {
	IdDirection uint `json:"id_direction"`
}

type AddCompetitiveGroup struct {
	CompetitiveGroup         CompetitiveGroup                 `json:"competitive"`
	CompetitiveGroupPrograms []digest.CompetitiveGroupProgram `json:"education_programs"`
	EntranceTests            []digest.EntranceTest            `json:"entrance_tests"`
}
type EditNumberCompetitive struct {
	IdCompetitive uint
	Number        int64 `json:"number"`
}
type AddEntrance struct {
	EntranceTests []digest.EntranceTest `json:"entrance_tests"`
}

type AddCompetitiveGroupPrograms struct {
	CompetitiveGroupPrograms []digest.CompetitiveGroupProgram `json:"education_programs"`
}
type AddEntranceTestDate struct {
	IdEntranceTest       uint
	EntranceTestCalendar []digest.EntranceTestCalendar `json:"entrance_test_calendar"`
}

var sortByArrayCompetitive = []string{
	`id`,
	`uid`,
	`name`,
	`id_direction`,
	`code_specialty`,
	`name_specialty`,
	`id_level_budget`,
	`id_education_source`,
	`id_education_level`,
	`id_education_form`,
	`name_level_budget`,
	`name_education_source`,
	`name_education_level`,
	`name_education_form`,
}

func (result *ResultCheck) CheckNumberAddCompetitive() {
	result.Done = true
	return
}
func (result *Result) GetListCompetitiveGroupsByCompanyId(campaignId uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var competitiveGroups []digest.CompetitiveGroup
	var db *gorm.DB
	if service.SearchStringInSliceString(result.Sort.Field, sortByArrayCompetitive) >= 0 {
		order := `asc`
		if service.SearchStringInSliceString(result.Sort.Order, orderArray) >= 0 {
			order = result.Sort.Order
		}
		db = conn.Order(result.Sort.Field + ` ` + order)
	} else {
		db = conn.Order(`created asc `)
	}
	for _, search := range result.Search {
		if service.SearchStringInSliceString(search[0], CompetitiveSearchArray) >= 0 {
			db = db.Where(`(UPPER(`+search[0]+`) LIKE ? OR ( UPPER( code_specialty || ' ' || name_specialty) LIKE ?))`, `%`+strings.ToUpper(search[1])+`%`, `%`+strings.ToUpper(search[1])+`%`)
		}
	}
	db = db.Table(`cmp.v_competitive_groups `).Where(`id_campaign=? AND actual IS TRUE`, campaignId)
	//if result.Search != `` {
	//	db = db.Where(`UPPER(name) LIKE ?`, `%`+strings.ToUpper(result.Search)+`%`)
	//}
	dbCount := db.Model(&competitiveGroups).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()

	db = db.Preload(`EducationLevel`).Preload(`Campaign`).Preload(`LevelBudget`).Preload(`EducationSource`).Preload(`EducationForm`).Preload(`Direction`).Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&competitiveGroups)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Конкурсы не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, competitveGroup := range competitiveGroups {
			//number := competitveGroup.BudgetO+competitveGroup.BudgetOz+competitveGroup.BudgetZ+competitveGroup.PaidO+competitveGroup.PaidOz+competitveGroup.PaidZ+competitveGroup.PaidO+competitveGroup.PaidOz+competitveGroup.PaidZ+competitveGroup.TargetO+competitveGroup.TargetOz+competitveGroup.TargetZ
			var number int64
			switch competitveGroup.IdEducationSource {
			case 1: // Бюджетные места
				switch competitveGroup.IdEducationForm {
				case 1: // Очная форма
					number = competitveGroup.BudgetO
					break
				case 2: // Очно-заочная(вечерняя)
					number = competitveGroup.BudgetOz
					break
				case 3: // Заочная
					number = competitveGroup.BudgetZ
					break
				}
				break
			case 2: // Квота приема лиц, имеющих особое право
				switch competitveGroup.IdEducationForm {
				case 1: // Очная форма
					number = competitveGroup.QuotaO
					break
				case 2: // Очно-заочная(вечерняя)
					number = competitveGroup.QuotaOz
					break
				case 3: // Заочная
					number = competitveGroup.QuotaZ
					break
				}
				break
			case 3: // С оплатой обучения
				switch competitveGroup.IdEducationForm {
				case 1: // Очная форма
					number = competitveGroup.PaidO
					break
				case 2: // Очно-заочная(вечерняя)
					number = competitveGroup.PaidOz
					break
				case 3: // Заочная
					number = competitveGroup.PaidZ
					break
				}
				break
			case 4: // Целевой прием
				switch competitveGroup.IdEducationForm {
				case 1: // Очная форма
					number = competitveGroup.TargetO
					break
				case 2: // Очно-заочная(вечерняя)
					number = competitveGroup.TargetOz
					break
				case 3: // Заочная
					number = competitveGroup.TargetZ
					break
				}
			}
			c := map[string]interface{}{
				`id`:                    competitveGroup.Id,
				`name`:                  competitveGroup.Name,
				`id_campaign`:           competitveGroup.IdCampaign,
				`name_campaign`:         competitveGroup.Campaign.Name,
				`year_start_campaign`:   competitveGroup.Campaign.YearStart,
				`year_end_campaign`:     competitveGroup.Campaign.YearEnd,
				`id_education_form`:     competitveGroup.EducationForm.Id,
				`name_education_form`:   competitveGroup.EducationForm.Name,
				`id_education_source`:   competitveGroup.EducationSource.Id,
				`name_education_source`: competitveGroup.EducationSource.Name,
				`id_education_level`:    competitveGroup.EducationLevel.Id,
				`name_education_level`:  competitveGroup.EducationLevel.Name,
				`id_direction`:          competitveGroup.Direction.Id,
				`name_direction`:        competitveGroup.Direction.Name,
				`code_specialty`:        competitveGroup.Direction.Code,
				`id_level_budget`:       competitveGroup.LevelBudget.Id,
				`name_level_budget`:     competitveGroup.LevelBudget.Name,
				`id_author`:             competitveGroup.IdAuthor,
				`actual`:                competitveGroup.Actual,
				`uid`:                   competitveGroup.UID,
				`id_organization`:       competitveGroup.IdOrganization,
				`created`:               competitveGroup.Created,
				`comment`:               competitveGroup.Comment,
				`number`:                number,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Конкурсы не найдены.`
		result.Message = &message
		result.Items = []digest.CompetitiveGroup{}
		return
	}
}

func (result *ResultInfo) GetDirectionByEntrant(keys map[string][]string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var db *gorm.DB
	var directions []RowsCls
	db = conn.Select(`id, name`).Table(`cls.v_direction_specialty`)
	if len(keys[`id_entrant`]) > 0 {
		if v, ok := strconv.Atoi(keys[`id_entrant`][0]); ok == nil {
			var idDirection []DirectionCompetitiveGroups
			_ = conn.Raw(`select  cg.id_direction 
						from app.applications a 
						join cmp.competitive_groups cg on cg.id=a.id_competitive_group 
						where id_entrant = ?  AND cg.actual IS TRUE AND a.actual IS TRUE
						group by cg.id_direction `, v).Scan(&idDirection)
			if len(idDirection) >= 3 {
				var existDirection []uint
				for _, val := range idDirection {
					existDirection = append(existDirection, val.IdDirection)
				}
				if len(existDirection) > 0 {
					db = db.Where(`id IN (?)`, existDirection)
				}
			}
		} else {
			message := `Неверный идентификатор абитуриента.`
			result.Message = &message
			return
		}
	}
	if len(keys[`search`]) > 0 {
		db = db.Where(`UPPER(name) LIKE ?`, `%`+strings.ToUpper(keys[`search`][0])+`%`)
	}

	db = db.Where(`actual IS TRUE`).Find(&directions)
	if db.Error != nil {
		message := db.Error.Error()
		result.Message = &message
		return
	}
	result.Done = true
	result.Items = directions
	return
}
func (result *ResultInfo) GetListCompetitiveGroups(keys map[string][]string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var competitiveGroups []digest.CompetitiveGroup
	var db *gorm.DB

	db = conn.Where(`id_organization=? AND actual IS TRUE`, result.User.CurrentOrganization.Id)
	if len(keys[`id_entrant`]) > 0 {
		if v, ok := strconv.Atoi(keys[`id_entrant`][0]); ok == nil {
			var idDirection []DirectionCompetitiveGroups
			_ = conn.Raw(`select  cg.id_direction 
						from app.applications a 
						join cmp.competitive_groups cg on cg.id=a.id_competitive_group 
						where id_entrant = ? AND a.id_organization = ? AND cd.actual IS TRUE AND a.actual IS TRUE
						group by cg.id_direction `, v, result.User.CurrentOrganization.Id).Scan(&idDirection)
			if len(idDirection) >= 3 {
				var existDirection []uint
				for _, val := range idDirection {
					existDirection = append(existDirection, val.IdDirection)
				}
				if len(existDirection) > 0 {
					db = db.Where(`id_direction IN (?)`, existDirection)
				}
			}
			var existApplication []digest.Application
			_ = conn.Where(`id_entrant=? AND id_organization=? AND actual IS TRUE`, v, result.User.CurrentOrganization.Id).Find(&existApplication)

			if len(existApplication) >= 0 {
				var existCompetitiveIds []uint
				for _, val := range existApplication {
					existCompetitiveIds = append(existCompetitiveIds, val.IdCompetitiveGroup)
				}
				if len(existCompetitiveIds) > 0 {
					db = db.Where(`id NOT IN (?)`, existCompetitiveIds)
				}
			}
		} else {
			message := `Неверный идентификатор абитуриента.`
			result.Message = &message
			return
		}
	} else {
		message := `Нет идентификатора абитуриента.`
		result.Message = &message
		return
	}

	for key, value := range keys {
		if len(value) > 0 {
			switch key {
			case `id_campaign`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_campaign=?`, v)
				}
				break
			case `id_education_form`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_education_form=?`, v)
				}
				break
			case `id_education_level`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_education_level=?`, v)
				}
				break
			case `id_education_source`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_education_source=?`, v)
				}
				break
			case `id_direction`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_direction=?`, v)
				}
				break
			case `id_level_budget`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_level_budget=?`, v)
				}
				break
			case `search`:
				db = db.Where(`UPPER(name) LIKE ?`, `%`+strings.ToUpper(value[0])+`%`)
				break
			}
		}
	}

	db = db.Preload(`EducationLevel`).Preload(`Campaign`).Preload(`LevelBudget`).Preload(`EducationSource`).Preload(`EducationForm`).Preload(`Direction`).Find(&competitiveGroups)

	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Конкурсные группы не найдены.`
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
				`id`:           competitveGroup.Id,
				`name`:         competitveGroup.Name,
				`uid`:          competitveGroup.UID,
				`id_direction`: competitveGroup.IdDirection,
				`created`:      competitveGroup.Created,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Конкурсные группы не найдены.`
		result.Message = &message
		result.Items = []digest.CompetitiveGroup{}
		return
	}
}
func (result *ResultInfo) GetShortListCompetitiveGroups(keys map[string][]string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var competitiveGroups []digest.CompetitiveGroup
	var db *gorm.DB
	db = conn.Where(`id_organization=? AND actual IS TRUE`, result.User.CurrentOrganization.Id)

	for key, value := range keys {
		if len(value) > 0 {
			switch key {
			case `id_campaign`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_campaign=?`, v)
				}
				break
			case `id_education_form`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_education_form=?`, v)
				}
				break
			case `id_education_level`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_education_level=?`, v)
				}
				break
			case `id_education_source`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_education_source=?`, v)
				}
				break
			case `id_direction`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_direction=?`, v)
				}
				break
			case `id_level_budget`:
				if v, ok := strconv.Atoi(value[0]); ok == nil {
					db = db.Where(`id_level_budget=?`, v)
				}
				break
			case `search`:
				db = db.Where(`UPPER(name) LIKE ?`, `%`+strings.ToUpper(value[0])+`%`)
				break
			}
		}
	}

	db = db.Find(&competitiveGroups)

	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Конкурсные группы не найдены.`
			result.Message = &message
			result.Items = []digest.CompetitiveGroup{}
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, competitveGroup := range competitiveGroups {
			c := map[string]interface{}{
				`id`:   competitveGroup.Id,
				`name`: competitveGroup.Name,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Конкурсные группы не найдены.`
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
	competitive.IdEducationLevel = data.CompetitiveGroup.IdEducationLevel
	competitive.Name = strings.TrimSpace(data.CompetitiveGroup.Name)
	err := CheckAddCompetitive(data.CompetitiveGroup.IdCampaign)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	var campaign digest.Campaign
	db := tx.Where(`actual IS TRUE`).Preload(`CampaignType`).Find(&campaign, data.CompetitiveGroup.IdCampaign)
	if campaign.Id < 1 {
		result.SetErrorResult(`Компания не найдена`)
		tx.Rollback()
		return
	}
	row := conn.Table(`cls.edu_levels_campaign_types`).Where(`id_campaign_types=? AND id_education_level=? `, campaign.CampaignType.Id, data.CompetitiveGroup.IdEducationLevel).Select(`id`).Row()
	var idEducLevelCampaignType uint
	err = row.Scan(&idEducLevelCampaignType)
	if err != nil || idEducLevelCampaignType <= 0 {
		result.SetErrorResult(`Данный уровень образования не соответствует типу приемной компании`)
		tx.Rollback()
		return
	}
	competitive.IdCampaign = data.CompetitiveGroup.IdCampaign
	var direction digest.Direction
	db = tx.Where(`id_education_level=?`, data.CompetitiveGroup.IdEducationLevel).Find(&direction, data.CompetitiveGroup.IdDirection)
	if direction.Id < 1 {
		result.SetErrorResult(`Специальность не найдена`)
		tx.Rollback()
		return
	}
	competitive.IdDirection = data.CompetitiveGroup.IdDirection

	var educForm digest.EducationForm
	db = tx.Find(&educForm, data.CompetitiveGroup.IdEducationForm)
	if direction.Id < 1 {
		result.SetErrorResult(`Форма образования не найдена`)
		tx.Rollback()
		return
	}
	competitive.IdEducationForm = data.CompetitiveGroup.IdEducationForm

	var educSource digest.EducationSource
	db = tx.Find(&educSource, data.CompetitiveGroup.IdEducationSource)
	if direction.Id < 1 {
		result.SetErrorResult(`Форма образования не найдена`)
		tx.Rollback()
		return
	}
	competitive.IdEducationSource = data.CompetitiveGroup.IdEducationSource
	if competitive.IdEducationSource != 3 && data.CompetitiveGroup.IdLevelBudget == nil { // 3 - платка.
		result.SetErrorResult(`Не указан уровень бюджета`)
		tx.Rollback()
		return
	}
	if data.CompetitiveGroup.IdLevelBudget != nil {
		var levelBudget digest.LevelBudget
		db = tx.Find(&levelBudget, *data.CompetitiveGroup.IdLevelBudget)
		if direction.Id < 1 {
			result.SetErrorResult(`Уровень бюджета не найден`)
			tx.Rollback()
			return
		}
		competitive.IdLevelBudget = data.CompetitiveGroup.IdLevelBudget
	}

	if data.CompetitiveGroup.UID != nil {
		var exist digest.CompetitiveGroup
		tx.Where(`id_organization=? AND uid=? AND actual IS TRUE`, result.User.CurrentOrganization.Id, *data.CompetitiveGroup.UID).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Конкурсная группа с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		competitive.UID = data.CompetitiveGroup.UID
	}
	competitive.Actual = true
	number := data.CompetitiveGroup.Number
	switch competitive.IdEducationSource {
	case 1: // Бюджетные места
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			competitive.BudgetO = *number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.BudgetOz = *number
			break
		case 3: // Заочная
			competitive.BudgetZ = *number
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
			competitive.QuotaO = *number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.QuotaOz = *number
			break
		case 3: // Заочная
			competitive.QuotaZ = *number
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
			competitive.PaidO = *number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.PaidOz = *number
			break
		case 3: // Заочная
			competitive.PaidZ = *number
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
			competitive.TargetO = *number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.TargetOz = *number
			break
		case 3: // Заочная
			competitive.TargetZ = *number
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
	err = CheckNumberCompetitive(data.CompetitiveGroup, *number)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	if data.CompetitiveGroup.Comment != nil && strings.TrimSpace(*data.CompetitiveGroup.Comment) != `` {
		competitive.Comment = data.CompetitiveGroup.Comment
	} else {
		competitive.Comment = nil
	}
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&competitive)
	var idsPrograms []uint
	var idsEntrance []uint
	if db.Error == nil {
		if len(data.CompetitiveGroupPrograms) > 0 {
			for _, value := range data.CompetitiveGroupPrograms {
				var program digest.CompetitiveGroupProgram
				program = value
				program.IdOrganization = result.User.CurrentOrganization.Id
				program.IdAuthor = result.User.Id
				program.IdCompetitiveGroup = competitive.Id
				program.Actual = true
				program.Created = time.Now()
				if value.Uid != nil {
					var exist digest.CompetitiveGroupProgram
					tx.Where(`id_organization=? AND uid=? AND actual IS TRUE`, result.User.CurrentOrganization.Id, *value.Uid).Find(&exist)
					if exist.Id > 0 {
						result.SetErrorResult(`Образовательная программа с данным uid уже существует у выбранной организации`)
						tx.Rollback()
						return
					}
					program.Uid = value.Uid
				}
				db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&program)
				if db.Error != nil {
					result.SetErrorResult(db.Error.Error())
					tx.Rollback()
					return
				}
				idsPrograms = append(idsPrograms, program.Id)
			}
		}
		if len(data.EntranceTests) > 0 {
			for _, value := range data.EntranceTests {
				var entrance digest.EntranceTest
				entrance = value
				entrance.IdOrganization = result.User.CurrentOrganization.Id
				entrance.IdAuthor = result.User.Id
				entrance.IdCompetitiveGroup = competitive.Id
				entrance.Actual = true
				entrance.Created = time.Now()
				if value.Uid != nil {
					var exist digest.EntranceTest
					tx.Where(`id_organization=? AND uid=? AND actual IS TRUE`, result.User.CurrentOrganization.Id, *value.Uid).Find(&exist)
					if exist.Id > 0 {
						result.SetErrorResult(`Вступительный тест с данным uid уже существует у выбранной организации`)
						tx.Rollback()
						return
					}
					entrance.Uid = value.Uid
				}
				var entranceTestType digest.EntranceTestType
				db = tx.Find(&entranceTestType, entrance.IdEntranceTestType)
				if entranceTestType.Id < 1 {
					result.SetErrorResult(`Тип вступительного теста не найден`)
					tx.Rollback()
					return
				}
				db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&entrance)
				if db.Error != nil {
					result.SetErrorResult(db.Error.Error())
					tx.Rollback()
					return
				}
				idsEntrance = append(idsEntrance, entrance.Id)
			}
		}

		result.Items = map[string]interface{}{
			`id_competitive_group`:  competitive.Id,
			`id_education_programs`: idsPrograms,
			`id_entrance_tests`:     idsEntrance,
		}
		result.Done = true
		tx.Commit()
	} else {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}

	tx.Rollback()
	return

}
func (result *ResultInfo) EditCompetitive(data CompetitiveGroup) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var competitive digest.CompetitiveGroup

	db := tx.Where(`id=?  AND actual IS TRUE`, data.Id).Find(&competitive)
	if competitive.Id <= 0 {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}
	data.IdCampaign = competitive.IdCampaign
	err := CheckEditCompetitiveGroup(data.Id)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	if competitive.IdOrganization != result.User.CurrentOrganization.Id {
		result.SetErrorResult(`Конкурсная группа принадлежит другой организации.`)
		tx.Rollback()
		return
	}

	competitive.Name = strings.TrimSpace(data.Name)
	if data.UID != nil && data.UID != competitive.UID && *data.UID != `` {
		var exist digest.CompetitiveGroup
		db = tx.Where(`uid=? AND id_organization=? AND id!=?  AND actual IS TRUE`, data.UID, result.User.CurrentOrganization.Id, competitive.Id).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Конкурсная группа с данным uid уже существует у данной организации.`)
			tx.Rollback()
			return
		}
		competitive.UID = data.UID
	}
	if competitive.IdEducationLevel != data.IdEducationLevel {
		var category digest.EducationLevel
		db = tx.Find(&category, data.IdEducationLevel)
		if category.Id < 1 {
			result.SetErrorResult(`Уровень образования не найден`)
			tx.Rollback()
			return
		}
		competitive.IdEducationLevel = data.IdEducationLevel
	}
	if competitive.IdEducationForm != data.IdEducationForm {
		var category digest.EducationForm
		db = tx.Find(&category, data.EducationForm)
		if category.Id < 1 {
			result.SetErrorResult(`Форма образования не найдена`)
			tx.Rollback()
			return
		}
		competitive.IdEducationForm = data.IdEducationForm
	}
	if data.IdLevelBudget != nil && competitive.IdLevelBudget != data.IdLevelBudget {
		var category digest.LevelBudget
		db = tx.Find(&category, *data.IdLevelBudget)
		if category.Id < 1 {
			result.SetErrorResult(`Уровень бюджета не найден`)
			tx.Rollback()
			return
		}
		competitive.IdLevelBudget = data.IdLevelBudget
	}
	if competitive.IdEducationSource != data.IdEducationSource {
		var category digest.EducationSource
		db = tx.Find(&category, data.IdEducationSource)
		if category.Id < 1 {
			result.SetErrorResult(`Источник финансирования не найден`)
			tx.Rollback()
			return
		}
		competitive.IdEducationSource = data.IdEducationSource
	}
	t := time.Now()

	competitive.Changed = &t
	competitive.IdAuthor = result.User.Id
	var campaign digest.Campaign
	db = tx.Where(`actual IS TRUE`).Preload(`CampaignType`).Find(&campaign, competitive.IdCampaign)
	if campaign.Id < 1 {
		result.SetErrorResult(`Компания не найдена`)
		tx.Rollback()
		return
	}
	row := conn.Table(`cls.edu_levels_campaign_types`).Where(`id_campaign_types=? AND id_education_level=?`, campaign.CampaignType.Id, data.IdEducationLevel).Select(`id`).Row()
	var idEducLevelCampaignType uint
	err = row.Scan(&idEducLevelCampaignType)
	if err != nil || idEducLevelCampaignType <= 0 {
		result.SetErrorResult(`Данный уровень образования не соответствует типу приемной компании`)
		tx.Rollback()
		return
	}
	if competitive.IdDirection != data.IdDirection {
		var direction digest.Direction
		db = tx.Where(`id_education_level=?`, data.IdEducationLevel).Find(&direction, data.IdDirection)
		if direction.Id < 1 {
			result.SetErrorResult(`Специальность не найдена`)
			tx.Rollback()
			return
		}
		competitive.IdDirection = data.IdDirection
	}
	if data.Number != nil {
		number := data.Number
		switch competitive.IdEducationSource {
		case 1: // Бюджетные места
			switch competitive.IdEducationForm {
			case 1: // Очная форма
				competitive.BudgetO = *number
				break
			case 2: // Очно-заочная(вечерняя)
				competitive.BudgetOz = *number
				break
			case 3: // Заочная
				competitive.BudgetZ = *number
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
				competitive.QuotaO = *number
				break
			case 2: // Очно-заочная(вечерняя)
				competitive.QuotaOz = *number
				break
			case 3: // Заочная
				competitive.QuotaZ = *number
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
				competitive.PaidO = *number
				break
			case 2: // Очно-заочная(вечерняя)
				competitive.PaidOz = *number
				break
			case 3: // Заочная
				competitive.PaidZ = *number
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
				competitive.TargetO = *number
				break
			case 2: // Очно-заочная(вечерняя)
				competitive.TargetOz = *number
				break
			case 3: // Заочная
				competitive.TargetZ = *number
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
		err = CheckNumberCompetitive(data, *number)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
	}

	if data.Comment != nil && strings.TrimSpace(*data.Comment) != `` {
		competitive.Comment = data.Comment
	} else {
		competitive.Comment = nil
	}
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&competitive)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
	result.Done = true
	result.Items = map[string]interface{}{
		`id_competitive_group`: competitive.Id,
	}
	tx.Commit()
	return

}
func (result *ResultInfo) EditCompetitiveComment(idCompetitiveGroup uint, comment string) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var competitive digest.CompetitiveGroup

	db := tx.Where(`id=?  AND actual IS TRUE`, idCompetitiveGroup).Find(&competitive)
	if competitive.Id <= 0 {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}

	t := time.Now()
	competitive.Changed = &t
	competitive.Comment = &comment

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&competitive)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
	result.Done = true
	result.Items = map[string]interface{}{
		`id_competitive_group`: competitive.Id,
	}
	tx.Commit()
	return

}
func (result *ResultInfo) EditCompetitiveName(idCompetitiveGroup uint, name string) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var competitive digest.CompetitiveGroup

	db := tx.Where(`id=?  AND actual IS TRUE`, idCompetitiveGroup).Find(&competitive)
	if competitive.Id <= 0 {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}

	t := time.Now()
	competitive.Changed = &t
	competitive.Name = name

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&competitive)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
	result.Done = true
	result.Items = map[string]interface{}{
		`id_competitive_group`: competitive.Id,
	}
	tx.Commit()
	return

}
func (result *ResultInfo) EditNumberCompetitive(data EditNumberCompetitive) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var competitive digest.CompetitiveGroup

	db := tx.Where(`id=?  AND actual IS TRUE`, data.IdCompetitive).Find(&competitive)
	if competitive.Id <= 0 {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}
	var campaign digest.Campaign
	conn.Where(`id=? AND actual is true`, competitive.IdCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		result.SetErrorResult(`Компания конкурсной группы не найдена `)
		tx.Rollback()
		return
	}
	if campaign.IdCampaignStatus == 3 { // 3 - статус завершена
		result.SetErrorResult(`Редактирование невозможно. Приемная компания завершена. `)
		tx.Rollback()
		return
	}
	if competitive.IdOrganization != result.User.CurrentOrganization.Id {
		result.SetErrorResult(`Конкурсная группа принадлежит другой организации.`)
		tx.Rollback()
		return
	}

	number := data.Number
	switch competitive.IdEducationSource {
	case 1: // Бюджетные места
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			competitive.BudgetO = number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.BudgetOz = number
			break
		case 3: // Заочная
			competitive.BudgetZ = number
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
			competitive.QuotaO = number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.QuotaOz = number
			break
		case 3: // Заочная
			competitive.QuotaZ = number
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
			competitive.PaidO = number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.PaidOz = number
			break
		case 3: // Заочная
			competitive.PaidZ = number
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
			competitive.TargetO = number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.TargetOz = number
			break
		case 3: // Заочная
			competitive.TargetZ = number
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
	err := CheckNumberCompetitiveById(competitive.Id, number)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	t := time.Now()
	competitive.Changed = &t
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&competitive)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
	result.Done = true
	result.Items = map[string]interface{}{
		`id_competitive_group`: competitive.Id,
	}
	tx.Commit()
	return

}

func (result *ResultInfo) AddProgram(idCompetitive uint, data AddCompetitiveGroup) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var competitive digest.CompetitiveGroup
	db := conn.Where(`actual IS TRUE`).Find(&competitive, idCompetitive)
	if competitive.Id == 0 {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}
	err := CheckEditProgramsCompetitiveGroup(competitive.Id)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	var idsPrograms []uint

	if len(data.CompetitiveGroupPrograms) > 0 {
		for _, value := range data.CompetitiveGroupPrograms {
			var program digest.CompetitiveGroupProgram
			program = value
			program.IdOrganization = result.User.CurrentOrganization.Id
			program.IdAuthor = result.User.Id
			program.IdCompetitiveGroup = competitive.Id
			program.Actual = true
			program.Created = time.Now()
			if value.Uid != nil {
				var exist digest.CompetitiveGroupProgram
				tx.Where(`id_organization=? AND uid=? AND actual IS TRUE`, result.User.CurrentOrganization.Id, *value.Uid).Find(&exist)
				if exist.Id > 0 {
					result.SetErrorResult(`Образовательная программа с данным uid уже существует у выбранной организации`)
					tx.Rollback()
					return
				}
				program.Uid = value.Uid
			}
			db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&program)
			if db.Error != nil {
				result.SetErrorResult(db.Error.Error())
				tx.Rollback()
				return
			}
			idsPrograms = append(idsPrograms, program.Id)
		}
	} else {
		result.SetErrorResult(`Образовательные программы не найдены`)
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_competitive_group`:  competitive.Id,
		`id_education_programs`: idsPrograms,
	}
	result.Done = true
	tx.Commit()
	return
}

func (result *ResultInfo) AddEntrance(idCompetitive uint, data AddEntrance) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var competitive digest.CompetitiveGroup
	db := conn.Where(`actual IS TRUE`).Find(&competitive, idCompetitive)
	if competitive.Id == 0 {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}
	err := CheckEditEntranceCompetitiveGroup(competitive.Id)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	var idsEntrance []uint
	if len(data.EntranceTests) > 0 {
		for _, value := range data.EntranceTests {
			var entrance digest.EntranceTest
			entrance = value
			entrance.IdOrganization = result.User.CurrentOrganization.Id
			entrance.IdAuthor = result.User.Id
			entrance.IdCompetitiveGroup = competitive.Id
			entrance.Actual = true
			entrance.Created = time.Now()
			if value.Uid != nil {
				var exist digest.EntranceTest
				tx.Where(`id_organization=? AND uid=? AND actual IS TRUE`, result.User.CurrentOrganization.Id, *value.Uid).Find(&exist)
				if exist.Id > 0 {
					result.SetErrorResult(`Вступительный тест с данным uid уже существует у выбранной организации`)
					tx.Rollback()
					return
				}
				entrance.Uid = value.Uid
			}
			var entranceTestType digest.EntranceTestType
			db = tx.Find(&entranceTestType, entrance.IdEntranceTestType)
			if entranceTestType.Id < 1 {
				result.SetErrorResult(`Тип вступительного теста не найден`)
				tx.Rollback()
				return
			}
			db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&entrance)
			if db.Error != nil {
				result.SetErrorResult(db.Error.Error())
				tx.Rollback()
				return
			}
			idsEntrance = append(idsEntrance, entrance.Id)
		}
	} else {
		result.SetErrorResult(`Вступительные тесты не найдены`)
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_competitive_group`: competitive.Id,
		`id_entrance_tests`:    idsEntrance,
	}
	result.Done = true
	tx.Commit()
	return
}

func (result *ResultInfo) GetEntranceTestsCalendarByEntrance(idEntranceTest uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var entranceTestCalendar []digest.EntranceTestCalendar
	db := conn.Where(`id_entrance_test=? AND actual IS TRUE`, idEntranceTest).Order(`created desc`).Find(&entranceTestCalendar)

	if db.RowsAffected > 0 {
		var entranceTestDates []interface{}
		for index, _ := range entranceTestCalendar {
			//date := entranceTestCalendar[index].EntranceTestDate.Format("2006-01-02 15:04:05")
			entranceTestDates = append(entranceTestDates, map[string]interface{}{
				"id":                 entranceTestCalendar[index].Id,
				"id_entrance_test":   entranceTestCalendar[index].IdEntranceTest,
				"entrance_test_date": entranceTestCalendar[index].EntranceTestDate,
				"exam_location":      entranceTestCalendar[index].ExamLocation,
				"uid":                entranceTestCalendar[index].Uid,
				"uid_epgu":           entranceTestCalendar[index].UidEpgu,
				"created":            entranceTestCalendar[index].Created,
				"count_c":            entranceTestCalendar[index].CountC,
			})
		}
		result.Done = true
		result.Items = entranceTestDates
		return
	} else {
		result.Done = true
		message := `Даты вступительных испытаний не найдены.`
		result.Message = &message
		result.Items = []digest.CompetitiveGroup{}
		return
	}
}
func (result *ResultInfo) GetListDatesByEntranceTest(idEntranceTest uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var ec []digest.EntranceTestCalendar
	var db *gorm.DB

	db = conn.Where(`id_entrance_test=? AND actual IS TRUE`, idEntranceTest).Find(&ec)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Даты не найдены не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for index, _ := range ec {
			date := ec[index].EntranceTestDate.Format("2006-01-02 15:04:05")
			c := map[string]interface{}{
				"id":                 ec[index].Id,
				"id_entrance_test":   ec[index].IdEntranceTest,
				"entrance_test_date": date,
				"exam_location":      ec[index].ExamLocation,
				"uid":                ec[index].Uid,
				"uid_epgu":           ec[index].UidEpgu,
				"created":            ec[index].Created,
				"count_c":            ec[index].CountC,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Даты не найдены.`
		result.Message = &message
		result.Items = []digest.CompetitiveGroup{}
		return
	}
}
func (result *ResultInfo) AddEntranceTestCalendar(data AddEntranceTestDate) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var entranceTest digest.EntranceTest
	db := conn.Where(`id=? AND actual IS TRUE`, data.IdEntranceTest).Find(&entranceTest)
	if entranceTest.Id == 0 {
		result.SetErrorResult(`Вступительное испытание не найдено`)
		tx.Rollback()
		return
	}
	err := CheckAddRemoveEntranceTestCalendar(entranceTest.IdCompetitiveGroup)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	if entranceTest.IdOrganization != result.User.CurrentOrganization.Id {
		result.SetErrorResult(`Не совпадает организация`)
		tx.Rollback()
		return
	}
	//err := CheckEditEntranceCompetitiveGroup(competitive.Id)
	//if err != nil {
	//	result.SetErrorResult(err.Error())
	//	tx.Rollback()
	//	return
	//}
	var idsEntranceCalendar []uint
	if len(data.EntranceTestCalendar) > 0 {
		for _, value := range data.EntranceTestCalendar {
			var entranceCalendar digest.EntranceTestCalendar
			if value.EntranceTestDate.Year() < 2006 {
				result.SetErrorResult(`Неверное значение даты`)
				tx.Rollback()
				return
			}
			if strings.TrimSpace(value.ExamLocation) == `` {
				result.SetErrorResult(`Неверное значение места проведения`)
				tx.Rollback()
				return
			}
			entranceCalendar = value
			if value.CountC != nil {
				entranceCalendar.CountC = value.CountC
			} else {
				entranceCalendar.CountC = nil
			}

			entranceCalendar.IdOrganization = result.User.CurrentOrganization.Id
			entranceCalendar.IdAuthor = result.User.Id
			entranceCalendar.IdEntranceTest = data.IdEntranceTest
			entranceCalendar.Actual = true
			entranceCalendar.Created = time.Now()
			entranceCalendar.UidEpgu = nil
			if value.Uid != nil {
				var exist digest.EntranceTestCalendar
				tx.Where(`id_organization=? AND uid=? AND actual IS TRUE`, result.User.CurrentOrganization.Id, *value.Uid).Find(&exist)
				if exist.Id > 0 {
					result.SetErrorResult(`Дата проведения вступительного испытания с данным uid уже существует у выбранной организации`)
					tx.Rollback()
					return
				}
				entranceCalendar.Uid = value.Uid
			}

			db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&entranceCalendar)
			if db.Error != nil {
				result.SetErrorResult(db.Error.Error())
				tx.Rollback()
				return
			}
			idsEntranceCalendar = append(idsEntranceCalendar, entranceCalendar.Id)
		}
	} else {
		result.SetErrorResult(`Даты проведения втсупительных  испытаний не найдены`)
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrance_test`:          data.IdEntranceTest,
		`id_entrance_test_calendar`: idsEntranceCalendar,
	}
	result.Done = true
	tx.Commit()
	return
}
func (result *ResultInfo) RemoveEntranceTestCalendar(idEntranceTestDate uint) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var entranceCalendar digest.EntranceTestCalendar
	db := conn.Where(`id=? AND actual IS TRUE`, idEntranceTestDate).Find(&entranceCalendar)
	if entranceCalendar.Id == 0 {
		result.SetErrorResult(`Дата вступительного испытания не найдена`)
		tx.Rollback()
		return
	}
	if entranceCalendar.UidEpgu != nil {
		result.SetErrorResult(`Невозможно удалить дату вступительных испытаний, полученную с ЕПГУ`)
		tx.Rollback()
		return
	}
	if db.RowsAffected > 0 {
		var entranceTest digest.EntranceTest
		db = conn.Where(`id=? AND actual IS TRUE`, entranceCalendar.IdEntranceTest).Find(&entranceTest)
		if entranceTest.Id == 0 {
			result.SetErrorResult(`Вступительное испытание не найдено`)
			tx.Rollback()
			return
		}
		err := CheckAddRemoveEntranceTestCalendar(entranceTest.IdCompetitiveGroup)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		if entranceTest.IdOrganization != result.User.CurrentOrganization.Id {
			result.SetErrorResult(`Не совпадает организация`)
			tx.Rollback()
			return
		}
		t := time.Now()
		entranceCalendar.Changed = &t
		entranceCalendar.Actual = false
		conn.Exec(`DELETE FROM app.entrance_test_agreed WHERE id_entrance_test_calendar=?`, idEntranceTestDate)
		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&entranceCalendar)
		if db.Error != nil {
			result.SetErrorResult(db.Error.Error())
			tx.Rollback()
			return
		}

	} else {
		result.Done = true
		message := `Дата вступительного испытания не найдена.`
		result.Message = &message
		result.Items = []digest.CompetitiveGroup{}
		return
	}
	result.Items = map[string]interface{}{
		`id_entrance_test`:          entranceCalendar.IdEntranceTest,
		`id_entrance_test_calendar`: entranceCalendar.Id,
	}
	result.Done = true
	tx.Commit()
	return
}

func (result *ResultInfo) GetInfoCompetitiveGroup(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var competitive digest.CompetitiveGroup
	db := conn.Where(`actual IS TRUE`).Preload(`Campaign`).Preload(`EducationForm`).Preload(`EducationLevel`).Preload(`EducationSource`).Preload(`LevelBudget`).Preload(`Direction`).Find(&competitive, ID)
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Конкурсная группа не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		var number int64
		//var campEducForms []digest.CampaignEducForm
		//db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducForms)
		//var campEducLevels []digest.CampaignEducLevel
		//db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducLevels)

		switch competitive.IdEducationSource {
		case 1: // Бюджетные места
			switch competitive.IdEducationForm {
			case 1: // Очная форма
				number = competitive.BudgetO
				break
			case 2: // Очно-заочная(вечерняя)
				number = competitive.BudgetOz
				break
			case 3: // Заочная
				number = competitive.BudgetZ
				break
			default:
				result.SetErrorResult(`Ошибка`)
				return
			}
			break
		case 2: // Квота приема лиц, имеющих особое право
			switch competitive.IdEducationForm {
			case 1: // Очная форма
				number = competitive.QuotaO
				break
			case 2: // Очно-заочная(вечерняя)
				number = competitive.QuotaOz
				break
			case 3: // Заочная
				number = competitive.QuotaZ
				break
			default:
				result.SetErrorResult(`Ошибка`)
				return
			}
			break
		case 3: // С оплатой обучения
			switch competitive.IdEducationForm {
			case 1: // Очная форма
				number = competitive.PaidO
				break
			case 2: // Очно-заочная(вечерняя)
				number = competitive.PaidOz
				break
			case 3: // Заочная
				number = competitive.PaidZ
				break
			default:
				result.SetErrorResult(`Ошибка`)
				return
			}
			break
		case 4: // Целевой прием
			switch competitive.IdEducationForm {
			case 1: // Очная форма
				number = competitive.TargetO
				break
			case 2: // Очно-заочная(вечерняя)
				number = competitive.TargetOz
				break
			case 3: // Заочная
				number = competitive.TargetZ
				break
			default:
				result.SetErrorResult(`Ошибка`)
				return
			}
		default:
			result.SetErrorResult(`Ошибка`)
			return
		}
		c := map[string]interface{}{
			"id":                    competitive.Id,
			"id_campaign":           competitive.IdCampaign,
			"name_campaign":         competitive.Campaign.Name,
			"id_direction":          competitive.Direction.Id,
			"code_direction":        competitive.Direction.Code,
			"name_direction":        competitive.Direction.Name,
			"name":                  competitive.Name,
			"uid":                   competitive.UID,
			"id_education_form":     competitive.EducationForm.Id,
			"name_education_form":   competitive.EducationForm.Name,
			"id_education_level":    competitive.EducationLevel.Id,
			"name_education_level":  competitive.EducationLevel.Name,
			"id_education_source":   competitive.EducationSource.Id,
			"name_education_source": competitive.EducationSource.Name,
			"id_level_budget":       competitive.LevelBudget.Id,
			"name_level_budget":     competitive.LevelBudget.Name,
			"number":                number,
			`comment`:               competitive.Comment,
		}
		//for _, campEducLevel := range campEducLevels {
		//	var educLevel digest.EducationLevel
		//	db = conn.Find(&educLevel, campEducLevel.IdEducationLevel)
		//	c.EducationLevels = append(c.EducationLevels, educLevel.Id)
		//	c.EducationLevelsName = append(c.EducationLevelsName, educLevel.Name)
		//}
		result.Done = true
		result.Items = c
		return
	} else {
		result.Done = true
		message := `Конкурсная группа не найдена.`
		result.Message = &message
		result.Items = []digest.CompetitiveGroup{}
		return
	}
}
func (result *ResultInfo) GetEducationProgramsCompetitiveGroup(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var competitive digest.CompetitiveGroup
	db := conn.Where(`actual IS TRUE`).Preload(`Campaign`).Find(&competitive, ID)
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Конкурсная группа не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		var idsPrograms []uint
		db = conn.Where(`id_competitive_group=? AND actual IS TRUE`, ID).Table(`cmp.competitive_group_programs`).Pluck(`id`, &idsPrograms)
		var programs []interface{}
		for _, id := range idsPrograms {
			var program digest.CompetitiveGroupProgram
			db = conn.Find(&program, id)
			programs = append(programs, map[string]interface{}{
				`id`:                 program.Id,
				`id_subdivision_org`: program.IdSubdivisionOrg,
				`name`:               program.Name,
				`uid`:                program.Uid,
			})
		}
		result.Done = true
		if len(programs) <= 0 {
			result.Items = []digest.CompetitiveGroupProgram{}
		} else {
			result.Items = programs
		}

		return
	} else {
		result.Done = true
		message := `Конкурсная группа не найдена.`
		result.Message = &message
		result.Items = []digest.CompetitiveGroup{}
		return
	}
}
func (result *ResultInfo) GetEntranceTestsCompetitiveGroup(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var competitive digest.CompetitiveGroup
	db := conn.Where(`actual IS TRUE`).Preload(`Campaign`).Find(&competitive, ID)
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Конкурсная группа не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		var idsEntrance []uint
		db = conn.Where(`id_competitive_group=? AND actual IS TRUE`, ID).Table(`cmp.entrance_test`).Pluck(`id`, &idsEntrance)
		var entranceTests []interface{}
		for _, id := range idsEntrance {
			var entrance digest.EntranceTest
			db = conn.Preload(`Subject`).Preload(`EntranceTestType`).Find(&entrance, id)
			e := map[string]interface{}{
				"id":                      entrance.Id,
				"id_entrance_test_type":   entrance.EntranceTestType.Id,
				"name_entrance_test_type": entrance.EntranceTestType.Name,
				"id_subject":              entrance.Subject.Id,
				"name_subject":            entrance.Subject.Name,
				"priority":                entrance.Priority,
				"uid":                     entrance.Uid,
				"test_name":               entrance.TestName,
				"min_score":               entrance.MinScore,
				"is_ege":                  entrance.IsEge,
			}
			var entranceTestCalendar []digest.EntranceTestCalendar
			db := conn.Where(`id_entrance_test=? AND actual IS TRUE`, entrance.Id).Order(`created desc`).Find(&entranceTestCalendar)
			var entranceTestDates []interface{}
			if db.RowsAffected > 0 {
				for index, _ := range entranceTestCalendar {
					//date := entranceTestCalendar[index].EntranceTestDate.Format("2006-01-02 15:04:05")
					entranceTestDates = append(entranceTestDates, map[string]interface{}{
						"id":                 entranceTestCalendar[index].Id,
						"id_entrance_test":   entranceTestCalendar[index].IdEntranceTest,
						"entrance_test_date": entranceTestCalendar[index].EntranceTestDate,
						"exam_location":      entranceTestCalendar[index].ExamLocation,
						"uid":                entranceTestCalendar[index].Uid,
						"uid_epgu":           entranceTestCalendar[index].UidEpgu,
						"created":            entranceTestCalendar[index].Created,
						"count_c":            entranceTestCalendar[index].CountC,
					})
				}
			}
			if len(entranceTestDates) > 0 {
				e["entrance_test_calendar"] = entranceTestDates
			} else {
				e["entrance_test_calendar"] = []digest.EntranceTestCalendar{}
			}

			entranceTests = append(entranceTests, e)

		}
		result.Done = true
		r := map[string]interface{}{
			"id_campaign": competitive.Campaign.Id,
		}
		if len(entranceTests) <= 0 {
			r[`tests`] = []digest.EntranceTest{}
		} else {
			r[`tests`] = entranceTests
		}
		result.Items = r
		return
	} else {
		result.Done = true
		message := `Конкурсная группа не найдена.`
		result.Message = &message
		result.Items = []digest.CompetitiveGroup{}
		return
	}
}

func (result *ResultInfo) RemoveEntranceCompetitive(idCompetitive uint, idEntrance uint) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var old digest.CompetitiveGroup
	db := tx.Where(`actual IS TRUE`).Find(&old, idCompetitive)
	if old.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}
	err := CheckEditEntranceCompetitiveGroup(old.Id)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	var entrance digest.EntranceTest
	db = conn.Where(`id_competitive_group=? AND id=? AND actual IS TRUE`, idCompetitive, idEntrance).Find(&entrance)
	if entrance.Id > 0 {
		db = tx.Exec(`DELETE FROM app.entrance_test WHERE id_entrance_test=?`, idEntrance)
		db = tx.Exec(`UPDATE cmp.entrance_test SET actual=false, changed=? WHERE id_competitive_group=? AND id=? `, time.Now(), idCompetitive, idEntrance)
		if db.Error == nil {
			result.Done = true
			tx.Commit()
			result.Items = map[string]interface{}{
				`id_competitive_group`: idCompetitive,
				`id_entrance`:          idEntrance,
			}
			return
		} else {
			result.SetErrorResult(db.Error.Error())
			tx.Rollback()
			return
		}
	} else {
		result.SetErrorResult(`Не найден вступительный тест`)
		tx.Rollback()
		return
	}
}

func (result *ResultInfo) RemoveProgramCompetitive(idCompetitive uint, idProgram uint) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var old digest.CompetitiveGroup
	db := tx.Where(`actual IS TRUE`).Find(&old, idCompetitive)
	if old.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}
	err := CheckEditProgramsCompetitiveGroup(old.Id)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	var program digest.CompetitiveGroupProgram
	db = conn.Where(`id_competitive_group=? AND id=? AND actual IS TRUE`, idCompetitive, idProgram).Find(&program)
	if program.Id > 0 {
		db = tx.Exec(`UPDATE cmp.competitive_group_programs SET actual=false, changed=? WHERE id_competitive_group=? AND id=?`, time.Now(), idCompetitive, idProgram)
		if db.Error == nil {
			result.Done = true
			tx.Commit()
			result.Items = map[string]interface{}{
				`id_competitive_group`: idCompetitive,
				`id_program`:           idProgram,
			}
			return
		} else {
			result.SetErrorResult(db.Error.Error())
			tx.Rollback()
			return
		}
	} else {
		result.SetErrorResult(`Не найдена образовательная программа`)
		tx.Rollback()
		return
	}
}

func (result *ResultInfo) RemoveCompetitive(idCompetitive uint) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var old digest.CompetitiveGroup
	db := tx.Where(`actual IS TRUE`).Find(&old, idCompetitive)
	if old.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}
	err := CheckAddCompetitive(old.IdCampaign)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	var programs []digest.CompetitiveGroupProgram
	db = conn.Where(`id_competitive_group=? AND actual IS TRUE`, idCompetitive).Find(&programs)
	if len(programs) > 0 {
		db = tx.Exec(`UPDATE cmp.competitive_group_programs SET actual=false, changed=? WHERE id_competitive_group=?`, time.Now(), idCompetitive)
	}
	var entrance []digest.EntranceTest
	db = conn.Where(`id_competitive_group=? AND actual IS TRUE`, idCompetitive).Find(&entrance)
	if len(entrance) > 0 {
		for _, value := range entrance {
			db = tx.Exec(`UPDATE cmp.entrance_test_calendar SET actual=false, changed=? WHERE id_entrance_test=?`, time.Now(), value.Id)
		}
		db = tx.Exec(`UPDATE cmp.entrance_test SET actual=false, changed=? WHERE id_competitive_group=?`, time.Now(), idCompetitive)
	}

	db = tx.Exec(`UPDATE app.applications SET actual=false, changed=? WHERE id_competitive_group=?`, time.Now(), idCompetitive)
	db = tx.Exec(`UPDATE cmp.competitive_groups SET actual=false, changed=? WHERE id=?`, time.Now(), idCompetitive)
	if db.Error == nil {
		result.Done = true
		tx.Commit()
		result.Items = map[string]interface{}{
			`id_competitive_group`: idCompetitive,
		}
		return
	} else {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}

}

func (result *ResultList) GetEntranceTestsSelectListByCompetitive(idCompetitive uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)

	var items []digest.EntranceTest
	sortField := `created`
	sortOrder := `desc`
	db := conn.Where(`id_competitive_group=? AND actual`, idCompetitive).Order(sortField + ` ` + sortOrder)

	if result.Search != `` {
		//db = db.Where(`UPPER(name) like ?`, `%`+strings.ToUpper(result.Search)+`%`)
	}

	db = db.Preload(`EntranceTestType`).Preload(`Subject`).Find(&items)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Вступительные испытания не найдены.`
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
				`id`:                      item.Id,
				`name_entrance_test_type`: item.EntranceTestType.Name,
				`min_score`:               item.MinScore,
				`is_ege`:                  item.IsEge,
				`priority`:                item.Priority,
			}
			if item.Subject.Id > 0 {
				c[`name_subject`] = item.Subject.Name
			} else {
				c[`name_subject`] = item.TestName
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
