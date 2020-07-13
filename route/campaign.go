package route

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"persons/digest"
	"persons/handlers"
	"persons/service"
	"strconv"
	"time"
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
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `cmp.campaign`,
			Action:      "view",
			Result:      res.Done,
			OldData:     nil,
			NewData:     nil,
			Errors:      res.Errors,
			Created:     time.Now(),
			Source:      "cabinet",
			Route:       &r.URL.Path,
			IdAuthor:    res.User.Id,
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// конкурсы приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/competitive", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.CompetitiveSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		err = handlers.CheckCampaignByUser(uint(id), res.User)
		if err != nil {
			message := err.Error()
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		res.GetListCompetitiveGroupsByCompanyId(uint(id))
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// редактирование приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.CampaignMain
		vars := mux.Vars(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `cmp.campaign`,
			Action:      "edit",
			Created:     time.Now(),
			Source:      "cabinet",
			Route:       &r.URL.Path,
			IdAuthor:    res.User.Id,
		}
		res.User = *handlers.CheckAuthCookie(r)
		res.PrimaryLogging.IdAuthor = res.User.Id
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			res.PrimaryLogging.Result = false
			res.PrimaryLogging.Errors = res.Message
			service.ReturnJSON(w, &res)
			return
		}
		idCampaign := uint(id)
		res.PrimaryLogging.IdObject = &idCampaign
		err = handlers.CheckCampaignByUser(idCampaign, res.User)
		if err != nil {
			message := err.Error()
			res.Message = &message
			res.PrimaryLogging.Result = false
			res.PrimaryLogging.Errors = res.Message
			service.ReturnJSON(w, &res)
			return
		}
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &cmp)
		if err != nil {
			m := err.Error()
			res.Message = &m
			res.PrimaryLogging.Result = false
			res.PrimaryLogging.Errors = res.Message
			service.ReturnJSON(w, &res)
			return
		}
		cmp.Id = idCampaign
		res.EditCampaign(cmp)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// короткий список компаний для выпадаек
	r.HandleFunc("/campaign/short", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultList
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetShortListCampaign()
		service.ReturnJSON(w, &res)
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
		service.ReturnJSON(w, &res)
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
		service.ReturnJSON(w, &res)
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
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// добавить уровни образования у приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/education_levels/add", func(w http.ResponseWriter, r *http.Request) {
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
				var cmp handlers.AddEducationLevels
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &cmp)
				if err == nil {
					cmp.IdCampaign = uint(id)
					res.AddEducationLevelsCampaign(cmp)
				} else {
					m := err.Error()
					res.Message = &m
				}
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// добавить формы обучения у приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/education_forms/add", func(w http.ResponseWriter, r *http.Request) {
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
				var cmp handlers.AddEducationForms
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &cmp)
				if err == nil {
					cmp.IdCampaign = uint(id)
					res.AddEducationFormsCampaign(cmp)
				} else {
					m := err.Error()
					res.Message = &m
				}
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// это что? какой то спсиок предметов по вузам . ааа, выбор предметов для егэ по годам приемной компании! Женя, КОля - привет!
	// Аня сказала "это список предметов ЕГЭ с их минимальными баллами в зависимости от года, каждый год разный минимальный балл"
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
		service.ReturnJSON(w, &res)
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
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// удаление приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.RemoveCampaign(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// меняем статус приемной компании POST
	r.HandleFunc("/campaign/{id:[0-9]+}/status/set", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		var cmp handlers.ChangeStatusCampaign
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err == nil {
				cmp.IdCampaign = uint(id)
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &cmp)
				if err == nil {
					res.SetStatusCampaign(cmp)
				} else {
					message := err.Error()
					res.Message = &message
				}
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")

	// даты окончания приемной компании
	r.HandleFunc("/campaign/{id:[0-9]+}/enddate", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		err = handlers.CheckCampaignByUser(uint(id), res.User)
		if err != nil {
			message := err.Error()
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}

		res.GetEndDateCampaign(uint(id))
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// редактирование даты окончания
	r.HandleFunc("/campaign/{id:[0-9]+}/enddate/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.AddEndData
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		err = handlers.CheckCampaignByUser(uint(id), res.User)
		if err != nil {
			message := err.Error()
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &cmp)
		if err != nil {
			m := err.Error()
			res.Message = &m
			service.ReturnJSON(w, &res)
			return
		}
		cmp.IdCampaign = uint(id)
		res.EditEndDateCampaign(cmp)
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// обнуление даты окончания
	r.HandleFunc("/campaign/{id:[0-9]+}/enddate/{id_end_date:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				idEndDate, err := strconv.ParseInt(vars[`id_end_date`], 10, 32)
				if err == nil {
					res.RemoveEndDateCampaign(uint(id), uint(idEndDate))
				} else {
					message := `Неверный параметр id_end_date.`
					res.Message = &message
				}
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// добавление дат по волнам набора
	r.HandleFunc("/campaign/{id:[0-9]+}/accept_phases/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		err = handlers.CheckCampaignByUser(uint(id), res.User)
		if err != nil {
			message := err.Error()
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}

		var cmp handlers.AddAppAcceptPhases
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &cmp)
		if err != nil {
			m := err.Error()
			res.Message = &m
			service.ReturnJSON(w, &res)
			return
		}
		cmp.IdCampaign = uint(id)
		res.AddAppAcceptPhases(cmp)
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// даты окончания приемной компании по волнам набора
	r.HandleFunc("/campaign/{id:[0-9]+}/accept_phases", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		err = handlers.CheckCampaignByUser(uint(id), res.User)
		if err != nil {
			message := err.Error()
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}

		res.GetAppAcceptPhasesCampaign(uint(id))
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// редактирование дат по волнам набора
	r.HandleFunc("/campaign/{id:[0-9]+}/accept_phases/{id_end_application:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.CampaignAppAcceptPhases
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		idEndApplication, err := strconv.ParseInt(vars[`id_end_application`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		err = handlers.CheckCampaignByUser(uint(id), res.User)
		if err != nil {
			message := err.Error()
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &cmp)
		if err != nil {
			m := err.Error()
			res.Message = &m
			service.ReturnJSON(w, &res)
			return
		}

		res.EditAppAcceptPhasesCampaign(uint(id), uint(idEndApplication), cmp)
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// редактирование дат по волнам набора
	r.HandleFunc("/campaign/{id:[0-9]+}/accept_phases/{id_end_application:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		idEndApplication, err := strconv.ParseInt(vars[`id_end_application`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		err = handlers.CheckCampaignByUser(uint(id), res.User)
		if err != nil {
			message := err.Error()
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		res.RemoveAppAcceptPhasesCampaign(uint(id), uint(idEndApplication))
		service.ReturnJSON(w, &res)
	}).Methods("POST")

}
