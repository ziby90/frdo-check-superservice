package route

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"persons/digest"
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

	r.HandleFunc("/docs/{id:[0-9]+}/file", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetFileDoc(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")

	r.HandleFunc("/entrants/{id_entrant:[0-9]+}/docs/{table_name}/add", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		tableName := vars[`table_name`]
		idEntrant, err := strconv.ParseInt(vars[`id_entrant`], 10, 32)
		if err == nil {
			err = r.ParseMultipartForm(0)
			decoder := schema.NewDecoder()
			file, header, fileErr := r.FormFile("file")
			switch tableName {
			case `compatriot`:
				cmp := digest.Compatriot{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddCompatriot(uint(idEntrant), cmp)
				break
			case `composition`:
				cmp := digest.Composition{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddComposition(uint(idEntrant), cmp)
				break
			case `disability`:
				cmp := digest.Disability{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddDisability(uint(idEntrant), cmp)
				break
			case `ege`:
				cmp := digest.Ege{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddEge(uint(idEntrant), cmp)
				break
			case `educations`:
				cmp := digest.Educations{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddEducations(uint(idEntrant), cmp)
				break
			case `identifications`:
				cmp := digest.Identifications{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddIdentifications(uint(idEntrant), cmp)
				break
			case `militaries`:
				cmp := digest.Militaries{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddMilitaries(uint(idEntrant), cmp)
				break
			case `olympics`:
				cmp := digest.OlympicsDocs{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddOlympicsDocs(uint(idEntrant), cmp)
				break
			case `orphans`:
				cmp := digest.Orphans{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddOrphans(uint(idEntrant), cmp)
				break
			case `other`:
				cmp := digest.Other{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddOther(uint(idEntrant), cmp)
				break
			case `parents_lost`:
				cmp := digest.ParentsLost{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddParentsLost(uint(idEntrant), cmp)
				break
			case `radiation_work`:
				cmp := digest.RadiationWork{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddRadiationWork(uint(idEntrant), cmp)
				break
			case `veteran`:
				cmp := digest.Veteran{}
				err := decoder.Decode(&cmp, r.Form)
				if err != nil {
					res.SetErrorResult(err.Error())
					break
				}
				if fileErr != nil && fileErr.Error() != `http: no such file` {
					res.SetErrorResult(fileErr.Error())
					break
				}
				if fileErr == nil {
					cmp.File = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				}
				res.AddVeteran(uint(idEntrant), cmp)
				break
			default:
				message := `Неверный параметр table_name.`
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")

}
