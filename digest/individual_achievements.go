package digest

import (
	"errors"
	"persons/config"
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
	Changed             *time.Time          `json:"changed"`                               // Дата изменения
	IdAuthor            uint                `gorm:"foreignkey:id_author" json:"id_author"` // Идентификатор автора
	Actual              bool                `json:"actual"`                                // Актуальность
	Organization        Organization        `gorm:"foreignkey:IdOrganization"`
	IdOrganization      uint                // Идентификатор организации
}

func GetIndividualAchievements(id uint) (*IndividualAchievements, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item IndividualAchievements
	db := conn.Find(&item, id)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, db.Error
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if item.Id <= 0 {
		return nil, errors.New(`Достижение не найдено. `)
	}
	return &item, nil
}

func (i IndividualAchievements) CheckUid(uid string, primary PrimaryDataDigest) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var exist DistributedAdmissionVolume
	conn.Where(`id_organization=? AND uid=? AND actual IS TRUE`, primary.IdOrganization, uid).Find(&exist)
	if exist.Id > 0 {
		m := `Индивидуальные достижения с данным uid уже существует у выбранной организации. `
		return errors.New(m)
	}
	return nil
}
func (i IndividualAchievements) GetById(id uint) (*PrimaryDataDigest, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	db := conn.Where(`actual IS TRUE`).Find(&i, id)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Индивидуальные достижения не найдены. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if i.Id <= 0 {
		return nil, errors.New(`Индивидуальные достижения не найдены. `)
	}
	primary := PrimaryDataDigest{
		Id:             i.Id,
		Uid:            i.Uid,
		Actual:         i.Actual,
		IdOrganization: i.IdOrganization,
		Created:        i.Created,
		TableName:      i.TableName(),
	}
	return &primary, nil
}

func (IndividualAchievements) TableName() string {
	return "cmp.achievements"
}

func (i *IndividualAchievements) Init(action string) {

}

func (i *IndividualAchievements) Add() error {
	return nil
}
func (i *IndividualAchievements) Edit() error {
	return nil
}
func (i *IndividualAchievements) Remove() error {
	return nil
}
