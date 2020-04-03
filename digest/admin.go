package digest

import "time"

type Organization struct {
	Id         uint      `xml:"id"  json:"id"`
	Ogrn       string    `xml:"ogrn" json:"ogrn"`
	Kpp        string    `xml:"kpp" json:"kpp"`
	Inn        string    `xml:"inn" json:"inn"`
	ShortTitle string    `xml:"short_title" json:"short_title"`
	FullTitle  string    `xml:"full_title" json:"full_title"`
	IDRegion   string    `xml:"id_region" json:"id_region"`
	Address    string    `xml:"address" json:"address"`
	Phone      string    `xml:"phone" json:"phone"`
	Created    time.Time `xml:"created" json:"created"`
	Actual     bool      `xml:"actual" json:"actual"`
	Filial     bool      `xml:"filial" json:"filial"`
}

// TableNames
func (Organization) TableName() string {
	return "admin.organizations"
}
