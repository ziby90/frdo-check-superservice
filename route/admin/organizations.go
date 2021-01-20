package route_admin

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"persons/digest"
	"persons/handlers"
	handlers_admin "persons/handlers/admin"
	"persons/service"
	"strconv"
	"time"
)

func AddOrgsAddminHandler(r *mux.Router) {
	//список организаций
	r.HandleFunc("/organizations/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers_admin.OrganizationSearchArray)
		res.GetListOrganization()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	//список заявок на связь с организацией
	r.HandleFunc("/organizations/requests/links/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		//res.MakeUrlParamsSearch(keys, handlers_admin.OrganizationSearchArray)
		res.GetListRequestsLinks()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// получаем файл заявки
	r.HandleFunc("/organizations/requests/links/{id:[0-9]+}/file", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		res.GetFileRequestLink(uint(id))
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
	// принятие заявки на связь от пользователя
	r.HandleFunc("/organizations/requests/links/{id:[0-9]+}/accept", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo
		vars := mux.Vars(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.request_links`,
			Action:      "accept_request_link",
			Created:     time.Now(),
			Source:      "cabinet",
			Route:       &r.URL.Path,
			IdAuthor:    res.User.Id,
		}
		res.User = *handlers.CheckAuthCookie(r)
		res.PrimaryLogging.IdAuthor = res.User.Id
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			res.SetErrorResult(`Неверный параметр id.`)
			service.ReturnJSON(w, &res)
			return
		}
		idUser := uint(id)
		res.PrimaryLogging.IdObject = &idUser
		res.AcceptRequestLink(idUser)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// отклонение заявки на связь от пользователя
	r.HandleFunc("/organizations/requests/links/{id:[0-9]+}/decline", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo
		var cmp struct {
			Comment string `json:"comment"`
		}

		vars := mux.Vars(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.request_links`,
			Action:      "decline_request_link",
			Created:     time.Now(),
			Source:      "cabinet",
			Route:       &r.URL.Path,
			IdAuthor:    res.User.Id,
		}
		res.User = *handlers.CheckAuthCookie(r)
		res.PrimaryLogging.IdAuthor = res.User.Id
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			res.SetErrorResult(`Неверный параметр id.`)
			service.ReturnJSON(w, &res)
			return
		}
		idUser := uint(id)
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &cmp)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnJSON(w, &res)
			return
		}
		res.PrimaryLogging.IdObject = &idUser
		res.DeclineRequestLink(idUser, cmp.Comment)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	//список организаций для выпадайки
	//r.HandleFunc("/organizations/short", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers_admin.ResultInfo
	//	res.GetListOrganizationShort()
	//	service.ReturnJSON(w, &res)
	//}).Methods("GET")
}
