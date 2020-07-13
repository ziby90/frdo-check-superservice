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

func AddEntrantHandler(r *mux.Router) {
	r.HandleFunc("/entrants/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.EntrantsSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetListEntrants()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	r.HandleFunc("/entrants/{id:[0-9]+}/short", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetInfoEntrantApp(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	r.HandleFunc("/entrants/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var data handlers.AddEntrantData
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &data)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.AddEntrant(data)
		} else {
			m := err.Error()
			res.Message = &m
		}
		service.ReturnJSON(w, &res)
	}).Methods("Post")
	r.HandleFunc("/entrants/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return

		}
		res.GetInfoEntrant(uint(id))
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	r.HandleFunc("/entrants/{id:[0-9]+}/others", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.GetDocsEntrant(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	r.HandleFunc("/entrants/{id:[0-9]+}/idents", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.GetDocsIdentsEntrant(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	r.HandleFunc("/entrants/{id:[0-9]+}/idents/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.GetListDocsIdentsEntrant(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	r.HandleFunc("/entrants/{id:[0-9]+}/docs/short", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		keys := r.URL.Query()
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		res.User = *handlers.CheckAuthCookie(r)
		if err == nil {
			res.GetShortListDocsEntrant(uint(id), keys)
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	// получаем файл фотки абитуриента.
	r.HandleFunc("/entrants/photo/{id:[0-9]+}/file", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		// TODO а если чужие спиздят? Утечка! надо замутить проверку на доступ, а как?
		res.GetFilePhotoPersons(uint(id))
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
	// добавляем  файл фотки абитуриента.
	r.HandleFunc("/entrants/{id:[0-9]+}/photo/file/add", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		var f *digest.File
		if err == nil {
			err = r.ParseMultipartForm(0)
			file, header, fileErr := r.FormFile("file")
			if fileErr != nil && fileErr.Error() != `http: no such file` {
				res.SetErrorResult(fileErr.Error())
				service.ReturnJSON(w, &res)
				return
			}

			if fileErr == nil {
				f = &digest.File{
					MultFile: file,
					Header:   *header,
				}
			} else {
				f = nil
			}
			res.AddFilePhotoPersons(uint(id), f)
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// изменяем  файл фотки абитуриента.
	r.HandleFunc("/entrants/photo/{id:[0-9]+}/file/edit", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		var f *digest.File
		if err == nil {
			err = r.ParseMultipartForm(0)
			file, header, fileErr := r.FormFile("file")
			if fileErr != nil && fileErr.Error() != `http: no such file` {
				res.SetErrorResult(fileErr.Error())
				service.ReturnJSON(w, &res)
				return
			}

			if fileErr == nil {
				f = &digest.File{
					MultFile: file,
					Header:   *header,
				}
			} else {
				f = nil
			}
			res.EditFilePhotoPersons(uint(id), f)
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// удаляем  файл фотки абитуриента.
	r.HandleFunc("/entrants/photo/{id:[0-9]+}/file/remove", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			// TODO а если чужие спиздят? Утечка! надо замутить проверку на доступ, а как?
			res.RemoveFilePhotoPersons(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")
}
