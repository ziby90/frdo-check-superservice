package handlers

import (
	"net/http"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
)

func GetAuthUser() {

}

func CheckAuthBase(login string, hash string) ResultAuth {
	var user digest.User
	res := ResultAuth{
		User:    nil,
		Done:    false,
		Message: "",
	}
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	conn.Where(`UPPER(login)=?`, strings.ToUpper(login)).Find(&user)
	if hash == user.Password {
		conn.Model(&user).Related(&user.Role, `IdRole`)
		res.Done = true
		res.Links = GetOrganizationsLinks(&user)
		res.User = &user
	} else {
		res.Message = `Неверный логин или пароль`
	}
	return res
}
func CheckAuthCookie(r *http.Request) *digest.User {
	login := ``
	pass := ``
	cookieLogin, err := r.Cookie(`login`)
	if err == nil {
		login = strings.ToUpper(cookieLogin.Value)
	}
	cookiePass, err := r.Cookie(`password`)
	if err == nil {
		pass = cookiePass.Value
	}

	if login != `` && pass != `` {
		var user digest.User
		conn := config.Db.ConnGORM
		conn.LogMode(config.Conf.Dblog)
		conn.Where(`UPPER(login)=UPPER(?)`, login).Find(&user)
		conn.Model(&user).Related(&user.Role, `IdRole`)
		userAgent := r.Header.Get(`User-agent`)
		if pass == service.GetHash(user.Password+userAgent, true) {
			orgId := CheckOrgCookie(user, r)
			if orgId > 0 {
				org := digest.Organization{}
				conn.Find(&org, orgId)
				organization := digest.CurrentOrganization{
					Id:         org.Id,
					ShortTitle: org.ShortTitle,
					Ogrn:       org.Ogrn,
					Kpp:        org.Kpp,
					IdEiis:     org.IdEiis,
					IsOOVO:     org.IsOOVO,
				}
				user.CurrentOrganization = &organization
			}
			return &user
		}
	}
	return nil
}
