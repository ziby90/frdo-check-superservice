package service

import (
	"crypto/sha256"
	_ "database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"persons/config"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

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

func ConvertInterfaceToUint(v interface{}) (uint, error) {
	u, err := strconv.ParseUint(fmt.Sprintf(`%v`, v), 10, 32)
	if err == nil {
		return uint(u), nil
	}
	return 0, err
}

func ConvertInterfaceToInt(v interface{}) (int64, error) {
	i, err := strconv.ParseInt(fmt.Sprintf(`%v`, v), 10, 64)
	if err == nil {
		return i, nil
	}
	return 0, err
}

func GetHash(str string, salt bool) string {
	hasher := sha256.New()
	s := str
	if salt {
		s += config.Conf.Salt
	}
	hasher.Write([]byte(s))
	h := hex.EncodeToString(hasher.Sum(nil))
	return h
}

func SplitNameCell(name string) (string, int) {
	char := regexp.MustCompile(`\d`).ReplaceAllString(name, ``)
	stringIndex := regexp.MustCompile(`[A-z]*`).ReplaceAllString(name, ``)
	index, err := strconv.Atoi(stringIndex)
	if err != nil {
		index = -1
	}
	return char, index
}
func GetValueFromInt64ToPoint(f *int64) interface{} {
	if f == nil {
		return nil
	} else {
		return *f
	}
}

func GetTimeFromString(value string, layout *string) interface{} {
	if layout == nil {
		s := `02-01-2006`
		layout = &s
	}
	t, err := time.Parse(*layout, fmt.Sprintf("%v", value))
	if err == nil {
		return t
	} else {
		//fmt.Println(value, err)
		return nil
	}
}
func StringInSliceString(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func SearchStringInSliceString(a string, list []string) int {
	for index, b := range list {
		if b == a {
			return index
		}
	}
	return -1
}

func NumberInSliceNumber(a int64, list []int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CutString(text string, limit int) string {
	runes := []rune(text)
	if len(runes) >= limit {
		return string(runes[:limit])
	}
	return text
}

func ReturnJSON(w http.ResponseWriter, object interface{}) {
	ansB, _ := json.Marshal(object)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ansB)
}

func ReturnXml(w http.ResponseWriter, object string) {
	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(object))
}

func ParseTimeStringToString(s string) string {
	layout := "2006-01-02T15:04:05Z07:00" // RFC3339
	tm, _ := time.Parse(layout, s)
	res := tm.Format("2006-01-02 15:04:05")
	return res
}

func GetNameStringTrim(s string) string {
	cutset := " \t\n\r"
	s = strings.ToUpper(s)
	s = strings.Replace(s, "Ё", "Е", -1)
	s = strings.Trim(s, cutset)
	return s
}

func GetBaseString(v interface{}) string {
	b, _ := json.Marshal(v)
	res := base64.StdEncoding.EncodeToString([]byte(string(b)))
	return res
}

func WaitGormDB() {
	DbStatus := config.Db.ConnGORM.DB().Ping()
	if DbStatus != nil {
	WaitDB:
		for {
			time.Sleep(time.Second)
			DbStatus = config.Db.ConnGORM.DB().Ping()
			if DbStatus != nil {
				continue
			}
			break WaitDB
		}
	}
}
