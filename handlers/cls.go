package handlers

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
)

type RowsCls struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type SysRows struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	NameTable string `json:"name_table"`
}

type Subjects struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Year     int    `json:"year"`
	MinScore int    `json:"min_score"`
}

var ListClsTableName = []string{
	`achievement_categories`,
	`appeal_statuses`,
	`benefits`,
	`campaign_statuses`,
	`campaign_types`,
	`composition_themes`,
	`compatriot_categories`,
	`directions`,
	`disability_types`,
	`document_categories`,
	`document_sys_category`,
	`v_document_types`,
	`document_types`,
	`education_forms`,
	`education_levels`,
	`education_sources`,
	`entrance_test_types`,
	`gender`,
	`level_budget`,
	`military_categories`,
	`olympic_diploma_types`,
	`olympic_levels`,
	`v_okcm`,
	`orphan_categories`,
	`parents_lost_categories`,
	`radiation_work_categories`,
	`regions`,
	`subjects`,
	`veteran_categories`,
	`v_direction_specialty`,
	`v_edu_levels_campaign_types`,
}

var ListFilterColumns = []string{
	`name_table`,
	`id_campaign_types`,
	`id_education_level`,
	`id_parent`,
	`is_ege`,
}

func (result *ResultCls) GetClsResponse(clsName string) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var db *gorm.DB
	var r []RowsCls
	var fields []string
	if service.SearchStringInSliceString(clsName, ListClsTableName) >= 0 {
		fields = []string{`id`, `name`}
		db = conn.Table(`cls.` + clsName)
	} else {
		switch clsName {
		case `v_okso_enlarged_group`:
			fields = []string{`id`, `(code || ' ' || name)`}
			db = conn.Table(`cls.` + clsName)
			break
		case `v_okso_specialty`:
			fields = []string{`id`, `(code || ' ' || name)`}
			db = conn.Table(`cls.` + clsName)
			break
		default:
			message := `Неизвестный справочник.`
			result.Message = &message
			return
		}

	}
	if result.Search != `` {
		db = db.Where(`UPPER(`+fields[1]+`) like ?`, `%`+strings.ToUpper(result.Search)+`%`)
	}
	fmt.Println(result.Filter, result.Value)
	fmt.Println(service.SearchStringInSliceString(result.Filter, ListFilterColumns) >= 0)

	if result.Filter != `` && result.Value != `` {

		if service.SearchStringInSliceString(result.Filter, ListFilterColumns) >= 0 {
			if strings.HasPrefix(result.Filter, `id`) {
				db = db.Where(`(`+result.Filter+`) = ?`, result.Value)
			} else {
				if clsName == `subjects` && result.Filter == `is_ege` {
					if result.Value == `true` {
						db = db.Where(`ege is TRUE`)
					} else {
						db = db.Where(`ege is FALSE`)
					}

				} else {
					db = db.Where(`UPPER(`+result.Filter+`) like ?`, `%`+strings.ToUpper(result.Value)+`%`)
				}

			}

		}
	}
	db.Where(`actual`).Select(strings.Join(fields, `,`) + ` as name`).Scan(&r)
	result.Items = r

	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Формы образования не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	result.Done = true
	return
}

func (result *ResultInfo) GetSubjectsNoEge(idCampaign uint) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	db := conn.Find(&campaign, idCampaign)
	if campaign.Id <= 0 || db.Error != nil {
		message := `Компания не найдена.`
		result.Message = &message
		return
	}
	var subjects []Subjects
	conn.Raw(`SELECT 
				s.id
				, s.name
				, COALESCE(ss.min_score, 0) AS min_score 
				, coalesce(ss.ege_year , ? ) as year
				FROM cls.subjects s
				LEFT JOIN cls.min_score_subjects ss ON s.id = ss.id_subject AND ss.ege_year = ? WHERE s.ege IS TRUE
				`, campaign.YearStart, campaign.YearStart).Scan(&subjects)
	result.Items = subjects
	result.Done = true
	return

}

func (result *ResultCls) GetClsSysCategoryResponse() {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var r []SysRows

	db := conn.Table(`cls.document_sys_categories`).Select("id, name, name_table").Scan(&r)
	result.Items = r

	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Справочник не найден.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	result.Done = true
	return
}
