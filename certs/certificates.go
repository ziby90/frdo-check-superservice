package certs

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/req"
	_ "go.nanomsg.org/mangos/v3/transport/all"
	"persons/config"
	"persons/service"
	"regexp"
	"strings"
	"time"
)

type File struct {
	Content []byte `json:"content,omitempty"`
	Title   string `json:"title,omitempty"`
	Size    int64  `json:"size,omitempty"`
	Type    string `json:"type, omitempty"`
}
type Request struct {
	Action      string `json:"action"`
	Content     string `json:"content"`
	Certificate string `json:"cert64"`
	Signature   string `json:"signature"`
	UseCertName string `json:"certname"`
}

type Result struct {
	Result     string       `json:"result"`
	ResultSign string       `json:"result_sign,omitempty"`
	ResultId   int          `json:"result_id,omitempty"`
	ResError   *ErrorCrypta `json:"error,omitempty"`
	ResCert    *CertInfo    `json:"certificate,omitempty"`
	Content    string       `json:"content,omitempty"`
}

type ErrorCrypta struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e ErrorCrypta) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Message)
}

type CertInfo struct {
	Subject            string   `json:"Subject,omitempty"`
	Issuer             string   `json:"Issuer,omitempty"`
	SerialNumber       string   `json:"SerialNumber,omitempty"`
	SHA1Hash           string   `json:"SHA1Hash,omitempty"`
	SubjKeyID          string   `json:"SubjKeyID,omitempty"`
	NotBefore          string   `json:"NotBefore,omitempty"`
	NotAfter           string   `json:"NotAfter,omitempty"`
	SignatureAlgorithm string   `json:"SignatureAlgorithm,omitempty"`
	UsageKeys          []string `json:"UsageKey,omitempty"`
	Cert64             string   `json:"Cert64,omitempty"`
}

func (c *Certificate) parseSubject() {
	sub := Subject{}
	//	fmt.Println(c.Issuer)
	subjkeys := []string{`SNILS`, `OGRN`, `STREET`, `E`, `INN`, `C`, `L`, `O`, `CN`, `SN`, `G`, `T`, `S`, `OU`}
	if len(c.Subject) > 0 {
		re, _ := regexp.Compile(`,\s([A-Z]|\d\.\d\.\d{1,6})`)
		replasedstr := string(re.ReplaceAll([]byte(c.Subject), []byte(`\t$1`)))
		reg := strings.Split(replasedstr, `\t`)
		subject := map[string]string{}
		for _, elem := range reg {
			row := strings.Split(elem, "=")
			if service.SearchStringInSliceString(row[0], subjkeys) >= 0 {
				value := strings.Replace(row[1], `, `, ``, 1)
				value = strings.Trim(value, `\t\n\r`)
				if row[0] == "G" {
					split := strings.Split(string(value), " ")
					fmt.Println(len(split))
					switch len(split) {
					case 1:
						subject[`NAME`] = split[0]
					case 2:
						subject[`NAME`] = split[0]
						subject[`PATRONYMIC`] = split[1]
					default:
						subject[`NAME`] = "Unknown"
						subject[`PATRONYMIC`] = "Unknown"
					}
				} else {
					subject[row[0]] = value
				}
			}
			mapstructure.Decode(subject, &sub)
		}
	}
	c.ParsedSubject = sub
}
func SignCheckRequest(content []byte) (operationResult string, cert *Certificate, err error) {
	Req := &Request{}
	Req.Action = `simpleVerify`
	Req.Content = base64.StdEncoding.EncodeToString(content)
	res, err := Req.sendRequest()
	if err != nil {
		return res.Result, nil, errors.New(`Сертификат не прошел проверку. ` + err.Error())
	}
	if res.Result == "success" {
		var c Certificate
		c.Subject = res.ResCert.Subject
		c.Issuer = res.ResCert.Issuer
		c.SerialNumber = res.ResCert.SerialNumber
		c.SHA1Hash = res.ResCert.SHA1Hash
		c.SubjKey = res.ResCert.SubjKeyID
		const AeneasTimeFormat = "02/01/2006 15:04:05"
		NotBefore, err := time.Parse(AeneasTimeFormat, res.ResCert.NotBefore)
		if err != nil { // Always check errors even if they should not happen.
			return res.Result, nil, errors.New(`Ошибка времени. ` + err.Error())
		}
		c.NotBefore = NotBefore
		NotAfter, err := time.Parse(AeneasTimeFormat, res.ResCert.NotAfter)
		if err != nil { // Always check errors even if they should not happen.
			return res.Result, nil, errors.New(`Ошибка времени. ` + err.Error())
		}
		c.NotAfter = NotAfter
		c.SubjKey = res.ResCert.SubjKeyID
		c.Cert64 = res.ResCert.Cert64
		c.parseSubject()
		return res.Result, &c, nil
	}
	return res.Result, nil, errors.New(`Сертификат не прошел проверку. `)

}
func (sReq *Request) sendRequest() (res Result, err error) {
	request, _ := json.Marshal(sReq)
	var sock mangos.Socket
	res = Result{}
	if sock, err = req.NewSocket(); err != nil {
		return
	}
	defer sock.Close()
	sock.SetOption("MAX-RCV-SIZE", 25*1048576)
	c := make(chan int)
	stop := time.After(3 * time.Second)
	go func() {
		if err = sock.Dial(config.Conf.Crypta.Tcp); err != nil {
			return
		}
		c <- 5
	}()
F:
	for {
		select {
		case <-c:
			break F
		case <-stop:
			err = ErrorCrypta{
				"100",
				"Can't conection to Aeneas",
			}
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
	if err = sock.Send(request); err != nil {
		return
	}
	if msg, err := sock.Recv(); err == nil {
		json.Unmarshal(msg, &res)
		// fmt.Println(string(msg))
	}
	return
}
