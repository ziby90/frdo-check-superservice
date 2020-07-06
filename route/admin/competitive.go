package route_admin

import (
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/handlers/admin"
	"persons/service"
)

func AddCompetitiveGroupsHandler(r *mux.Router) {
	// просто список конкурсов для выбора при подаче заявления
	r.HandleFunc("/competitive/list", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo
		keys := r.URL.Query()
		res.User = *handlers.CheckAuthCookie(r)
		// TODO ограничение на пять вузов
		res.GetListCompetitiveGroups(keys)
		service.ReturnJSON(w, res)
	}).Methods("GET")
	//// просто список конкурсов для выбора при подаче заявления
	//r.HandleFunc("/competitive/tests/{id:[0-9]+}/date/short", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	vars := mux.Vars(r)
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	// TODO ограничение на пять вузов
	//	res.GetListDatesByEntranceTest(uint(id))
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
	//// добавление конкурсной группы
	//r.HandleFunc("/competitive/add", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	var cmp handlers_admin.AddCompetitiveGroup
	//	b, _ := ioutil.ReadAll(r.Body)
	//	err := json.Unmarshal(b, &cmp)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	err = handlers_admin.CheckCampaignByUser(cmp.CompetitiveGroup.IdCampaign, res.User)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//
	//	res.AddCompetitive(cmp)
	//	service.ReturnJSON(w, res)
	//}).Methods("Post")
	//// изменение конкурсной группы
	//r.HandleFunc("/competitive/{id:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	var cmp handlers_admin.CompetitiveGroup
	//	vars := mux.Vars(r)
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	b, _ := ioutil.ReadAll(r.Body)
	//	err = json.Unmarshal(b, &cmp)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	cmp.Id = uint(id)
	//	res.EditCompetitive(cmp)
	//	service.ReturnJSON(w, res)
	//}).Methods("Post")
	//// изменение только количества мест конкурсной группы
	//r.HandleFunc("/competitive/{id:[0-9]+}/number/edit", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	var cmp handlers_admin.EditNumberCompetitive
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//
	//	b, _ := ioutil.ReadAll(r.Body)
	//	err = json.Unmarshal(b, &cmp)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	cmp.IdCompetitive = uint(id)
	//	res.EditNumberCompetitive(cmp)
	//	service.ReturnJSON(w, res)
	//}).Methods("Post")
	//// добавление вступительного испытания для конкурсной группы
	//r.HandleFunc("/competitive/{id:[0-9]+}/entrance/add", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	var cmp handlers_admin.AddEntrance
	//	b, _ := ioutil.ReadAll(r.Body)
	//	err := json.Unmarshal(b, &cmp)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	err = handlers_admin.CheckCompetitiveGroupByUser(uint(id), res.User)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.AddEntrance(uint(id), cmp)
	//	service.ReturnJSON(w, res)
	//}).Methods("Post")
	//// добавление образовательной программы для конкурсной группы
	//r.HandleFunc("/competitive/{id:[0-9]+}/program/add", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	var cmp handlers_admin.AddCompetitiveGroup
	//	b, _ := ioutil.ReadAll(r.Body)
	//	err := json.Unmarshal(b, &cmp)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	err = handlers_admin.CheckCompetitiveGroupByUser(uint(id), res.User)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.AddProgram(uint(id), cmp)
	//	service.ReturnJSON(w, res)
	//}).Methods("Post")
	//// проверка цифры конкурсной группы
	//r.HandleFunc("/competitive/check", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultCheck
	//	var cmp handlers_admin.AddCompetitiveGroup
	//	b, _ := ioutil.ReadAll(r.Body)
	//	err := json.Unmarshal(b, &cmp)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	// TODO заглушка. не забудь. UPD СатАня сказала, что не сипользует эту функцию.
	//	err = handlers_admin.CheckCampaignByUser(cmp.CompetitiveGroup.IdCampaign, res.User)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.CheckNumberAddCompetitive()
	//	service.ReturnJSON(w, res)
	//}).Methods("Post")
	//// инфа по конкусрной группе
	//r.HandleFunc("/competitive/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
	//	res := handlers_admin.ResultInfo{}
	//	vars := mux.Vars(r)
	//
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	err = handlers_admin.CheckCompetitiveGroupByUser(uint(id), res.User)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.GetInfoCompetitiveGroup(uint(id))
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
	//// Образовательные программы конкурсной группы
	//r.HandleFunc("/competitive/{id:[0-9]+}/programs", func(w http.ResponseWriter, r *http.Request) {
	//	res := handlers_admin.ResultInfo{}
	//	vars := mux.Vars(r)
	//
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	err = handlers_admin.CheckCompetitiveGroupByUser(uint(id), res.User)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.GetEducationProgramsCompetitiveGroup(uint(id))
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
	//// Вступительные испытания конкурсной группы
	//r.HandleFunc("/competitive/{id:[0-9]+}/tests", func(w http.ResponseWriter, r *http.Request) {
	//	res := handlers_admin.ResultInfo{}
	//	vars := mux.Vars(r)
	//
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	err = handlers_admin.CheckCompetitiveGroupByUser(uint(id), res.User)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.GetEntranceTestsCompetitiveGroup(uint(id))
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
	//// удаление вступительного испытания у конкусрной группы
	//r.HandleFunc("/competitive/{id:[0-9]+}/entrance/{id_entrance:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	vars := mux.Vars(r)
	//
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	idEntrance, err := strconv.ParseInt(vars[`id_entrance`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id_entrance.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	err = handlers_admin.CheckCompetitiveGroupByUser(uint(id), res.User)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.RemoveEntranceCompetitive(uint(id), uint(idEntrance))
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
	//// удаление конкусрной группы вместе с испытаниями и программами
	//r.HandleFunc("/competitive/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	err = handlers_admin.CheckCompetitiveGroupByUser(uint(id), res.User)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.RemoveCompetitive(uint(id))
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
	//// удаление образовательной программы у конкурсной группы
	//r.HandleFunc("/competitive/{id:[0-9]+}/program/{id_program:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	vars := mux.Vars(r)
	//
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	idProgram, err := strconv.ParseInt(vars[`id_program`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id_program.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	err = handlers_admin.CheckCompetitiveGroupByUser(uint(id), res.User)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.RemoveProgramCompetitive(uint(id), uint(idProgram))
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
	//// выпадайка на список вступительных испытаний у конкурсной группы
	//r.HandleFunc("/competitive/{id:[0-9]+}/tests/select", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultList
	//	keys := r.URL.Query()
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	res.MakeUrlParams(keys)
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	err = handlers_admin.CheckCompetitiveGroupByUser(uint(id), res.User)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.GetEntranceTestsSelectListByCompetitive(uint(id))
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
	//
	//// добавление дат вступительных испытаний
	//r.HandleFunc("/competitive/entrance/{id:[0-9]+}/date/add", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	var cmp handlers_admin.AddEntranceTestDate
	//	b, _ := ioutil.ReadAll(r.Body)
	//	err := json.Unmarshal(b, &cmp)
	//	vars := mux.Vars(r)
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	if err != nil {
	//		message := err.Error()
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	cmp.IdEntranceTest = uint(id)
	//	res.AddEntranceTestCalendar(cmp)
	//	service.ReturnJSON(w, res)
	//}).Methods("Post")
	//// удаление даты вступительных испытаний
	//r.HandleFunc("/competitive/entrance/date/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	vars := mux.Vars(r)
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.RemoveEntranceTestCalendar(uint(id))
	//	service.ReturnJSON(w, res)
	//}).Methods("Post")
	//// список дат вступительных испытаний
	//r.HandleFunc("/competitive/entrance/{id:[0-9]+}/date", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	vars := mux.Vars(r)
	//	res.User = *handlers_admin.CheckAuthCookie(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err != nil {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//		service.ReturnJSON(w, res)
	//		return
	//	}
	//	res.GetEntranceTestsCalendarByEntrance(uint(id))
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")

}
