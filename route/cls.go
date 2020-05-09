package route

import (
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/service"
)

func AddClsHandler(r *mux.Router) {

	r.HandleFunc("/cls/list/{cls_name}", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultCls
		vars := mux.Vars(r)
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.GetClsResponse(vars[`cls_name`])
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/cls/sys_category", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultCls
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.GetClsSysCategoryResponse()
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/directions/list", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		keys := r.URL.Query()
		res.GetDirectionByEntrant(keys)
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
