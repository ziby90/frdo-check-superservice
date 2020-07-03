package digest

import (
	"errors"
	"persons/config"
	"time"
)

type CompetitiveGroup struct {
	Id                uint      `gorm:"primary_key";json:"id"` // Идентификатор
	UID               *string   `xml:"UID" json:"uid" `        // Идентификатор от организации
	Direction         Direction `gorm:"foreignkey:IdDirection"`
	IdDirection       uint      // Идентификатор направления
	Name              string
	EducationForm     EducationForm `gorm:"foreignkey:IdEducationForm"`
	IdEducationForm   uint
	EducationLevel    EducationLevel `gorm:"foreignkey:IdEducationLevel"`
	IdEducationLevel  uint
	EducationSource   EducationSource `gorm:"foreignkey:IdEducationSource"`
	IdEducationSource uint
	LevelBudget       LevelBudget `gorm:"foreignkey:IdLevelBudget"`
	IdLevelBudget     *uint
	Campaign          Campaign `gorm:"foreignkey:IdCampaign"`
	IdCampaign        uint
	BudgetO           int64        `json:"budget_o,omitempty"`
	BudgetOz          int64        `json:"budget_oz,omitempty"`
	BudgetZ           int64        `json:"budget_z,omitempty"`
	QuotaO            int64        `json:"quota_o,omitempty"`
	QuotaOz           int64        `json:"quota_oz,omitempty"`
	QuotaZ            int64        `json:"quota_z,omitempty"`
	PaidO             int64        `json:"paid_o,omitempty"`
	PaidOz            int64        `json:"paid_oz,omitempty"`
	PaidZ             int64        `json:"paid_z,omitempty"`
	TargetO           int64        `json:"target_o,omitempty"`
	TargetOz          int64        `json:"target_oz,omitempty"`
	TargetZ           int64        `json:"target_z,omitempty"`
	Created           time.Time    `json:"created"`
	Changed           *time.Time   `json:"changed"`
	Comment           *string      `json:"comment"`
	IdAuthor          uint         `gorm:"foreignkey:id_author" json:"id_author"` // Идентификатор автора
	Actual            bool         `json:"actual"`
	Organization      Organization `gorm:"foreignkey:IdOrganization"`
	IdOrganization    uint         // Идентификатор организации
	XmlPath           XmlPath      `json:"-" xml:"-" gorm:"-"`
}

type CompetitiveGroupProgram struct {
	Id                 uint      `json:"id"`
	IdCompetitiveGroup uint      `json:"id_competitive_group"`
	IdSubdivisionOrg   *uint     `json:"id_subdivision_org"`
	Name               string    `json:"name"`
	Created            time.Time `json:"created"`
	IdAuthor           uint      `json:"id_author"` // Идентификатор автора
	Actual             bool      `json:"actual"`
	Uid                *string   `json:"uid"`
	IdOrganization     uint      `json:"id_organization"`
}

// TableNames
func (CompetitiveGroup) TableName() string {
	return "cmp.competitive_groups"
}

func (CompetitiveGroupProgram) TableName() string {
	return "cmp.competitive_group_programs"
}

func (c CompetitiveGroup) CheckUid(uid string, primary PrimaryDataDigest) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var exist CompetitiveGroup
	conn.Where(`id_organization=? AND uid=? AND actual IS TRUE`, primary.IdOrganization, uid).Find(&exist)
	if exist.Id > 0 {
		m := `Конкурсная группа с данным uid уже существует у выбранной организации. `
		return errors.New(m)
	}
	return nil
}

func (c CompetitiveGroup) GetById(id uint) (*PrimaryDataDigest, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	db := conn.Where(`actual IS TRUE`).Find(&c, id)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Конкурсная группа не найдена. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if c.Id <= 0 {
		return nil, errors.New(`Конкурсная группа не найдена. `)
	}
	primary := PrimaryDataDigest{
		Id:             c.Id,
		Uid:            c.UID,
		Actual:         c.Actual,
		IdOrganization: c.IdOrganization,
		Created:        c.Created,
		TableName:      c.TableName(),
	}
	return &primary, nil
}

func (p CompetitiveGroupProgram) CheckUid(uid string, primary PrimaryDataDigest) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var exist CompetitiveGroupProgram
	conn.Where(`id_organization=? AND uid=? AND actual IS TRUE`, primary.IdOrganization, uid).Find(&exist)
	if exist.Id > 0 {
		m := `Образовательные программы с данным uid уже существует у выбранной организации. `
		return errors.New(m)
	}
	return nil
}
func (p CompetitiveGroupProgram) GetById(id uint) (*PrimaryDataDigest, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	db := conn.Where(`actual IS TRUE`).Find(&p, id)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Образовательные программы не найдены. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if p.Id <= 0 {
		return nil, errors.New(`Образовательные программы не найдены. `)
	}
	primary := PrimaryDataDigest{
		Id:             p.Id,
		Uid:            p.Uid,
		Actual:         p.Actual,
		IdOrganization: p.IdOrganization,
		Created:        p.Created,
		TableName:      p.TableName(),
	}
	return &primary, nil
}
