package route

import (
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddOrgsHandler(r *mux.Router) {
	r.HandleFunc("/organizations/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.GetListOrganization()
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/organizations/{id:[0-9]+}/info", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{
			Done:    false,
			Message: nil,
			Items:   nil,
		}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetInfoCampaign(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
}
