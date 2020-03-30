package service

import (
	"persons/config"

	"github.com/jinzhu/gorm"
)

type Logging struct {
	gorm.Model
	IdOrganization int64
	Ipaddress      string `gorm:"column:ipaddress";json:"ipaddress,omitempty"`
	Machine        string `json:"machine"`
	IdAction       int64  `json:"id_action"`
	IdObject       int64  `json:"id_object"`
	Result         string `json:"result"`
	Additional     string `json:"additional"`
	Done           bool   `json:"done"`
}

func (*Logging) TableName() string {
	return `logging.actions`
}

func Migrate() {
	config.Db.ConnGORM.AutoMigrate(&Logging{})
}
func Drop() {
	config.Db.ConnGORM.DropTable(&Logging{})
}
