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

func AddEntrantHandler(r *mux.Router) {
	r.HandleFunc("/entrants/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.EntrantsSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetListEntrants()
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/entrants/{id:[0-9]+}/short", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetInfoEntrantApp(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	r.HandleFunc("/entrants/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var data handlers.AddEntrantData
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &data)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.AddEntrant(data)
		} else {
			m := err.Error()
			res.Message = &m
		}
		service.ReturnJSON(w, res)
	}).Methods("Post")

	r.HandleFunc("/entrants/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.GetInfoEntrant(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/entrants/{id:[0-9]+}/others", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.GetDocsEntrant(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/entrants/{id:[0-9]+}/idents", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.GetDocsIdentsEntrant(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/entrants/{id:[0-9]+}/idents/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.GetListDocsIdentsEntrant(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/entrants/{id:[0-9]+}/docs/short", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		keys := r.URL.Query()
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.GetShortListDocsEntrant(uint(id), keys)
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

}
