package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"persons/config"
	"persons/route"
	"strings"
)

var configuration = config.GetConfiguration("conf.json")

func loggingMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here

		fmt.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		//fmt.Println(`after handler`)
	})
}

func main() {
	fmt.Println(`start`)
	config.GetDbConnection()
	fmt.Println(`db`)
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./super-service-frontend/assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", neuter(fs)))
	route.GetApiHandler(r)
	r.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		if HasSuffix(r.URL.Path, []string{"js", "css", "gif", "jpeg", "woff2", "woff", "ttf", "ico"}) == false {
			http.ServeFile(w, r, "super-service-frontend/index.html")
		} else {
			http.ServeFile(w, r, "super-service-frontend/"+r.URL.Path)
		}
	})
	r.Use(loggingMw)

	fmt.Print("Server Listen...")
	fmt.Println(configuration.Port)
	err := http.ListenAndServe(":"+configuration.Port, r)
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
