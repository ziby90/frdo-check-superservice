package digest

import (
	"time"
)

// Приемная компания
type Campaign struct {
	Id               uint           `gorm:"primary_key" json:"id" xml:"CampaignID"` // Идентификатор
	Uid              string         `xml:"UID" json:"uid"`                          // Идентификатор от организации
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
}
type CampaignEducLevel struct {
	Id               uint `gorm:"primary_key";json:"id"`
	IdCampaign       uint
	IdEducationLevel uint
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



