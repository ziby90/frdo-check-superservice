package handlers

import (
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
		rows, err := conn.Table(`admin.organizations_users`).Where(`id_user=?`, user.Id).Select(`id_organization`).Rows()
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id_organization uint
				err := rows.Scan(&id_organization)
				if err != nil {
					log.Fatal(err)
				}
				org := digest.Organization{}
				conn.Find(&org, id_organization)
				links = append(links, map[string]interface{}{
					`id`:          org.Id,
					`short_title`: org.ShortTitle,
					`ogrn`:        org.Ogrn,
					`kpp`:         org.Kpp,
				})
			}
		}
	}
	return links
}

func SetCurrentOrganization(currentOrgId uint, user *digest.User) uint {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	row := conn.Table(`admin.organizations_users`).Where(`id_user=? AND id_organization=?`, user.Id, currentOrgId).Select(`id_organization`).Row()
	var organizationId uint
	err := row.Scan(&organizationId)
	if err == nil && organizationId > 0 {
		return organizationId
	}
	return 0
}
