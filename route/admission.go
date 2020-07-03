package route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"persons/digest"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddAdmissionHandler(r *mux.Router) {
	//два запроса на вывод кцп. оставили пока  оба, Аня просила
	r.HandleFunc("/campaign/{id:[0-9]+}/admission", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.AdmissionVolumeSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.GetListAdmissionVolume(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	r.HandleFunc("/campaign/{id:[0-9]+}/admission2", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.AdmissionVolumeSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.GetListAdmissionVolumeBySpec(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// список уровней бюджета кцп
	r.HandleFunc("/admission/{id:[0-9]+}/budget", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckAdmissionVolumeByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.GetLevelBudgetAdmission(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// добавление бюджета к кцп
	r.HandleFunc("/admission/{id:[0-9]+}/budget/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			var cmp handlers.AddBudget
			b, _ := ioutil.ReadAll(r.Body)
			err := json.Unmarshal(b, &cmp)
			if err == nil {
				err = handlers.CheckAdmissionVolumeByUser(uint(id), res.User)
				if err != nil {
					message := err.Error()
					res.Message = &message
				} else {
					res.AddAdmissionBudget(uint(id), cmp)
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
	}).Methods("POST")
	// Добавление кцп с нулевыми значениями
	r.HandleFunc("/admission/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.AddAdmission
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			err = handlers.CheckCampaignByUser(cmp.IdCampaign, res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.AddAdmission(cmp)
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("Post")
	// редактирование цифр кцп
	r.HandleFunc("/admission/{id:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp digest.AdmissionVolume
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckAdmissionVolumeByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &cmp)
				if err == nil {
					fmt.Print(cmp.QuotaZ)
					res.EditAdmission(uint(id), cmp)
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
	}).Methods("Post")

	// удаление кцп с уровнями бюджета
	r.HandleFunc("/admission/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckAdmissionVolumeByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.RemoveAdmission(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// удаление конкретного уровня бюджета у кцп
	r.HandleFunc("/admission/{id:[0-9]+}/budget/{id_budget:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			idBudget, err := strconv.ParseInt(vars[`id_budget`], 10, 32)
			if err == nil {
				err = handlers.CheckAdmissionVolumeByUser(uint(id), res.User)
				if err != nil {
					message := err.Error()
					res.Message = &message
				} else {
					res.RemoveBudgetAdmission(uint(id), uint(idBudget))
				}
			} else {
				message := `Неверный параметр id_budget.`
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// редактирование конкретного уровня бюджета у кцп
	r.HandleFunc("/admission/{id:[0-9]+}/budget/{id_budget:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.EditDistributedAdmissionVolume
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			idBudget, err := strconv.ParseInt(vars[`id_budget`], 10, 32)
			if err == nil {
				err = handlers.CheckAdmissionVolumeByUser(uint(id), res.User)
				if err != nil {
					message := err.Error()
					res.Message = &message
				} else {
					b, _ := ioutil.ReadAll(r.Body)
					err := json.Unmarshal(b, &cmp)
					if err == nil {
						cmp.IdAdmissionVolume = uint(id)
						cmp.IdLevelBudget = uint(idBudget)
						res.EditAdmissionLevelBudget(cmp)
					} else {
						m := err.Error()
						res.Message = &m
					}
				}
			} else {
				message := `Неверный параметр id_budget.`
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")
	// информация по кцп
	r.HandleFunc("/admission/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.Result{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckAdmissionVolumeByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.GetAdmissionVolumeById(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
}
