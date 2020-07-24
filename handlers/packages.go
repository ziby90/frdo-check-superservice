package handlers

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
)

var PackageSearchArray = []string{
	`name`,
}

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
	for _, search := range result.Search {
		if service.SearchStringInSliceString(search[0], PackageSearchArray) >= 0 {
			db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
		}
	}
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
	db = db.Limit(result.Paginator.Limit).Preload(`Region`).Offset(result.Paginator.Offset).Find(&elements)
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
				"name_region":     elements[index].Region.Name,
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
	path := getPath(result.User.CurrentOrganization.Id, doc.TableName(), time.Now())
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

func (result *Result) GetRatingApplicationsPackages() {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var packages []digest.RatingApplicationsPackages
	if result.Sort.Field == `` {
		result.Sort.Field = `created`
	}
	if result.Sort.Order == `` {
		result.Sort.Order = `asc`
	}

	db := conn.Order(result.Sort.Field + ` ` + result.Sort.Order)
	db = db.Where(`id_organization=?`, result.User.CurrentOrganization.Id)
	for _, search := range result.Search {
		if service.SearchStringInSliceString(search[0], PackageSearchArray) >= 0 {
			db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
		}
	}
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
		result.Items = []digest.RatingApplicationsPackages{}
		return
	}

}
func (result *Result) GetRatingApplicationsElements(idPackage uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var pack digest.RatingApplicationsPackages
	db := conn.Where(`id=?`, idPackage).Find(&pack)
	if pack.Id <= 0 {
		result.Done = false
		m := `Пакет не найден`
		result.Message = &m
		return
	}
	var elements []digest.RatingApplicationsElement
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
	answer := make(map[string]interface{})
	answer[`package_info`] = nil
	answer[`items`] = nil
	var response []interface{}
	if db.RowsAffected > 0 {
		answer[`package_info`] = map[string]interface{}{
			"id_competitive_group": elements[0].IdCompetitiveGroup,
			"competitive_group":    elements[0].CompetitiveGroup,
			"id_organization":      elements[0].IdOrganization,
			"organization":         elements[0].Organization,
			"admission_volume":     elements[0].AdmissionVolume,
			"count_first_step":     elements[0].CountFirstStep,
			"count_second_step":    elements[0].CountSecondStep,
			"changed":              elements[0].Changed,
		}
		for index, _ := range elements {
			response = append(response, map[string]interface{}{
				"id":                    elements[index].Id,
				"orderid":               elements[index].Orderid,
				"fio":                   elements[index].Fio,
				"rating":                elements[index].Rating,
				"without_tests":         elements[index].WithoutTests,
				"reason_without_tests":  elements[index].ReasonWithoutTests,
				"entrance_test1":        elements[index].EntranceTest1,
				"result1":               elements[index].Result1,
				"entrance_test2":        elements[index].EntranceTest2,
				"result2":               elements[index].Result2,
				"entrance_test3":        elements[index].EntranceTest3,
				"result3":               elements[index].Result3,
				"entrance_test4":        elements[index].EntranceTest4,
				"result4":               elements[index].Result4,
				"entrance_test5":        elements[index].EntranceTest5,
				"result5":               elements[index].Result5,
				"mark":                  elements[index].Mark,
				"benefit":               elements[index].Benefit,
				"reason_benefit":        elements[index].ReasonBenefit,
				"sum_mark":              elements[index].SumMark,
				"agreed":                elements[index].Agreed,
				"original":              elements[index].Original,
				"addition":              elements[index].Addition,
				"enlisted":              elements[index].Enlisted,
				"id_package":            elements[index].IdPackage,
				"checked":               elements[index].Checked,
				"error":                 elements[index].Error,
				"created":               elements[index].Created,
				"id_rating_application": elements[index].IdRatingApplication,
			})
		}
		answer[`items`] = response
		result.Done = true
		result.Items = answer
		return
	} else {
		result.Done = true
		message := `Элементы не найдены.`
		result.Message = &message
		result.Items = answer
		return
	}

}
func (result *ResultInfo) AddFileRatingApplicationsPackage(packageName string, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var doc digest.RatingApplicationsPackages
	path := getPath(result.User.CurrentOrganization.Id, doc.TableName(), time.Now())
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
