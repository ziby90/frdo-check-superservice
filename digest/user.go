package digest

import "time"

type User struct {
	Id                  uint                 `json:"id"`
	Login               string               `json:"login"`
	Password            string               `json:"-"`
	Patronymic          *string              `json:"patronymic"`
	Surname             string               `json:"surname"`
	Name                string               `json:"name"`
	IdRole              uint                 `json:"-"`
	IdRegion            *uint                `json:"id_region"`
	Region              Region               `json:"region" gorm:"foreignkey:IdRegion"`
	Role                Role                 `json:"role" gorm:"foreignkey:IdRole"`
	CurrentOrganization *CurrentOrganization `json:"organization,omitempty" gorm:"-"`
	RegistrationDate    time.Time            `json:"registration_date"`
	Post                *string              `json:"post"`
	Work                *string              `json:"work"`
	Adress              *string              `json:"adress"`
	Phone               *string              `json:"phone"`
	Email               *string              `json:"email"`
	IdAuthor            *uint                `json:"id_Author"`
	Actual              bool                 `json:"actual"` // Актуальность
	UpdatedAt           *time.Time           `json:"update_at"`
	Snils               *string              `json:"snils"`
}

type CurrentOrganization struct {
	Id         uint   `json:"id"`
	ShortTitle string `json:"short_title"`
	Ogrn       string `json:"ogrn"`
	Kpp        string `json:"kpp"`
	IdEiis     string `json:"id_eiis"`
	IsOOVO     bool   `json:"is_oovo"`
}

type Role struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Actual      bool      `json:"actual"` // Актуальность
}

// TableNames
func (Role) TableName() string {
	return "admin.roles"
}
func (User) TableName() string {
	return "admin.users"
}
