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
	// просто список конкурсов для выбора при подаче заявления
	r.HandleFunc("/competitive/list", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		keys := r.URL.Query()
		res.User = *handlers.CheckAuthCookie(r)
		// TODO ограничение на пять вузов
		res.GetListCompetitiveGroups(keys)
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// добавление конкурсной группы
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
	// изменение конкурсной группы
	r.HandleFunc("/competitive/{id:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.CompetitiveGroup
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			b, _ := ioutil.ReadAll(r.Body)
			err := json.Unmarshal(b, &cmp)
			res.User = *handlers.CheckAuthCookie(r)
			if err == nil {
				cmp.Id = uint(id)
				res.EditCompetitive(cmp)
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("Post")
	// редактирование конкурсной группы
	r.HandleFunc("/competitive/{id:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.CompetitiveGroup
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			b, _ := ioutil.ReadAll(r.Body)
			err := json.Unmarshal(b, &cmp)
			res.User = *handlers.CheckAuthCookie(r)
			if err == nil {
				err = handlers.CheckCampaignByUser(uint(id), res.User)
				if err != nil {
					message := err.Error()
					res.Message = &message
				} else {
					cmp.Id = uint(id)
					//res.AddCompetitive(cmp)
				}
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}

		service.ReturnJSON(w, res)
	}).Methods("Post")
	// добавление вступительного испытания для конкурсной группы
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
	// добавление образовательной программы для конкурсной группы
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
	// проверка цифры конкурсной группы
	r.HandleFunc("/competitive/check", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultCheck
		var cmp handlers.AddCompetitiveGroup
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			// TODO заглушка. не забудь
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
	// инфа по конкусрной группе
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
	// удаление вступительного испытания у конкусрной группы
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
	// удаление конкусрной группы вместе с испытаниями и программами
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
	// удаление образовательной программы у конкурсной группы
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
	// выпадайка на список вступительных испытаний у конкурсной группы
	r.HandleFunc("/competitive/{id:[0-9]+}/tests/select", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultList
		keys := r.URL.Query()
		res.User = *handlers.CheckAuthCookie(r)
		res.MakeUrlParams(keys)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCompetitiveGroupByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.GetEntranceTestsSelectListByCompetitive(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
}
