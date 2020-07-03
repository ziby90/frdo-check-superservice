package handlers

import (
	"fmt"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
)

var EntrantsSearchArray = []string{
	`surname`,
	`name`,
	`patronymic`,
	`snils`,
}
var PriorityTable = []string{
	`identification`,
	`educations`,
	`ege`,
	`orphans`,
	`veteran`,
	`olympics`,
	`militaries`,
	`other`,
	`disability`,
	`compatriot`,
	`parents_lost`,
	`radiation_work`,
	`composition`,
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
		result.Done = true
		birthday := entrant.Birthday.Format("2006-01-02")
		response := map[string]interface{}{
			`id`:                     entrant.Id,
			`created`:                entrant.Created,
			"snils":                  entrant.Snils,
			"surname":                entrant.Surname,
			"name":                   entrant.Name,
			"patronymic":             entrant.Patronymic,
			"id_gender":              entrant.Gender.Id,
			"name_gender":            entrant.Gender.Name,
			"birthday":               birthday,
			"birthplace":             entrant.Birthplace,
			"phone":                  entrant.Phone,
			"email":                  entrant.Email,
			"id_okcm":                entrant.Okcm.Id,
			`name_okcm`:              entrant.Okcm.ShortName,
			`registration_addr_full`: entrant.RegistrationAddrFull,
			`fact_addr_full`:         entrant.FactAddrFull,
		}
		fmt.Println(entrant.IdFactAddr)
		var factAdrr digest.Address
		if entrant.IdFactAddr != nil {
			db = conn.Preload(`Region`).Where(`id=?`, entrant.IdFactAddr).Find(&factAdrr)
			if factAdrr.Id > 0 {
				response[`fact_address`] = map[string]interface{}{
					"id":                factAdrr.Id,
					"full_addr":         factAdrr.FullAddr,
					"index_addr":        factAdrr.IndexAddr,
					"id_region":         factAdrr.Region.Id,
					"name_region":       factAdrr.Region.Name,
					"area":              factAdrr.Area,
					"city":              factAdrr.City,
					"city_area":         factAdrr.CityArea,
					"place":             factAdrr.Place,
					"street":            factAdrr.Street,
					"additional_area":   factAdrr.AdditionalArea,
					"additional_street": factAdrr.AdditionalStreet,
					"house":             factAdrr.House,
					"building1":         factAdrr.Building1,
					"building2":         factAdrr.Building2,
					"apartment":         factAdrr.Apartment,
					"id_author":         factAdrr.IdAuthor,
				}
			}
		} else {
			response[`fact_address`] = nil
		}

		var regAdrr digest.Address
		if entrant.IdFactAddr != nil {
			db = conn.Preload(`Region`).Where(`id=?`, entrant.IdRegistrationAddr).Find(&regAdrr)
			if regAdrr.Id > 0 {
				response[`registration_address`] = map[string]interface{}{
					"id":                regAdrr.Id,
					"full_addr":         regAdrr.FullAddr,
					"index_addr":        regAdrr.IndexAddr,
					"id_region":         regAdrr.Region.Id,
					"name_region":       regAdrr.Region.Name,
					"area":              regAdrr.Area,
					"city":              regAdrr.City,
					"city_area":         regAdrr.CityArea,
					"place":             regAdrr.Place,
					"street":            regAdrr.Street,
					"additional_area":   regAdrr.AdditionalArea,
					"additional_street": regAdrr.AdditionalStreet,
					"house":             regAdrr.House,
					"building1":         regAdrr.Building1,
					"building2":         regAdrr.Building2,
					"apartment":         regAdrr.Apartment,
					"id_author":         regAdrr.IdAuthor,
				}
			}
		} else {
			response[`registration_address`] = nil
		}
		result.Items = response
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
		result.Done = false
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
			issueDate := identifications[index].IssueDate.Format(`2006-01-02`)
			series := ``
			if identifications[index].DocSeries != nil {
				series = *identifications[index].DocSeries
			}
			name := identifications[index].DocumentType.Name + ` ` + series + ` ` + identifications[index].DocNumber + ` от ` + issueDate
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

func (result *ResultInfo) GetShortListDocsEntrant(idEntrant uint, keys map[string][]string) {
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
		var allDocuments []digest.AllDocuments
		cmd := `
					with a as (select id as id_entrant from persons.entrants where id=?),
					b as (SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, doc_number, id_document_type, doc_series, NULL::integer as mark, NULL::character varying as name_subject, issue_date, 'educations' as name_table  FROM documents.educations educ WHERE EXISTS(SELECT 1 FROM a WHERE educ.id_entrant=a.id_entrant)
					UNION
					SELECT ege.id,path_file,id_entrant,uid_epgu,id_organization,checked, doc_number, id_document_type, NULL::character varying as doc_series, mark, sbj.name as name_subject,  issue_date, 'ege' as name_table
						FROM documents.ege ege
						join cls.subjects sbj ON sbj.id = ege.id_subject WHERE EXISTS(SELECT 1 FROM a WHERE ege.id_entrant=a.id_entrant)
					UNION
					SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'orphans' as name_table FROM documents.orphans orph WHERE EXISTS(SELECT 1 FROM a WHERE orph.id_entrant=a.id_entrant)
					UNION
					SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'veteran' as name_table FROM documents.veteran vet WHERE EXISTS(SELECT 1 FROM a WHERE vet.id_entrant=a.id_entrant)
					UNION
					SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'olympics' as name_table FROM documents.olympics olymp WHERE EXISTS(SELECT 1 FROM a WHERE olymp.id_entrant=a.id_entrant)
					UNION
					SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'militaries' as name_table FROM documents.militaries mil WHERE EXISTS(SELECT 1 FROM a WHERE mil.id_entrant=a.id_entrant)
					UNION
					SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'other' as name_table FROM documents.other oth WHERE EXISTS(SELECT 1 FROM a WHERE oth.id_entrant=a.id_entrant)
					UNION
					SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, doc_number, id_document_type, NULL::character varying as doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'disability' as name_table FROM documents.disability dis WHERE EXISTS(SELECT 1 FROM a WHERE dis.id_entrant=a.id_entrant)
					UNION
					SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, NULL::character varying as doc_number, id_document_type,  NULL::character varying as doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, NULL::timestamp with time zone as issue_date, 'compatriot' as name_table
					FROM documents.compatriot compar WHERE EXISTS(SELECT 1 FROM a WHERE compar.id_entrant=a.id_entrant)
					UNION
					SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, doc_number, id_document_type, doc_series, NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'parents_lost' as name_table FROM documents.parents_lost par WHERE EXISTS(SELECT 1 FROM a WHERE par.id_entrant=a.id_entrant)
					UNION
					SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, doc_number, id_document_type, doc_series, NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'radiation_work' as name_table FROM documents.radiation_work rad WHERE EXISTS(SELECT 1 FROM a WHERE rad.id_entrant=a.id_entrant)
					UNION
					SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, NULL::character varying as doc_number, id_document_type, NULL::character varying as doc_series, NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'composition' as name_table
					FROM documents.composition compos WHERE EXISTS(SELECT 1 FROM a WHERE compos.id_entrant=a.id_entrant)
					UNION
					SELECT id,path_file,id_entrant,uid_epgu,id_organization,checked, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'identification' as name_table
					FROM documents.identification ident WHERE EXISTS(SELECT 1 FROM a WHERE ident.id_entrant=a.id_entrant))
					SELECT b.*, sys.id as id_sys_categories, sys."name" as name_sys_categories, dt."name" as name_document_type
					from b  
					join cls.document_sys_categories sys on b.name_table = sys.name_table
					join cls.document_types dt on dt.id = b.id_document_type			
					Where b.id IS NOT NULL
`
		if len(keys[`categories`]) > 0 {
			var categories []string
			fmt.Println()
			for _, v := range keys[`categories`] {
				if service.SearchStringInSliceString(v, PriorityTable) >= 0 {
					categories = append(categories, `'`+v+`'`)
				}
			}
			if len(categories) > 0 {
				cmd += ` AND b.name_table IN (` + strings.Join(categories, `,`) + `) `
			}
		}
		if len(keys[`no_categories`]) > 0 {
			var categories []string
			for _, v := range keys[`no_categories`] {
				if service.SearchStringInSliceString(v, PriorityTable) >= 0 {
					categories = append(categories, `'`+v+`'`)
				}
			}
			if len(categories) > 0 {
				cmd += ` AND b.name_table NOT IN (` + strings.Join(categories, `,`) + `) `
			}
		}
		db = conn.Raw(cmd, idEntrant).Scan(&allDocuments)
		if db.Error != nil {
			result.Done = false
			message := db.Error.Error()
			result.Message = &message
			return
		}
		sysCategory := make(map[string]CategoryDocs)
		for index := range allDocuments {
			var category CategoryDocs
			if val, ok := sysCategory[allDocuments[index].NameTable]; ok {
				category = val
			} else {
				category.Name = allDocuments[index].NameSysCategories
				category.Code = allDocuments[index].NameTable
			}
			var issueDate *string
			if allDocuments[index].IssueDate != nil {
				date := *allDocuments[index].IssueDate
				dateF := date.Format(`2006-01-02`)
				issueDate = &dateF
			} else {
				issueDate = nil
			}
			canEdit := false
			if allDocuments[index].IdOrganization != nil && *allDocuments[index].IdOrganization == result.User.CurrentOrganization.Id && allDocuments[index].UidEpgu == nil {
				canEdit = true
			}
			file := false
			if allDocuments[index].PathFile != nil {
				file = true
			}
			category.Docs = append(category.Docs, map[string]interface{}{
				"id":               allDocuments[index].Id,
				"doc_number":       allDocuments[index].DocNumber,
				"doc_series":       allDocuments[index].DocSeries,
				"uid_epgu":         allDocuments[index].UidEpgu,
				"id_document_type": allDocuments[index].IdDocumentType,
				"checked":          allDocuments[index].Checked,
				"id_organization":  allDocuments[index].IdOrganization,
				"id_entrant":       allDocuments[index].IdEntrant,
				//"id_sys_categories": 		allDocuments[index].IdSysCategories,
				"issue_date":         issueDate,
				"can_edit":           canEdit,
				"mark":               allDocuments[index].Mark,
				"name_document_type": allDocuments[index].NameDocumentType,
				"name_subject":       allDocuments[index].NameSubject,
				"has_file":           file,
				//"name_sys_categories": 		allDocuments[index].NameSysCategories,
				//"name_table": 				allDocuments[index].NameTable,
			})
			sysCategory[allDocuments[index].NameTable] = category
		}
		for index, _ := range PriorityTable {
			if val, ok := sysCategory[PriorityTable[index]]; ok {
				items = append(items, val)
			}
		}
		if len(items) < 1 {
			result.Items = []digest.AllDocuments{}
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
	if entrant.Patronymic != nil {
		s := strings.TrimSpace(*entrant.Patronymic)
		entrant.Patronymic = &s
	}

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

	if entrantData.Entrant.RegistrationAddr != nil {
		var registrAddr digest.Address
		registrAddr = *entrant.RegistrationAddr
		registrAddr.IdAuthor = &result.User.Id
		entrant.RegistrationAddr.IdAuthor = &result.User.Id

		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&registrAddr)
		if db.Error != nil {
			m := `Ошибка при добавлении регистрационного адреса: ` + db.Error.Error()
			result.Message = &m
			tx.Rollback()
			return
		}
		entrant.IdRegistrationAddr = &registrAddr.Id
		entrant.RegistrationAddr.Id = registrAddr.Id
	}

	if entrantData.Entrant.FactAddr != nil {
		var factAddr digest.Address
		factAddr = *entrant.FactAddr
		factAddr.IdAuthor = &result.User.Id

		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&factAddr)
		if db.Error != nil {
			m := `Ошибка при добавлении фактического адреса: ` + db.Error.Error()
			result.Message = &m
			tx.Rollback()
			return
		}

		entrant.IdFactAddr = &factAddr.Id
		entrant.FactAddr.Id = factAddr.Id
		entrant.FactAddr.IdAuthor = &result.User.Id
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
	identification.IdOrganization = &result.User.CurrentOrganization.Id
	identification.EntrantsId = entrant.Id
	identification.Created = time.Now()
	identification.Name = strings.TrimSpace(identification.Name)
	identification.Surname = strings.TrimSpace(identification.Surname)
	if identification.Patronymic != nil {
		s := strings.TrimSpace(*identification.Patronymic)
		identification.Patronymic = &s
	}

	var existIdent digest.Identifications
	if identification.DocSeries != nil {
		db = tx.Where(`doc_series is null AND UPPER(doc_number)=? AND issue_date::date=?::date AND id_document_type=?`, strings.ToUpper(identification.DocNumber), identification.IssueDate, identification.IdDocumentType).Find(&existIdent)
	} else {
		db = tx.Where(`UPPER(doc_series)=? AND UPPER(doc_number)=? AND issue_date::date=?::date AND id_document_type=?`, strings.ToUpper(*identification.DocSeries), strings.ToUpper(identification.DocNumber), identification.IssueDate, identification.IdDocumentType).Find(&existIdent)
	}

	if existIdent.Id > 0 {
		result.SetErrorResult(`Удостоверяющий документ с указанными серией, номером, типом и датой выдачи уже существует`)
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
