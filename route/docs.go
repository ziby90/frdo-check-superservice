package route

import (
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddDocsHandler(r *mux.Router) {
	r.HandleFunc("/docs/{table_name}/{id:[0-9]+}/info", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		tableName := vars[`table_name`]
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetInfoEDocs(uint(id), tableName)
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

}
