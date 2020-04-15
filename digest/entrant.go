package digest

import "time"

type Entrants struct {
	Id                 uint      `json:"id"` // Идентификатор
	Snils              string    `json:"snils"`
	Surname            string    `json:"surname"`
	Name               string    `json:"name"`
	Patronymic         string    `json:"patronymic"`
	Gender             Gender    `json:"gender" gorm:"foreignkey:IdGender"`
	IdGender           uint      `json:"id_gender"`
	Birthday           time.Time `json:"birthday"`
	Birthplace         string    `json:"birthplace"`
	Phone              string    `json:"phone"`
	Email              string    `json:"email"`
	IdRegistrationAddr uint      `json:"id_registration_address"`
	RegistrationAddr   Address   `json:"registration_addr" gorm:"foreignkey:IdRegistrationAddr"`
	IdFactAddr         uint      `json:"id_fact_address"`
	FactAddr           Address   `json:"fact_addr" gorm:"foreignkey:IdFactAddr"`
	Okcm               Okcm      `json:"okcm" gorm:"foreignkey:IdOkcm"`
	IdOkcm             uint      `json:"id_okcm"`
	Created            time.Time `json:"created"` // Дата создания
}

type Address struct {
	Id               uint   `json:"id"`
	FullAddr         string `json:"full_addr"`
	IndexAddr        string `json:"index_addr"`
	IdRegion         uint   `json:"id_region"`
	Region           Region `json:"region" gorm:"foreignkey:IdRegion"`
	Area             string `json:"area"`
	City             string `json:"city"`
	CityArea         string `json:"city_area"`
	Place            string `json:"place"`
	Street           string `json:"street"`
	AdditionalArea   string `json:"additional_area"`
	AdditionalStreet string `json:"additional_street"`
	House            string `json:"house"`
	Building1        string `json:"building1"`
	Building2        string `json:"building2"`
	Apartment        string `json:"apartment"`
	IdAuthor         uint   `json:"id_author" gorm:"column:id_author"` // Идентификатор автора
}

func (Entrants) TableName() string {
	return "persons.entrants"
}
func (Address) TableName() string {
	return "persons.addresses"
}
