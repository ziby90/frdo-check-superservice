package route

import (
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/service"
)


func AddClsHandler(r *mux.Router) {

	r.HandleFunc("/api/cls/list/{cls_name}", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultCls
		vars := mux.Vars(r)
		switch vars[`cls_name`]{
			case `education_levels`:
				res.GetEducLevelResponse()
				break
			case `education_forms`:
				res.GetEducFormResponse()
				break
			default:
				message := `Неизвестный справочник.`
				res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")


}
