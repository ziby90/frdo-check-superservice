package digest

import (
	"errors"
	"persons/config"
	"time"
)

type Application struct {
	Id                       uint                `json:"id" schema:"id"` // Идентификатор
	IdOrganization           uint                `json:"id_organization"`
	Entrants                 Entrants            `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId               uint                `json:"id_entrant" gorm:"column:id_entrant" schema:"id_entrant"`
	CompetitiveGroup         CompetitiveGroup    `gorm:"foreignkey:IdCompetitiveGroup"`
	IdCompetitiveGroup       uint                `json:"id_competitive_group"`
	IdCompetitiveGroupTarget *uint               `json:"id_competitive_group_target"`
	AppNumber                string              `json:"app_number"`
	RegistrationDate         time.Time           `json:"registration_date" schema:"registration_date"`
	Rating                   float32             `json:"rating" schema:"rating"`
	Status                   ApplicationStatuses `gorm:"foreignkey:id_status"`
	IdStatus                 uint                `json:"id_status" schema:"id_status"`
	Priority                 int64               `json:"priority" schema:"priority"`
	FirstHigherEducation     bool                `json:"first_higher_education" schema:"first_higher_education"`
	NeedHostel               bool                `json:"need_hostel" schema:"need_hostel"`
	IdDisabledDocument       *uint               `json:"id_disabled_document" schema:"id_disabled_document"`
	SpecialConditions        bool                `json:"special_conditions" schema:"special_conditions"`
	DistanceTest             bool                `json:"distance_test" schema:"distance_test"`
	DistancePlace            *string             `json:"distance_place" schema:"distance_place"`
	ViolationTypes           ViolationTypes      `json:"violation_type" gorm:"foreignkey:id_violation"`
	IdViolation              *uint               `json:"id_violation" schema:"id_violation"`
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
	IdBenefit                *uint               `json:"id_benefit" schema:"id_benefit"`
	Uid                      *string             `json:"uid" schema:"uid"`
	UidEpgu                  *int64              `json:"uid_epgu" schema:"uid_epgu"`
	Created                  time.Time           `json:"created"` // Дата создания
	StatusComment            *string             `json:"status_comment" schema:"status_comment"`
}

type ApplicationsAgreedHistory struct {
	Id             uint      `json:"id" schema:"id"` // Идентификатор
	IdApplication  uint      `json:"id_application"`
	Agreed         bool      `json:"agreed"`
	Date           time.Time `json:"date"`
	IdOrganization *uint     `json:"id_organization"`
	Created        time.Time `json:"created"` // Дата создания
}

type Documents struct {
	Id                    uint                  `json:"id" schema:"id"` // Идентификатор
	Application           Application           `gorm:"foreignkey:IdApplication"`
	IdApplication         uint                  `json:"id_application"`
	Document              VDocuments            `gorm:"foreignkey:IdDocument"`
	IdDocument            uint                  `json:"id_document"`
	DocumentSysCategory   DocumentSysCategories `gorm:"foreignkey:IdDocumentSysCategory"`
	IdDocumentSysCategory uint                  `json:"id_document_sys_category"`
	Uid                   *string               `json:"uid" schema:"uid"`
	OriginalDoc           *time.Time            `json:"original_doc"`
	Created               time.Time             `json:"created"` // Дата создания
}
type AppEntranceTest struct {
	Id                         uint                          `json:"id" schema:"id"` // Идентификатор
	Application                Application                   `gorm:"foreignkey:IdApplication"`
	IdApplication              uint                          `json:"id_application"`
	EntranceTest               EntranceTest                  `gorm:"foreignkey:IdEntranceTest"`
	IdEntranceTest             uint                          `json:"id_entrance_test"`
	Document                   VDocuments                    `gorm:"foreignkey:IdDocument"`
	IdDocument                 *uint                         `json:"id_document"`
	EntranceTestResultSource   EntranceTestTypeResultSources `gorm:"foreignkey:IdEntranceTestResultSource"`
	IdResultSource             uint                          `json:"id_result_source"`
	EntranceTestDocumentType   EntranceTestTypeDocumentTypes `gorm:"foreignkey:IdEntranceTestDocumentType"`
	IdEntranceTestDocumentType *uint                         `json:"id_entrance_test_document_type"`
	//DocumentSysCategory   	DocumentSysCategories `gorm:"foreignkey:IdDocumentSysCategory"`
	//IdDocumentSysCategory 	uint                  `json:"id_document_sys_category"`
	ResultValue    int64          `json:"result_value"`
	Benefit        Benefit        `gorm:"foreignkey:IdBenefit"`
	IdBenefit      *uint          `json:"id_benefit"`
	AppealStatus   AppealStatuses `gorm:"foreignkey:IdAppealStatus"`
	IdAppealStatus *uint          `json:"id_appeal_status"`
	HasEge         *bool          `json:"has_ege"`
	EgeValue       *int64         `json:"ege_value"`
	Uid            *string        `json:"uid" schema:"uid"`
	IssueDate      time.Time      `json:"issue_date"`
	Created        time.Time      `json:"created"` // Дата создания
}

type VApplications struct {
	Id                   uint      `json:"id" schema:"id"` // Идентификатор
	IdOrganization       uint      `json:"id_organization"`
	IdEntrant            uint      `json:"id_entrant"`
	EntrantFullname      string    `json:"entrant_fullname"`
	EntrantSnils         string    `json:"entrant_snils"`
	IdCompetitiveGroup   uint      `json:"id_competitive_group"`
	NameCompetitiveGroup string    `json:"name_competitive_group"`
	AppNumber            string    `json:"app_number"`
	RegistrationDate     time.Time `json:"registration_date"`
	Rating               float32   `json:"rating" schema:"rating"`
	IdStatus             uint      `json:"id_status" schema:"id_status"`
	NameStatus           string    `json:"name_status" schema:"name_status"`
	Agreed               *bool     `json:"agreed" schema:"agreed"`
	Original             bool      `json:"original" schema:"original"`
	Uid                  *string   `json:"uid" schema:"uid"`
	Created              time.Time `json:"created"` // Дата создания
}
type AppAchievements struct {
	Id                  uint                   `json:"id" schema:"id"` // Идентификатор
	Application         Application            `gorm:"foreignkey:IdApplication"`
	IdApplication       uint                   `json:"id_application"`
	AchievementCategory AchievementCategory    `gorm:"foreignkey:IdCategory" json:"category"`
	IdCategory          uint                   `json:"id_category"` // Идентификатор наименования категории индивидуального достижения
	Achievement         IndividualAchievements `gorm:"foreignkey:IdAchievement"`
	IdAchievement       *uint                  `json:"id_achievement"`
	Document            VDocuments             `gorm:"foreignkey:IdDocument"`
	IdDocument          *uint                  `json:"id_document"`
	Name                string                 `json:"name"`
	Mark                *int64                 `json:"mark"`
	Uid                 *string                `json:"uid" schema:"uid"`
	UidEpgu             *string                `json:"uid_epgu" schema:"uid_epgu"`
	PathFile            *string                `json:"path_file" schema:"file" schema:"file"`
	Created             time.Time              `json:"created"` // Дата создания
}

type OrderAdmission struct {
}

func GetApplication(id uint) (*Application, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item Application
	db := conn.Preload(`CompetitiveGroup`).Find(&item, id)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Заявление не найдено. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if item.Id <= 0 {
		return nil, errors.New(`Заявление не найдено. `)
	}
	return &item, nil
}

func (Application) TableName() string {
	return "app.applications"
}
func (ApplicationsAgreedHistory) TableName() string {
	return "app.applications_agreed_history"
}
func (AppAchievements) TableName() string {
	return "app.achievements"
}
func (AppEntranceTest) TableName() string {
	return "app.entrance_test"
}
func (OrderAdmission) TableName() string {
	return "app.order_admission"
}
func (Documents) TableName() string {
	return "app.documents"
}
func (VApplications) TableName() string {
	return "app.v_applications"
}
