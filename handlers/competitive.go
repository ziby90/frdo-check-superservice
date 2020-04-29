package handlers

import (
	"errors"
	"github.com/jinzhu/gorm"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
)

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
	Number            int64                  `json:"number"`
}

type AddCompetitiveGroup struct {
	CompetitiveGroup         CompetitiveGroup                 `json:"competitive"`
	CompetitiveGroupPrograms []digest.CompetitiveGroupProgram `json:"education_programs"`
	EntranceTests            []digest.EntranceTest            `json:"entrance_tests"`
}

type AddEntrance struct {
	EntranceTests []digest.EntranceTest `json:"entrance_tests"`
}

type AddCompetitiveGroupPrograms struct {
	CompetitiveGroupPrograms []digest.CompetitiveGroupProgram `json:"education_programs"`
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

	db = db.Preload(`EducationLevel`).Preload(`Campaign`).Preload(`LevelBudget`).Preload(`EducationSource`).Preload(`EducationForm`).Preload(`Direction`).Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&competitiveGroups)
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
	competitive.IdEducationLevel = data.CompetitiveGroup.IdEducationLevel
	competitive.Name = strings.TrimSpace(data.CompetitiveGroup.Name)

	var campaign digest.Campaign
	db := tx.Preload(`CampaignType`).Find(&campaign, data.CompetitiveGroup.IdCampaign)
	if campaign.Id < 1 {
		result.SetErrorResult(`Компания не найдена`)
		tx.Rollback()
		return
	}
	row := conn.Table(`cls.edu_levels_campaign_types`).Where(`id_campaign_types=? AND id_education_level=?`, campaign.CampaignType.Id, data.CompetitiveGroup.IdEducationLevel).Select(`id`).Row()
	var idEducLevelCampaignType uint
	err := row.Scan(&idEducLevelCampaignType)
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
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *data.CompetitiveGroup.UID).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Конкурсная группа с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		competitive.UID = data.CompetitiveGroup.UID
	}
	competitive.Actual = true

	switch competitive.IdEducationSource {
	case 1: // Бюджетные места
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			competitive.BudgetO = data.CompetitiveGroup.Number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.BudgetOz = data.CompetitiveGroup.Number
			break
		case 3: // Заочная
			competitive.BudgetZ = data.CompetitiveGroup.Number
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
			competitive.QuotaO = data.CompetitiveGroup.Number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.QuotaOz = data.CompetitiveGroup.Number
			break
		case 3: // Заочная
			competitive.QuotaZ = data.CompetitiveGroup.Number
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
			competitive.PaidO = data.CompetitiveGroup.Number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.PaidOz = data.CompetitiveGroup.Number
			break
		case 3: // Заочная
			competitive.PaidZ = data.CompetitiveGroup.Number
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
			competitive.TargetO = data.CompetitiveGroup.Number
			break
		case 2: // Очно-заочная(вечерняя)
			competitive.TargetOz = data.CompetitiveGroup.Number
			break
		case 3: // Заочная
			competitive.TargetZ = data.CompetitiveGroup.Number
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
					tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *value.Uid).Find(&exist)
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
					tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *value.Uid).Find(&exist)
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

func (result *ResultInfo) AddProgram(idCompetitive uint, data AddCompetitiveGroup) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var competitive digest.CompetitiveGroup
	db := conn.Find(&competitive, idCompetitive)
	if competitive.Id == 0 {
		result.SetErrorResult(`Конкурсная группа не найдена`)
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
				tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *value.Uid).Find(&exist)
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
	db := conn.Find(&competitive, idCompetitive)
	if competitive.Id == 0 {
		result.SetErrorResult(`Конкурсная группа не найдена`)
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
				tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *value.Uid).Find(&exist)
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

func (result *ResultInfo) GetInfoCompetitiveGroup(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var competitive digest.CompetitiveGroup
	db := conn.Preload(`Campaign`).Preload(`EducationForm`).Preload(`EducationLevel`).Preload(`EducationSource`).Preload(`LevelBudget`).Preload(`Direction`).Find(&competitive, ID)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
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
			"id_direction":          competitive.Direction.Id,
			"code_direction":        competitive.Direction.Code,
			"name_direction":        competitive.Direction.Name,
			"name":                  competitive.Name,
			"id_education_form":     competitive.EducationForm.Id,
			"name_education_form":   competitive.EducationForm.Name,
			"id_education_level":    competitive.EducationLevel.Id,
			"name_education_level":  competitive.EducationLevel.Name,
			"id_education_source":   competitive.EducationSource.Id,
			"name_education_source": competitive.EducationSource.Name,
			"id_level_budget":       competitive.LevelBudget.Id,
			"name_level_budget":     competitive.LevelBudget.Name,
			"number":                number,
		}
		var idsPrograms []uint
		db = conn.Where(`id_competitive_group=?`, ID).Table(`cmp.competitive_group_programs`).Pluck(`id`, &idsPrograms)
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
		var idsEntrance []uint
		db = conn.Where(`id_competitive_group=?`, ID).Table(`cmp.entrance_test`).Pluck(`id`, &idsEntrance)
		var entranceTests []interface{}
		for _, id := range idsEntrance {
			var entrance digest.EntranceTest
			db = conn.Preload(`Subject`).Preload(`EntranceTestType`).Find(&entrance, id)
			entranceTests = append(entranceTests, map[string]interface{}{
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
			})
		}
		//for _, campEducLevel := range campEducLevels {
		//	var educLevel digest.EducationLevel
		//	db = conn.Find(&educLevel, campEducLevel.IdEducationLevel)
		//	c.EducationLevels = append(c.EducationLevels, educLevel.Id)
		//	c.EducationLevelsName = append(c.EducationLevelsName, educLevel.Name)
		//}
		result.Done = true
		result.Items = map[string]interface{}{
			`competitive`:        c,
			`education_programs`: programs,
			`entrance_tests`:     entranceTests,
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

func (result *ResultInfo) RemoveEntranceCompetitive(idCompetitive uint, idEntrance uint) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var old digest.CompetitiveGroup
	db := tx.Find(&old, idCompetitive)
	if old.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}
	var entrance digest.EntranceTest
	db = conn.Where(`id_competitive_group=? AND id=?`, idCompetitive, idEntrance).Find(&entrance)
	if entrance.Id > 0 {
		db = tx.Where(`id_competitive_group=? AND id=?`, idCompetitive, idEntrance).Delete(&entrance)
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
	db := tx.Find(&old, idCompetitive)
	if old.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}
	var program digest.CompetitiveGroupProgram
	db = conn.Where(`id_competitive_group=? AND id=?`, idCompetitive, idProgram).Find(&program)
	if program.Id > 0 {
		db = tx.Where(`id_competitive_group=? AND id=?`, idCompetitive, idProgram).Delete(&program)
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
	db := tx.Find(&old, idCompetitive)
	if old.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}

	var programs []digest.CompetitiveGroupProgram
	db = conn.Where(`id_competitive_group=?`, idCompetitive).Find(&programs)
	if len(programs) > 0 {
		db = tx.Where(`id_competitive_group=?`, idCompetitive).Delete(&programs)
	}
	var entrance []digest.EntranceTest
	db = conn.Where(`id_competitive_group=?`, idCompetitive).Find(&entrance)
	if len(entrance) > 0 {
		db = tx.Where(`id_competitive_group=?`, idCompetitive).Delete(&entrance)
	}

	db = tx.Where(`id=?`, idCompetitive).Delete(&old)
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

func CheckCompetitiveGroupByUser(idCompetitiveGroup uint, user digest.User) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var count int
	if user.Role.Code == `administrator` {
		return nil
	}
	db := conn.Table(`cmp.competitive_groups`).Select(`id`).Where(`id_organization=? AND id=?`, user.CurrentOrganization.Id, idCompetitiveGroup).Count(&count)
	if db.Error != nil {
		return db.Error
	}
	if count > 0 {
		return nil
	} else {
		return errors.New(`У пользователя нет доступа к данной конкурсной группе `)
	}
}
