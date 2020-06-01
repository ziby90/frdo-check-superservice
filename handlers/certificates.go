package handlers

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"persons/certs"
	"persons/config"
	"time"
)

func (result *ResultInfo) CheckFile(file certs.File) {
	conn := &config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	operationResult, cert, err := certs.SignCheckRequest(file.Content)
	if err == nil && operationResult == `success` {
		var oldCert certs.Certificate
		db := conn.Where(`id_organization=? AND deleted_at IS NULL`, result.User.CurrentOrganization.Id).Find(&oldCert)
		if cert.ParsedSubject.OGRN != result.User.CurrentOrganization.Ogrn {
			result.SetErrorResult(`ОГРН выбранной организации не совпадает с ОГРН , указанным в сертификате`)
			tx.Rollback()
			return
		}
		if oldCert.ID > 0 {
			db = tx.Where(`id=?`, oldCert.ID).Delete(&oldCert)
			if db.Error != nil {
				result.SetErrorResult(`Ошибка при удалении предыдущего сертфииката`)
				tx.Rollback()
				return
			}
		}
		cert.IdOrganization = result.User.CurrentOrganization.Id

		path := getPath(result.User.CurrentOrganization.Id, cert.TableName(), time.Now())
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		// write the whole body at once
		ext := filepath.Ext(path + `/` + file.Title)
		sha1FileName := sha1.Sum([]byte(cert.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		err = ioutil.WriteFile(path+name, file.Content, 0644)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		cert.Link = name

		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&cert)
		if db.Error != nil {
			result.SetErrorResult(`Ошибка при добавлении сертификата ` + db.Error.Error())
			tx.Rollback()
			return
		}
		var oldSubject certs.Subject
		db = conn.Where(`id_certificate=? AND deleted_at IS NULL`, oldCert.ID).Find(&oldSubject)
		if oldSubject.ID > 0 {
			db = tx.Where(`id=?`, oldSubject.ID).Delete(&oldSubject)
			if db.Error != nil {
				result.SetErrorResult(`Ошибка при удалении предыдущего обьекта сертификата`)
				tx.Rollback()
				return
			}
		}

		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&cert.ParsedSubject)
		if db.Error != nil {
			result.SetErrorResult(`Ошибка при добавлении обьекта сертификата ` + db.Error.Error())
			tx.Rollback()
			return
		}
		tx.Commit()

		result.Done = true
		result.Items = map[string]interface{}{
			`id_certificate`:  cert.ID,
			`id_cert_subject`: cert.ParsedSubject.ID,
		}
		return
	}
	if err != nil {
		result.SetErrorResult(err.Error())
		return
	}
	result.SetErrorResult(`Неудача почему то`)
	return

}

func (result *ResultInfo) GetOrganizationCertificate() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var cert certs.Certificate
	conn.Where(`id_organization=? AND deleted_at IS NULL`, result.User.CurrentOrganization.Id).Find(&cert)
	if cert.ID <= 0 {
		result.SetErrorResult(`Не найден сертификат организации`)
		return
	}
	response := map[string]interface{}{
		`id`:            cert.ID,
		`subject`:       cert.Subject,
		`issuer`:        cert.Issuer,
		`serial_number`: cert.SerialNumber,
		`sha1hash`:      cert.SHA1Hash,
		`subjKey_id`:    cert.SubjKey,
		`not_before`:    cert.NotBefore,
		`not_after`:     cert.NotAfter,
		`ogrn`:          nil,
		`email`:         nil,
		`surname`:       nil,
		`name`:          nil,
		`patronymic`:    nil,
	}
	var subject certs.Subject
	conn.Where(`id_certificate=? AND deleted_at IS NULL`, cert.ID).Find(&subject)
	if subject.ID > 0 {
		response[`ogrn`] = subject.OGRN
		response[`email`] = subject.E
		response[`surname`] = subject.SN
		response[`name`] = subject.NAME
		response[`patronymic`] = subject.PATRONYMIC
	}

	result.Done = true
	result.Items = response

}

func CheckHandlerPost(w http.ResponseWriter, r *http.Request) {

}
