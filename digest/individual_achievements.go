package digest

import (
	"time"
)

// Индивидуальные достижения
type IndividualAchievements struct {
	Id                  uint                `gorm:"primary_key" json:"id"`                     // Идентификатор
	Uid                 *string             `xml:"UID" json:"uid"`                             // Идентификатор от организации
	IdCampaign          uint                `gorm:"foreignkey:id_campaign" json:"id_campaign"` // Идентификатор приемной компании
	AchievementCategory AchievementCategory `gorm:"foreignkey:IdCategory" json:"category"`
	IdCategory          uint                `json:"id_category"`                           // Идентификатор наименования категории индивидуального достижения
	Name                string              `json:"name"`                                  // Наименование
	MaxValue            int64               `json:"max_value"`                             // Максимальное значение
	Created             time.Time           `json:"created"`                               // Дата создания
	IdAuthor            uint                `gorm:"foreignkey:id_author" json:"id_author"` // Идентификатор автора
	Actual              bool                `json:"actual"`                                // Актуальность
	Organization        Organization        `gorm:"foreignkey:IdOrganization"`
	IdOrganization      uint                // Идентификатор организации
}

func (IndividualAchievements) TableName() string {
	return "cmp.achievements"
}

func (r *IndividualAchievements) Init(action string) {

}

func (r *IndividualAchievements) Add() error {
	return nil
}
func (r *IndividualAchievements) Edit() error {
	return nil
}
func (r *IndividualAchievements) Remove() error {
	return nil
}
