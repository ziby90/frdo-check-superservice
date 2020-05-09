package handlers

import (
	"fmt"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
)

var ApplicationSearchArray = []string{
	`hz`,
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
		info = map[string]interface{}{
			"id":                     application.Id,
			"app_number":             application.AppNumber,
			"name_competitive_group": application.CompetitiveGroup.Name,
			"entrant_fullname":       application.Entrants.Surname + ` ` + application.Entrants.Name + ` ` + application.Entrants.Patronymic,
			"entrant_snils":          application.Entrants.Snils,
			"id_status":              application.Status.Id,
			"name_status":            application.Status.Name,
			"registration_date":      application.RegistrationDate,
			"agreed":                 application.Agreed,
			"original":               application.Original,
			"rating":                 application.Rating,
		}
		var docs []digest.Documents
		var responseDocs []interface{}
		db = conn.Preload(`DocumentSysCategory`).Where(`id_application=?`, idApplication).Find(&docs)
		for index, _ := range docs {
			if docs[index].DocumentSysCategory.NameTable == `identification` {
				var doc digest.Identifications
				db = conn.Where(`id=?`, docs[index].IdDocument).Find(&doc)
				issueDate := doc.IssueDate.Format(`02-01-2006`)
				responseDocs = append(responseDocs, map[string]interface{}{
					"id":            doc.Id,
					"name_type":     doc.DocumentType.Name,
					"doc_series":    doc.DocSeries,
					"name_category": `identification`,
					"doc_number":    doc.DocNumber,
					"issue_date":    issueDate,
				})
			} else {
				var doc digest.VDocuments
				db = conn.Where(`id_document=?`, docs[index].IdDocument).Find(&doc)
				responseDocs = append(responseDocs, map[string]interface{}{
					"id":            doc.IdDocument,
					"name_type":     doc.DocumentName,
					"doc_name":      doc.DocName,
					"name_category": doc.NameTable,
				})
			}
		}
		response = map[string]interface{}{
			"application": info,
			"docs":        responseDocs,
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
				"id": applications[index].Id,
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
	//existDb := conn.Where(`id_organization=? AND uid=?`)

	//
	//
	//if len(existApplication)>=0 {
	//	var existCompetitiveIds []uint
	//	for _, val := range existApplication{
	//		existCompetitiveIds = append(existCompetitiveIds, val.IdCompetitiveGroup)
	//	}
	//	if len(existCompetitiveIds)>0 {
	//		db = db.Where(`id NOT IN (?)`, existCompetitiveIds)
	//	}
	//}
	//

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
						where id_entrant = ?
						group by cg.id_direction `, data.IdEntrant).Scan(&idDirection)
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
						result.SetErrorResult(db.Error.Error())
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
						result.SetErrorResult(db.Error.Error())
						tx.Rollback()
						return
					}
					idsDocs = append(idsDocs, newDoc.Id)
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
