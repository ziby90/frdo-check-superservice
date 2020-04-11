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
}
