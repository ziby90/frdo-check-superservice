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

func AddOrgsHandler(r *mux.Router) {
	//список организаций
	r.HandleFunc("/organizations/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.GetListOrganization()
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// инфа по организации
	r.HandleFunc("/organizations/{id:[0-9]+}/info", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{
			Done:    false,
			Message: nil,
			Items:   nil,
		}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetInfoCampaign(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// добавление направления обучения организации
	r.HandleFunc("/organizations/directions/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.IdsDirectionOrganization
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.AddOrganizationDirection(cmp)
		} else {
			message := err.Error()
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")
	// удаление направления обучения организации
	r.HandleFunc("/organizations/directions/remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.IdsDirectionOrganization
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.RemoveOrganizationDirection(cmp)
		} else {
			message := err.Error()
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")
	// список всех направлений
	r.HandleFunc("/directions/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.CampaignSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetDirections()
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// список всех направлений
	r.HandleFunc("/organizations/directions/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.DirectionsSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetDirectionsByOrganization(keys)
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// направления обучения по вузам
	r.HandleFunc("/organizations/directions/select", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultList
		keys := r.URL.Query()
		res.User = *handlers.CheckAuthCookie(r)

		res.MakeUrlParams(keys)
		res.GetDirectionsSelectListByOrg(keys)
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// группы направлений обучения по вузам
	r.HandleFunc("/organizations/directions/parents/select", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultList
		keys := r.URL.Query()
		res.User = *handlers.CheckAuthCookie(r)

		res.MakeUrlParams(keys)
		res.GetDirectionsParentsSelectListByOrg(keys)
		service.ReturnJSON(w, res)
	}).Methods("GET")
}
