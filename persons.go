package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"persons/config"
	"persons/digest"
	"persons/handlers"
	"persons/route"
	"persons/service"
	"strings"
)

var configuration = config.GetConfiguration("conf.json")
var mainUser *digest.User

func authMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(`authMw`)
		fmt.Println(r.RequestURI)
		mainUser = handlers.CheckAuthCookie(r)
		if mainUser != nil {
			next.ServeHTTP(w, r)
		} else {
			m := `Для выполнения данного действия необходима авторизация`
			service.ReturnJSON(w, handlers.ResultInfo{
				Done:    false,
				Message: &m,
			})
		}
	})
}

func organizationMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(`organizationMw`)
		fmt.Println(r.RequestURI)

		if handlers.CheckOrgCookie(*handlers.CheckAuthCookie(r), r) > 0 {
			next.ServeHTTP(w, r)
		} else {
			m := `Для выполнения данного действия необходимо выбрать организацию.`
			service.ReturnJSON(w, handlers.ResultInfo{
				Done:    false,
				Message: &m,
			})
		}
	})
}

func main() {
	config.GetDbConnection()

	// общие маршруты
	routeAll := mux.NewRouter()
	route.GetApiHandlerNoAuth(routeAll)

	// маршруты с организацией выбранной и авторизацией
	routeWithAuth := routeAll.PathPrefix(`/api`).Subrouter()
	route.GetApiHandlerAuth(routeWithAuth)
	routeWithAuth.Use(authMw)
	routeWithAuth.Use(organizationMw)

	// маршруты с авторизацией
	routeWithoutOrg := routeAll.PathPrefix(`/api`).Subrouter()
	route.GetApiHandlerNoOrg(routeWithoutOrg)
	routeWithoutOrg.Use(authMw)

	fs := http.FileServer(http.Dir("./super-service-frontend/assets"))
	routeAll.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", neuter(fs)))

	routeAll.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		if HasSuffix(r.URL.Path, []string{"js", "css", "gif", "jpeg", "woff2", "woff", "ttf", "ico"}) == false {
			http.ServeFile(w, r, "super-service-frontend/index.html")
		} else {
			http.ServeFile(w, r, "super-service-frontend/"+r.URL.Path)
		}
	})

	fmt.Print("Server Listen...")
	fmt.Println(configuration.Port)

	err := http.ListenAndServe(":"+configuration.Port, routeAll)

	fmt.Println(err)
}

func HasSuffix(path string, parts []string) bool {
	for _, part := range parts {
		if strings.HasSuffix(path, part) == true {
			return true
		}
	}
	return false
}

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
