package route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"persons/digest"
	"persons/handlers"
	"persons/service"
	"strings"
	"time"
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
		var res handlers.ResultInfo
		data := make(map[string]string)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.users`,
			Action:      "login",
			Created:     time.Now(),
			Source:      "cabinet",
			Route:       &r.URL.Path,
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			fmt.Println(`ошибка ` + err.Error())
		}

		if data[`login`] == `` || data[`password`] == `` {
			res.SetErrorResult(`Неверные параметры`)
			service.ReturnJSON(w, &res)
		}
		userAgent := r.Header.Get(`User-agent`)
		pass := data[`password`]
		mail := strings.ToUpper(data[`login`])
		resNext := handlers.CheckAuthBase(mail, pass)
		resNext.PrimaryLogging = res.PrimaryLogging
		if resNext.Done {
			resNext.PrimaryLogging.Result = true
			resNext.PrimaryLogging.IdObject = &resNext.User.Id
			resNext.PrimaryLogging.IdAuthor = resNext.User.Id
			http.SetCookie(w, &http.Cookie{
				Name:     "login",
				Value:    resNext.User.Login,
				HttpOnly: true,
				Path:     `/`,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     "password",
				Value:    service.GetHash(resNext.User.Password+userAgent, true),
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
			resNext.PrimaryLogging.Result = false
			resNext.SetNewData(data)
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
		service.ReturnJSON(w, &resNext)
	}).Methods("POST")
	// уже уходите? нам тоже не понравилось. сами вы лагаете.
	r.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultAuth
		res.User = handlers.CheckAuthCookie(r)
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.users`,
			Action:      "logout",
			Created:     time.Now(),
			Source:      "cabinet",
			Route:       &r.URL.Path,
		}
		if res.User != nil {
			res.PrimaryLogging.IdObject = &res.User.Id
			res.PrimaryLogging.IdAuthor = res.User.Id
		}

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
		res.User = nil
		res.Done = true
		res.PrimaryLogging.Result = true
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// регистрация пользователя
	r.HandleFunc("/api/registration", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp handlers.AddUser
		res.PrimaryLogging = digest.PrimaryLogging{
			TableObject: `admin.users`,
			Action:      "registration",
			Created:     time.Now(),
			Source:      "cabinet",
			Route:       &r.URL.Path,
			IdAuthor:    res.User.Id,
		}
		err := r.ParseMultipartForm(0)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnJSON(w, &res)
			return
		}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&cmp, r.Form)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnJSON(w, &res)
			return
		}
		file, header, fileErr := r.FormFile("file")
		if fileErr != nil && fileErr.Error() != "http: no such file" {
			res.SetErrorResult(fileErr.Error())
			service.ReturnJSON(w, &res)
			return
		}
		var f *digest.File
		if fileErr == nil {
			f = &digest.File{
				MultFile: file,
				Header:   *header,
			}
		}
		res.RegistrationUser(cmp, f)
		res.PrimaryLogging.Result = res.Done
		res.PrimaryLogging.Errors = res.Message
		service.ReturnJSON(w, &res)
	}).Methods("POST")
}
