package handlers

import "time"

type CmpLogs struct {
	Id         uint      `json:"id"`
	IdObject   uint      `json:"id_object"`
	TypeObject string    `json:"type_object"`
	Action     string    `json:"action"`
	Result     bool      `json:"result"`
	OldData    *string   `json:"old_data"`
	NewData    *string   `json:"new_data"`
	Errors     *string   `json:"errors"`
	Created    time.Time `json:"created"`
	IdAuthor   *uint     `json:"id_author"`
}

func (CmpLogs) TableName() string {
	return `logging.cmp_logs`
}
func GetNewCmpLogs(idObject uint, typeObject string, action string) CmpLogs {
	newCmpLogs := CmpLogs{
		IdObject:   idObject,
		TypeObject: typeObject,
		Action:     action,
	}
	return newCmpLogs
}

func (*CmpLogs) NewLog() {

}
