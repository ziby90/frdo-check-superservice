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
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	r.HandleFunc("/api/analytics/dima_sokolin/test", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		res.GetAnalyticsTest()
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	r.HandleFunc("/api/generate/applications", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		auth := handlers.CheckAuthCookie(r)
		if auth == nil {
			res.Done = false
			m := "Can't read current_organization"
			res.Message = &m
			service.ReturnJSON(w, &res)
			return
		}
		res.User = *auth
		res.GetAnalyticsListApplications()
		fmt.Println(res.Items)
		if res.Done {
			path := fmt.Sprintf(`%v`, res.Items)
			file, err := ioutil.ReadFile(path)
			if err != nil {
				res.Done = false
				m := "Can't open file: " + path
				res.Message = &m
				service.ReturnErrorJSON(w, &res, 400)
				return
			} else {
				filename := "attachment; filename=" + time.Now().Format("2006-01-02 15:04:05") + ".xlsx"
				w.Header().Set("Content-Disposition", filename)
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Write(file)
				return
			}
		}
	}).Methods("GET")
	r.HandleFunc("/api/analytics/applications", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		auth := handlers.CheckAuthCookie(r)
		if auth == nil {
			res.Done = false
			m := "Can't read current_organization"
			res.Message = &m
			service.ReturnJSON(w, &res)
			return
		}
		res.User = *auth
		res.GetAnalyticsListApplications()
		fmt.Println(res.Items)
		if res.Done {
			path := fmt.Sprintf(`%v`, res.Items)
			file, err := ioutil.ReadFile(path)
			if err != nil {
				res.Done = false
				m := "Can't open file: " + path
				res.Message = &m
				service.ReturnErrorJSON(w, &res, 400)
				return
			} else {
				filename := "attachment; filename=" + time.Now().Format("2006-01-02 15:04:05") + ".xlsx"
				w.Header().Set("Content-Disposition", filename)
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Write(file)
				return
			}
		}
	}).Methods("GET")
	r.HandleFunc("/api/generate/competitive_groups", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		auth := handlers.CheckAuthCookie(r)
		if auth == nil {
			res.Done = false
			m := "Can't read current_organization"
			res.Message = &m
			service.ReturnJSON(w, &res)
			return
		}
		res.User = *auth
		res.GetAnalyticsListCompetitiveGroup()

		if res.Done {
			path := fmt.Sprintf(`%v`, res.Items)
			file, err := ioutil.ReadFile(path)
			if err != nil {
				res.Done = false
				m := "Can't open file: " + path
				res.Message = &m
				service.ReturnErrorJSON(w, &res, 400)
				return
			} else {
				filename := "attachment; filename=" + time.Now().Format("2006-01-02 15:04:05") + ".xlsx"
				w.Header().Set("Content-Disposition", filename)
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Write(file)
				return
			}
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	r.HandleFunc("/api/analytics/competitive_groups", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		auth := handlers.CheckAuthCookie(r)
		if auth == nil {
			res.Done = false
			m := "Can't read current_organization"
			res.Message = &m
			service.ReturnJSON(w, &res)
			return
		}
		res.User = *auth
		res.GetAnalyticsListCompetitiveGroup()

		if res.Done {
			path := fmt.Sprintf(`%v`, res.Items)
			file, err := ioutil.ReadFile(path)
			if err != nil {
				res.Done = false
				m := "Can't open file: " + path
				res.Message = &m
				service.ReturnErrorJSON(w, &res, 400)
				return
			} else {
				filename := "attachment; filename=" + time.Now().Format("2006-01-02 15:04:05") + ".xlsx"
				w.Header().Set("Content-Disposition", filename)
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Write(file)
				return
			}
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	r.HandleFunc("/api/analytics/dima_sokolin/test", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		res.GetAnalyticsTest()
		service.ReturnJSON(w, &res)
	}).Methods("GET")
}
