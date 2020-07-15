package route_admin

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"io/ioutil"
	"net/http"
	"persons/digest"
	"persons/handlers"
	handlers_admin "persons/handlers/admin"
	"persons/service"
	"strconv"
	"time"
)

func AddUserHandler(r *mux.Router) {
	// список пользователей
	r.HandleFunc("/users/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.NewResult()
		res.Sort = handlers_admin.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers_admin.UserSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetListUsers()
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.users`,
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
	// информация по конкретному пользователю
	r.HandleFunc("/users/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.ResultInfo{}
		vars := mux.Vars(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.users`,
			Action:      "view",
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
		res.GetInfoUser(idUser)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	// список связей кокнретного пользователя
	r.HandleFunc("/users/{id:[0-9]+}/links", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.ResultInfo{}
		vars := mux.Vars(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.users`,
			Action:      "view",
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
		res.GetLinksUser(idUser)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// добавление связи с организацией пользователю
	r.HandleFunc("/users/{id:[0-9]+}/links/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo
		var cmp struct {
			IdOrganization uint    `json:"id_organization" schema:"id_organization"`
			Comment        *string `json:"comment" schema:"comment"`
		}
		vars := mux.Vars(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.organizations_users`,
			Action:      "add",
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
		err = r.ParseMultipartForm(0)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnJSON(w, &res)
			return
		}
		idUser := uint(id)
		decoder := schema.NewDecoder()
		err = decoder.Decode(&cmp, r.Form)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnJSON(w, &res)
			return
		}
		file, header, fileErr := r.FormFile("file")
		if fileErr != nil {
			res.SetErrorResult(fileErr.Error())
			service.ReturnJSON(w, &res)
			return
		}
		f := &digest.File{
			MultFile: file,
			Header:   *header,
		}
		res.AddLinksToUser(idUser, cmp.IdOrganization, f, cmp.Comment)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// получаем файл связи у пользователя
	r.HandleFunc("/users/links/{id:[0-9]+}/file", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.ResultInfo{}
		vars := mux.Vars(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.organizations_users`,
			Action:      "view",
			Created:     time.Now(),
			Source:      "cabinet",
			Route:       &r.URL.Path,
			IdAuthor:    res.User.Id,
		}
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			res.SetErrorResult(`Неверный параметр id.`)
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		idLink := uint(id)
		res.PrimaryLogging.IdObject = &idLink

		res.PrimaryLogging.IdAuthor = res.User.Id
		res.GetFileLink(idLink)
		if !res.Done {
			res.PrimaryLogging.Result = res.Done
			res.PrimaryLogging.Errors = res.Message
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		path := fmt.Sprintf(`%v`, res.Items)
		file, err := ioutil.ReadFile(path)
		if err != nil {
			res.Done = false
			res.SetErrorResult("Can't open file: " + path)
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		w.Write(file)
		return
	}).Methods("GET")

	// меняем пароль пользователю
	r.HandleFunc("/users/{id:[0-9]+}/password/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo
		var cmp struct {
			Password string `json:"password"`
		}
		vars := mux.Vars(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.users`,
			Action:      "password-reset",
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
		res.ResetPasswordUser(idUser, cmp.Password)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// добавление связи с организацией пользователю
	r.HandleFunc("/users/create", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo
		var cmp handlers_admin.AddUser
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.users`,
			Action:      "create",
			Created:     time.Now(),
			Source:      "cabinet",
			Route:       &r.URL.Path,
			IdAuthor:    res.User.Id,
		}
		res.User = *handlers.CheckAuthCookie(r)
		res.PrimaryLogging.IdAuthor = res.User.Id
		err := r.ParseMultipartForm(0)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnJSON(w, &res)
			return
		}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&cmp, r.Form)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnJSON(w, &res)
			return
		}
		file, header, fileErr := r.FormFile("file")
		if fileErr != nil && fileErr.Error() != "http: no such file" {
			res.SetErrorResult(fileErr.Error())
			service.ReturnJSON(w, &res)
			return
		}
		var f *digest.File
		if fileErr == nil {
			f = &digest.File{
				MultFile: file,
				Header:   *header,
			}
		}
		res.CreateUser(cmp, f)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// блокировка пользователя
	r.HandleFunc("/users/{id:[0-9]+}/block", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo
		vars := mux.Vars(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.users`,
			Action:      "block",
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
		res.BlockUser(idUser)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// блокировка пользователя
	r.HandleFunc("/users/{id:[0-9]+}/unblock", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo
		vars := mux.Vars(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.users`,
			Action:      "unblock",
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
		res.UnblockUser(idUser)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("POST")

}
