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

func AddAchievementsHandler(r *mux.Router) {
	// спсиок достижекний
	r.HandleFunc("/achievements/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.GetListAchievement()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// добавление достижения в заявление
	r.HandleFunc("/applications/{id:[0-9]+}/achievements/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		var cmp digest.AppAchievements
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				cmp.IdApplication = uint(id)
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &cmp)
				if err == nil {
					res.AddAppAchievement(cmp)
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
	// получаем файл у достижения заявления
	r.HandleFunc("/applications/achievements/file/{id:[0-9]+}/get", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		// TODO а если чужие спиздят? Утечка! надо замутить проверку на доступ, а как?
		res.GetFileAppAchievement(uint(id))
		if !res.Done {
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		path := fmt.Sprintf(`%v`, res.Items)
		file, err := ioutil.ReadFile(path)
		if err != nil {
			res.Done = false
			m := "Can't open file: " + path
			res.Message = &m
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		w.Write(file)
		return
	}).Methods("GET")
	// добавление достижений
	r.HandleFunc("/achievements/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp digest.IndividualAchievements
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		res.User = *handlers.CheckAuthCookie(r)
		res.Items = cmp

		err = handlers.CheckCampaignByUser(cmp.IdCampaign, res.User)
		if err != nil {
			message := err.Error()
			res.Message = &message
		} else {
			res.AddAchievement(cmp, *handlers.CheckAuthCookie(r))
		}

		service.ReturnJSON(w, &res)
	}).Methods("Post")
	// изменение достижения
	r.HandleFunc("/achievements/{id:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.AchievementMain
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			b, _ := ioutil.ReadAll(r.Body)
			err := json.Unmarshal(b, &cmp)
			res.User = *handlers.CheckAuthCookie(r)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				cmp.Id = uint(id)
				res.EditAchievement(cmp)
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("Post")
	// достижения по компании
	r.HandleFunc("/campaign/{id:[0-9]+}/achievements", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetListAchievementByCompanyId(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// достижения конкурсной группы
	r.HandleFunc("/competitive/{id:[0-9]+}/achievements/select", func(w http.ResponseWriter, r *http.Request) {
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
				res.GetAchievementsSelectListByCompetitive(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// достижения по компании, выпадайка
	r.HandleFunc("/campaign/{id:[0-9]+}/achievements/select", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultList
		keys := r.URL.Query()
		res.User = *handlers.CheckAuthCookie(r)
		res.MakeUrlParams(keys)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckCampaignByUser(uint(id), res.User)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.GetAchievementsSelectListByCampaign(uint(id))
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// информация по достижениям
	r.HandleFunc("/achievements/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{
			Done:    false,
			Message: nil,
			Items:   nil,
		}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetInfoAchievement(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// удаление достижения
	r.HandleFunc("/achievements/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{
			Done:    false,
			Message: nil,
			Items:   nil,
		}

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.RemoveAchievement(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// удаление достижения
	r.HandleFunc("/applications/achievements/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{
			Done:    false,
			Message: nil,
			Items:   nil,
		}
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.RemoveApplicationAchievement(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
}
