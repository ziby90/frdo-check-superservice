package route

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"persons/digest"
	"persons/handlers"
	"persons/service"
	"strconv"
)

func AddEntityHandler(r *mux.Router) {
	// изменение uid
	r.HandleFunc("/{entity_name}/{id:[0-9]+}/uid/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var cmp digest.EditUid
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &cmp)
		if err != nil {
			message := err.Error()
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		res.User = *handlers.CheckAuthCookie(r)
		cmp.Id = uint(id)
		entityName := vars[`entity_name`]
		switch entityName {
		case `achievements`:
			var c digest.IndividualAchievements
			cmp.Entity = c
			break
		case `admission_volume`:
			var c digest.AdmissionVolume
			cmp.Entity = c
			break
		case `competitive`:
			var c digest.CompetitiveGroup
			cmp.Entity = c
			break
		case `campaign`:
			var c digest.Campaign
			cmp.Entity = c
			break
		case `distributed_admission_volume`:
			var c digest.DistributedAdmissionVolume
			cmp.Entity = c
			break
		case `entrance_test`:
			var c digest.EntranceTest
			cmp.Entity = c
			break
		case `entrance_test_calendar`:
			var c digest.EntranceTestCalendar
			cmp.Entity = c
			break
		case `enddate`:
			var c digest.EndApplication
			cmp.Entity = c
			break
		case `programs`:
			var c digest.CompetitiveGroupProgram
			cmp.Entity = c
			break
		default:
			message := `Неверный параметр entity_name.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		res.EditUid(cmp)
		service.ReturnJSON(w, &res)
	}).Methods("Post")
}
