package digest

import (
	"errors"
	"persons/config"
	"time"
)

// Вступительные испытания
type EntranceTest struct {
	Id                    uint             `gorm:"primary_key";json:"id"` // Идентификатор
	Uid                   *string          `json:"uid" `                  // Идентификатор от организации
	TestName              *string          `json:"test_name"`
	Priority              int64            `json:"priority"`
	IsEge                 bool             `json:"is_ege"`
	MinScore              int64            `json:"min_score"`
	Subject               Subject          `gorm:"foreignkey:IdSubject"`
	IdSubject             *uint            `json:"id_subject"`
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

// Даты проведения проведения вступительных испытаний
type EntranceTestCalendar struct {
	Id               uint         `gorm:"primary_key";json:"id"` // Идентификатор
	EntranceTest     EntranceTest `gorm:"foreignkey:IdEntranceTest"`
	IdEntranceTest   uint         `json:"id_entrance_test"`
	Uid              *string      `json:"uid" `      // Идентификатор от организации
	UidEpgu          *string      `json:"uid_epgu" ` // Идентификатор от организации ЕПГУ
	ExamLocation     string       `json:"exam_location"`
	Created          time.Time    `json:"created"`
	EntranceTestDate time.Time    `json:"entrance_test_date"`
	Changed          *time.Time   `json:"changed"`
	IdAuthor         uint         `gorm:"foreignkey:id_author" json:"id_author"` // Идентификатор автора
	Actual           bool         `json:"actual"`
	Organization     Organization `gorm:"foreignkey:IdOrganization"`
	IdOrganization   uint         `json:"id_organization"`
	CountC           *int64       `json:"count_c"`
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

func (e EntranceTest) CheckUid(uid string, primary PrimaryDataDigest) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var exist EntranceTest
	conn.Where(`id_organization=? AND uid=? AND actual IS TRUE`, primary.IdOrganization, uid).Find(&exist)
	if exist.Id > 0 {
		m := `Вступительные испытания с данным uid уже существует у выбранной организации. `
		return errors.New(m)
	}
	return nil
}
func (e EntranceTest) GetById(id uint) (*PrimaryDataDigest, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	db := conn.Where(`actual IS TRUE`).Find(&e, id)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Вступительные испытания не найдены. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if e.Id <= 0 {
		return nil, errors.New(`Вступительные испытания не найдены. `)
	}
	primary := PrimaryDataDigest{
		Id:             e.Id,
		Uid:            e.Uid,
		Actual:         e.Actual,
		IdOrganization: e.IdOrganization,
		Created:        e.Created,
		TableName:      e.TableName(),
	}
	return &primary, nil
}
func (ec EntranceTestCalendar) CheckUid(uid string, primary PrimaryDataDigest) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var exist EntranceTest
	conn.Where(`id_organization=? AND uid=? AND actual IS TRUE`, primary.IdOrganization, uid).Find(&exist)
	if exist.Id > 0 {
		m := `Даты вступительных испытаний с данным uid уже существует у выбранной организации. `
		return errors.New(m)
	}
	return nil
}
func (ec EntranceTestCalendar) GetById(id uint) (*PrimaryDataDigest, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	db := conn.Where(`actual IS TRUE`).Find(&ec, id)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Даты вступительных испытаний не найдены. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if ec.Id <= 0 {
		return nil, errors.New(`Даты вступительных испытаний не найдены. `)
	}
	primary := PrimaryDataDigest{
		Id:             ec.Id,
		Uid:            ec.Uid,
		Actual:         ec.Actual,
		IdOrganization: ec.IdOrganization,
		Created:        ec.Created,
		TableName:      ec.TableName(),
	}
	return &primary, nil
}

// TableNames
func (EntranceTest) TableName() string {
	return "cmp.entrance_test"
}
func (EntranceTestCalendar) TableName() string {
	return "cmp.entrance_test_calendar"
}
