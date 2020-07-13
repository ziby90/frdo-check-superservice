package route

import (
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/service"
)

func AddOlympHandler(r *mux.Router) {
	r.HandleFunc("/olympics/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.GetListOlympics()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
}
