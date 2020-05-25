package handlers

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
)

var ApplicationSearchArray = []string{
	`hz`,
}

type ChangeStatusApplication struct {
	Application         digest.Application         `gorm:"foreignkey:IdApplication"`
	IdApplication       uint                       `json:"id_application"`
	ApplicationStatus   digest.ApplicationStatuses `gorm:"foreignkey:IdApplicationStatus"`
	IdApplicationStatus *uint                      `json:"id_application_status"`
	CodeStatus          string                     `json:"code"`
}
type AddApplication struct {
	IdCompetitiveGroup uint      `json:"id_competitive_group"`
	IdEntrant          uint      `json:"id_entrant" schema:"id_entrant"`
	AppNumber          string    `json:"app_number"`
	RegistrationDate   time.Time `json:"registration_date" schema:"registration_date"`
	//Rating                   	float32             `json:"rating" schema:"rating"`
	Priority             int64 `json:"priority" schema:"priority"`
	FirstHigherEducation bool  `json:"first_higher_education" schema:"first_higher_education"`
	NeedHostel           bool  `json:"need_hostel" schema:"need_hostel"`
	//IdDisabledDocument       	uint                `json:"id_disabled_document" schema:"id_disabled_document"`
	SpecialConditions bool    `json:"special_conditions" schema:"special_conditions"`
	DistanceTest      bool    `json:"distance_test" schema:"distance_test"`
	DistancePlace     *string `json:"distance_place" schema:"distance_place"`
	//IdViolation              	uint                `json:"id_violation" schema:"id_violation"`
	//EgeCheck                 	*time.Time          `json:"ege_check" schema:"ege_check"`
	Agreed        *bool      `json:"agreed" schema:"agreed"`
	Disagreed     *bool      `json:"disagreed" schema:"disagreed"`
	AgreedDate    *time.Time `json:"agreed_date" schema:"agreed_date"`
	DisagreedDate *time.Time `json:"disagreed_date" schema:"disagreed_date"`
	//IdOrderAdmission         	*uint               `json:"id_order_admission" schema:"id_order_admission"`
	//OrderAdmissionDate       	*time.Time          `json:"order_admission_date" schema:"order_admission_date"`
	//IdReturnType             	*uint               `json:"id_return_type" schema:"id_return_type"`
	//ReturnDate               	*time.Time          `json:"return_date" schema:"return_date"`
	Original bool `json:"original" schema:"original"`
	//IdBenefit                	uint                `json:"id_benefit" schema:"id_benefit"`
	Uid           *string           `json:"uid" schema:"uid"`
	StatusComment *string           `json:"status_comment" schema:"status_comment"`
	Docs          []DocsApplication `json:"docs" schema:"docs"`
}

type DocsApplication struct {
	Id   uint   `json:"id"`
	Type string `json:"type"`
}

//func CheckHandlerPost(w http.ResponseWriter, r *http.Request) {
//	data := AddNewData{
//		Params: map[string]interface{}{},
//		Files:  []*multipart.FileHeader{},
//	}
//	err := r.ParseMultipartForm(200000)
//	if err != nil {
//		//fmt.Println(err)
//	} else {
//		for _, file := range r.MultipartForm.File["files"] {
//			data.Files = append(data.Files, file)
//		}
//	}
//	for key, value := range r.Form {
//		data.Params[key] = value[0]
//	}
//	json.NewDecoder(r.Body).Decode(&data)
//
//	var check = CheckFiles(data)
//	service.ReturnJSON(w, check)
//
//}

func (result *Result) GetApplications() {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var applications []digest.Application
	sortField := `created`
	sortOrder := `desc`
	if result.Sort.Field != `` {
		sortField = result.Sort.Field
	} else {
		result.Sort.Field = sortField
	}

	fmt.Print(result.Sort.Field, sortField)
	db := conn.Where(`id_organization=?`, result.User.CurrentOrganization.Id).Preload(`Status`).Preload(`Entrants`).Preload(`CompetitiveGroup`)
	for _, search := range result.Search {
		if service.SearchStringInSliceString(search[0], ApplicationSearchArray) >= 0 {
			db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
		}
	}

	dbCount := db.Model(&applications).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Order(sortField + ` ` + sortOrder).Find(&applications)
	var response []interface{}
	if db.RowsAffected > 0 {
		for index, _ := range applications {
			response = append(response, map[string]interface{}{
				"id":                     applications[index].Id,
				"app_number":             applications[index].AppNumber,
				"name_competitive_group": applications[index].CompetitiveGroup.Name,
				"entrant_fullname":       applications[index].Entrants.Surname + ` ` + applications[index].Entrants.Name + ` ` + applications[index].Entrants.Patronymic,
				"entrant_snils":          applications[index].Entrants.Snils,
				"id_status":              applications[index].Status.Id,
				"name_status":            applications[index].Status.Name,
				"registration_date":      applications[index].RegistrationDate,
				"agreed":                 applications[index].Agreed,
				"original":               applications[index].Original,
				"rating":                 applications[index].Rating,
			})
		}
		result.Done = true
		result.Items = response
		return
	} else {
		result.Done = true
		message := `Заявления не найдены.`
		result.Message = &message
		result.Items = []digest.Application{}
		return
	}

}
func (result *ResultInfo) GetApplicationsByEntrant(idEntrant uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var applications []digest.Application

	db := conn.Where(`id_organization=? AND id_entrant=?`, result.User.CurrentOrganization.Id, idEntrant).Preload(`Status`).Preload(`Entrants`).Preload(`CompetitiveGroup`).Find(&applications)

	var response []interface{}
	if db.RowsAffected > 0 {
		for index, _ := range applications {
			response = append(response, map[string]interface{}{
				"id":                     applications[index].Id,
				"app_number":             applications[index].AppNumber,
				"name_competitive_group": applications[index].CompetitiveGroup.Name,
				"entrant_fullname":       applications[index].Entrants.Surname + ` ` + applications[index].Entrants.Name + ` ` + applications[index].Entrants.Patronymic,
				"entrant_snils":          applications[index].Entrants.Snils,
				"id_status":              applications[index].Status.Id,
				"name_status":            applications[index].Status.Name,
				"registration_date":      applications[index].RegistrationDate,
				"agreed":                 applications[index].Agreed,
				"original":               applications[index].Original,
				"rating":                 applications[index].Rating,
			})
		}
		result.Done = true
		result.Items = response
		return
	} else {
		result.Done = true
		message := `Заявления не найдены.`
		result.Message = &message
		result.Items = []digest.Application{}
		return
	}

}
func (result *ResultInfo) GetApplicationById(idApplication uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var application digest.Application

	db := conn.Where(`id_organization=? AND id=?`, result.User.CurrentOrganization.Id, idApplication).Preload(`Status`).Preload(`Entrants`).Preload(`CompetitiveGroup`).Find(&application)

	if db.RowsAffected > 0 {
		var response interface{}
		var info interface{}
		birthday := application.Entrants.Birthday.Format(`2006-01-02`)
		entrant := map[string]interface{}{
			"id":         application.Entrants.Id,
			"surname":    application.Entrants.Surname,
			"name":       application.Entrants.Name,
			"patronymic": application.Entrants.Patronymic,
			"snils":      application.Entrants.Snils,
			"birthday":   birthday,
		}
		competitive := map[string]interface{}{
			"name": application.CompetitiveGroup.Name,
			"id":   application.CompetitiveGroup.Id,
		}
		info = map[string]interface{}{
			"id":                application.Id,
			"app_number":        application.AppNumber,
			"id_status":         application.Status.Id,
			"name_status":       application.Status.Name,
			"code_status":       application.Status.Code,
			"registration_date": application.RegistrationDate,
		}
		response = map[string]interface{}{
			"application": info,
			"entrant":     entrant,
			"competitive": competitive,
		}
		result.Done = true
		result.Items = response
		return
	} else {
		result.Done = false
		message := `Заявление не найдено.`
		result.Message = &message
		return
	}

}
func (result *ResultInfo) GetApplicationInfoById(idApplication uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var application digest.Application

	db := conn.Where(`id_organization=? AND id=?`, result.User.CurrentOrganization.Id, idApplication).Preload(`Status`).Find(&application)

	if db.RowsAffected > 0 {
		var info interface{}
		info = map[string]interface{}{
			"id":                     application.Id,
			"id_entrant":             application.EntrantsId,
			"app_number":             application.AppNumber,
			"id_status":              application.Status.Id,
			"name_status":            application.Status.Name,
			"registration_date":      application.RegistrationDate,
			"first_higher_education": application.FirstHigherEducation,
			"need_hostel":            application.NeedHostel,
			"distance_test":          application.DistanceTest,
			"distance_place":         application.DistancePlace,
			"special_conditions":     application.SpecialConditions,
			"agreed":                 application.Agreed,
			"original":               application.Original,
			"rating":                 application.Rating,
		}
		result.Done = true
		result.Items = info
		return
	} else {
		result.Done = false
		message := `Заявление не найдено.`
		result.Message = &message
		return
	}

}
func (result *ResultInfo) GetApplicationEntranceTestsById(idApplication uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var application digest.Application

	db := conn.Where(`id_organization=? AND id=?`, result.User.CurrentOrganization.Id, idApplication).Find(&application)

	if db.RowsAffected > 0 {
		var tests []interface{}
		var appEntranceTests []digest.AppEntranceTest
		db = conn.Preload(`EntranceTestDocumentType`).Where(`id_application=?`, idApplication).Find(&appEntranceTests)
		for index, _ := range appEntranceTests {
			db = conn.Model(&appEntranceTests[index]).Related(&appEntranceTests[index].EntranceTest, `IdEntranceTest`)
			db = conn.Model(&appEntranceTests[index].EntranceTest).Related(&appEntranceTests[index].EntranceTest.EntranceTestType, `IdEntranceTestType`)
			db = conn.Model(&appEntranceTests[index].EntranceTest).Related(&appEntranceTests[index].EntranceTest.Subject, `IdSubject`)

			r := make(map[string]interface{})
			r[`test_info`] = map[string]interface{}{
				"id":                      appEntranceTests[index].EntranceTest.Id,
				"id_entrance_test_type":   appEntranceTests[index].EntranceTest.EntranceTestType.Id,
				"name_entrance_test_type": appEntranceTests[index].EntranceTest.EntranceTestType.Name,
				"id_subject":              appEntranceTests[index].EntranceTest.Subject.Id,
				"name_subject":            appEntranceTests[index].EntranceTest.Subject.Name,
				"priority":                appEntranceTests[index].EntranceTest.Priority,
				"uid":                     appEntranceTests[index].EntranceTest.Uid,
				"test_name":               appEntranceTests[index].EntranceTest.TestName,
				"min_score":               appEntranceTests[index].EntranceTest.MinScore,
				"is_ege":                  appEntranceTests[index].EntranceTest.IsEge,
			}
			var doc map[string]interface{}
			if appEntranceTests[index].IdDocument != nil {
				var d digest.Ege
				db = conn.Preload(`DocumentType`).Preload(`Subject`).Where("id=?", *appEntranceTests[index].IdDocument).Find(&d)
				issueDate := d.IssueDate.Format(`2006-01-02`)
				doc = map[string]interface{}{
					"id":            d.Id,
					"code_category": `ege`,
					"name_category": `Свидетельство о результатах ЕГЭ`,
					"subject": map[string]interface{}{
						"name": d.Subject.Name,
						"id":   d.Subject.Id,
					},
					"mark":       d.Mark,
					"issue_date": issueDate,
				}
			}
			issueDate := appEntranceTests[index].IssueDate.Format(`2006-01-02`)
			r[`data`] = map[string]interface{}{
				"id":               appEntranceTests[index].Id,
				"id_entrance_test": appEntranceTests[index].IdEntranceTest,
				"issue_date":       issueDate,
				"uid":              appEntranceTests[index].Uid,
				"doc":              doc,
			}
			tests = append(tests, r)
		}
		result.Items = []digest.AppEntranceTest{}
		if len(tests) > 0 {
			result.Items = tests
		}
		result.Done = true
		return
	} else {
		result.Done = true
		result.Items = []digest.AppEntranceTest{}
		return
	}

}
func (result *ResultInfo) GetApplicationDocsById(idApplication uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var application digest.Application

	db := conn.Where(`id_organization=? AND id=?`, result.User.CurrentOrganization.Id, idApplication).Find(&application)

	if db.RowsAffected > 0 {
		var response interface{}
		var responseDocs []interface{}
		var allDocuments []digest.AllDocuments
		cmd := `
					 with a as(SELECT id_document, id_document_sys_category FROM app.documents WHERE id_application = ?),
					b as (SELECT id, doc_number, id_document_type, doc_series, NULL::integer as mark, NULL::character varying as name_subject, issue_date, 'educations' as name_table  FROM documents.educations educ WHERE EXISTS(SELECT 1 FROM a WHERE educ.id =a.id_document)
					UNION
					SELECT ege.id, doc_number, id_document_type, NULL::character varying as doc_series, mark, sbj.name as name_subject,  issue_date, 'ege' as name_table
						FROM documents.ege ege
						join cls.subjects sbj ON sbj.id = ege.id_subject WHERE EXISTS(SELECT 1 FROM a WHERE ege.id =a.id_document)
					UNION
					SELECT id, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'orphans' as name_table FROM documents.orphans orph WHERE EXISTS(SELECT 1 FROM a WHERE orph.id =a.id_document)
					UNION
					SELECT id, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'veteran' as name_table FROM documents.veteran vet WHERE EXISTS(SELECT 1 FROM a WHERE vet.id =a.id_document)
					UNION
					SELECT id, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'olympics' as name_table FROM documents.olympics olymp WHERE EXISTS(SELECT 1 FROM a WHERE olymp.id =a.id_document)
					UNION
					SELECT id, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'militaries' as name_table FROM documents.militaries mil WHERE EXISTS(SELECT 1 FROM a WHERE mil.id =a.id_document)
					UNION
					SELECT id, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'other' as name_table FROM documents.other oth WHERE EXISTS(SELECT 1 FROM a WHERE oth.id =a.id_document)
					UNION
					SELECT id, doc_number, id_document_type, NULL::character varying as doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'disability' as name_table FROM documents.disability dis WHERE EXISTS(SELECT 1 FROM a WHERE dis.id =a.id_document)
					UNION
					SELECT id, NULL::character varying as doc_number, id_document_type,  NULL::character varying as doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, NULL::timestamp with time zone as issue_date, 'compatriot' as name_table
					FROM documents.compatriot compar WHERE EXISTS(SELECT 1 FROM a WHERE compar.id =a.id_document)
					UNION
					SELECT id, doc_number, id_document_type, doc_series, NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'parents_lost' as name_table FROM documents.parents_lost par WHERE EXISTS(SELECT 1 FROM a WHERE par.id =a.id_document)
					UNION
					SELECT id, doc_number, id_document_type, doc_series, NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'radiation_work' as name_table FROM documents.radiation_work rad WHERE EXISTS(SELECT 1 FROM a WHERE rad.id =a.id_document)
					UNION
					SELECT id, NULL::character varying as doc_number, id_document_type, NULL::character varying as doc_series, NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'composition' as name_table
					FROM documents.composition compos WHERE EXISTS(SELECT 1 FROM a WHERE compos.id =a.id_document)
					UNION
					SELECT id, doc_number, id_document_type, doc_series,  NULL::integer  as mark, NULL::character varying as name_subject, issue_date, 'identification' as name_table
					FROM documents.identification ident WHERE EXISTS(SELECT 1 FROM a WHERE ident.id =a.id_document))
					SELECT b.*, sys.id as id_sys_categories, sys."name" as name_sys_categories, dt."name" as name_document_type
					from b  
					join cls.document_sys_categories sys on b.name_table = sys.name_table
					join cls.document_types dt on dt.id = b.id_document_type			
					Where b.id IS NOT NULL
`
		db = conn.Raw(cmd, idApplication).Scan(&allDocuments)
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

			category.Docs = append(category.Docs, map[string]interface{}{
				"id":               allDocuments[index].Id,
				"doc_number":       allDocuments[index].DocNumber,
				"doc_series":       allDocuments[index].DocSeries,
				"id_document_type": allDocuments[index].IdDocumentType,
				"id_entrant":       application.EntrantsId,
				//"id_sys_categories": 		allDocuments[index].IdSysCategories,
				"issue_date":         issueDate,
				"mark":               allDocuments[index].Mark,
				"name_document_type": allDocuments[index].NameDocumentType,
				"name_subject":       allDocuments[index].NameSubject,
				//"name_sys_categories": 		allDocuments[index].NameSysCategories,
				//"name_table": 				allDocuments[index].NameTable,
			})
			sysCategory[allDocuments[index].NameTable] = category
		}
		for index, _ := range PriorityTable {
			if val, ok := sysCategory[PriorityTable[index]]; ok {
				responseDocs = append(responseDocs, val)
			}
		}
		if len(responseDocs) == 0 {
			result.SetErrorResult(`Документы у заявления не найдены`)
			return
		}
		response = map[string]interface{}{
			"id_application": idApplication,
			"docs":           responseDocs,
		}
		result.Done = true
		result.Items = response
		return
	} else {
		result.Done = false
		message := `Заявление не найдено.`
		result.Message = &message
		return
	}

}
func (result *ResultInfo) GetApplicationsById() {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)

	var response []interface{}

	var applications []digest.Application
	db := conn.Where(`id_organization=?`, result.User.CurrentOrganization.Id).Preload(`ViolationTypes`).Preload(`Entrants`).Preload(`CompetitiveGroup`).Where(``).Find(&applications)
	fmt.Print(len(applications))

	if db.RowsAffected > 0 {
		for index, _ := range applications {
			response = append(response, map[string]interface{}{
				"id":         applications[index].Id,
				"id_entrant": applications[index].EntrantsId,
				"entrant": map[string]interface{}{
					"id":         applications[index].Entrants.Id,
					"surname":    applications[index].Entrants.Surname,
					"name":       applications[index].Entrants.Name,
					"patronymic": applications[index].Entrants.Patronymic,
					"snils":      applications[index].Entrants.Snils,
				},
				"id_organization":             applications[index].IdOrganization,
				"id_competitive_group":        applications[index].CompetitiveGroup.Id,
				"id_competitive_group_target": applications[index].IdCompetitiveGroupTarget,
				"app_number":                  applications[index].AppNumber,
				"registration_date":           applications[index].RegistrationDate,
				"rating":                      applications[index].Rating,
				"id_status":                   applications[index].IdStatus,
				"priority":                    applications[index].Priority,
				"first_higher_education":      applications[index].FirstHigherEducation,
				"need_hostel":                 applications[index].NeedHostel,
				"id_disabled_document":        applications[index].IdDisabledDocument,
				"special_conditions":          applications[index].SpecialConditions,
				"distance_test":               applications[index].DistanceTest,
				"distance_place":              applications[index].DistancePlace,
				"id_violation":                applications[index].ViolationTypes.Id,
				"ege_check":                   applications[index].EgeCheck,
				"agreed":                      applications[index].Agreed,
				"disagreed":                   applications[index].Disagreed,
				"agreed_date":                 applications[index].AgreedDate,
				"disagreed_date":              applications[index].DisagreedDate,
				"id_order_admission":          applications[index].IdOrderAdmission,
				"order_admission_date":        applications[index].OrderAdmissionDate,
				"id_return_type":              applications[index].IdReturnType,
				"return_date":                 applications[index].ReturnDate,
				"original":                    applications[index].Original,
				"original_doc":                applications[index].OriginalDoc,
				"id_benefit":                  applications[index].IdBenefit,
				"uid":                         applications[index].Uid,
				"created":                     applications[index].Created,
				"status_comment":              applications[index].StatusComment,
			})
		}
		result.Done = true
		result.Items = response
		return
	} else {
		result.Done = true
		message := `Завяления не найдены.`
		result.Message = &message
		return
	}

}
func (result *ResultInfo) AddApplication(data AddApplication) {
	conn := &config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var application digest.Application
	application.IdOrganization = result.User.CurrentOrganization.Id
	application.Created = time.Now()

	//
	var existApplication []digest.Application
	_ = conn.Where(`id_entrant=? AND id_competitive_group=?`, data.IdEntrant, data.IdCompetitiveGroup).Find(&existApplication)
	if len(existApplication) > 0 {
		result.SetErrorResult(`Данный абитуриент уже подавал заявление на указанную конкусрную группу`)
		tx.Rollback()
		return
	}

	var competitive digest.CompetitiveGroup
	db := tx.Where(`id_organization=?`, result.User.CurrentOrganization.Id).Find(&competitive, data.IdCompetitiveGroup)
	if competitive.Id < 1 {
		result.SetErrorResult(`Конкурсная группа не найдена`)
		tx.Rollback()
		return
	}
	var idDirection []DirectionCompetitiveGroups
	_ = conn.Raw(`select  cg.id_direction
						from app.applications a
						join cmp.competitive_groups cg on cg.id=a.id_competitive_group
						where id_entrant = ? AND id_organization=?
						group by cg.id_direction `, data.IdEntrant, result.User.CurrentOrganization.Id).Scan(&idDirection)
	if len(idDirection) >= 3 {
		var existDirection []uint
		for _, val := range idDirection {
			existDirection = append(existDirection, val.IdDirection)
		}
		if service.SearchUintInSliceUint(competitive.IdDirection, existDirection) < 0 {
			result.SetErrorResult(`Абитуриент уже подал заявления по трем направлениям.`)
			tx.Rollback()
			return
		}
	}
	application.IdCompetitiveGroup = data.IdCompetitiveGroup

	var entrant digest.Entrants
	db = tx.Find(&entrant, data.IdEntrant)
	if entrant.Id < 1 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	application.EntrantsId = data.IdEntrant
	if data.AppNumber == `` {
		result.SetErrorResult(`Поле app_number не может быть пустым`)
		tx.Rollback()
		return
	}
	application.AppNumber = data.AppNumber
	application.RegistrationDate = data.RegistrationDate
	application.IdStatus = 1
	application.Priority = data.Priority
	application.FirstHigherEducation = data.FirstHigherEducation
	application.NeedHostel = data.NeedHostel
	//application.IdDisabledDocument = data.IdDisabledDocument
	application.SpecialConditions = data.SpecialConditions
	application.DistanceTest = data.DistanceTest
	application.DistanceTest = data.DistanceTest
	application.DistancePlace = data.DistancePlace

	application.Agreed = data.Agreed
	application.AgreedDate = data.AgreedDate
	application.Disagreed = data.Disagreed
	application.DisagreedDate = data.DisagreedDate
	application.Original = data.Original
	application.StatusComment = data.StatusComment

	if data.Uid != nil {
		var exist digest.Application
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *data.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Заявление с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		application.Uid = data.Uid
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&application)
	if db.Error != nil {
		result.SetErrorResult(`Ошибка добавлении заявления ` + db.Error.Error())
		tx.Rollback()
		return
	}
	var idsDocs []uint
	ident := false
	educ := false
	if db.Error == nil {
		if len(data.Docs) > 0 {
			for _, value := range data.Docs {
				if value.Type == `identification` {
					var doc digest.Identifications
					db = conn.Where(`id=? AND id_entrant=?`, value.Id, data.IdEntrant).Find(&doc)
					if doc.Id <= 0 {
						result.SetErrorResult(`Документ ` + fmt.Sprintf(`%v`, value.Id) + ` не найден`)
						tx.Rollback()
						return
					}
					ident = true
					var newDoc digest.Documents
					newDoc.IdDocument = doc.Id
					newDoc.IdApplication = application.Id
					newDoc.IdDocumentSysCategory = 10 // sys_category for identification
					newDoc.Created = time.Now()
					db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&newDoc)
					if db.Error != nil {
						result.SetErrorResult(`Ошибка при добавлении доумента, удостоверяющего личность, у заявления ` + db.Error.Error())
						tx.Rollback()
						return
					}
					idsDocs = append(idsDocs, newDoc.Id)
				} else {
					var doc digest.VDocuments
					db = conn.Where(`id_document=? AND id_entrant=?`, value.Id, data.IdEntrant).Find(&doc)
					if doc.IdDocument <= 0 {
						result.SetErrorResult(`Документ ` + fmt.Sprintf(`%v`, value.Id) + ` не найден`)
						tx.Rollback()
						return
					}
					if doc.NameTable == `educations` {
						educ = true
					}
					var newDoc digest.Documents
					newDoc.IdDocument = doc.IdDocument
					newDoc.IdApplication = application.Id
					newDoc.IdDocumentSysCategory = doc.IdSysCategories
					newDoc.Created = time.Now()

					db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&newDoc)
					if db.Error != nil {
						result.SetErrorResult(`Ошибка при добавлении документа у заявления ` + db.Error.Error())
						tx.Rollback()
						return
					}
					idsDocs = append(idsDocs, newDoc.Id)
					if value.Type == `ege` {
						var ege digest.Ege
						db = conn.Preload(`Subject`).Where(`id=?`, value.Id).Find(&ege)
						// TODO add app.entrance_test
						var entranceTest digest.EntranceTest
						db = conn.Where(`id_competitive_group=? AND id_subject=?`, data.IdCompetitiveGroup, ege.IdSubject).Find(&entranceTest)
						if entranceTest.Id <= 0 {
							result.SetErrorResult(`Вступительный тест у конкурсной группы с предметом "` + fmt.Sprintf(`%v`, ege.Subject.Name) + `" не найден`)
							tx.Rollback()
							return
						}
						rTrue := true
						appEntranceTest := digest.AppEntranceTest{
							IdApplication:  application.Id,
							IdEntranceTest: entranceTest.Id,
							IdDocument:     &value.Id,
							IdResultSource: 1,
							ResultValue:    ege.Mark,
							HasEge:         &rTrue,
							EgeValue:       &ege.Mark,
							IssueDate:      ege.IssueDate,
							Created:        time.Now(),
						}
						db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&appEntranceTest)
						if db.Error != nil {
							result.SetErrorResult(`Ошибка при добавлении вступительного сипытания у заявления ` + db.Error.Error())
							tx.Rollback()
							return
						}
					}
				}
			}
		} else {
			result.SetErrorResult(`Не найдено ни одного документа`)
			tx.Rollback()
			return
		}
		if !educ || !ident {
			if !educ {
				result.SetErrorResult(`Заявление должно содержать хотя бы один документ об образовании`)
				tx.Rollback()
				return
			}
			if !ident {
				result.SetErrorResult(`Заявление должно содержать хотя бы один удостоверяющий документ`)
				tx.Rollback()
				return
			}
		}

		result.Items = map[string]interface{}{
			`id_application`: application.Id,
			`id_documents`:   idsDocs,
		}
		result.Done = true
		tx.Commit()
		return
	} else {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
}
func (result *ResultInfo) AddAppAchievement(data digest.AppAchievements) {
	conn := &config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var achievement digest.IndividualAchievements
	if data.IdAchievement != nil {
		ach, err := digest.GetIndividualAchievements(*data.IdAchievement)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		achievement = *ach
	} else {
		result.SetErrorResult(`Нет достижения`)
		tx.Rollback()
		return
	}

	application, err := digest.GetApplication(data.IdApplication)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}

	var new digest.AppAchievements

	if achievement.IdCampaign != application.CompetitiveGroup.IdCampaign {
		result.SetErrorResult(`Приемная компания достижения не совпадает с приемной компанией заявления`)
		tx.Rollback()
		return
	}

	if data.Mark != nil {
		if *data.Mark > achievement.MaxValue {
			result.SetErrorResult(`Оценка не может превышать максимально допустимое значение достижения`)
			tx.Rollback()
			return
		}
		new.Mark = data.Mark
	}

	if data.Uid != nil {
		var exist digest.AppAchievements
		_ = tx.Where(`upper(uid)=upper(?) AND id_application=?`, data.Uid, data.IdApplication).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Достижекние с данным uid уже существуют`)
			tx.Rollback()
			return
		} else {
			new.Uid = data.Uid
		}
	}

	new.Created = time.Now()
	new.IdApplication = application.Id
	new.IdAchievement = &achievement.Id
	new.IdCategory = achievement.IdCategory
	new.Name = achievement.Name
	if data.IdDocument != nil {
		doc, err := digest.GetVDocuments(*data.IdDocument)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		if doc.EntrantsId != application.EntrantsId {
			result.SetErrorResult(`Абитуриент документа и заявления не совпадает`)
			tx.Rollback()
			return
		}
		new.IdDocument = &doc.IdDocument
	} else {
		result.SetErrorResult(`Где ваши докУменты?`)
		tx.Rollback()
		return
	}
	db := tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&new)
	if db.Error != nil {
		result.SetErrorResult(`Ошибка при добавлении достижения к заявлению ` + db.Error.Error())
		tx.Rollback()
		return
	}
	result.Items = map[string]interface{}{
		`id_application`:  application.Id,
		`id_achievements`: new.Id,
	}
	result.Done = true
	tx.Commit()
	return

}
func (result *ResultInfo) SetStatusApplication(data ChangeStatusApplication) {
	conn := &config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	if data.CodeStatus == `` {
		result.SetErrorResult(`Пустой статус`)
		tx.Rollback()
		return
	}
	status, err := GetApplicationStatusByCode(data.CodeStatus)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	application, err := digest.GetApplication(data.IdApplication)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}

	if application.IdStatus == status.Id {
		result.SetErrorResult(`Заявление уже в этом статусе`)
		tx.Rollback()
		return
	}
	application.IdStatus = status.Id
	db := tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&application)
	if db.Error != nil {
		result.SetErrorResult(`Ошибка при изменении статуса заявления ` + db.Error.Error())
		tx.Rollback()
		return
	}
	result.Items = map[string]interface{}{
		`id_application`: application.Id,
		`new_status`:     application.IdStatus,
	}
	result.Done = true
	tx.Commit()
	return
}
func (result *ResultInfo) GetApplicationStatuses(keys map[string][]string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var db *gorm.DB
	var statuses []digest.ApplicationStatuses
	//db = conn.Select(`id, name`).Table(`cls.v_direction_specialty`)
	db = conn.Where(`actual`)
	if len(keys[`search`]) > 0 {
		db = db.Where(`UPPER(name) LIKE ?`, `%`+strings.ToUpper(keys[`search`][0])+`%`)
	}

	db = db.Find(&statuses)
	if db.Error != nil {
		message := db.Error.Error()
		result.Message = &message
		return
	}
	var response []interface{}
	for index, _ := range statuses {
		response = append(response, map[string]interface{}{
			"id":   statuses[index].Id,
			"name": statuses[index].Name,
			"code": statuses[index].Code,
		})
	}
	result.Done = true
	result.Items = response
	return
}
func (result *ResultInfo) RemoveApplicationAchievement(idAchievement uint) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var old digest.AppAchievements
	db := tx.Find(&old, idAchievement)
	if old.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Индивидуальное достижение не найдено`)
		tx.Rollback()
		return
	}
	var application digest.Application
	db = tx.Find(&application, old.IdApplication)
	if application.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Заявление не найдено`)
		tx.Rollback()
		return
	}
	if application.IdOrganization != result.User.CurrentOrganization.Id {
		result.SetErrorResult(`Организация заявления не совпадает с выбранной вами`)
		tx.Rollback()
		return
	}
	db = tx.Where(`id=?`, idAchievement).Delete(&old)
	if db.Error == nil {
		result.Done = true
		tx.Commit()
		result.Items = map[string]interface{}{
			`id_achievement`: idAchievement,
		}
		return
	} else {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
}
func GetApplicationStatusByCode(code string) (*digest.ApplicationStatuses, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item digest.ApplicationStatuses
	db := conn.Where(`code=?`, code).Find(&item)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Статус не найден. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if item.Id <= 0 {
		return nil, errors.New(`Статус не найден. `)
	}
	return &item, nil
}
func CheckApplicationByUser(idApplication uint, user digest.User) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var count int
	if user.Role.Code == `administrator` {
		return nil
	}
	db := conn.Table(`app.applications`).Select(`id`).Where(`id_organization=? AND id=?`, user.CurrentOrganization.Id, idApplication).Count(&count)
	if db.Error != nil {
		return db.Error
	}
	if count > 0 {
		return nil
	} else {
		return errors.New(`У пользователя нет доступа к данному заявлению `)
	}
}
