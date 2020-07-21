package service

import (
	"crypto/sha256"
	_ "database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"persons/config"
	"persons/digest"
	"regexp"
	"strconv"
)

var ErrorDbNotFound = `record not found`

type Result struct {
	Done    bool    `json:"done,omitempty"`
	Message *string `json:"message,omitempty"`
	Code    int64   `json:"code,omitempty"`
}

type File struct {
	Content []byte `json:"content,omitempty"`
	Title   string `json:"title,omitempty"`
	Size    int64  `json:"size,omitempty"`
	Type    string `json:"type, omitempty"`
}

func GetParamKeyFilterUintArray(key []string) []uint {
	var filter []uint
	for _, value := range key {
		if id, err := strconv.ParseUint(value, 10, 32); err == nil {
			filter = append(filter, uint(id))
		}
	}
	return filter
}

//
//func ConvertInterfaceToUint(v interface{}) (uint, error) {
//	u, err := strconv.ParseUint(fmt.Sprintf(`%v`, v), 10, 32)
//	if err == nil {
//		return uint(u), nil
//	}
//	return 0, err
//}
//
//func ConvertInterfaceToInt(v interface{}) (int64, error) {
//	i, err := strconv.ParseInt(fmt.Sprintf(`%v`, v), 10, 64)
//	if err == nil {
//		return i, nil
//	}
//	return 0, err
//}

func GetHash(str string, salt bool) string {
	hash := sha256.New()
	s := str
	if salt {
		s += config.Conf.Salt
	}
	hash.Write([]byte(s))
	h := hex.EncodeToString(hash.Sum(nil))
	return h
}

//
//func SplitNameCell(name string) (string, int) {
//	char := regexp.MustCompile(`\d`).ReplaceAllString(name, ``)
//	stringIndex := regexp.MustCompile(`[A-z]*`).ReplaceAllString(name, ``)
//	index, err := strconv.Atoi(stringIndex)
//	if err != nil {
//		index = -1
//	}
//	return char, index
//}
//func GetValueFromInt64ToPoint(f *int64) interface{} {
//	if f == nil {
//		return nil
//	} else {
//		return *f
//	}
//}
//
//func GetTimeFromString(value string, layout *string) interface{} {
//	if layout == nil {
//		s := `02-01-2006`
//		layout = &s
//	}
//	t, err := time.Parse(*layout, fmt.Sprintf("%v", value))
//	if err == nil {
//		return t
//	} else {
//		//fmt.Println(value, err)
//		return nil
//	}
//}
//func StringInSliceString(a string, list []string) bool {
//	for _, b := range list {
//		if b == a {
//			return true
//		}
//	}
//	return false
//}

func CheckSnils(snils string) error {
	if len(snils) != 11 {
		return errors.New(`Неверное число символов в строке снилс. `)
	}
	matched, _ := regexp.Match(`\D`, []byte(snils))
	if matched {
		return errors.New(`Недопустимый символ в строке снилс. `)
	}
	control := snils[len(snils)-2:]
	if controlNumber, ok := strconv.Atoi(control); ok == nil {
		s := snils[:len(snils)-2]
		result := 0
		for i := 0; i < len(s); i++ {
			num, err := strconv.Atoi(string(s[i]))
			if err != nil {
				return err
			}
			result += (len(s) - i) * num
		}
		if result == 100 || result == 101 {
			result = 0
		}
		if result > 101 {
			result %= 101
		}
		fmt.Println(`result`, result)
		if result == controlNumber {
			return nil
		} else {
			return errors.New(`Снилс некорректен `)
		}
	}
	return errors.New(`Ошибка проверки снилс. `)
}

func SearchStringInSliceString(a string, list []string) int {
	for index, b := range list {
		if b == a {
			return index
		}
	}
	return -1
}
func SearchUintInSliceUint(a uint, list []uint) int {
	for index, b := range list {
		if b == a {
			return index
		}
	}
	return -1
}

//
//func NumberInSliceNumber(a int64, list []int64) bool {
//	for _, b := range list {
//		if b == a {
//			return true
//		}
//	}
//	return false
//}
//
//func CutString(text string, limit int) string {
//	runes := []rune(text)
//	if len(runes) >= limit {
//		return string(runes[:limit])
//	}
//	return text
//}

func ReturnJSON(w http.ResponseWriter, object digest.Logging) {
	ansB, _ := json.Marshal(object)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(ansB)
	if err != nil {
		fmt.Println(`ошибка ` + err.Error())
	}
	if object.Check() {
		err = object.SaveLogs()
	}

}
func ReturnErrorJSON(w http.ResponseWriter, object digest.Logging, statusCode int) {
	ansB, _ := json.Marshal(object)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(ansB)
	if err != nil {
		fmt.Println(`ошибка ` + err.Error())
	}
	if object.Check() {
		err = object.SaveLogs()
	}
}

//
//func ReturnXml(w http.ResponseWriter, object string) {
//	w.Header().Set("Content-Type", "text/xml")
//	_, err := w.Write([]byte(object))
//	if err!=nil{
//		fmt.Println(`ошибка `+err.Error())
//	}
//}

//func ParseTimeStringToString(s string) string {
//	layout := "2006-01-02T15:04:05Z07:00" // RFC3339
//	tm, _ := time.Parse(layout, s)
//	res := tm.Format("2006-01-02 15:04:05")
//	return res
//}
//
//func GetNameStringTrim(s string) string {
//	cutSet := " \t\n\r"
//	s = strings.ToUpper(s)
//	s = strings.Replace(s, "Ё", "Е", -1)
//	s = strings.Trim(s, cutSet)
//	return s
//}
//
//func GetBaseString(v interface{}) string {
//	b, _ := json.Marshal(v)
//	res := base64.StdEncoding.EncodeToString([]byte(string(b)))
//	return res
//}

//func WaitGormDB() {
//	DbStatus := config.Db.ConnGORM.DB().Ping()
//	if DbStatus != nil {
//	WaitDB:
//		for {
//			time.Sleep(time.Second)
//			DbStatus = config.Db.ConnGORM.DB().Ping()
//			if DbStatus != nil {
//				continue
//			}
//			break WaitDB
//		}
//	}
//}
