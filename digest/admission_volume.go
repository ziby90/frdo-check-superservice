package digest

import (
	"errors"
	"persons/config"
	"time"
)

// КЦП
type AdmissionVolume struct {
	Id                 uint      `gorm:"primary_key";json:"id"` // Идентификатор
	Uid                *string   `xml:"UID" json:"uid" `        // Идентификатор от организации
	Direction          Direction `gorm:"foreignkey:IdDirection"`
	IdDirection        uint      `json:"id_specialty"` // Идентификатор направления
	IdCampaign         uint      `json:"id_campaign"`  // Идентификатор направления
	ReceptionCampaign  Campaign
	CampaignUID        CampaignUID    `json:"Campaign"` // Идентификатор направления
	EducationLevel     EducationLevel `gorm:"foreignkey:IdEducationLevel"`
	IdEducationLevel   uint           `json:"id_education_level"`
	NameEducationLevel string         `gorm:"-"`
	CodeEducationLevel string         `gorm:"-"`
	BudgetO            int64          `json:"budget_o,omitempty"`
	BudgetOz           int64          `json:"budget_oz,omitempty"`
	BudgetZ            int64          `json:"budget_z,omitempty"`
	QuotaO             int64          `json:"quota_o,omitempty"`
	QuotaOz            int64          `json:"quota_oz,omitempty"`
	QuotaZ             int64          `json:"quota_z,omitempty"`
	PaidO              int64          `json:"paid_o,omitempty"`
	PaidOz             int64          `json:"paid_oz,omitempty"`
	PaidZ              int64          `json:"paid_z,omitempty"`
	TargetO            int64          `json:"target_o,omitempty"`
	TargetOz           int64          `json:"target_oz,omitempty"`
	TargetZ            int64          `json:"target_z,omitempty"`
	Created            time.Time      `json:"created"`
	Changed            *time.Time     `json:"changed"`
	IdAuthor           *uint          `gorm:"foreignkey:id_author" json:"id_author"` // Идентификатор автора
	Actual             bool           `json:"actual"`
	Organization       Organization   `gorm:"foreignkey:IdOrganization"`
	IdOrganization     uint           // Идентификатор организации
	CodeSpecialty      string         `json:"code_specialty" gorm:"-"`
	NameSpecialty      string         `json:"name_specialty" gorm:"-"`
	IdGroups           uint           `json:"id_groups" gorm:"-"`
	CodeGroups         string         `json:"code_groups" gorm:"-"`
	NameGroups         string         `json:"name_groups" gorm:"-"`
	HasDistributed     bool           `json:"has_distributed" gorm:"-"`
	SumDistributed     int64          `json:"sum_distributed" gorm:"-"`
	SumCompetitive     int64          `json:"sum_competitive" gorm:"-"`
}

type DistributedAdmissionVolume struct {
	Id                uint            `gorm:"primary_key";json:"id"` // Идентификатор
	AdmissionVolume   AdmissionVolume `gorm:"association_foreignkey:Id"`
	AdmissionVolumeId uint            `json:"id_admission_volume"  gorm:"column:id_admission_volume"`
	LevelBudget       LevelBudget     `gorm:"foreignkey:IdLevelBudget"`
	IdLevelBudget     uint            `json:"id_level_budget"`
	Uid               *string         `xml:"UID" json:"uid" ` // Идентификатор от организации
	BudgetO           int64           `json:"budget_o,omitempty"`
	BudgetOz          int64           `json:"budget_oz,omitempty"`
	BudgetZ           int64           `json:"budget_z,omitempty"`
	QuotaO            int64           `json:"quota_o,omitempty"`
	QuotaOz           int64           `json:"quota_oz,omitempty"`
	QuotaZ            int64           `json:"quota_z,omitempty"`
	PaidO             int64           `json:"paid_o,omitempty"`
	PaidOz            int64           `json:"paid_oz,omitempty"`
	PaidZ             int64           `json:"paid_z,omitempty"`
	TargetO           int64           `json:"target_o,omitempty"`
	TargetOz          int64           `json:"target_oz,omitempty"`
	TargetZ           int64           `json:"target_z,omitempty"`
	Created           time.Time       `json:"created"`
	Changed           *time.Time      `json:"changed"`
	IdAuthor          *uint           `gorm:"foreignkey:id_author" json:"id_author"` // Идентификатор автора
	Actual            bool            `json:"actual"`
	Organization      Organization    `gorm:"foreignkey:IdOrganization"`
	IdOrganization    uint            // Идентификатор организации
}

type CampaignUID struct {
	UID string `xml:"CampaignUID,omitempty"`
}

// TableNames
func (AdmissionVolume) TableName() string {
	return "cmp.v_admission_volume_specialty_groups"
}
func (DistributedAdmissionVolume) TableName() string {
	return "cmp.distributed_admission_volume"
}

func (c AdmissionVolume) GetById(id uint) (*PrimaryDataDigest, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	db := conn.Where(`actual IS TRUE`).Find(&c, id)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`КЦП не найдена. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if c.Id <= 0 {
		return nil, errors.New(`КЦП не найдена. `)
	}
	primary := PrimaryDataDigest{
		Id:             c.Id,
		Uid:            c.Uid,
		Actual:         c.Actual,
		IdOrganization: c.IdOrganization,
		Created:        c.Created,
		TableName:      "cmp.admission_volume",
	}
	return &primary, nil
}

func (c AdmissionVolume) CheckUid(uid string, primary PrimaryDataDigest) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var exist AdmissionVolume
	conn.Where(`id_organization=? AND uid=? AND actual IS TRUE`, primary.IdOrganization, uid).Find(&exist)
	if exist.Id > 0 {
		m := `КЦП с данным uid уже существует у выбранной организации. `
		return errors.New(m)
	}
	return nil
}

func (d DistributedAdmissionVolume) CheckUid(uid string, primary PrimaryDataDigest) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var exist DistributedAdmissionVolume
	conn.Where(`id_organization=? AND uid=? AND actual IS TRUE`, primary.IdOrganization, uid).Find(&exist)
	if exist.Id > 0 {
		m := `Уровни бюджета с данным uid уже существует у выбранной организации. `
		return errors.New(m)
	}
	return nil
}
func (d DistributedAdmissionVolume) GetById(id uint) (*PrimaryDataDigest, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	db := conn.Where(`actual IS TRUE`).Find(&d, id)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Уровни бюджета не найдены. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if d.Id <= 0 {
		return nil, errors.New(`Уровни бюджета не найдены. `)
	}
	primary := PrimaryDataDigest{
		Id:             d.Id,
		Uid:            d.Uid,
		Actual:         d.Actual,
		IdOrganization: d.IdOrganization,
		Created:        d.Created,
		TableName:      d.TableName(),
	}
	return &primary, nil
}
