package digest

import (
	"errors"
	"persons/config"
	"time"
)

// Приемная компания
type Campaign struct {
	Id               uint           `gorm:"primary_key" json:"id" xml:"CampaignID"` // Идентификатор
	Uid              *string        `xml:"UID" json:"uid"`                          // Идентификатор от организации
	Name             string         `json:"name";xml:"Name"`                        // Наименование
	CampaignType     CampaignType   `gorm:"foreignkey:IdCampaignType"`              // Идентификатор типа компании
	CampaignStatus   CampaignStatus `gorm:"foreignkey:IdCampaignStatus"`            // Идентификатор типа компании
	IdCampaignType   uint
	IdCampaignStatus uint
	EducationForms   EducationForms
	EducationLevels  EducationLevels
	YearStart        int64        `json:"year_start";xml:"YearStart"`            // Год начала компании
	YearEnd          int64        `json:"year_end";xml:"YearEnd"`                // Год окончания компании
	IdAuthor         uint         `gorm:"foreignkey:id_author";json:"id_author"` // Идентификатор автора
	Organization     Organization `gorm:"foreignkey:IdOrganization"`
	IdOrganization   uint         // Идентификатор организации
	Created          time.Time    `json:"created"` // Дата создания
	Changed          *time.Time   `json:"changed"` // Дата создания
}

type EndApplication struct {
	Id               uint       `gorm:"primary_key" json:"id"` // Идентификатор
	IdCampaign       uint       `json:"id_campaign"`
	IdAppAcceptPhase *uint      `json:"id_app_accept_phase"`
	IdEducationLevel uint       `json:"id_education_level"`
	IdEducationForm  uint       `json:"id_education_form"`
	EndDate          time.Time  `json:"end_date"`
	OrderEndApp      *string    `json:"order_end_app"`
	Actual           bool       `json:"actual"`
	IdOrganization   uint       // Идентификатор организации
	Created          time.Time  `json:"created"`       // Дата создания
	Changed          *time.Time `json:"changed"`       // Дата создания
	Uid              *string    `xml:"UID" json:"uid"` // Идентификатор от организации
}
type VEndApplication struct {
	Id               uint       `gorm:"primary_key" json:"id"` // Идентификатор
	Campaign         string     `json:"name_campaign"`
	IdAppAcceptPhase *uint      `json:"id_app_accept_phase"`
	AppAcceptPhase   *string    `json:"name_app_accept_phase"`
	EducationLevel   string     `json:"name_education_level"`
	IdEducationLevel uint       `json:"id_education_level"`
	EducationForm    string     `json:"name_education_form"`
	IdEducationForm  uint       `json:"id_education_form"`
	IdEndApplication *uint      `json:"id_end_application"`
	EndDate          *time.Time `json:"end_date"`
	OrderEndApp      *string    `json:"order_end_app"`
	Actual           *bool      `json:"actual"`
	IdOrganization   *uint      // Идентификатор организации
	Created          *time.Time `json:"created"`       // Дата создания
	Changed          *time.Time `json:"changed"`       // Дата создания
	Uid              *string    `xml:"UID" json:"uid"` // Идентификатор от организации
}

type EducationForms struct {
	EducationFormID []uint `xml:"EducationFormID"` // Идентификаторы формы образования
}
type EducationLevels struct {
	EducationLevelID []uint `xml:"EducationLevelID"` // Идентификаторы уровня образования
}
type CampaignEducForm struct {
	Id              uint `gorm:"primary_key";json:"id"`
	IdCampaign      uint
	IdEducationForm uint
	IdOrganization  uint
}
type CampaignEducLevel struct {
	Id               uint `gorm:"primary_key";json:"id"`
	IdCampaign       uint
	IdEducationLevel uint
	IdOrganization   uint
}

func GetCampaign(id uint) (*Campaign, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item Campaign
	db := conn.Find(&item, id)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Приемная компания не найдена. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if item.Id <= 0 {
		return nil, errors.New(`Приемная компания не найдена. `)
	}
	return &item, nil
}

// TableNames
func (Campaign) TableName() string {
	return "cmp.campaigns"
}
func (CampaignEducForm) TableName() string {
	return "cmp.campaigns_educ_form"
}
func (CampaignEducLevel) TableName() string {
	return "cmp.campaigns_educ_level"
}
func (EndApplication) TableName() string {
	return "cmp.end_application"
}
func (VEndApplication) TableName() string {
	return "cmp.v_end_application"
}

func (r *Campaign) Check(payload string) error {
	return nil
}

func (r *Campaign) Add() error {
	return nil
}
func (r *Campaign) Edit() error {
	return nil
}

func (r *Campaign) Remove() error {
	return nil
}

func (r *Campaign) GetTestJson() interface{} {
	campaign := map[string]interface{}{
		"id":                 1,
		"full_title":         `Name`,
		"name":               `Name`,
		"id_campaign_type":   1,
		"id_education_level": 1,
		"year_start":         2019,
		"year_end":           2020,
		"id_author ":         1,
		"id_organization":    1,
	}
	return campaign
}
