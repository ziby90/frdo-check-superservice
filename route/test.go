package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"persons/config"
	"persons/handlers"
	"persons/service"
)

func AddTestHandler(r *mux.Router) {
	r.HandleFunc("/test2/smev_connect", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		conn := config.Db.ConnSmevGorm
		conn.LogMode(config.Conf.Dblog)
		var db *gorm.DB
		var rows handlers.RowsCls
		db = conn.Select(`id, request_data as name`).Table(`data.input_query`).Scan(&rows)
		if db.Error != nil {
			fmt.Println(db.Error)
		}
		res.Items = rows

		service.ReturnJSON(w, res)
	}).Methods("GET")
}
