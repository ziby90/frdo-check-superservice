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

func AddCampaignHandler(r *mux.Router) {
	// список приемных компаний
	r.HandleFunc("/campaign/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.CampaignSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetListCampaign()
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// конкурсы приемной компании
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
	// конкурсы приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/enddate", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetEndDateCampaign(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// редактирование даты окончания
	r.HandleFunc("/campaign/{id:[0-9]+}/enddate/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.AddEndData
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &cmp)
				if err == nil {
					cmp.IdCampaign = uint(id)
					res.EditEndDateCampaign(cmp)
				} else {
					m := err.Error()
					res.Message = &m
				}
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")
	// редактирование приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.CampaignMain
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &cmp)
				if err == nil {
					cmp.Id = uint(id)
					res.EditCampaign(cmp)
				} else {
					m := err.Error()
					res.Message = &m
				}
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")
	// короткий список компаний для выпадаек
	r.HandleFunc("/campaign/short", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultList
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetShortListCampaign()
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// добавление приемной компании
	r.HandleFunc("/campaign/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.CampaignMain
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		if err == nil {
			res.AddCampaign(cmp, *handlers.CheckAuthCookie(r))
		} else {
			m := err.Error()
			res.Message = &m
		}
		service.ReturnJSON(w, res)
	}).Methods("Post")
	// информация по приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.GetInfoCampaign(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// уровни образования у приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/education_levels", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.GetEducationLevelCampaign(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// это что? какой то спсиок предметов по вузам . ааа, выбор предметов для егэ по годам приемной компании! Женя, КОля - привет!
	r.HandleFunc("/campaign/{id:[0-9]+}/subjects", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.GetSubjectsNoEge(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// ФОрмы образования приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/education_forms", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.GetEducationFormCampaign(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// удаление приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {

		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		_, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			//TODO сделать удаление
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
}
