package digest

import (
	"time"
)

// КЦП
type AdmissionVolume struct {
	Id                 uint      `gorm:"primary_key";json:"id"` // Идентификатор
	Uid                string    `xml:"UID" json:"uid" `        // Идентификатор от организации
	Direction          Direction `gorm:"foreignkey:IdDirection"`
	IdDirection        uint      // Идентификатор направления
	IdCampaign         uint      // Идентификатор направления
	ReceptionCampaign  Campaign
	CampaignUID        CampaignUID    `xml:"Campaign"` // Идентификатор направления
	EducationLevel     EducationLevel `gorm:"foreignkey:IdEducationLevel"`
	IdEducationLevel   uint
	NameEducationLevel string
	CodeEducationLevel string
	BudgetO            int64        `xml:"budget_o,omitempty"`
	BudgetOz           int64        `xml:"budget_oz,omitempty"`
	BudgetZ            int64        `xml:"budget_z,omitempty"`
	QuotaO             int64        `xml:"quota_o,omitempty"`
	QuotaOz            int64        `xml:"quota_oz,omitempty"`
	QuotaZ             int64        `xml:"quota_z,omitempty"`
	PaidO              int64        `xml:"paid_o,omitempty"`
	PaidOz             int64        `xml:"paid_oz,omitempty"`
	PaidZ              int64        `xml:"paid_z,omitempty"`
	TargetO            int64        `xml:"target_o,omitempty"`
	TargetOz           int64        `xml:"target_oz,omitempty"`
	TargetZ            int64        `xml:"target_z,omitempty"`
	Created            time.Time    `xml:"created"`
	IdAuthor           uint         `gorm:"foreignkey:id_author" json:"id_author"` // Идентификатор автора
	Actual             bool         `xml:"actual"`
	Organization       Organization `gorm:"foreignkey:IdOrganization"`
	CodeSpecialty      string       `json:"code_specialty"`
	NameSpecialty      string       `json:"name_specialty"`
	IdGroups           uint         `json:"id_groups"`
	CodeGroups         string       `json:"code_groups"`
	NameGroups         string       `json:"name_groups"`
	Distributed        bool         `json:"distributed"`
	IdOrganization     uint         // Идентификатор организации
}
type DistributedAdmissionVolume struct {
	Id                uint            `gorm:"primary_key";json:"id"` // Идентификатор
	AdmissionVolume   AdmissionVolume `gorm:"association_foreignkey:Id"`
	AdmissionVolumeId uint            `json:"id_admission_volume"  gorm:"column:id_admission_volume"`
	LevelBudget       LevelBudget     `gorm:"foreignkey:IdLevelBudget"`
	IdLevelBudget     uint            `json:"id_level_budget"`
	Uid               string          `xml:"UID" json:"uid" ` // Идентификатор от организации
	BudgetO           int64           `xml:"budget_o,omitempty"`
	BudgetOz          int64           `xml:"budget_oz,omitempty"`
	BudgetZ           int64           `xml:"budget_z,omitempty"`
	QuotaO            int64           `xml:"quota_o,omitempty"`
	QuotaOz           int64           `xml:"quota_oz,omitempty"`
	QuotaZ            int64           `xml:"quota_z,omitempty"`
	PaidO             int64           `xml:"paid_o,omitempty"`
	PaidOz            int64           `xml:"paid_oz,omitempty"`
	PaidZ             int64           `xml:"paid_z,omitempty"`
	TargetO           int64           `xml:"target_o,omitempty"`
	TargetOz          int64           `xml:"target_oz,omitempty"`
	TargetZ           int64           `xml:"target_z,omitempty"`
	Created           time.Time       `xml:"created"`
	IdAuthor          uint            `gorm:"foreignkey:id_author" json:"id_author"` // Идентификатор автора
	Actual            bool            `xml:"actual"`
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
