package digest

import (
	"encoding/json"
	"persons/config"
	"time"
)

type Logging interface {
	SaveLogs() error
	GetPrimaryLogging() *PrimaryLogging
	Check() bool
	SetError(string)
	SetNewData(interface{})
}

type PrimaryLogging struct {
	Id          uint      `json:"id"`
	IdObject    *uint     `json:"id_object"`
	TableObject string    `json:"table_object"`
	Action      string    `json:"action"`
	Result      bool      `json:"result"`
	OldData     *string   `json:"old_data"`
	NewData     *string   `json:"new_data"`
	Errors      *string   `json:"errors"`
	Created     time.Time `json:"created"`
	Source      string    `json:"source"`
	Route       *string   `json:"route"`
	IdAuthor    uint      `json:"id_author"` // Идентификатор автора
	TypeLogging *string   `gorm:"-"`
}

func (p *PrimaryLogging) GetTableNameLogging() string {
	if p.Action != `view` {
		return `logging.digest_crud_logs`
	} else {
		return `logging.digest_view_logs`
	}
}

func (p *PrimaryLogging) SaveLogs() error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	db := conn.Table(p.GetTableNameLogging()).Create(p)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (p *PrimaryLogging) GetPrimaryLogging() *PrimaryLogging {
	return p
}
func (p *PrimaryLogging) Check() bool {
	if p.Action != `` {
		return true
	}
	return false
}
func (p *PrimaryLogging) SetError(m string) {
	p.Result = false
	p.Errors = &m
}
func (p *PrimaryLogging) SetNewData(object interface{}) {
	newData, _ := json.Marshal(object)
	strNewData := string(newData)
	p.NewData = &strNewData

}
