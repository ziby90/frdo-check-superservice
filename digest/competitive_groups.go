package digest

import "time"

type CompetitiveGroup struct {
	Id                uint      `gorm:"primary_key";json:"id"` // Идентификатор
	UID               string    `xml:"UID" json:"uid" `        // Идентификатор от организации
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
	ReceptionCampaign Campaign    `gorm:"foreignkey:IdCampaign"`
	CampaignUID       CampaignUID `xml:"Campaign"` // Идентификатор направления
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
	IdAuthor          uint         `gorm:"foreignkey:id_author" json:"id_author"` // Идентификатор автора
	Actual            bool         `json:"actual"`
	Organization      Organization `gorm:"foreignkey:IdOrganization"`
	IdOrganization    uint         // Идентификатор организации
	XmlPath           XmlPath      `json:"-" xml:"-" gorm:"-"`
}

// TableNames
func (CompetitiveGroup) TableName() string {
	return "cmp.competitive_groups"
}
