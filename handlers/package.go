package handlers

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"persons/config"
	"persons/digest"
	"time"
)

func (result *Result) GetMarkEgePackages() {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var packages []digest.MarkEgePackages
	if result.Sort.Field == `` {
		result.Sort.Field = `created`
	}
	if result.Sort.Order == `` {
		result.Sort.Order = `asc`
	}
	db := conn.Order(result.Sort.Field + ` ` + result.Sort.Order)
	db = db.Where(`id_organization=?`, result.User.CurrentOrganization.Id)

	dbCount := db.Model(&packages).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Preload(`Status`).Find(&packages)
	var response []interface{}
	if db.RowsAffected > 0 {
		for index, _ := range packages {
			response = append(response, map[string]interface{}{
				"id":              packages[index].Id,
				"name":            packages[index].Name,
				"id_organization": packages[index].IdOrganization,
				"error":           packages[index].Error,
				"id_author":       packages[index].IdAuthor,
				"id_status":       packages[index].IdStatus,
				"name_status":     packages[index].Status.Name,
				"code_status":     packages[index].Status.Code,
				"created":         packages[index].Created,
				"count_all":       packages[index].CountAll,
				"count_add":       packages[index].CountAdd,
			})
		}
		result.Done = true
		result.Items = response
		return
	} else {
		result.Done = true
		message := `Пакеты не найдены.`
		result.Message = &message
		result.Items = []digest.MarkEgePackages{}
		return
	}

}
func (result *Result) GetMarkEgeElements(idPackage uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var pack digest.MarkEgePackages
	db := conn.Where(`id=?`, idPackage).Find(&pack)
	if pack.Id <= 0 {
		result.Done = false
		m := `Пакет не найден`
		result.Message = &m
		return
	}
	var elements []digest.MarkEgeElement
	if result.Sort.Field == `` {
		result.Sort.Field = `created`
	}
	if result.Sort.Order == `` {
		result.Sort.Order = `asc`
	}
	db = conn.Order(result.Sort.Field + ` ` + result.Sort.Order)
	db = db.Where(`id_package=?`, idPackage)

	dbCount := db.Model(&elements).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&elements)
	var response []interface{}
	if db.RowsAffected > 0 {
		for index, _ := range elements {
			response = append(response, map[string]interface{}{
				"id":              elements[index].Id,
				"id_entrant":      elements[index].Id,
				"surname":         elements[index].Surname,
				"name":            elements[index].Name,
				"patronymic":      elements[index].Patronymic,
				"id_document":     elements[index].IdDocument,
				"doc_series":      elements[index].DocSeries,
				"doc_number":      elements[index].DocNumber,
				"id_subject":      elements[index].IdSubject,
				"name_subject":    elements[index].Subject,
				"mark":            elements[index].Mark,
				"year":            elements[index].Year,
				"id_region":       elements[index].IdRegion,
				"status":          elements[index].Status,
				"app_per":         elements[index].AppPer,
				"cert_number":     elements[index].CertNumber,
				"typ_number":      elements[index].TypNumber,
				"app_status":      elements[index].AppStatus,
				"checked":         elements[index].Checked,
				"error":           elements[index].Error,
				"created":         elements[index].Created,
				"id_document_ege": elements[index].IdDocumentEge,
			})
		}
		result.Done = true
		result.Items = response
		return
	} else {
		result.Done = true
		message := `Элементы не найдены.`
		result.Message = &message
		result.Items = []digest.MarkEgeElement{}
		return
	}

}
func (result *ResultInfo) AddFileMarkEgePackage(packageName string, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var doc digest.MarkEgePackages
	path := getPath(0, doc.TableName(), time.Now())
	ext := filepath.Ext(path + `/` + f.Header.Filename)
	sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
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
	doc.Name = packageName
	doc.PathFile = name
	doc.Created = time.Now()
	doc.IdStatus = 1
	doc.IdAuthor = result.User.Id
	doc.IdOrganization = result.User.CurrentOrganization.Id

	db := conn.Create(&doc)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		return
	}
	result.Items = map[string]interface{}{
		`id_package`: doc.Id,
	}
	result.Done = true
	return
}
