package route

import (
	"bytes"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"persons/certs"
	"persons/handlers"
	"persons/service"
)

func AddCertificatesHandler(r *mux.Router) {
	// смена сертификата организации
	r.HandleFunc("/organizations/certificates/change", func(w http.ResponseWriter, r *http.Request) {
		res := handlers.ResultInfo{}
		res.User = *handlers.CheckAuthCookie(r)
		err := r.ParseMultipartForm(0)
		_, header, fileErr := r.FormFile("file")
		if err != nil {
			res.SetErrorResult(err.Error())
		}
		if fileErr != nil && fileErr.Error() != `http: no such file` {
			res.SetErrorResult(fileErr.Error())
		}
		if fileErr == nil {
			f := certs.File{}
			fileContent, err := header.Open()
			if err != nil {
				res.SetErrorResult(`CheckFilesError ` + err.Error())
			} else {
				defer fileContent.Close()
			}
			buf := bytes.NewBuffer(nil)

			if _, err := io.Copy(buf, fileContent); err != nil {
				res.SetErrorResult(`CheckFilesError2 ` + err.Error())
			}
			f.Title = header.Filename
			f.Size = header.Size
			f.Content = buf.Bytes()

			res.CheckFile(f)

		} else {
			res.SetErrorResult(`fileErr ` + err.Error())
		}
		service.ReturnJSON(w, res)
	}).Methods("POST")
	// конкретная информация по сертификату организации
	r.HandleFunc("/organizations/certificates/info", func(w http.ResponseWriter, r *http.Request) {
		var res handlers.ResultInfo
		res.User = *handlers.CheckAuthCookie(r)
		res.GetOrganizationCertificate()
		service.ReturnJSON(w, res)
	}).Methods("GET")
}
