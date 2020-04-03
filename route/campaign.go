package route

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddCampaignHandler(r *mux.Router) {
	r.HandleFunc("/campaign/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetListCampaign()
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/campaign/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo

		var cmp handlers.CampaignMain
		b, _ := ioutil.ReadAll(r.Body)

		err := json.Unmarshal(b, &cmp)
		if err == nil {
			res.AddCampaign(cmp, *handlers.CheckAuthCookie(r))
		} else {
			m := err.Error()
			res.Message = &m
		}
		service.ReturnJSON(w, res)
	}).Methods("Post")

	r.HandleFunc("/campaign/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
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
