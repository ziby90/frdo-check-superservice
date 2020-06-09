package route

import (
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/service"
)

func AddClsHandler(r *mux.Router) {
	// список справочников для чики-выпадаек. ВНутри есть массив с перечислением таблиц
	r.HandleFunc("/cls/list/{cls_name}", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultCls
		vars := mux.Vars(r)
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.GetClsResponse(vars[`cls_name`])
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// отдельно наши основные категории документов
	r.HandleFunc("/cls/sys_category", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultCls
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.GetClsSysCategoryResponse()
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// и отдельно список направлений
	r.HandleFunc("/directions/short", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		keys := r.URL.Query()
		res.GetDirectionByEntrant(keys)
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// и отдельно список статусов заявлений с кодами
	r.HandleFunc("/cls/appstatuses/short", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		keys := r.URL.Query()
		res.GetApplicationStatuses(keys)
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// и отдельно список статусов заявлений с кодами
	r.HandleFunc("/cls/cmpstatuses/short", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		keys := r.URL.Query()
		res.GetCampaignStatuses(keys)
		service.ReturnJSON(w, res)
	}).Methods("GET")

	//r.HandleFunc("/cls/score/{id_subject:[0-9]+}/{year:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers.ResultCls
	//	vars := mux.Vars(r)
	//	idSubject, err := strconv.ParseInt(vars[`id_subject`], 10, 32)
	//	if err == nil {
	//		year, err := strconv.ParseInt(vars[`year`], 10, 32)
	//		if err == nil {
	//			res.GetMinScoreSubject(uint(idSubject), uint(year))
	//		} else {
	//			message := `Неверный параметр year.`
	//			res.Message = &message
	//		}
	//	} else {
	//		message := `Неверный параметр id_subject.`
	//		res.Message = &message
	//	}
	//
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
}
