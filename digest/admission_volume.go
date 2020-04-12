package digest

import (
	"time"
)

// КЦП
type AdmissionVolume struct {
	Id                  uint      `gorm:"primary_key";json:"id"` // Идентификатор
	Uid                 string    `xml:"UID" json:"uid" `        // Идентификатор от организации
	Direction           Direction `gorm:"foreignkey:IdDirection"`
	IdDirection         uint      // Идентификатор направления
	IdCampaign          uint      // Идентификатор направления
	ReceptionCampaign   Campaign
	CampaignUID         CampaignUID    `xml:"Campaign"` // Идентификатор направления
	EducationLevel      EducationLevel `gorm:"foreignkey:IdEducationLevel"`
	IdEducationLevel    uint
	BudgetO             int64        `xml:"budget_o,omitempty"`
	BudgetOz            int64        `xml:"budget_oz,omitempty"`
	BudgetZ             int64        `xml:"budget_z,omitempty"`
	QuotaO              int64        `xml:"quota_o,omitempty"`
	QuotaOz             int64        `xml:"quota_oz,omitempty"`
	QuotaZ              int64        `xml:"quota_z,omitempty"`
	PaidO               int64        `xml:"paid_o,omitempty"`
	PaidOz              int64        `xml:"paid_oz,omitempty"`
	PaidZ               int64        `xml:"paid_z,omitempty"`
	TargetO             int64        `xml:"target_o,omitempty"`
	TargetOz            int64        `xml:"target_oz,omitempty"`
	TargetZ             int64        `xml:"target_z,omitempty"`
	Created             time.Time    `xml:"created"`
	IdAuthor            uint         `gorm:"foreignkey:id_author" json:"id_author"` // Идентификатор автора
	Actual              bool         `xml:"actual"`
	Organization        Organization `gorm:"foreignkey:IdOrganization"`
	IdOrganization      uint         // Идентификатор организации
	IdCampaignEducLevel uint
}

type CampaignUID struct {
	UID string `xml:"CampaignUID,omitempty"`
}

// TableNames
func (AdmissionVolume) TableName() string {
	return "cmp.admission_volume"
}
