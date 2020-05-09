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

	r.HandleFunc("/competitive/list", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		keys := r.URL.Query()
		res.User = *handlers.CheckAuthCookie(r)
		res.GetListCompetitiveGroups(keys)
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/competitive/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.AddCompetitiveGroup
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			err = handlers.CheckCampaignByUser(cmp.CompetitiveGroup.IdCampaign, res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.AddCompetitive(cmp)
			}
		} else {
			message := err.Error()
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("Post")

	r.HandleFunc("/competitive/{id:[0-9]+}/entrance/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.AddEntrance
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			id, err := strconv.ParseInt(vars[`id`], 10, 32)
			if err == nil {
				err = handlers.CheckCompetitiveGroupByUser(uint(id), res.User)
				if err != nil {
					message := err.Error()
					res.Message = &message
				} else {
					res.AddEntrance(uint(id), cmp)
				}
			} else {
				message := `Неверный параметр id.`
				res.Message = &message
			}
		} else {
			message := err.Error()
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("Post")

	r.HandleFunc("/competitive/{id:[0-9]+}/program/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.AddCompetitiveGroup
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			id, err := strconv.ParseInt(vars[`id`], 10, 32)
			if err == nil {
				err = handlers.CheckCompetitiveGroupByUser(uint(id), res.User)
				if err != nil {
					message := err.Error()
					res.Message = &message
				} else {
					res.AddProgram(uint(id), cmp)
				}
			} else {
				message := `Неверный параметр id.`
				res.Message = &message
			}
		} else {
			message := err.Error()
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("Post")

	r.HandleFunc("/competitive/check", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultCheck
		var cmp handlers.AddCompetitiveGroup
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			err = handlers.CheckCampaignByUser(cmp.CompetitiveGroup.IdCampaign, res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.CheckNumberAddCompetitive()
			}
		} else {
			message := err.Error()
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("Post")

	r.HandleFunc("/organizations/directions", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultList
		keys := r.URL.Query()
		idEducationLevel := uint(0)
		res.User = *handlers.CheckAuthCookie(r)
		if len(keys[`id_education_level`]) > 0 {
			if v, ok := strconv.Atoi(keys[`id_education_level`][0]); ok == nil {
				idEducationLevel = uint(v)
			}
		}
		res.MakeUrlParams(keys)
		res.GetDirectiontsListByOrg(idEducationLevel)
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/competitive/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCompetitiveGroupByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.GetInfoCompetitiveGroup(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/competitive/{id:[0-9]+}/entrance/{id_entrance:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			idEntrance, err := strconv.ParseInt(vars[`id_entrance`], 10, 32)
			if err == nil {
				err = handlers.CheckCompetitiveGroupByUser(uint(id), res.User)
				if err != nil {
					message := err.Error()
					res.Message = &message
				} else {
					res.RemoveEntranceCompetitive(uint(id), uint(idEntrance))
				}
			} else {
				message := `Неверный параметр id_entrance.`
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/competitive/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCompetitiveGroupByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.RemoveCompetitive(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/competitive/{id:[0-9]+}/program/{id_program:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			idProgram, err := strconv.ParseInt(vars[`id_program`], 10, 32)
			if err == nil {
				err = handlers.CheckCompetitiveGroupByUser(uint(id), res.User)
				if err != nil {
					message := err.Error()
					res.Message = &message
				} else {
					res.RemoveProgramCompetitive(uint(id), uint(idProgram))
				}
			} else {
				message := `Неверный параметр id_program.`
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
}
