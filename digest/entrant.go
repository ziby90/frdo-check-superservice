package digest

import "time"

type Entrants struct {
	Id                  uint      `json:"id"` // Идентификатор
	Snils               string    `json:"snils"`
	Surname             string    `json:"surname"`
	Name                string    `json:"name"`
	Patronymic          string    `json:"patronymic"`
	Gender              Gender    `json:"gender" gorm:"foreignkey:IdGender"`
	IdGender            uint      `json:"id_gender"`
	Birthday            time.Time `json:"birthday"`
	Birthplace          string    `json:"birthplace"`
	Phone               string    `json:"phone"`
	Email               string    `json:"email"`
	RegistrationAddress string    `json:"registration_address"`
	FactAddress         string    `json:"fact_address"`
	Okcm                Okcm      `json:"okcm" gorm:"foreignkey:IdOkcm"`
	IdOkcm              uint      `json:"id_okcm"`
	Created             time.Time `json:"created"` // Дата создания
}

func (Entrants) TableName() string {
	return "persons.entrants"
}
