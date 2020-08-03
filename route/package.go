package route

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"persons/digest"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddPackageHandler(r *mux.Router) {
	// добавляем файл с баллами егэ
	r.HandleFunc("/packages/mark-ege/file/add", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		res.User = *handlers.CheckAuthCookie(r)
		var f *digest.File
		r.ParseMultipartForm(0)
		file, header, fileErr := r.FormFile("file")
		if fileErr != nil && fileErr.Error() != `http: no such file` {
			res.SetErrorResult(fileErr.Error())
			service.ReturnJSON(w, &res)
			return
		}
		cmp := struct {
			Name string `json:"name" schema:"name"`
		}{}
		decoder := schema.NewDecoder()
		err := decoder.Decode(&cmp, r.Form)
		if err != nil {
			res.SetErrorResult(err.Error())
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
		res.AddFileMarkEgePackage(cmp.Name, f)
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// список пакетов с баллами егэ
	r.HandleFunc("/packages/mark-ege/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.PackageSearchArray)
		//res.MakeUrlParamsSearch(keys, handlers.ApplicationSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetMarkEgePackages()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// список элементов пакета с баллами егэ
	r.HandleFunc("/packages/mark-ege/{id:[0-9]+}/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		vars := mux.Vars(r)
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		//res.MakeUrlParamsSearch(keys, handlers.ApplicationSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		res.GetMarkEgeElements(uint(id))
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	// добавляем файл с рейтинговыми списками заявлений
	r.HandleFunc("/packages/rating-applications/file/add", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		res.User = *handlers.CheckAuthCookie(r)
		var f *digest.File
		r.ParseMultipartForm(0)
		file, header, fileErr := r.FormFile("file")
		if fileErr != nil && fileErr.Error() != `http: no such file` {
			res.SetErrorResult(fileErr.Error())
			service.ReturnJSON(w, &res)
			return
		}
		cmp := struct {
			Name string `json:"name" schema:"name"`
		}{}
		decoder := schema.NewDecoder()
		err := decoder.Decode(&cmp, r.Form)
		if err != nil {
			res.SetErrorResult(err.Error())
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
		res.AddFileRatingApplicationsPackage(cmp.Name, f)
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// список пакетов с рейтинговыми списками заявлений
	r.HandleFunc("/packages/rating-applications/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.PackageSearchArray)
		//res.MakeUrlParamsSearch(keys, handlers.ApplicationSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetRatingApplicationsPackages()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// список элементов пакета с рейтинговыми списками заявлений
	r.HandleFunc("/packages/rating-applications/{id:[0-9]+}/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		vars := mux.Vars(r)
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		//res.MakeUrlParamsSearch(keys, handlers.ApplicationSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		res.GetRatingApplicationsElements(uint(id))
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	//добавляем файл с рейтингами по конкурсу
	r.HandleFunc("/packages/rating-competitive/file/add", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		res.User = *handlers.CheckAuthCookie(r)
		var f *digest.File
		r.ParseMultipartForm(0)
		file, header, fileErr := r.FormFile("file")
		if fileErr != nil && fileErr.Error() != `http: no such file` {
			res.SetErrorResult(fileErr.Error())
			service.ReturnJSON(w, &res)
			return
		}
		cmp := struct {
			Name string `json:"name" schema:"name"`
		}{}
		decoder := schema.NewDecoder()
		err := decoder.Decode(&cmp, r.Form)
		if err != nil {
			res.SetErrorResult(err.Error())
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
		res.AddFileRatingCompetitivePackage(cmp.Name, f)
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	//список пакетов с Рейтингами по конкурсу
	r.HandleFunc("/packages/rating-competitive/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.MakeUrlParamsSearch(keys, handlers.PackageSearchArray)
		//res.MakeUrlParamsSearch(keys, handlers.ApplicationSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetRatingCompetitivePackages()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	//детализация пакета с рейтингами по конкурсу
	r.HandleFunc("/packages/rating-competitive/{id:[0-9]+}/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		vars := mux.Vars(r)
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		//res.MakeUrlParamsSearch(keys, handlers.ApplicationSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		res.GetRatingCompetitiveElements(uint(id))
		service.ReturnJSON(w, &res)
	}).Methods("GET")

}
