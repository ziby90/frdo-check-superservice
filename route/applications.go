package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"persons/handlers"
	"persons/service"
	"strconv"
	"strings"
	"time"
)

func AddApplicationHandler(r *mux.Router) {
	// список заявлений
	r.HandleFunc("/applications/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		res.Sort = handlers.Sort{
			Field: "",
			Order: "",
		}
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		//res.MakeUrlParamsSearch(keys, handlers.ApplicationSearchArray)
		res.User = *handlers.CheckAuthCookie(r)
		res.GetApplications(keys)
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// добавление заявления
	r.HandleFunc("/applications/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		var cmp handlers.AddApplication
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &cmp)
		if err == nil {
			res.AddApplication(cmp)
		} else {
			message := err.Error()
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// заявления конкретного абитуриента
	r.HandleFunc("/entrants/{id:[0-9]+}/applications", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		keys := r.URL.Query()
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.GetApplicationsByEntrant(uint(id), keys)
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// меняем статус заявления GET
	r.HandleFunc("/applications/{id:[0-9]+}/status/set/{code_status}", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		var cmp handlers.ChangeStatusApplication
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				cmp.IdApplication = uint(id)
				codeStatus := fmt.Sprintf(`%v`, vars[`code_status`])
				cmp.CodeStatus = codeStatus
				res.SetStatusApplication(cmp)
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// меняем статус заявления POST
	r.HandleFunc("/applications/{id:[0-9]+}/status/set", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		var cmp handlers.ChangeStatusApplication
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				cmp.IdApplication = uint(id)
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &cmp)
				if err == nil {
					res.SetStatusApplication(cmp)
				} else {
					message := err.Error()
					res.Message = &message
				}
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// меняем статус заявления POST
	r.HandleFunc("/applications/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		var cmp struct {
			StatusComment string `json:"status_comment"`
		}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &cmp)
				if err == nil {
					res.RemoveApplication(uint(id), cmp.StatusComment)
				} else {
					message := err.Error()
					res.Message = &message
				}
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("POST")
	// основная информация по заявлению
	r.HandleFunc("/applications/{id:[0-9]+}/main", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				res.GetApplicationById(uint(id))
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// конкретная информация по заявлению
	r.HandleFunc("/applications/{id:[0-9]+}/info", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				res.GetApplicationInfoById(uint(id))
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// редактирование конкретной информации по заявлению
	r.HandleFunc("/applications/{id:[0-9]+}/info/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var data handlers.EditApplicationInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &data)
				data.IdApplication = uint(id)
				if err == nil {
					res.EditApplicationInfoById(data)
					//res.Ed(data)
				} else {
					m := err.Error()
					res.Message = &m
				}
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("Post")
	// редактирование конкретной информации по заявлению
	r.HandleFunc("/applications/{id:[0-9]+}/test/{id_test:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var data handlers.EditApplicationTest
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &data)
				data.IdApplication = uint(id)
				if err == nil {
					idTest, err := strconv.ParseInt(vars[`id_test`], 10, 32)
					if err == nil {
						data.IdApplication = uint(id)
						data.IdEntranceTest = uint(idTest)
						res.EditApplicationTestById(data)
					} else {
						message := `Неверный параметр id_test.`
						res.Message = &message
					}
				} else {
					m := err.Error()
					res.Message = &m
				}
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("Post")
	// документы заявления
	r.HandleFunc("/applications/{id:[0-9]+}/docs/list", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				res.GetApplicationDocsById(uint(id))
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// документы заявления кратакая запись чисто айдишников
	r.HandleFunc("/applications/{id:[0-9]+}/docs/short", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				res.GetApplicationDocsByIdShort(uint(id))
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// вступительные тесты заявления
	r.HandleFunc("/applications/{id:[0-9]+}/tests/list", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				res.GetApplicationEntranceTestsById(uint(id))
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// список достижений к заявлению
	r.HandleFunc("/applications/{id:[0-9]+}/achievements/list", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.NewResult()
		keys := r.URL.Query()
		res.MakeUrlParams(keys)
		res.User = *handlers.CheckAuthCookie(r)
		res.MakeUrlParamsSearch(keys, handlers.AppAchievementsSearchArray)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				res.GetAchievementsByApplicationId(uint(id))
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// добавление вступительного теста к заявлению
	r.HandleFunc("/applications/{id:[0-9]+}/tests/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var data handlers.AddApplicationEntranceTest
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &data)
				data.IdApplication = uint(id)
				if err == nil {
					res.AddEntranceTestApplication(data)
				} else {
					m := err.Error()
					res.Message = &m
				}
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("Post")
	// выбор даты ви
	r.HandleFunc("/applications/{id:[0-9]+}/calendar/choose", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var data handlers.ChooseCalendar
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		res.User = *handlers.CheckAuthCookie(r)
		err = handlers.CheckApplicationByUser(uint(id), res.User)
		if err != nil {
			m := err.Error()
			res.Message = &m
			service.ReturnJSON(w, &res)
			return
		}
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &data)
		data.IdApplication = uint(id)
		if err != nil {
			m := err.Error()
			res.Message = &m
			service.ReturnJSON(w, &res)
			return
		}
		res.ChooseCalendarEntranceTest(data)
		service.ReturnJSON(w, &res)
	}).Methods("Post")
	// удаление даты ви
	r.HandleFunc("/applications/{id:[0-9]+}/calendar/{id_entrance_test_agreed:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			message := `Неверный параметр id.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		res.User = *handlers.CheckAuthCookie(r)
		err = handlers.CheckApplicationByUser(uint(id), res.User)
		if err != nil {
			m := err.Error()
			res.Message = &m
			service.ReturnJSON(w, &res)
			return
		}
		idCalendar, err := strconv.ParseInt(vars[`id_entrance_test_agreed`], 10, 32)
		if err != nil {
			message := `Неверный параметр id_entrance_test_agreed.`
			res.Message = &message
			service.ReturnJSON(w, &res)
			return
		}
		res.RemoveCalendarEntranceTest(uint(id), uint(idCalendar))
		service.ReturnJSON(w, &res)
	}).Methods("Post")
	// добавление документов к заявлению
	r.HandleFunc("/applications/{id:[0-9]+}/docs/add", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var data handlers.AddApplicationDocs
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			err = handlers.CheckApplicationByUser(uint(id), res.User)
			if err == nil {
				b, _ := ioutil.ReadAll(r.Body)
				err := json.Unmarshal(b, &data)
				data.IdApplication = uint(id)
				if err == nil {
					res.AddDocsApplication(data)
				} else {
					m := err.Error()
					res.Message = &m
				}
			} else {
				message := err.Error()
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("Post")
	// генерация файла пдф согласия/отзыва согласия
	r.HandleFunc("/applications/{id:[0-9]+}/generate/agreed/pdf", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var data handlers.PDFApplicationParams
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			res.SetErrorResult(`Неверный параметр id.`)
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		err = handlers.CheckApplicationByUser(uint(id), res.User)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		data.IdApplication = uint(id)
		//data.Docs = append(data.Docs, handlers.DocsApplication{
		//	Id:   156,
		//	Type: "identification",
		//})
		//data.Docs = append(data.Docs, handlers.DocsApplication{
		//	Id:   1306,
		//	Type: "education",
		//})
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &data)

		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		resultData := res.GeneratePDFApplicationAgreedData(data)
		if !res.Done {
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		tmpl, err := template.ParseFiles(*resultData.TmplPath)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, resultData)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		pdfg, err := wkhtml.NewPDFGenerator()
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		pdfg.AddPage(wkhtml.NewPageReader(strings.NewReader(tpl.String())))
		err = pdfg.Create()
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		filename := `attachment; filename="` + time.Now().Format(`2006-01-02 15:04:05`) + `.pdf"`
		w.Header().Set("Content-Disposition", filename)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

		w.Write(pdfg.Bytes())
	}).Methods("Post")
	// генерация файла пдф заявления
	r.HandleFunc("/applications/{id:[0-9]+}/generate/pdf", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		var data handlers.PDFApplicationParams
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			res.SetErrorResult(`Неверный параметр id.`)
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		err = handlers.CheckApplicationByUser(uint(id), res.User)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		data.IdApplication = uint(id)
		//data.Docs = append(data.Docs, handlers.DocsApplication{
		//	Id:   156,
		//	Type: "identification",
		//})
		//data.Docs = append(data.Docs, handlers.DocsApplication{
		//	Id:   1306,
		//	Type: "education",
		//})
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &data)

		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		resultData := res.GeneratePDFApplicationData(data)
		if !res.Done {
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		tmpl, err := template.ParseFiles(*resultData.TmplPath)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, resultData)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		//pdfg, err := wkhtml.NewPDFGenerator()
		//if err != nil {
		//	res.SetErrorResult(err.Error())
		//	service.ReturnErrorJSON(w, &res, 400)
		//	return
		//}
		//pdfg.AddPage(wkhtml.NewPageReader(strings.NewReader(tpl.String())))
		//err = pdfg.Create()
		//if err != nil {
		//	res.SetErrorResult(err.Error())
		//	service.ReturnErrorJSON(w, &res, 400)
		//	return
		//}
		//filename := `attachment; filename="` + time.Now().Format(`2006-01-02 15:04:05`) + `.pdf"`
		//w.Header().Set("Content-Disposition", filename)
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		//w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		//
		//w.Write(pdfg.Bytes())
	}).Methods("GET")
	// получение информации о согласии - несогласии
	r.HandleFunc("/applications/{id:[0-9]+}/info/agreed", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err != nil {
			res.SetErrorResult(`Неверный параметр id.`)
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		err = handlers.CheckApplicationByUser(uint(id), res.User)
		if err != nil {
			res.SetErrorResult(err.Error())
			service.ReturnErrorJSON(w, &res, 400)
			return
		}
		res.GetAgreedApplicationsById(uint(id))
		service.ReturnJSON(w, &res)
	}).Methods("Get")
	//// генерация файла пдф заявления
	//r.HandleFunc("/applications/{id:[0-9]+}/generate/pdf", func(w http.ResponseWriter, r *http.Request) {
	//	var res handlers.ResultInfo
	//	var data handlers.PDFApplicationParams
	//	res.GeneratePDFApplication(data)
	//}).Methods("GET")
	// удаление вступительного теста к заявлению
	r.HandleFunc("/applications/tests/{id:[0-9]+}/remove", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{
			Done:    false,
			Message: nil,
			Items:   nil,
		}
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars[`id`], 10, 32)
		if err == nil {
			res.RemoveApplicationTest(uint(id))
		} else {
			message := `Неверный параметр id.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")
	// удаление документа к заявлению
	r.HandleFunc("/applications/{id_application:[0-9]+}/docs/{id_document:[0-9]+}/{code_category}/remove", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{
			Done:    false,
			Message: nil,
			Items:   nil,
		}
		res.User = *handlers.CheckAuthCookie(r)
		vars := mux.Vars(r)
		idApplication, err := strconv.ParseInt(vars[`id_application`], 10, 32)
		if err == nil {
			idDocument, err := strconv.ParseInt(vars[`id_document`], 10, 32)
			if err == nil {
				codeCategory := fmt.Sprintf(`%v`, vars[`code_category`])
				res.RemoveApplicationDocuments(uint(idApplication), uint(idDocument), codeCategory)
			} else {
				message := `Неверный параметр id_document.`
				res.Message = &message
			}
		} else {
			message := `Неверный параметр id_application.`
			res.Message = &message
		}
		service.ReturnJSON(w, &res)
	}).Methods("GET")

	//
	//r.HandleFunc("/entrants/{id:[0-9]+}/others", func(w http.ResponseWriter, r *http.Request) {
	//	res := handlers.ResultInfo{}
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err == nil {
	//		res.GetDocsEntrant(uint(id))
	//	} else {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//	}
	//	service.ReturnJSON(w, &res)
	//}).Methods("GET")
	//
	//r.HandleFunc("/entrants/{id:[0-9]+}/idents", func(w http.ResponseWriter, r *http.Request) {
	//	res := handlers.ResultInfo{}
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err == nil {
	//		res.GetDocsIdentsEntrant(uint(id))
	//	} else {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//	}
	//	service.ReturnJSON(w, &res)
	//}).Methods("GET")
	//
	//r.HandleFunc("/entrants/{id:[0-9]+}/idents/list", func(w http.ResponseWriter, r *http.Request) {
	//	res := handlers.ResultInfo{}
	//
	//	vars := mux.Vars(r)
	//	id, err := strconv.ParseInt(vars[`id`], 10, 32)
	//	if err == nil {
	//		res.GetListDocsIdentsEntrant(uint(id))
	//	} else {
	//		message := `Неверный параметр id.`
	//		res.Message = &message
	//	}
	//	service.ReturnJSON(w, &res)
	//}).Methods("GET")

}
