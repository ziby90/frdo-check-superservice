package route

import (
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddNewHandler(r *mux.Router) {

	// получаем список новостей
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
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	// детализация новости
	r.HandleFunc("/api/new/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
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

	// получаем файл новости
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
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	// страдаем херней
	r.HandleFunc("/api/new/stradat/hernya", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		res.Items = map[string]interface{}{
			`stradat`: true,
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
}
