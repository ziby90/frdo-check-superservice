package digest

import (
	"io/ioutil"
	"persons/config"
	"persons/error_handler"
	"time"
)

// Индивидуальные достижения
type IndividualAchievements struct {
	Id                  uint                `gorm:"primary_key";json:"id"`                     // Идентификатор
	Uid                 string              `xml:"UID" json:"uid"`                             // Идентификатор от организации
	IdCampaign          uint                `gorm:"foreignkey:id_campaign";json:"id_campaign"` // Идентификатор приемной компании
	AchievementCategory AchievementCategory `gorm:"foreignkey:IdCategory"`
	IdCategory          uint                // Идентификатор наименования категории индивидуального достижения
	Name                string              `json:"name"`                                  // Наименование
	MaxValue            int64               `json:"max_value"`                             // Максимальное значение
	Created             time.Time           `json:"created"`                               // Дата создания
	IdAuthor            uint                `gorm:"foreignkey:id_author";json:"id_author"` // Идентификатор автора
	Actual              bool                `json:"actual"`                                // Актуальность
	XmlPath             XmlPath             `json:"-" xml:"-" gorm:"-"`
}

func (IndividualAchievements) TableName() string {
	return "cmp.achievements"
}

func (r *IndividualAchievements) Init(action string) {
	switch action {
	case `add`:
		r.XmlPath.TrueResultAction = `Элемент добавлен.`
		r.XmlPath.PathXml = config.Conf.RootDir + "/schema/example/xml/add/individual_achievements.xml"
		r.XmlPath.PathXsd = config.Conf.RootDir + "/schema/xsd/add/individual_achievements_schema.xsd"
	case `edit`:
		r.XmlPath.TrueResultAction = `Элемент изменен.`
		r.XmlPath.PathXml = config.Conf.RootDir + "/schema/example/xml/add/individual_achievements.xml"
		r.XmlPath.PathXsd = config.Conf.RootDir + "/schema/xsd/add/individual_achievements_schema.xsd"
		break
	case `remove`:
		r.XmlPath.TrueResultAction = `Элемент удален`
		r.XmlPath.PathXml = config.Conf.RootDir + "/schema/example/xml/remove/individual_achievements.xml"
		r.XmlPath.PathXsd = config.Conf.RootDir + "/schema/xsd/remove/individual_achievements_schema.xsd"
		break
	}
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

func (r *IndividualAchievements) GetTestXml() (string, error) {
	content, err := ioutil.ReadFile(r.XmlPath.PathXml)
	if err != nil {
		return ``, error_handler.ErrorType.New(validXmlErr, err.Error())
	}
	return string(content), nil
}
