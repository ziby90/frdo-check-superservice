package digest

import "time"

type Application struct {
	Id                       uint                `json:"id" schema:"id"` // Идентификатор
	IdOrganization           uint                `json:"id_organization"`
	Entrants                 Entrants            `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId               uint                `json:"id_entrant" gorm:"column:id_entrant" schema:"id_entrant"`
	CompetitiveGroup         CompetitiveGroup    `gorm:"foreignkey:IdCompetitiveGroup"`
	IdCompetitiveGroup       uint                `json:"id_competitive_group"`
	IdCompetitiveGroupTarget uint                `json:"id_competitive_group_target"`
	AppNumber                string              `json:"app_number"`
	RegistrationDate         time.Time           `json:"registration_date" schema:"registration_date"`
	Rating                   float32             `json:"rating" schema:"rating"`
	Status                   ApplicationStatuses `gorm:"foreignkey:id_status"`
	IdStatus                 uint                `json:"id_status" schema:"id_status"`
	Priority                 int64               `json:"priority" schema:"priority"`
	FirstHigherEducation     bool                `json:"first_higher_education" schema:"first_higher_education"`
	NeedHostel               bool                `json:"need_hostel" schema:"need_hostel"`
	IdDisabledDocument       uint                `json:"id_disabled_document" schema:"id_disabled_document"`
	SpecialConditions        bool                `json:"special_conditions" schema:"special_conditions"`
	DistanceTest             bool                `json:"distance_test" schema:"distance_test"`
	DistancePlace            *string             `json:"distance_place" schema:"distance_place"`
	ViolationTypes           ViolationTypes      `json:"violation_type" gorm:"foreignkey:id_violation"`
	IdViolation              uint                `json:"id_violation" schema:"id_violation"`
	EgeCheck                 *time.Time          `json:"ege_check" schema:"ege_check"`
	Agreed                   *bool               `json:"agreed" schema:"agreed"`
	Disagreed                *bool               `json:"disagreed" schema:"disagreed"`
	AgreedDate               *time.Time          `json:"agreed_date" schema:"agreed_date"`
	DisagreedDate            *time.Time          `json:"disagreed_date" schema:"disagreed_date"`
	OrderAdmission           *OrderAdmission     `json:"order_admission" gorm:"foreignkey:id_order_admission"`
	IdOrderAdmission         *uint               `json:"id_order_admission" schema:"id_order_admission"`
	OrderAdmissionDate       *time.Time          `json:"order_admission_date" schema:"order_admission_date"`
	ReturnType               *ReturnTypes        `json:"return_type" gorm:"foreignkey:id_return_type"`
	IdReturnType             *uint               `json:"id_return_type" schema:"id_return_type"`
	ReturnDate               *time.Time          `json:"return_date" schema:"return_date"`
	Original                 bool                `json:"original" schema:"original"`
	OriginalDoc              *time.Time          `json:"original_doc" schema:"original_doc"`
	IdBenefit                uint                `json:"id_benefit" schema:"id_benefit"`
	Uid                      *string             `json:"uid" schema:"uid"`
	Created                  time.Time           `json:"created"` // Дата создания
	StatusComment            *string             `json:"status_comment" schema:"status_comment"`
}

type OrderAdmission struct {
}

func (Application) TableName() string {
	return "app.applications"
}
func (OrderAdmission) TableName() string {
	return "app.order_admission"
}
