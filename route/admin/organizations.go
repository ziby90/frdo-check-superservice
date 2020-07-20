package route_admin

import (
	"github.com/gorilla/mux"
	"net/http"
	handlers_admin "persons/handlers/admin"
	"persons/service"
)

func AddOrgsAddminHandler(r *mux.Router) {
	//список организаций
	r.HandleFunc("/organizations/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers_admin.OrganizationSearchArray)
		res.GetListOrganization()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	//список организаций для выпадайки
	r.HandleFunc("/organizations/short", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo
		res.GetListOrganizationShort()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
}
