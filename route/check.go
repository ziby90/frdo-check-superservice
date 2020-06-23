package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddChecksHandler(r *mux.Router) {
	// можно ли редактировать ПК
	r.HandleFunc("/campaign/{id:[0-9]+}/check/edit-remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			can := false
			err = handlers.CheckEditCampaign(uint(id))
			if err == nil {
				can = true
			}
			res.Items = map[string]interface{}{
				`can`: can,
			}
			res.Done = true
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// можно ли редактировать и обнулять контролльные даты
	r.HandleFunc("/campaign/{id:[0-9]+}/enddate/check/edit-remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			can := false
			err = handlers.CheckEditEndDate(uint(id))
			if err == nil {
				can = true
			}
			res.Items = map[string]interface{}{
				`can`: can,
			}
			res.Done = true
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// можно ли добавлять кцп
	r.HandleFunc("/campaign/{id:[0-9]+}/admission/check/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			can := false
			err = handlers.CheckAddAdmission(uint(id))
			if err == nil {
				can = true
			}
			res.Items = map[string]interface{}{
				`can`: can,
			}
			res.Done = true
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// можно ли редактировать кцп
	r.HandleFunc("/campaign/{id:[0-9]+}/admission/check/edit-remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			can := false
			err = handlers.CheckEditAdmission(uint(id))
			if err == nil {
				can = true
			}
			res.Items = map[string]interface{}{
				`can`: can,
			}
			res.Done = true
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// можно ли добавлять кцп
	r.HandleFunc("/campaign/{id:[0-9]+}/competitive/check/add-remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			can := false
			err = handlers.CheckAddCompetitive(uint(id))
			if err == nil {
				can = true
			}
			res.Items = map[string]interface{}{
				`can`: can,
			}
			res.Done = true
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// можно ли удалять кцп
	//r.HandleFunc("/competitive/{id:[0-9]+}/check/remove", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers.ResultInfo
	//	vars := mux.Vars(r)
	//	res.User = *handlers.CheckAuthCookie(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err == nil {
	//		can := false
	//		err = handlers.CheckRemoveCompetitive(uint(id))
	//		if err == nil {
	//			can = true
	//		}
	//		res.Items = map[string]interface{}{
	//			`can`: can,
	//		}
	//		res.Done = true
	//	} else {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//	}
	//	service.ReturnJSON(w, res)
	//}).Methods("GET")
	// можно ли редактировать кг
	r.HandleFunc("/competitive/{id:[0-9]+}/check/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			can := false
			err = handlers.CheckEditCompetitiveGroup(uint(id))
			if err == nil {
				can = true
			}
			res.Items = map[string]interface{}{
				`can`: can,
			}
			res.Done = true
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// можно ли редактировать образовательные программы у кг
	r.HandleFunc("/competitive/{id:[0-9]+}/programs/check/add-remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			can := false
			err = handlers.CheckEditProgramsCompetitiveGroup(uint(id))
			if err == nil {
				can = true
			}
			res.Items = map[string]interface{}{
				`can`: can,
			}
			res.Done = true
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// можно ли редактировать втсупительные испытания у кг
	r.HandleFunc("/competitive/{id:[0-9]+}/tests/check/add-remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			can := false
			err = handlers.CheckEditEntranceCompetitiveGroup(uint(id))
			if err == nil {
				can = true
			}
			res.Items = map[string]interface{}{
				`can`: can,
			}
			res.Done = true
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// можно ли добавлять достижение
	r.HandleFunc("/campaign/{id:[0-9]+}/achievements/check/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			can := false
			err = handlers.CheckAddAchievements(uint(id))
			if err == nil {
				can = true
			}
			res.Items = map[string]interface{}{
				`can`: can,
			}
			res.Done = true
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// можно ли редактировать достижение
	r.HandleFunc("/campaign/{id:[0-9]+}/achievements/check/edit-remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			can := false
			err = handlers.CheckEditAchievements(uint(id))
			if err == nil {
				can = true
			}
			res.Items = map[string]interface{}{
				`can`: can,
			}
			res.Done = true
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// цифра в кг
	r.HandleFunc("/competitive/{id:[0-9]+}/number/{number:[0-9]+}/check", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			number, err := strconv.ParseInt(vars[`number`], 10, 32)
			if err == nil {
				can := false
				err = handlers.CheckNumberCompetitiveById(uint(id), number)
				if err == nil {
					can = true
				}
				res.Items = map[string]interface{}{
					`can`:   can,
					`error`: err,
				}
				res.Done = true
				fmt.Println(fmt.Sprintf(`can : %v`, can))
				fmt.Println(fmt.Sprintf(`done : %v`, res.Done))
			} else {
				message := `Неверный параметр number.`
				res.Message = &message
			}

		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
}
