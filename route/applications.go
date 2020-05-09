package route

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddApplicationHandler(r *mux.Router) {
	r.HandleFunc("/applications/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		//res.MakeUrlParamsSearch(keys, handlers.ApplicationSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetApplications()
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/applications/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		var cmp handlers.AddApplication
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		if err == nil {
			res.AddApplication(cmp)
		} else {
			message := err.Error()
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")

	r.HandleFunc("/entrants/{id:[0-9]+}/applications", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetApplicationsByEntrant(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/applications/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetApplicationById(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

	//
	//r.HandleFunc("/entrants/{id:[0-9]+}/others", func(w http.ResponseWriter, r *http.Request) {
	//	res := handlers.ResultInfo{}
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err == nil {
	//		res.GetDocsEntrant(uint(id))
	//	} else {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//	}
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
	//
	//r.HandleFunc("/entrants/{id:[0-9]+}/idents", func(w http.ResponseWriter, r *http.Request) {
	//	res := handlers.ResultInfo{}
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err == nil {
	//		res.GetDocsIdentsEntrant(uint(id))
	//	} else {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//	}
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
	//
	//r.HandleFunc("/entrants/{id:[0-9]+}/idents/list", func(w http.ResponseWriter, r *http.Request) {
	//	res := handlers.ResultInfo{}
	//
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err == nil {
	//		res.GetListDocsIdentsEntrant(uint(id))
	//	} else {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//	}
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")

}
