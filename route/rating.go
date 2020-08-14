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
	"time"
)

func AddRatingHandler(r *mux.Router) {
	// список всех рейтинговых списков заявлений по кг
	r.HandleFunc("/rating/competitive-groups/{id:[0-9]+}/applications/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		//res.MakeUrlParamsSearch(keys, handlers.CampaignSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetListRatingCompetitiveGroupsApplication(uint(id))
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `rating.completitive_groups_applications`,
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
	// отправляем кг на синхронизацию
	r.HandleFunc("/rating/competitive-groups/{id:[0-9]+}/applications/sync", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		//res.MakeUrlParamsSearch(keys, handlers.CampaignSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.AddSyncRatingCompetitiveGroupPackage(uint(id))
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `rating.completitive_groups_applications`,
			Action:      "sync",
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
	// проверяем статус отправки на синхронизацию
	r.HandleFunc("/rating/competitive-groups/{id:[0-9]+}/applications/sync-status", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		//res.MakeUrlParamsSearch(keys, handlers.CampaignSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetSyncRatingCompetitiveGroupPackage(uint(id))
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `rating.completitive_groups_applications`,
			Action:      "sync-status",
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
	// обновляем список всех рейтинговых списков заявлений по кг
	r.HandleFunc("/rating/competitive-groups/{id:[0-9]+}/applications/refresh", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		//res.MakeUrlParamsSearch(keys, handlers.CampaignSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.RefreshListRatingCompetitiveGroupsApplication(uint(id))
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `rating.completitive_groups_applications`,
			Action:      "refresh",
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
	// редактируем заявление в рейтинговом списке
	r.HandleFunc("/rating/competitive-groups/applications/{id:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.EditRatingCompetitiveGroupsApplication
		vars := mux.Vars(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `rating.completitive_groups_applications`,
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
		idRatingApplication := uint(id)
		res.PrimaryLogging.IdObject = &idRatingApplication
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
		fmt.Println(cmp)
		cmp.Id = idRatingApplication
		res.EditRatingCompetitiveGroupsApplication(cmp)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("POST")
}
