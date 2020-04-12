package handlers

import (
	"fmt"
	"persons/config"
	"persons/digest"
	"strings"
	"time"
)

var EntrantsSearchArray = []string{
	`surname`,
	`name`,
	`patronymic`,
	`snils`,
}

type DocsResponseByCategory struct {
	Name  string        `json:"name"`
	Id    uint          `json:"id"`
	Count int           `json:"count"`
	Docs  []interface{} `json:"docs"`
}

type AddEntrantData struct {
	Entrant        digest.Entrants        `json:"entrant"`
	Identification digest.Identifications `json:"identification"`
	Education      digest.Educations      `json:"education"`
}

//
//type Unmarshaler interface {
//	UnmarshalJSON([]byte) error
//}
//
//func (l *Entrants) UnmarshalJSON(j []byte) error {
//	fmt.Println(`****************`)
//	var rawStrings map[string]string
//
//	err := json.Unmarshal(j, &rawStrings)
//	if err != nil {
//		return err
//	}
//
//	for k, v := range rawStrings {
//		if strings.ToLower(k) == "id" {
//			l.Id, err = strconv.Atoi(v)
//			if err != nil {
//				return err
//			}
//		}
//		if strings.ToLower(k) == "birthday" {
//			t, err := time.Parse(time.RFC3339, v)
//			fmt.Println(t, err)
//			if err != nil {
//				return err
//			}
//			l.Birthday = t
//		}
//	}
//
//	return nil
//}

func (result *ResultInfo) GetInfoEntrant(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var entrant digest.Entrants
	db := conn.Find(&entrant, ID)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Абитуриент не найден.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		db = conn.Model(&entrant).Related(&entrant.Gender, `IdGender`)
		db = conn.Model(&entrant).Related(&entrant.Okcm, `IdOkcm`)
		result.Done = true
		result.Items = entrant
		return
	} else {
		result.Done = true
		message := `Абитуриент не найден.`
		result.Message = &message
		return
	}
}

func (result *ResultInfo) GetDocsEntrant(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var entrant digest.Entrants
	db := conn.Find(&entrant, ID)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Абитуриент не найден.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		res := make(map[string]DocsResponseByCategory)
		var documents []digest.VDocuments
		db = conn.Where(`name_table!='identification'`).Model(&entrant).Related(&documents)

		for index, doc := range documents {
			var docCategory DocsResponseByCategory
			db = conn.Model(&documents[index]).Related(&documents[index].DocumentType, `IdDocumentType`)
			val, ok := res[doc.NameTable]
			if ok {
				docCategory = val
			} else {
				docCategory.Name = doc.NameSysCategories
				docCategory.Id = doc.IdSysCategories
			}
			docCategory.Docs = append(docCategory.Docs, map[string]interface{}{
				"doc_name": documents[index].DocumentName,
				"id":       doc.IdDocument,
				"created":  doc.Created,
			})
			docCategory.Count = len(docCategory.Docs)
			res[doc.NameTable] = docCategory
		}
		result.Done = true
		result.Items = res
		return
	} else {
		result.Done = true
		message := `Абитуриент не найден.`
		result.Message = &message
		return
	}

}
func (result *ResultInfo) GetDocsIdentsEntrant(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var entrant digest.Entrants

	db := conn.Find(&entrant, ID)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Абитуриент не найден.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}

	if db.RowsAffected > 0 {
		res := make(map[string]DocsResponseByCategory)
		var identifications []digest.Identifications
		db = conn.Model(&entrant).Related(&identifications)
		db = conn.Model(&entrant).Related(&identifications)
		var docCategory DocsResponseByCategory
		for index := range identifications {
			db = conn.Model(&identifications[index]).Related(&identifications[index].DocumentType, `IdDocumentType`)
			var documentSysCategory digest.DocumentSysCategories
			db = conn.Where(`name_table=?`, `identification`).Find(&documentSysCategory)
			if documentSysCategory.Id <= 0 {
				message := `Не найдена катекгория удостоверяющих документы`
				result.Message = &message
				return
			}
			docCategory.Name = documentSysCategory.Name
			docCategory.Id = documentSysCategory.Id
			db = conn.Model(&identifications[index]).Related(&identifications[index].Okcm, `IdOkcm`)
			db = conn.Model(&identifications[index]).Related(&identifications[index].DocumentType, `IdDocumentType`)
			issueDate := identifications[index].IssueDate.Format(`2006-01-02`)
			docCategory.Docs = append(docCategory.Docs, map[string]interface{}{
				"doc_name": identifications[index].DocumentType.Name,
				"id":       identifications[index].Id,
				"created":  identifications[index].Created,
				"data": map[string]interface{}{
					"id_entrant":         entrant.Id,
					"id_document_type":   identifications[index].DocumentType.Id,
					"name_document_type": identifications[index].DocumentType.Name,
					"surname":            identifications[index].Surname,
					"name":               identifications[index].Name,
					"patronymic":         identifications[index].Patronymic,
					"doc_series":         identifications[index].DocSeries,
					"doc_number":         identifications[index].DocNumber,
					"doc_organization":   identifications[index].DocOrganization,
					"issue_date":         issueDate,
					"subdivision_code":   identifications[index].SubdivisionCode,
					"id_okcm":            identifications[index].IdOkcm,
					"name_okcm":          identifications[index].Okcm.ShortName,
					"checked":            identifications[index].Checked,
					"created":            identifications[index].Created,
					"name_sys_category":  documentSysCategory.Name,
				},
			})
			docCategory.Count = len(docCategory.Docs)
			res[`identification`] = docCategory
		}

		result.Done = true
		result.Items = res
		return
	} else {
		result.Done = true
		message := `Абитуриент не найден.`
		result.Message = &message
		return
	}

}

func (result *Result) GetListEntrants() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var entrants []digest.Entrants
	db := conn.Order(result.Sort.Field + ` ` + result.Sort.Order)
	for _, search := range result.Search {
		db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)

	}
	dbCount := db.Model(&entrants).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&entrants)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Абитуриенты не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.` + db.Error.Error()
		result.Message = &message
		return
	}
	var entrantsList []interface{}
	if db.RowsAffected > 0 {
		for _, entrant := range entrants {
			db = conn.Model(&entrant).Related(&entrant.Gender, `IdGender`)
			db = conn.Model(&entrant).Related(&entrant.Okcm, `IdOkcm`)
			e := map[string]interface{}{
				`id`:          entrant.Id,
				`surname`:     entrant.Surname,
				`name`:        entrant.Name,
				`patronymic`:  entrant.Patronymic,
				`snils`:       entrant.Snils,
				`birthday`:    entrant.Birthday,
				`id_gender`:   entrant.Gender.Id,
				`gender_name`: entrant.Gender.Name,
				`id_okcm`:     entrant.Okcm.Id,
				`okcm_name`:   entrant.Okcm.FullName,
				`created`:     entrant.Created,
			}
			entrantsList = append(entrantsList, e)
		}
		result.Done = true
		result.Items = entrantsList
		return
	} else {
		result.Done = true
		message := `Абитуриенты не найдены.`
		result.Message = &message
		result.Items = []digest.Entrants{}
		return
	}
}

func (result *ResultInfo) AddEntrant(entrantData AddEntrantData) {
	result.Items = entrantData.Entrant.Name
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)

	tx := conn.Begin()
	if entrantData.Entrant.Snils == `` {
		result.SetErrorResult(`Снилс обязательное поле`)
		return
	}

	var exist digest.Entrants
	db := tx.Where(`snils=?`, entrantData.Entrant.Snils).Find(&exist)
	if exist.Id > 0 {
		result.SetErrorResult(`Абитуриент с данным снилс уже существует`)
		return
	}

	var entrant digest.Entrants
	entrant = entrantData.Entrant
	entrant.Created = time.Now()
	entrant.Surname = strings.TrimSpace(entrant.Surname)
	entrant.Name = strings.TrimSpace(entrant.Name)
	entrant.Patronymic = strings.TrimSpace(entrant.Patronymic)
	fmt.Println(entrantData.Entrant.Birthplace)
	fmt.Println(entrant.Birthplace)
	db = tx.Find(&entrant.Gender, entrant.IdGender)
	if db.Error != nil || !entrant.Gender.Actual {
		result.SetErrorResult(`Не найден пол`)
		return
	}

	db = tx.Find(&entrant.Okcm, entrant.IdOkcm)
	if db.Error != nil || !entrant.Okcm.Actual {
		result.SetErrorResult(`Не найден оксм`)
		return
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&entrant)
	if db.Error != nil {
		tx.Rollback()
		m := `Ошибка при добавлении абитуриента: ` + db.Error.Error()
		result.Message = &m
		tx.Rollback()
		return
	}
	var identification digest.Identifications
	identification = entrantData.Identification
	identification.EntrantsId = entrant.Id
	identification.Created = time.Now()
	identification.Name = strings.TrimSpace(identification.Name)
	identification.Surname = strings.TrimSpace(identification.Surname)
	identification.Patronymic = strings.TrimSpace(identification.Patronymic)

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&identification)
	if db.Error != nil {
		tx.Rollback()
		m := `Ошибка при добавлении доумента, удостоверяющего личность: ` + db.Error.Error()
		result.Message = &m
		tx.Rollback()
		return
	}
	var education digest.Educations
	education = entrantData.Education
	education.IdEntrant = entrant.Id
	education.IdIdentDocument = identification.Id
	education.Created = time.Now()

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&education)
	if db.Error != nil {
		tx.Rollback()
		m := `Ошибка при добавлении доумента об образовании: ` + db.Error.Error()
		result.Message = &m
		tx.Rollback()
		return
	}
	result.Items = map[string]interface{}{
		`id_entrant`:        entrant.Id,
		`id_identification`: identification.Id,
		`id_education`:      education.Id,
	}
	result.Done = true
	tx.Commit()
}