package digest

import (
	"time"
)

type Olympics struct {
	Id          uint `xml:"OlympicsID"`
	Name        string
	OrderNumber int64
	OlympYear   int64
	Created     time.Time
	IdAuthor    uint // Идентификатор автора
	Actual      bool `json:"actual"` // Актуальность
}
type OlympicsProfiles struct {
	Id               uint     `xml:"ProfileID"`
	Olympic          Olympics `gorm:"foreignkey:IdOlympic"`
	IdOlympic        uint
	OrganizaerName   string
	OrganizaerAdress string
	Profile          Profile `gorm:"foreignkey:IdProfile"`
	IdProfile        uint
	OlympicLevel     OlympicLevel `gorm:"foreignkey:IdOlympicLevel"`
	IdOlympicLevel   uint
	Organization     Organization `gorm:"foreignkey:IdOrganization"`
	IdOrganization   uint         // Идентификатор организации
	Created          time.Time
	IdAuthor         uint // Идентификатор автора
	Actual           bool `json:"actual"` // Актуальность
}

type Profile struct {
	Id       uint `xml:"ProfileID"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}

// TableNames
func (Olympics) TableName() string {
	return "olymp.olympics"
}
func (OlympicsProfiles) TableName() string {
	return "olymp.olympics_profiles"
}
func (Profile) TableName() string {
	return "olymp.profiles"
}
