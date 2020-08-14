package handlers

import (
	"database/sql"
	"log"
	"persons/config"
	"persons/digest"
)

func (result *ResultInfo) GetUserInfoResponse() {
	if &result.User != nil {
		links := GetOrganizationsLinks(&result.User)
		items := make(map[string]interface{})
		items[`user_info`] = &result.User
		items[`organizations_list`] = links
		result.Items = items
		result.Done = true
	} else {
		m := `Авторизация не пройдена.`
		result.Message = &m
	}
}

func (result *ResultInfo) IsAuth(user *digest.User) {
	if user != nil {
		result.Done = true
	} else {
		m := `Авторизация не пройдена.`
		result.Message = &m
	}
}

func GetOrganizationsLinks(user *digest.User) interface{} {
	var links []interface{}
	if user != nil {
		conn := config.Db.ConnGORM
		conn.LogMode(config.Conf.Dblog)
		var rows *sql.Rows
		var err error
		if user.Role.Code == `administrator` {
			rows, err = conn.Table(`admin.organizations`).Select(`id, short_title, ogrn, kpp, id_eiis, is_oovo`).Where(`actual`).Rows()
		} else {
			rows, err = conn.Table(`admin.organizations_users ou`).Where(`id_user=? AND id_status=2`, user.Id).
				Select(`id_organization as id, short_title, ogrn, kpp, org.id_eiis, org.is_oovo`).
				Joins(`right join admin.organizations org ON org.id=ou.id_organization`).
				Rows()
		}
		if err == nil {
			defer func() {
				_ = rows.Close()
			}()
			for rows.Next() {
				var organization struct {
					IdEiis     string `json:"id_eiis"`
					IsOOVO     bool   `json:"is_oovo"`
					Id         uint   `json:"id"`
					ShortTitle string `json:"short_title"`
					Ogrn       string `json:"ogrn"`
					Kpp        string `json:"kpp"`
				}
				//var idOrganization uint
				err := rows.Scan(&organization.Id, &organization.ShortTitle, &organization.Ogrn, &organization.Kpp, &organization.IdEiis, &organization.IsOOVO)
				if err != nil {
					log.Fatal(err)
				}
				//org := digest.Organization{}
				//conn.Find(&org, idOrganization)
				links = append(links, map[string]interface{}{
					`id`:          organization.Id,
					`short_title`: organization.ShortTitle,
					`ogrn`:        organization.Ogrn,
					`kpp`:         organization.Kpp,
					`id_eiis`:     organization.IdEiis,
					`is_oovo`:     organization.IsOOVO,
				})
			}
		}
	}
	return links
}

func SetCurrentOrganization(currentOrgId uint, user *digest.User) uint {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var row *sql.Row
	if user.Role.Code == `administrator` {
		row = conn.Table(`admin.organizations`).Where(`id=?`, currentOrgId).Select(`id`).Row()
	} else {
		row = conn.Table(`admin.organizations_users`).Where(`id_user=? AND id_organization=?`, user.Id, currentOrgId).Select(`id_organization`).Row()
	}
	var organizationId uint
	err := row.Scan(&organizationId)
	if err == nil && organizationId > 0 {
		return organizationId
	}
	return 0
}
