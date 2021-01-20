package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strconv"
	"strings"
	"time"
)

type IdsDirectionOrganization struct {
	IdDirections []uint `json:"directions"`
}

var DirectionsSearchArray = []string{
	`code`,
	`name`,
}

func (result *ResultInfo) GetInfoOrganization() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var organization digest.Organization
	idOrganization := result.User.CurrentOrganization.Id
	db := conn.Find(&organization, idOrganization)
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Организация не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		//db = conn.Model(&campaign).Related(&campaign.CampaignType, `IdCampaignType`)
		//db = conn.Model(&campaign).Related(&campaign.CampaignStatus, `IdCampaignStatus`)

		////var educForms []EducationFormRespons
		//var campEducForms []digest.CampaignEducForm
		//db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducForms)
		//var campEducLevels []digest.CampaignEducLevel
		//db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducLevels)

		c := map[string]interface{}{
			`id`:          organization.Id,
			`short_title`: organization.ShortTitle,
			`created_at`:  organization.CreatedAt,
			`is_oovo`:     organization.IsOOVO,
			`ogrn`:        organization.Ogrn,
			`kpp`:         organization.Kpp,
			`address`:     organization.Address,
			`phone`:       organization.Phone,
		}
		result.Done = true
		result.Items = c
		return
	} else {
		result.Done = true
		message := `Организация не найдены.`
		result.Message = &message
		result.Items = make(map[string]string)
		return
	}
}

func (result *ResultInfo) AddOrganizationDirection(directions IdsDirectionOrganization) {
	result.Done = false
	conn := &config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var organization digest.Organization
	db := conn.Find(&organization, result.User.CurrentOrganization.Id)
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Организация не найдена.`
			result.Message = &message
			tx.Rollback()
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	var idsDirections []uint

	if len(directions.IdDirections) > 0 {
		for _, value := range directions.IdDirections {
			var r digest.Direction
			db = conn.Where(`id=?`, value).Find(&r)
			if r.Id <= 0 {
				result.SetErrorResult(`Направление на найдено`)
				tx.Rollback()
				return
			}
			var exist digest.OrgDirections
			db = conn.Where(`id_eiis=? AND id_direction=? AND actual`, organization.IdEiis, value).Find(&exist)

			if exist.Id > 0 {
				result.SetErrorResult(`Направление ` + fmt.Sprintf(`%v %v`, r.Code, r.Name) + ` у данной организации уже существует`)
				tx.Rollback()
				return
			}
			n := digest.OrgDirections{
				IdDirection: value,
				IdEiis:      organization.IdEiis,
				Code:        &r.Code,
				CreatedAt:   time.Now(),
				Actual:      true,
			}

			db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&n)
			if db.Error != nil {
				result.SetErrorResult(db.Error.Error())
				tx.Rollback()
				return
			}
			idsDirections = append(idsDirections, n.Id)
		}
	} else {
		result.SetErrorResult(`Направления не найдены`)
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_organization`: result.User.CurrentOrganization.Id,
		`id_directions`:   idsDirections,
	}
	result.Done = true
	tx.Commit()
	return
}
func (result *ResultInfo) RemoveOrganizationDirection(directions IdsDirectionOrganization) {
	result.Done = false
	conn := &config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var organization digest.Organization
	db := conn.Find(&organization, result.User.CurrentOrganization.Id)
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Организация не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	var idsDirections []uint

	if len(directions.IdDirections) > 0 {
		for _, value := range directions.IdDirections {
			var r digest.Direction
			db = conn.Where(`id=?`, value).Find(&r)
			if r.Id <= 0 {
				result.SetErrorResult(`Направление на найдено`)
				tx.Rollback()
				return
			}
			var exist digest.OrgDirections
			db = conn.Where(`id_eiis=? AND id_direction=? AND actual`, organization.IdEiis, value).Find(&exist)

			if exist.Id <= 0 {
				result.SetErrorResult(`Направление ` + fmt.Sprintf(`%v %v`, r.Code, r.Name) + ` у данной организации не найдено`)
				tx.Rollback()
				return
			}
			exist.Actual = false
			now := time.Now()
			exist.Changed = &now

			db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&exist)
			if db.Error != nil {
				result.SetErrorResult(db.Error.Error())
				tx.Rollback()
				return
			}
			idsDirections = append(idsDirections, exist.Id)
		}
	} else {
		result.SetErrorResult(`Направления не найдены`)
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_organization`: result.User.CurrentOrganization.Id,
		`id_directions`:   idsDirections,
	}
	result.Done = true
	tx.Commit()
	return
}

func (result *Result) GetDirectionsByOrganization(keys map[string][]string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var orgDirections []digest.VOrganizationsDirections
	db := conn.Where(`actual is TRUE AND id_organization=?`, result.User.CurrentOrganization.Id)
	if len(keys[`search_parent`]) > 0 {
		db = db.Where(`(UPPER(code_parent || ' ' || name_parent) LIKE ?)`, `%`+strings.ToUpper(keys[`search_parent`][0])+`%`)
	}
	if len(keys[`search_specialty`]) > 0 {
		db = db.Where(`(UPPER(code || ' ' || name) LIKE ?)`, `%`+strings.ToUpper(keys[`search_specialty`][0])+`%`)
	}
	if len(keys[`search_education_level`]) > 0 {
		var idsEducLevel []uint
		for _, val := range keys[`search_education_level`] {
			id, err := strconv.ParseInt(val, 10, 32)
			if err == nil {
				idsEducLevel = append(idsEducLevel, uint(id))
			}
		}
		if len(idsEducLevel) > 0 {
			db = db.Where(`id_education_level IN (?)`, idsEducLevel)
		}

	}
	dbCount := db.Model(&orgDirections).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {
	}

	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Order(`code_parent asc`).Find(&orgDirections)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Направления не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for index, _ := range orgDirections {
			c := map[string]interface{}{
				`id`:                   orgDirections[index].Id,
				`name`:                 orgDirections[index].Name,
				`code`:                 orgDirections[index].Code,
				`id_parent`:            orgDirections[index].IdParent,
				`name_parent`:          orgDirections[index].NameParent,
				`code_parent`:          orgDirections[index].CodeParent,
				`id_education_level`:   orgDirections[index].IdEducationLevel,
				`name_education_level`: orgDirections[index].NameEducationLevel,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Направления не найдены.`
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}

func (result *ResultList) GetDirectionsSelectListByOrg(keys map[string][]string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var items []RowsCls
	sortField := `id`
	sortOrder := `asc`
	db := conn.Where(`id_organization=? AND actual`, result.User.CurrentOrganization.Id).Order(sortField + ` ` + sortOrder).Table(`admin.v_org_directions`)
	if len(keys[`id_education_level`]) > 0 {
		if v, ok := strconv.Atoi(keys[`id_education_level`][0]); ok == nil {
			db = db.Where(`id_education_level=?`, v)
		}
	}
	if len(keys[`id_parent`]) > 0 {
		if v, ok := strconv.Atoi(keys[`id_parent`][0]); ok == nil {
			db = db.Where(`id_parent=?`, v)
		}
	}
	if result.Search != `` {
		db = db.Where(`UPPER((code || ' ' || name)) like ?`, `%`+strings.ToUpper(result.Search)+`%`)
	}
	db = db.Select(`id, (code || ' ' || name) as name`).Find(&items)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Направления не найдены.`
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
				`id`:   item.Id,
				`name`: item.Name,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Направления не найдены.`
		result.Message = &message
		result.Items = []digest.VOrganizationsDirections{}
		return
	}
}
func (result *ResultList) GetDirectionsParentsSelectListByOrg(keys map[string][]string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var items []RowsCls
	sortField := `id`
	sortOrder := `asc`
	db := conn.Where(`id_organization=? AND actual`, result.User.CurrentOrganization.Id).Order(sortField + ` ` + sortOrder).Table(`admin.v_org_directions`)
	if len(keys[`id_education_level`]) > 0 {
		if v, ok := strconv.Atoi(keys[`id_education_level`][0]); ok == nil {
			db = db.Where(`id_education_level=?`, v)
		}
	}
	if result.Search != `` {
		db = db.Where(`UPPER((code_parent || ' ' || name_parent)) like ?`, `%`+strings.ToUpper(result.Search)+`%`)
	}
	db = db.Select(`id_parent as id, (code_parent || ' ' || name_parent) as name`).Group(`id_parent, code_parent, name_parent`).Find(&items)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Направления не найдены.`
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
				`id`:   item.Id,
				`name`: item.Name,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Направления не найдены.`
		result.Message = &message
		result.Items = []digest.VOrganizationsDirections{}
		return
	}
}

func CheckOrgCookie(user digest.User, r *http.Request) uint {
	currentOrgId := uint(0)
	cookieOrg, err := r.Cookie(`current-org`)
	if err == nil {
		u64, err := strconv.ParseUint(cookieOrg.Value, 10, 32)
		if err == nil {
			currentOrgId = uint(u64)
		}
	}
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var row *sql.Row
	if user.Role.Code == `administrator` {
		if currentOrgId == 0 {
			return 0
		}
		row = conn.Table(`admin.organizations`).Where(`id=?`, currentOrgId).Select(`id`).Row()

	} else {
		row = conn.Table(`admin.organizations_users`).Where(`id_user=? AND id_organization=?`, user.Id, currentOrgId).Select(`id_organization`).Row()
	}

	if err == nil {
		var organizationId uint
		err = row.Scan(&organizationId)
		if organizationId > 0 {
			return organizationId
		}
	}
	return 0
}

func (result *ResultInfo) SetIsOOVOOrganization(isOOVO bool) {
	conn := &config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var organization digest.Organization
	db := conn.Where(`id=?`, result.User.CurrentOrganization.Id).Find(&organization)
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Организация не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	organization.IsOOVO = isOOVO
	t := time.Now()
	organization.Changed = &t
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&organization)
	if db.Error != nil {
		result.SetErrorResult(`Ошибка при изменении параметров отправки у организации ` + db.Error.Error())
		tx.Rollback()
		return
	}
	result.Items = map[string]interface{}{
		`id_organization`: organization.Id,
		`is_oovo`:         organization.IsOOVO,
	}
	result.Done = true
	tx.Commit()
	return
}

func (result *ResultInfo) GetListOrganizationShort() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var organizations []digest.Organization
	db := conn.Where(`actual is true`).Find(&organizations)
	var orgs []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Организации не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, organization := range organizations {
			c := map[string]interface{}{
				`id`:          organization.Id,
				`short_title`: organization.ShortTitle,
				`created_at`:  organization.CreatedAt,
				`is_oovo`:     organization.IsOOVO,
				`ogrn`:        organization.Ogrn,
				`kpp`:         organization.Kpp,
			}
			orgs = append(orgs, c)
		}
		result.Done = true
		result.Items = orgs
		return
	} else {
		result.Done = true
		message := `Организации не найдены.`
		result.Message = &message
		result.Items = []digest.Organization{}
		return
	}
}
