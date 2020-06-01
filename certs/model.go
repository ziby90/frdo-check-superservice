package certs

import (
	"persons/config"
	"persons/error_handler"
	"time"

	"github.com/jinzhu/gorm"
)

var UnknownSert = error_handler.ErrorType{Type: 11, ToUserType: 404}

type Certificate struct {
	gorm.Model
	IdOrganization uint
	Link           string
	Subject        string    `json:"subject,omitempty"`
	ParsedSubject  Subject   `gorm:"foreignkey:IDCertificate"`
	Issuer         string    `json:"issuer,omitempty"`
	SerialNumber   string    `json:"serial_number,omitempty"`
	SHA1Hash       string    `json:"sha1hash,omitempty"`
	SubjKey        string    `gorm:"column:subj_key"`
	NotBefore      time.Time `gorm:"type:timestamp"`
	NotAfter       time.Time `gorm:"type:timestamp"`
	Cert64         string
}

type Subject struct {
	gorm.Model
	IDCertificate uint
	SNILS         string `json:"snils, omitempty"`
	OGRN          string `json:"ogrn, omitempty"`
	STREET        string `json:"street, omitempty"`
	E             string `json:"mail, omitempty"`
	INN           string `json:"inn, omitempty"`
	L             string `json:"location, omitempty"`
	O             string `json:"organization, omitempty"`
	CN            string `json:"common_name, omitempty"`
	SN            string `json:"surname, omitempty"`
	NAME          string `json:"name, omitempty"`
	PATRONYMIC    string `json:"patronymic, omitempty"`
	G             string `json:"address, omitempty"`
	T             string `json:"post, omitempty"`
}

// type ResultCompare struct {
// 	service.Result
// 	r string
// }
// TableNames
func (Certificate) TableName() string {
	return "admin.certificates"
}

func (Subject) TableName() string {
	return "admin.cert_subject"
}

func Migrate() {
	config.Db.ConnGORM.AutoMigrate(&Certificate{})
	config.Db.ConnGORM.AutoMigrate(&Subject{})
	return
}

func DeleteMigrate() {
	config.Db.ConnGORM.DropTable(&Subject{}, &Certificate{})
}

func (c Certificate) SaveCert() {

	config.Db.ConnGORM.Create(&c)
	//service.Db.ConnGORM.Create(&T)

}

func (c Certificate) SelectCertWithMatchOrg(ogrn, kpp string) (IDorganization uint, err error) {
	var foundC int
	//config.Db.ConnGORM.LogMode(true)
	config.Db.ConnGORM.Where("serial_number=? AND sha1_hash=? AND subj_key=? AND id_organization=(?", c.SerialNumber, c.SHA1Hash, c.SubjKey, config.Db.ConnGORM.Table("admin.organizations").Select("id").Where("ogrn=? AND kpp=?)", ogrn, kpp).QueryExpr()).Find(&c).Count(&foundC)
	if foundC == 0 {
		return 0, UnknownSert.New("Отсутствует привязка сертификта к организации")
	}
	return c.IdOrganization, nil
}
