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
	// информация по документам на основании айди и названия таблицы
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
	// получаем файл у документа. кстати, рабоатет только на таблицу general
	r.HandleFunc("/docs/{id:[0-9]+}/file", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			// TODO а если чужие спиздят? Утечка! надо замутить проверку на доступ, а как?
			res.GetFileDoc(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("GET")
	// удаляем файл у документа. кстати, рабоатет только на таблицу general
	r.HandleFunc("/docs/{id:[0-9]+}/file/remove", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			// TODO а если чужие спиздят? Утечка! надо замутить проверку на доступ, а как?
			res.RemoveFileDoc(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")
	// добавление одного документа в профиль энтранта
	r.HandleFunc("/entrants/{id_entrant:[0-9]+}/docs/{table_name}/add", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		tableName := vars[`table_name`]
		idEntrant, err := strconv.ParseInt(vars[`id_entrant`], 10, 32)
		var f *digest.File
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}

				res.AddCompatriot(uint(idEntrant), cmp, f)
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddComposition(uint(idEntrant), cmp, f)
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddDisability(uint(idEntrant), cmp, f)
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddEge(uint(idEntrant), cmp, f)
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddEducations(uint(idEntrant), cmp, f)
				break
			case `identification`:
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddIdentifications(uint(idEntrant), cmp, f)
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddMilitaries(uint(idEntrant), cmp, f)
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddOlympicsDocs(uint(idEntrant), cmp, f)
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddOrphans(uint(idEntrant), cmp, f)
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddOther(uint(idEntrant), cmp, f)
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddParentsLost(uint(idEntrant), cmp, f)
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddRadiationWork(uint(idEntrant), cmp, f)
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
					f = &digest.File{
						MultFile: file,
						Header:   *header,
					}
				} else {
					f = nil
				}
				res.AddVeteran(uint(idEntrant), cmp, f)
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
	// добавляем файл у документа. кстати, рабоатет только на таблицу general
	r.HandleFunc("/docs/{id:[0-9]+}/file/add", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		vars := mux.Vars(r)
		res.User = *handlers.CheckAuthCookie(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		var f *digest.File
		if err == nil {
			err = r.ParseMultipartForm(0)
			file, header, fileErr := r.FormFile("file")
			if fileErr != nil && fileErr.Error() != `http: no such file` {
				res.SetErrorResult(fileErr.Error())
				service.ReturnJSON(w, res)
				return
			}

			if fileErr == nil {
				f = &digest.File{
					MultFile: file,
					Header:   *header,
				}
			} else {
				f = nil
			}
			res.AddFileDoc(uint(id), f)
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")

}
