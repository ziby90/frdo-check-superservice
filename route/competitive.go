package route

import (
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddCompetitiveGroupsHandler(r *mux.Router) {

	r.HandleFunc("/campaign/{id:[0-9]+}/competitive", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetListCompetitiveGroupsByCompanyId(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	//
	//r.HandleFunc("/achievements/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
	//	res := handlers.ResultInfo{
	//		Done:    false,
	//		Message: nil,
	//		Items:   nil,
	//	}
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err == nil {
	//		res.GetInfoAchievement(uint(id))
	//	} else {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//	}
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")

}
