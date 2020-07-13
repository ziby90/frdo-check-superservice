package route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/service"
	"strings"
)

func AddAuthHandler(r *mux.Router) {
	// авторизован ли?
	r.HandleFunc("/api/is-auth", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.IsAuth(handlers.CheckAuthCookie(r))
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// а роль? роль верная?
	r.HandleFunc("/api/is-role", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		data := make(map[string]string)
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			fmt.Println(`ошибка ` + err.Error())
		}

		if data[`code`] == `` {
			m := `Неверные параметры`
			res.Message = &m
			service.ReturnJSON(w, &res)
		}

		user := handlers.CheckAuthCookie(r)
		if strings.ToUpper(data[`code`]) == strings.ToUpper(user.Role.Code) {
			res.Done = true
		} else {
			res.Done = false
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// логин, ну что тут сказать
	r.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]string)
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			fmt.Println(`ошибка ` + err.Error())
		}

		if data[`login`] == `` || data[`password`] == `` {
			var res handlers.ResultInfo
			m := `Неверные параметры`
			res.Message = &m
			service.ReturnJSON(w, &res)
		}
		userAgent := r.Header.Get(`User-agent`)
		pass := data[`password`]
		mail := strings.ToUpper(data[`login`])
		res := handlers.CheckAuthBase(mail, pass)
		if res.Done {
			http.SetCookie(w, &http.Cookie{
				Name:     "login",
				Value:    res.User.Login,
				HttpOnly: true,
				Path:     `/`,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     "password",
				Value:    service.GetHash(res.User.Password+userAgent, true),
				HttpOnly: true,
				Path:     `/`,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     "current-org",
				Value:    ``,
				HttpOnly: true,
				Path:     `/`,
			})
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     "login",
				Value:    ``,
				HttpOnly: true,
				Path:     `/`,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     "password",
				Value:    ``,
				HttpOnly: true,
				Path:     `/`,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     "current-org",
				Value:    ``,
				HttpOnly: true,
				Path:     `/`,
			})
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// уже уходите? нам тоже не понравилось. сами вы лагаете.
	r.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "login",
			Value:    ``,
			HttpOnly: true,
			Path:     `/`,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "password",
			Value:    ``,
			HttpOnly: true,
			Path:     `/`,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "current-org",
			Value:    ``,
			HttpOnly: true,
			Path:     `/`,
		})
		service.ReturnJSON(w, &handlers.ResultAuth{
			User:    nil,
			Done:    true,
			Message: "Разлогинился",
		})
	}).Methods("GET")
}
