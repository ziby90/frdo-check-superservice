package handlers

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
)

type AddUser struct {
	Id             uint    `json:"id" schema:"id"`
	Login          string  `json:"login" schema:"login"`
	Password       string  `json:"password" schema:"password"`
	Patronymic     *string `json:"patronymic" schema:"patronymic"`
	Surname        string  `json:"surname" schema:"surname"`
	Name           string  `json:"name" schema:"name"`
	Phone          *string `json:"phone" schema:"phone"`
	Email          *string `json:"email" schema:"email"`
	IdOrganization *uint   `json:"id_organization" schema:"id_organization"`
	Comment        *string `json:"comment" schema:"comment"`
}

func GetAuthUser() {

}
func (result *ResultInfo) RegistrationUser(data AddUser, f *digest.File) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	result.PrimaryLogging.SetNewData(data)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var exist digest.User
	db := conn.Where(`upper(login)=?`, strings.ToUpper(data.Login)).Find(&exist)
	if exist.Id > 0 {
		result.SetErrorResult("Логин занят")
		return
	}

	user := digest.User{
		Login:            data.Login,
		Password:         data.Password,
		Patronymic:       data.Patronymic,
		Surname:          data.Surname,
		Name:             data.Name,
		IdRole:           1,
		RegistrationDate: time.Now(),
		Phone:            data.Phone,
		Email:            data.Email,
		Actual:           true,
	}
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&user)
	if db.Error != nil {
		result.SetErrorResult(`Ошибка при добавлении пользователя: ` + db.Error.Error())
		tx.Rollback()
		return
	}

	if f != nil {
		var organization digest.Organization
		conn.Where(`actual is true and id=?`, data.IdOrganization).Find(&organization)
		if organization.Id <= 0 {
			result.SetErrorResult("Организация не найдена")
			return
		}
		var link digest.RequestLinks
		path := getPath(user.Id, link.TableName(), time.Now())
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(link.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			return
		}
		link.IdOrganization = *data.IdOrganization
		link.IdUser = user.Id
		link.ConfirmingDoc = &name
		link.IdStatus = 1 // заявка
		link.CreatedAt = time.Now()
		link.Comment = data.Comment
		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&link)
		if db.Error != nil {
			result.SetErrorResult(`Ошибка при добавлении заявки на связь: ` + db.Error.Error())
			tx.Rollback()
			return
		}
	}

	result.Items = map[string]interface{}{
		`id_user`: user.Id,
	}
	result.Done = true
	tx.Commit()
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
