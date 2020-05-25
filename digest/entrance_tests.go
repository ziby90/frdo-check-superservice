package digest

import (
	"errors"
	"persons/config"
	"time"
)

// Конкурсные группы
type EntranceTest struct {
	Id                    uint             `gorm:"primary_key";json:"id"` // Идентификатор
	Uid                   *string          `json:"uid" `                  // Идентификатор от организации
	TestName              string           `json:"test_name"`
	Priority              int64            `json:"priority"`
	IsEge                 bool             `json:"is_ege"`
	MinScore              int64            `json:"min_score"`
	Subject               Subject          `gorm:"foreignkey:IdSubject"`
	IdSubject             uint             `json:"id_subject"`
	EntranceTestType      EntranceTestType `gorm:"foreignkey:IdEntranceTestType"`
	IdEntranceTestType    uint             `json:"id_entrance_test_type"`
	CompetitiveGroup      CompetitiveGroup `gorm:"foreignkey:IdCompetitiveGroup"`
	IdCompetitiveGroup    uint             `json:"id_competitive_group"`
	IdReplaceEntranceTest *uint            `json:"id_replace_entrance_test"`
	Created               time.Time        `json:"created"`
	IdAuthor              uint             `gorm:"foreignkey:id_author" json:"id_author"` // Идентификатор автора
	Actual                bool             `json:"actual"`
	Organization          Organization     `gorm:"foreignkey:IdOrganization"`
	IdOrganization        uint             `json:"id_organization"`
}

func GetEntranceTest(id uint) (*EntranceTest, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item EntranceTest
	db := conn.Preload(`CompetitiveGroup`).Find(&item, id)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Вступительный тест не найден. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if item.Id <= 0 {
		return nil, errors.New(`Вступительный тест не найден. `)
	}
	return &item, nil
}

// TableNames
func (EntranceTest) TableName() string {
	return "cmp.entrance_test"
}
