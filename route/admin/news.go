package route_admin

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"persons/digest"
	"persons/handlers"
	handlers_admin "persons/handlers/admin"
	"persons/service"
	"strconv"
)

func AddNewHandler(r *mux.Router) {

	//добавление новости
	r.HandleFunc("/new/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo

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
		var cmp digest.News
		err = decoder.Decode(&cmp, r.Form)
		if err != nil {
			message := err.Error()
			res.Message = &message
		} else {
			res.AddNews(cmp, files)
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")

	//список новостей
	r.HandleFunc("/new/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.NewResult()
		res.Sort = handlers_admin.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers_admin.NewsSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetListNews(keys)
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	//детализация новости
	r.HandleFunc("/new/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetInfoNew(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	//получить файл новости
	r.HandleFunc("/new/file/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			// TODO а если чужие спиздят? Утечка! надо замутить проверку на доступ, а как?
			res.GetFileNew(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	// добавляем файл к новости
	r.HandleFunc("/new/{id:[0-9]+}/file/add", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.ResultInfo{}
		vars := mux.Vars(r)
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
		service.ReturnJSON(w, &res)
	}).Methods("POST")

	//удалить файл новости
	r.HandleFunc("/new/file/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			// TODO а если чужие спиздят? Утечка! надо замутить проверку на доступ, а как?
			res.RemoveFileNew(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")

	//изменить новость
	r.HandleFunc("/new/{id:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers_admin.ResultInfo
		var cmp digest.News
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
		service.ReturnJSON(w, &res)
	}).Methods("Post")

	//удалить|восстановить новость
	r.HandleFunc("/new/{id:[0-9]+}/deleted", func(w http.ResponseWriter, r *http.Request) {
		res := handlers_admin.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			res.SetError(`Неверный параметр id.`)
			service.ReturnJSON(w, &res)
			return
		}
		var cmp struct {
			Deleted bool
		}
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &cmp)
		if err != nil {
			res.SetError(err.Error())
			service.ReturnJSON(w, &res)
			return
		}

		res.RemoveNew(uint(id), cmp.Deleted)
		service.ReturnJSON(w, &res)
	}).Methods("POST")

}
