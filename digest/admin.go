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
	CreatedAt  time.Time  `xml:"created_at" json:"created_at"`
	Actual     bool       `xml:"actual" json:"actual"`
	Filial     bool       `xml:"filial" json:"filial"`
	IsOOVO     bool       `json:"is_oovo"`
	Changed    *time.Time `json:"changed"`
}
type OrganizationsUsers struct {
	Id             uint `json:"id"`
	IdUser         uint `json:"id_user"`
	IdOrganization uint `json:"id_organization"`
	IdLink         uint `json:"id_link"`
}

type RequestLinks struct {
	Id                   uint                 `json:"id"`
	IdUser               uint                 `json:"id_user"`
	User                 User                 `gorm:"foreignkey:IdUser"`
	IdOrganization       uint                 `json:"id_organization"`
	Organization         Organization         `gorm:"foreignkey:IdOrganization"`
	IdStatus             uint                 `json:"id_status"`
	OrganizationStatuses OrganizationStatuses `gorm:"foreignkey:IdStatus"`
	ConfirmingDoc        *string              `json:"confirming_doc"`
	CreatedAt            time.Time            `xml:"created_at" json:"created_at"`
	UpdatedAt            *time.Time           `json:"updated_at"`
	Comment              *string              `json:"comment"`
	IdAuthor             *uint                `json:"id_author"`
}

type OrgDirections struct {
	Id          uint       `json:"id"`
	Direction   Direction  `json:"direction" gorm:"foreignkey:id_direction"`
	IdDirection uint       `json:"id_direction" schema:"id_direction"`
	IdLicense   *uint      `json:"id_license" schema:"id_license"`
	IdEiis      string     `json:"id_eiis" schema:"id_eiis"`
	Code        *string    `json:"code" schema:"code"`
	CreatedAt   time.Time  `xml:"created_at" json:"created_at"`
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
	CreatedAt          time.Time `xml:"created_at" json:"created_at"`
	Actual             bool      `json:"actual"`
}

// TableNames
func (Organization) TableName() string {
	return "admin.organizations"
}
func (OrgDirections) TableName() string {
	return "admin.org_directions"
}
func (OrganizationsUsers) TableName() string {
	return "admin.organizations_users"
}
func (RequestLinks) TableName() string {
	return "admin.request_links"
}
func (VOrganizationsDirections) TableName() string {
	return "admin.v_org_directions"
}
