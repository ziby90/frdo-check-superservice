package handlers

import (
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

type CategoryDocs struct {
	Name string        `json:"name_category"`
	Code string        `json:"code_category"`
	Docs []interface{} `json:"docs"`
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
		db = conn.Model(&entrant).Related(&entrant.FactAddr, `IdFactAddr`)
		db = conn.Model(&entrant).Related(&entrant.RegistrationAddr, `IdRegistrationAddr`)
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
func (result *ResultInfo) GetInfoEntrantApp(ID uint) {
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
		var countApp int64
		db = conn.Select("count(distinct(id_organization))").Table(`app.applications`).Where(`id_entrant=?`, ID).Count(&countApp)
		var applications []digest.Application
		db = conn.Table(`app.applications`).Where(`id_entrant=? AND id_organization=?`, ID, result.User.CurrentOrganization.Id).Find(&applications)
		birthday := entrant.Birthday.Format(`2006-01-02`)
		apps := []interface{}{} // ну вот надо так чикиной, че иде мне подчеркивает(
		for index, _ := range applications {
			apps = append(apps, map[string]interface{}{
				"id":         applications[index].Id,
				"app_number": applications[index].AppNumber,
			})
		}
		response := map[string]interface{}{
			`surname`:       entrant.Surname,
			`name`:          entrant.Name,
			`patronymic`:    entrant.Patronymic,
			`birthday`:      birthday,
			`snils`:         entrant.Snils,
			`count_org_app`: countApp,
			`app_org`:       apps,
		}
		result.Done = true
		result.Items = response
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

func (result *ResultInfo) GetListDocsIdentsEntrant(ID uint) {
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
	var items []interface{}
	if db.RowsAffected > 0 {
		var identifications []digest.Identifications
		db = conn.Preload(`DocumentType`).Model(&entrant).Related(&identifications)
		for index := range identifications {
			db = conn.Model(&identifications[index]).Related(&identifications[index].DocumentType, `IdDocumentType`)
			issueDate := identifications[index].IssueDate.Format(`02-01-2006`)
			name := identifications[index].DocumentType.Name + ` ` + identifications[index].DocSeries + ` ` + identifications[index].DocNumber + ` от ` + issueDate
			items = append(items, map[string]interface{}{
				"id":   identifications[index].Id,
				"name": name,
			})
		}

		result.Done = true
		result.Items = items
		return
	} else {
		result.Done = true
		message := `Абитуриент не найден.`
		result.Message = &message
		return
	}

}

func (result *ResultInfo) GetShortListDocsEntrant(idEntrant uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var entrant digest.Entrants

	db := conn.Find(&entrant, idEntrant)
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
	var items []interface{}
	if db.RowsAffected > 0 {
		var identifications []digest.Identifications
		db = conn.Preload(`DocumentType`).Where(`id_entrant=?`, idEntrant).Find(&identifications)
		var categoryIdents CategoryDocs
		var sysCategoryCls digest.DocumentSysCategories
		_ = conn.Where(`name_table=?`, `identification`).Find(&sysCategoryCls)
		categoryIdents.Code = sysCategoryCls.NameTable
		categoryIdents.Name = sysCategoryCls.Name
		for index := range identifications {
			issueDate := identifications[index].IssueDate.Format(`02-01-2006`)
			categoryIdents.Docs = append(categoryIdents.Docs, map[string]interface{}{
				"id":         identifications[index].Id,
				"name_type":  identifications[index].DocumentType.Name,
				"doc_series": identifications[index].DocSeries,
				"doc_number": identifications[index].DocNumber,
				"issue_date": issueDate,
			})
		}
		if len(categoryIdents.Docs) > 0 {
			items = append(items, categoryIdents)
		}

		var docs []digest.VDocuments
		db = conn.Where(`id_entrant=?`, idEntrant).Find(&docs)
		sysCategory := make(map[string]CategoryDocs)
		for index := range docs {
			var category CategoryDocs
			if val, ok := sysCategory[docs[index].NameTable]; ok {
				category = val
			} else {
				category.Name = docs[index].NameSysCategories
				category.Code = docs[index].NameTable
			}
			category.Docs = append(category.Docs, map[string]interface{}{
				"id":        docs[index].IdDocument,
				"name_type": docs[index].DocumentName,
				"doc_name":  docs[index].DocName,
			})
			sysCategory[docs[index].NameTable] = category
		}
		for index, _ := range sysCategory {
			items = append(items, sysCategory[index])
		}
		if len(items) < 1 {
			result.Items = []digest.VDocuments{}
		} else {
			result.Items = items
		}
		result.Done = true
		return
	} else {
		result.Done = false
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
	defer func() {
		tx.Rollback()
	}()
	if entrantData.Entrant.Snils == `` {
		result.SetErrorResult(`Снилс обязательное поле`)
		tx.Rollback()
		return
	}

	var exist digest.Entrants
	db := tx.Where(`snils=?`, entrantData.Entrant.Snils).Find(&exist)
	if exist.Id > 0 {
		result.SetErrorResult(`Абитуриент с данным снилс уже существует`)
		tx.Commit()
		return
	}
	//check := service.CheckSnils(entrantData.Entrant.Snils)
	//if check != nil {
	//	result.SetErrorResult(check.Error())
	//	return
	//}

	var entrant digest.Entrants
	entrant = entrantData.Entrant
	entrant.Created = time.Now()
	entrant.Surname = strings.TrimSpace(entrant.Surname)
	entrant.Name = strings.TrimSpace(entrant.Name)
	entrant.Patronymic = strings.TrimSpace(entrant.Patronymic)
	db = tx.Find(&entrant.Gender, entrant.IdGender)
	if db.Error != nil || !entrant.Gender.Actual {
		tx.Rollback()
		result.SetErrorResult(`Не найден пол`)
		return
	}

	db = tx.Find(&entrant.Okcm, entrant.IdOkcm)
	if db.Error != nil || !entrant.Okcm.Actual {
		tx.Rollback()
		result.SetErrorResult(`Не найден оксм`)
		return
	}

	if entrantData.Entrant.RegistrationAddr.IdRegion > 0 {
		var registrAddr digest.Address
		registrAddr = entrant.RegistrationAddr
		registrAddr.IdAuthor = result.User.Id

		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&registrAddr)
		if db.Error != nil {
			m := `Ошибка при добавлении регистрационного адреса: ` + db.Error.Error()
			result.Message = &m
			tx.Rollback()
			return
		}
		entrant.IdRegistrationAddr = registrAddr.Id
		entrant.RegistrationAddr.Id = registrAddr.Id
	}

	if entrantData.Entrant.FactAddr.IdRegion > 0 {
		var factAddr digest.Address
		factAddr = entrant.FactAddr
		factAddr.IdAuthor = result.User.Id

		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&factAddr)
		if db.Error != nil {
			m := `Ошибка при добавлении фактического адреса: ` + db.Error.Error()
			result.Message = &m
			tx.Rollback()
			return
		}

		entrant.IdFactAddr = factAddr.Id
		entrant.FactAddr.Id = factAddr.Id
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&entrant)
	if db.Error != nil {
		m := `Ошибка при добавлении абитуриента: ` + db.Error.Error()
		result.Message = &m
		tx.Rollback()
		return
	}

	var identification digest.Identifications
	identification = entrantData.Identification
	identification.IdOrganization = result.User.CurrentOrganization.Id
	identification.EntrantsId = entrant.Id
	identification.Created = time.Now()
	identification.Name = strings.TrimSpace(identification.Name)
	identification.Surname = strings.TrimSpace(identification.Surname)
	identification.Patronymic = strings.TrimSpace(identification.Patronymic)

	var existIdent digest.Identifications
	db = tx.Where(`UPPER(doc_series)=? AND UPPER(doc_number)=? AND issue_date::date=?::date`, strings.ToUpper(identification.DocSeries), strings.ToUpper(identification.DocNumber), identification.IssueDate).Find(&existIdent)
	if existIdent.Id > 0 {
		result.SetErrorResult(`Удостоверяющий документ с указанными серией, номером и датой выдачи уже существует`)
		return
	}
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&identification)
	if db.Error != nil {
		m := `Ошибка при добавлении документа, удостоверяющего личность: ` + db.Error.Error()
		result.Message = &m
		tx.Rollback()
		return
	}
	var education digest.Educations
	education = entrantData.Education
	education.IdEntrant = entrant.Id
	education.IdIdentDocument = identification.Id
	education.Created = time.Now()

	var existEduc digest.Educations
	db = tx.Where(`UPPER(doc_series)=? AND UPPER(doc_number)=? AND issue_date::date=?::date AND id_document_type=?`, strings.ToUpper(education.DocSeries), strings.ToUpper(education.DocNumber), education.IssueDate, education.IdDocumentType).Find(&existEduc)
	if existEduc.Id > 0 {
		result.SetErrorResult(`Документ об образовании с данными серией, номером, датой выдачи и типом документа уже существует`)
		return
	}
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&education)
	if db.Error != nil {
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
