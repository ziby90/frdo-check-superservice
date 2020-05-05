package handlers

import (
	"fmt"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
)

var ApplicationSearchArray = []string{
	`hz`,
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
	conn := config.Db.ConnGORM
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
		if service.SearchStringInSliceString(search[0], CampaignSearchArray) >= 0 {
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

func (result *ResultInfo) GetApplicationsById() {
	result.Done = false
	conn := config.Db.ConnGORM
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
