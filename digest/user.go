package digest

import "time"

type User struct {
	Id                  uint                 `json:"id"`
	Login               string               `json:"login"`
	Password            string               `json:"-"`
	Patronymic          string               `json:"patronymic"`
	Surname             string               `json:"surname"`
	Name                string               `json:"name"`
	IdRole              uint                 `json:"-"`
	Role                Role                 `json:"role" gorm:"foreignkey:IdRole"`
	CurrentOrganization *CurrentOrganization `json:"organization,omitempty" gorm:"-"`
}

type CurrentOrganization struct {
	Id         uint   `json:"id"`
	ShortTitle string `json:"short_title"`
	Ogrn       string `json:"ogrn"`
	Kpp        string `json:"kpp"`
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
