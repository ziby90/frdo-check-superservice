package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetApiHandler(r *mux.Router) {
	AddCampaignHandler(r)
	AddClsHandler(r)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(12334)
	}).Methods("GET")
}
