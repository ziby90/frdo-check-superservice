package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"persons/handlers"
	"persons/service"
	"time"
)

func AddAnalyticsHandler(r *mux.Router) {
	r.HandleFunc("/api/analytics/dima_sokolin", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		res.GetAnalytics()
		if res.Done {
			path := fmt.Sprintf(`%v`, res.Items)
			file, err := ioutil.ReadFile(path)
			if err != nil {
				res.Done = false
				m := "Cann't open file: " + path
				res.Message = &m
			} else {
				filename := "attachment; filename=" + time.Now().Format("2006-01-02 15:04:05") + ".xlsx"

				w.Header().Set("Content-Disposition", filename)
				w.Write(file)
				return
			}
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
}
