package digest

import "time"

type EditUid struct {
	Id     uint    `xml:"id"  json:"id"`
	Uid    *string `xml:"uid"  json:"uid"`
	Entity Entity
}
type Organization struct {
	Id         uint       `xml:"id"  json:"id"`
	Ogrn       string     `xml:"ogrn" json:"ogrn"`
	Kpp        string     `xml:"kpp" json:"kpp"`
	Inn        string     `xml:"inn" json:"inn"`
	ShortTitle string     `xml:"short_title" json:"short_title"`
	FullTitle  string     `xml:"full_title" json:"full_title"`
	IDRegion   string     `xml:"id_region" json:"id_region"`
	Address    string     `xml:"address" json:"address"`
	Phone      string     `xml:"phone" json:"phone"`
	IdEiis     string     `json:"id_eiis"`
	Created    time.Time  `xml:"created" json:"created"`
	Actual     bool       `xml:"actual" json:"actual"`
	Filial     bool       `xml:"filial" json:"filial"`
	IsOOVO     bool       `json:"is_oovo"`
	Changed    *time.Time `json:"changed"`
}

type OrgDirections struct {
	Id          uint       `json:"id"`
	Direction   Direction  `json:"direction" gorm:"foreignkey:id_direction"`
	IdDirection uint       `json:"id_direction" schema:"id_direction"`
	IdLicense   *uint      `json:"id_license" schema:"id_license"`
	IdEiis      string     `json:"id_eiis" schema:"id_eiis"`
	Code        *string    `json:"code" schema:"code"`
	Created     time.Time  `json:"created"`
	Changed     *time.Time `json:"changed"`
	Actual      bool       `json:"actual"`
}

type VOrganizationsDirections struct {
	Id                 uint      `json:"id"`
	IdOrganization     uint      `json:"id_organization"`
	Code               string    `json:"code"`
	Name               string    `json:"name"`
	CodeParent         string    `json:"code_parent"`
	NameParent         string    `json:"code_parent"`
	IdParent           uint      `json:"id_parent"`
	IdEducationLevel   uint      `json:"id_education_level"`
	NameEducationLevel string    `json:"name_education_level"`
	Created            time.Time `json:"created"`
	Actual             bool      `json:"actual"`
}

// TableNames
func (Organization) TableName() string {
	return "admin.organizations"
}
func (OrgDirections) TableName() string {
	return "admin.org_directions"
}
func (VOrganizationsDirections) TableName() string {
	return "admin.v_org_directions"
}
