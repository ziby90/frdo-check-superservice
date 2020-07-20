package handlers_admin

import (
	"crypto/sha1"
	"encoding/hex"
	sendToEpguPath "gitlab.com/unkal/sendtoepgu/path_files"
	"io"
	"os"
	"path/filepath"
	"persons/config"
	"persons/digest"
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

func getPath(idEntrant uint, category string, t time.Time) string {
	path, _ := sendToEpguPath.GetPath(idEntrant, t, category)
	//path := `./uploads/docs/` + fmt.Sprintf(`%v`, idEntrant) + `/` + category + `/` + t.Format(`02-01-2006`)
	return path
}

var UserSearchArray = []string{
	`name`,
	`surname`,
	`patronymic`,
	`login`,
}

func (result *Result) GetListUsers() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var items []digest.User

	db := conn.Preload(`Role`).Preload(`Region`)
	sortField := `registration_date`
	sortOrder := `desc`
	if result.Sort.Field != `` {
		sortField = result.Sort.Field
	} else {
		result.Sort.Field = sortField
	}
	db = db.Order(result.Sort.Field + ` ` + result.Sort.Order)
	for _, search := range result.Search {
		db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)

	}
	dbCount := db.Model(&items).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Order(sortField + ` ` + sortOrder).Find(&items)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Пользователи не найдены.`
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
				`id`:                item.Id,
				`login`:             item.Login,
				`name`:              item.Name,
				`surname`:           item.Surname,
				`patronymic`:        item.Patronymic,
				`id_role`:           item.Role.Id,
				`name_role`:         item.Role.Name,
				`code_role`:         item.Role.Code,
				`id_region`:         item.Region.Id,
				`name_region`:       item.Region.Name,
				`registration_date`: item.RegistrationDate,
				`actual`:            item.Actual,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Пользоватекли не найдены.`
		result.Message = &message
		result.Items = []digest.User{}
		return
	}
}
func (result *ResultInfo) GetInfoUser(idUser uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item digest.User
	db := conn.Preload(`Region`).Preload(`Role`).Find(&item, idUser)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Пользователь не найден.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {

		c := map[string]interface{}{
			`id`:                item.Id,
			`login`:             item.Login,
			`name`:              item.Name,
			`surname`:           item.Surname,
			`patronymic`:        item.Patronymic,
			`work`:              item.Work,
			`post`:              item.Post,
			`adress`:            item.Adress,
			`email`:             item.Email,
			`snils`:             item.Snils,
			`id_author`:         item.IdAuthor,
			`changed`:           item.Changed,
			`id_role`:           item.Role.Id,
			`name_role`:         item.Role.Name,
			`code_role`:         item.Role.Code,
			`id_region`:         item.Region.Id,
			`name_region`:       item.Region.Name,
			`registration_date`: item.RegistrationDate,
			`actual`:            item.Actual,
		}
		result.Done = true
		result.Items = c
		return
	} else {
		result.SetErrorResult(`Пользователь не найден.`)
		result.Items = []digest.User{}
		return
	}
}
func (result *ResultInfo) BlockUser(idUser uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item digest.User
	db := conn.Where(`actual is true`).Find(&item, idUser)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Пользователь не найден.`
			result.Message = &message
			return
		}
		result.SetErrorResult(`Ошибка подключения к БД. `)
		return
	}
	if db.RowsAffected > 0 {
		//item.Actual = false
		//t := time.Now()
		//item.Changed = &t
		//item.IdAuthor = &result.User.Id
		db := conn.Table(item.TableName()).Where(`id=?`, idUser).Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Updates(map[string]interface{}{"changed": time.Now(), "id_author": result.User.Id, "actual": false})
		if db.Error != nil {
			result.SetErrorResult(db.Error.Error())
			return
		}
		result.Done = true
		result.Items = map[string]interface{}{
			`id_user`: idUser,
		}
		return
	} else {
		result.Done = true
		message := `Пользователь не найден.`
		result.Message = &message
		result.Items = []digest.User{}
		return
	}
}
func (result *ResultInfo) UnblockUser(idUser uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item digest.User
	db := conn.Where(`actual is false`).Find(&item, idUser)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Пользователь не найден.`
			result.Message = &message
			return
		}
		result.SetErrorResult(`Ошибка подключения к БД. `)
		return
	}
	if db.RowsAffected > 0 {
		//item.Actual = true
		//t := time.Now()
		//item.Changed = &t
		//item.IdAuthor = &result.User.Id
		db := conn.Table(item.TableName()).Where(`id=?`, idUser).Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Updates(map[string]interface{}{"changed": time.Now(), "id_author": result.User.Id, "actual": true})

		db = conn.Model(&item).Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Update(&item)
		if db.Error != nil {
			result.SetErrorResult(db.Error.Error())
			return
		}
		result.Done = true
		result.Items = map[string]interface{}{
			`id_user`: idUser,
		}
		return
	} else {
		result.Done = true
		message := `Пользователь не найден.`
		result.Message = &message
		result.Items = []digest.User{}
		return
	}
}
func (result *Result) GetLinksUser(idUser uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item digest.User
	db := conn.Find(&item, idUser)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Пользователь не найден.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		var links []digest.OrganizationsUsers
		db := conn.Where(`id_user=? AND id_status=2`, idUser).Preload(`Organization`)
		dbCount := db.Model(&links).Count(&result.Paginator.TotalCount)
		if dbCount.Error != nil {

		}
		result.Paginator.Make()
		db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&links)

		if len(links) <= 0 {
			result.SetErrorResult(`Связи не найдены`)
			result.Done = true
			result.Items = []digest.OrganizationsUsers{}
			return
		}
		var c []interface{}
		for _, value := range links {
			hasFile := false
			if value.ConfirmingDoc != nil {
				hasFile = true
			}
			c = append(c, map[string]interface{}{
				`id`:              value.Id,
				`has_file`:        hasFile,
				`file`:            value.ConfirmingDoc,
				`id_organization`: value.Organization.Id,
				`ogrn`:            value.Organization.Ogrn,
				`kpp`:             value.Organization.Kpp,
				`short_title`:     value.Organization.ShortTitle,
			})
		}
		result.Done = true
		result.Items = c
		return
	} else {
		result.Done = true
		message := `Пользователь не найден.`
		result.Message = &message
		result.Items = []digest.User{}
		return
	}
}
func (result *ResultInfo) AddLinksToUser(idUser uint, idOrganization uint, f *digest.File, comment *string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)

	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var user digest.User
	db := conn.Where(`id=?`, idUser).Find(&user)
	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			message := "Пользователь не найден."
			result.Message = &message
			return
		}
		message := "Ошибка подключения к БД."
		result.Message = &message
		return
	}
	var organization digest.Organization
	conn.Where(`actual is true and id=?`, idOrganization).Find(&organization)
	if organization.Id <= 0 {
		message := "Организация не найдена"
		result.Message = &message
		return
	}
	var link digest.OrganizationsUsers
	var exist digest.OrganizationsUsers
	db = tx.Where(`id_user=? AND id_organization=? AND id_status IN (1,2) `, idUser, idOrganization).Find(&exist)
	if exist.Id > 0 {
		result.SetErrorResult(`Данная связь уже существует или на рассмотрении`)
		tx.Commit()
		return
	}
	path := getPath(idUser, link.TableName(), time.Now())
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
	link.IdOrganization = idOrganization
	link.IdUser = idUser
	link.ConfirmingDoc = &name
	link.IdStatus = 2
	link.IdAuthor = &result.User.Id
	link.Created = time.Now()
	link.Comment = comment
	result.PrimaryLogging.SetNewData(link)
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&link)
	if db.Error != nil {
		m := `Ошибка при добавлении связи: ` + db.Error.Error()
		result.Message = &m
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_user`:         idUser,
		`id_organization`: idOrganization,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) CreateUser(data AddUser, f *digest.File) {
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
		IdAuthor:         &result.User.Id,
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
		var link digest.OrganizationsUsers
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
		link.IdStatus = 2
		link.IdAuthor = &result.User.Id
		link.Created = time.Now()
		link.Comment = data.Comment
		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&link)
		if db.Error != nil {
			result.SetErrorResult(`Ошибка при добавлении связи: ` + db.Error.Error())
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
func (result *ResultInfo) ResetPasswordUser(idUser uint, password string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	if password == "" {
		result.SetErrorResult(`Аня, ты опять шлешь pass??`)
		tx.Rollback()
		return
	}
	var user digest.User
	conn.Where(`id=?`, idUser).Find(&user)
	if user.Id <= 0 {
		result.SetErrorResult(`Пользователь не найден`)
		tx.Rollback()
		return
	}
	db := conn.Table(user.TableName()).Where(`id=?`, idUser).Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).
		Updates(map[string]interface{}{"changed": time.Now(), "id_author": result.User.Id, "password": password})
	if db.Error != nil {
		result.SetErrorResult(`Ошибка при обновлении пароля пользователя: ` + db.Error.Error())
		tx.Rollback()
		return
	}
	result.Items = map[string]interface{}{
		`id_user`: user.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) BreakOff(idUser uint, idLink uint, comment string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var user digest.User
	conn.Where(`id=?`, idUser).Find(&user)
	if user.Id <= 0 {
		result.SetErrorResult(`Пользователь не найден`)
		tx.Rollback()
		return
	}
	var link digest.OrganizationsUsers
	conn.Where(`id=? and id_user=? and id_status=2`, idLink, idUser).Find(&link)
	if link.Id <= 0 {
		result.SetErrorResult(`Связь не найдена`)
		tx.Rollback()
		return
	}
	db := conn.Table(link.TableName()).Where(`id=?`, idLink).Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).
		Updates(map[string]interface{}{"changed": time.Now(), "id_author": result.User.Id, "comment": comment, "id_status": 4})
	if db.Error != nil {
		result.SetErrorResult(`Ошибка при разрые связи пользователя: ` + db.Error.Error())
		tx.Rollback()
		return
	}
	result.Items = map[string]interface{}{
		`id_user`: user.Id,
		`id_link`: link.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) GetFileLink(idLink uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var doc digest.OrganizationsUsers
	db := conn.Where(`id=?`, idLink).Find(&doc)
	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			result.SetErrorResult("Связь не найдена.")
			result.Items = []interface{}{}
			return
		}
		result.SetErrorResult("Ошибка подключения к БД.")
		return
	}
	if doc.ConfirmingDoc != nil && *doc.ConfirmingDoc != `` {
		filename := *doc.ConfirmingDoc
		path := getPath(doc.IdUser, doc.TableName(), doc.Created) + filename
		result.Items = path
	} else {
		result.SetErrorResult("Файл не найден.")
		return
	}
	result.Done = true
	return
}
