package route

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddNewHandler(r *mux.Router) {
	//
	r.HandleFunc("/api/new/add", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		auth := handlers.CheckAuthCookie(r)
		if auth == nil {
			message := "Необходима авторизация"
			res.Message = &message
		} else {
			res.User = *handlers.CheckAuthCookie(r)
			var files []*multipart.FileHeader
			err := r.ParseMultipartForm(0)
			if r.MultipartForm != nil {
				value, ok := r.MultipartForm.File["files"]
				if ok {
					for _, file := range value {
						files = append(files, file)
					}
				}
			}
			decoder := schema.NewDecoder()
			var cmp handlers.News
			err = decoder.Decode(&cmp, r.Form)
			if err != nil {
				message := err.Error()
				res.Message = &message
			} else {
				res.AddNews(cmp, files)
			}
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")

	//
	r.HandleFunc("/api/new/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.NewsSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetListNews()
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// изменение конкурсной группы
	r.HandleFunc("/api/new/{id:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.News
		auth := handlers.CheckAuthCookie(r)
		if auth == nil {
			message := "Необходима авторизация"
			res.Message = &message
		} else {
			res.User = *handlers.CheckAuthCookie(r)
			vars := mux.Vars(r)
			id, err := strconv.ParseInt(vars[`id`], 10, 32)
			if err == nil {
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &cmp)
				if err == nil {
					cmp.Id = uint(id)
					res.EditNew(cmp)
				} else {
					message := err.Error()
					res.Message = &message
				}
			} else {
				message := `Неверный параметр id.`
				res.Message = &message
			}
		}
		service.ReturnJSON(w, res)
	}).Methods("Post")
	// получаем файл у документа. кстати, рабоатет только на таблицу general
	r.HandleFunc("/api/new/file/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			// TODO а если чужие спиздят? Утечка! надо замутить проверку на доступ, а как?
			res.GetFileNew(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// добавляем файл у документа. кстати, рабоатет только на таблицу general
	r.HandleFunc("/api/new/{id:[0-9]+}/file/add", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		auth := handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		if auth == nil {
			message := "Необходима авторизация"
			res.Message = &message
		} else {
			id, err := strconv.ParseInt(vars[`id`], 10, 32)
			if err == nil {
				res.User = *handlers.CheckAuthCookie(r)
				var files []*multipart.FileHeader
				err := r.ParseMultipartForm(0)
				if r.MultipartForm != nil {
					value, ok := r.MultipartForm.File["files"]
					if ok {
						for _, file := range value {
							files = append(files, file)
						}
					}
				}
				if err != nil {
					message := err.Error()
					res.Message = &message
				} else {
					res.AddFileNew(uint(id), files)
				}
			} else {
				message := `Неверный параметр id.`
				res.Message = &message
			}
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")
	// удаляем файл у документа. кстати, рабоатет только на таблицу general
	r.HandleFunc("/api/new/file/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			// TODO а если чужие спиздят? Утечка! надо замутить проверку на доступ, а как?
			res.RemoveFileNew(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")

}
