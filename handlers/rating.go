package handlers

import (
	"encoding/json"
	"fmt"
	"persons/config"
	"persons/digest"
	"persons/service"
	"time"
)

type EditRatingCompetitiveGroupsApplication struct {
	Id               uint  `gorm:"primary_key"`
	CommonRating     int64 `json:"common_rating"`
	FirstRating      int64 `json:"first_rating"`
	AgreedRating     int64 `json:"agreed_rating"`
	ChangeRating     int64 `json:"change_rating"`
	CountFirstStep   int64 `json:"count_first_step"`
	CountSecondStep  int64 `json:"count_second_step"`
	CountApplication int64 `json:"count_application"`
	CountAgreed      int64 `json:"count_agreed"`
}

func (result *Result) GetListRatingCompetitiveGroupsApplication(idCompetitiveGroup uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var items []digest.RatingCompetitiveGroupsApplications
	sortField := `changed`
	sortOrder := `desc`
	if result.Sort.Field != `` {
		sortField = result.Sort.Field
	} else {
		result.Sort.Field = sortField
	}

	fmt.Print(result.Sort.Field, sortField)
	db := conn.Where(`id_organization=? AND id_competitive_group=?`, result.User.CurrentOrganization.Id, idCompetitiveGroup)

	dbCount := db.Model(&items).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Order(sortField + ` ` + sortOrder).Find(&items)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Заявления не найдены.`
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
				`id`:                   item.Id,
				`id_organization`:      item.IdOrganization,
				`id_competitive_group`: item.IdCompetitiveGroup,
				`id_application`:       item.IdApplication,
				`admission_volume`:     item.AdmissionVolume,
				`common_rating`:        item.CommonRating,
				`first_rating`:         item.FirstRating,
				`agreed_rating`:        item.AgreedRating,
				`change_rating`:        item.ChangeRating,
				`count_first_step`:     item.CountFirstStep,
				`count_second_step`:    item.CountSecondStep,
				`count_application`:    item.CountApplication,
				`count_agreed`:         item.CountAgreed,
				`changed`:              item.Changed,
				`surname`:              nil,
				`name`:                 nil,
				`patronymic`:           nil,
				`uid_epgu`:             nil,
				`app_number`:           nil,
			}
			application, _ := digest.GetApplication(item.IdApplication)
			if application.Id > 0 {
				c[`uid_epgu`] = application.UidEpgu
				c[`app_number`] = application.AppNumber
				var entrant digest.Entrants
				db = conn.Find(&entrant, application.EntrantsId)
				if entrant.Id > 0 {
					c[`surname`] = entrant.Surname
					c[`name`] = entrant.Name
					c[`patronymic`] = entrant.Patronymic
				}
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Заявления не найдены.`
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}
func (result *ResultInfo) EditRatingCompetitiveGroupsApplication(data EditRatingCompetitiveGroupsApplication) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item digest.RatingCompetitiveGroupsApplications

	db := conn.Where(`id_organization=? AND id=?`, result.User.CurrentOrganization.Id, data.Id).Find(&item)
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Заявление не найдено.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	oldData, _ := json.Marshal(item)
	strOldData := string(oldData)
	result.PrimaryLogging.OldData = &strOldData
	newData, _ := json.Marshal(data)
	strNewData := string(newData)
	result.PrimaryLogging.NewData = &strNewData

	if db.RowsAffected > 0 {
		item.CommonRating = data.CommonRating
		item.FirstRating = data.FirstRating
		item.AgreedRating = data.AgreedRating
		item.ChangeRating = data.ChangeRating
		item.CountFirstStep = data.CountFirstStep
		item.CountSecondStep = data.CountSecondStep
		item.CountApplication = data.CountApplication
		item.CountAgreed = data.CountAgreed
		item.Changed = time.Now()
		db = conn.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&item)
		result.Done = true
		return
	} else {
		result.Done = true
		message := `Заявление не найдено.`
		result.Message = &message
		return
	}
}
func (result *ResultInfo) RefreshListRatingCompetitiveGroupsApplication(idCompetitiveGroup uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)

	var competitive digest.CompetitiveGroup
	var admissionVolume digest.AdmissionVolume
	db := conn.Where(`id=? and id_organization=?`, idCompetitiveGroup, result.User.CurrentOrganization.Id).Find(&competitive)
	if competitive.Id <= 0 {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		return
	}
	db = conn.Table(`cmp.admission_volume`).Where(`id_education_level=? and id_campaign=? and id_direction=? and id_organization=?`, competitive.IdEducationLevel, competitive.IdCampaign, competitive.IdDirection, competitive.IdOrganization).Find(&admissionVolume)
	if admissionVolume.Id <= 0 {
		result.SetErrorResult(`КЦП не найдены`)
		return
	}
	var number int64
	switch competitive.IdEducationSource {
	case 1: // Бюджетные места
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			number = admissionVolume.BudgetO
			break
		case 2: // Очно-заочная(вечерняя)
			number = admissionVolume.BudgetOz
			break
		case 3: // Заочная
			number = admissionVolume.BudgetZ
			break
		default:
			result.SetErrorResult(`Ошибка`)
			return
		}
		break
	case 2: // Квота приема лиц, имеющих особое право
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			number = admissionVolume.QuotaO
			break
		case 2: // Очно-заочная(вечерняя)
			number = admissionVolume.QuotaOz
			break
		case 3: // Заочная
			number = admissionVolume.QuotaZ
			break
		default:
			result.SetErrorResult(`Ошибка`)
			return
		}
		break
	case 3: // С оплатой обучения
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			number = admissionVolume.PaidO
			break
		case 2: // Очно-заочная(вечерняя)
			number = admissionVolume.PaidOz
			break
		case 3: // Заочная
			number = admissionVolume.PaidZ
			break
		default:
			result.SetErrorResult(`Ошибка`)
			return
		}
		break
	case 4: // Целевой прием
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			number = admissionVolume.TargetO
			break
		case 2: // Очно-заочная(вечерняя)
			number = admissionVolume.TargetOz
			break
		case 3: // Заочная
			number = admissionVolume.TargetZ
			break
		default:
			result.SetErrorResult(`Ошибка`)
			return
		}
	default:
		result.SetErrorResult(`Ошибка`)
		return
	}
	cmdRemove := ` DELETE FROM rating.completitive_groups_applications cga WHERE cga.id_organization = ? AND
			  cga.id_competitive_group = ?
		  AND EXISTS(SELECT 1 FROM app.applications a WHERE a.id_organization = ? AND a.id_competitive_group = ?
		  AND actual IS false AND cga.id_application = a.id);`
	db = conn.Exec(cmdRemove, result.User.CurrentOrganization.Id, idCompetitiveGroup, result.User.CurrentOrganization.Id, idCompetitiveGroup)
	countRemove := db.RowsAffected
	cmd := `WITH a AS(SELECT a.id as id_application, a.id_organization, a.id_competitive_group FROM app.applications a WHERE a.id_organization = ? AND
		  a.id_competitive_group = ? AND id_status not in (1,2,3,4,10,17,21,22,24,25) AND actual IS true AND uid_epgu IS NOT NULL
		  AND NOT EXISTS(SELECT 1 FROM rating.completitive_groups_applications cga WHERE cga.id_organization = ? AND
		  cga.id_competitive_group = ? AND cga.id_application = a.id))
		  INSERT INTO rating.completitive_groups_applications(id_application, id_organization, id_competitive_group, admission_volume,
		  common_rating,
		  first_rating,
		  agreed_rating,
		  change_rating,
		  count_first_step,
		  count_second_step,
		  count_application,
		  count_agreed,
		  changed) SELECT id_application, id_organization, id_competitive_group, ? as admission_volume,
		  0 AS common_rating,
		  0 AS first_rating,
		  0 AS agreed_rating,
		  0 AS change_rating,
		  0 AS count_first_step,
		  0 AS count_second_step,
		  0 AS count_application,
		  0 AS count_agreed,
		  now() as changed FROM a;`
	db = conn.Exec(cmd, result.User.CurrentOrganization.Id, idCompetitiveGroup, result.User.CurrentOrganization.Id, idCompetitiveGroup, number)
	countAdd := db.RowsAffected
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		return
	} else {
		result.Items = map[string]interface{}{
			`count_add`:    countAdd,
			`count_remove`: countRemove,
		}
		result.Done = true
		return
	}

}
