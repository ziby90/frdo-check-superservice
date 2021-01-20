package route

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"persons/handlers"
	"persons/service"
)

func AddOrgsNoAuthHandler(r *mux.Router) {
	r.HandleFunc("/api/organizations/short", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.GetListOrganizationShort()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
}

func AddOrgsHandler(r *mux.Router) {
	// инфа по организации
	r.HandleFunc("/organizations/info", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{
			Done:    false,
			Message: nil,
			Items:   nil,
		}
		res.User = *handlers.CheckAuthCookie(r)
		res.GetInfoOrganization()

		service.ReturnJSON(w, &res)
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
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// изменение параметров отправки в ИС ООВО
	r.HandleFunc("/organizations/isoovo/change", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp struct {
			IsOOVO bool `json:"is_oovo"`
		}
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			//res.Items = cmp
			res.SetIsOOVOOrganization(cmp.IsOOVO)
		} else {
			message := err.Error()
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
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
		service.ReturnJSON(w, &res)
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
		service.ReturnJSON(w, &res)
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
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// направления обучения по вузам
	r.HandleFunc("/organizations/directions/select", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultList
		keys := r.URL.Query()
		res.User = *handlers.CheckAuthCookie(r)

		res.MakeUrlParams(keys)
		res.GetDirectionsSelectListByOrg(keys)
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// группы направлений обучения по вузам
	r.HandleFunc("/organizations/directions/parents/select", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultList
		keys := r.URL.Query()
		res.User = *handlers.CheckAuthCookie(r)

		res.MakeUrlParams(keys)
		res.GetDirectionsParentsSelectListByOrg(keys)
		service.ReturnJSON(w, &res)
	}).Methods("GET")
}
