package route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddUserHandler(r *mux.Router) {
	r.HandleFunc("/user/info", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		res.GetUserInfoResponse()
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/user/links", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.Items = handlers.GetOrganizationsLinks(handlers.CheckAuthCookie(r))
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/user/current-org", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			fmt.Println(`ошибка ` + err.Error())
		}
		paramId := fmt.Sprintf(`%v`, data[`id`])
		if paramId == `` {
			service.ReturnJSON(w, res)
		} else {
			u64, err := strconv.ParseUint(paramId, 10, 32)
			if err == nil {
				currentOrg := uint(u64)
				id := handlers.SetCurrentOrganization(currentOrg, handlers.CheckAuthCookie(r))
				fmt.Println(`id`, id)
				if id > 0 {
					http.SetCookie(w, &http.Cookie{
						Name:     "current-org",
						Value:    fmt.Sprintf(`%v`, id),
						HttpOnly: true,
						Path:     `/`,
					})
					res.Done = true
				} else {
					http.SetCookie(w, &http.Cookie{
						Name:     "current-org",
						Value:    ``,
						HttpOnly: true,
						Path:     `/`,
					})
				}
			}
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")

}
