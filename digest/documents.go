package digest

import (
	"time"
)

type Compatriot struct {
	Id                     uint               `json:"id"` // Идентификатор
	DocumentType           DocumentTypes      `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint               `json:"id_document_type"`
	CompatriotCategory     CompatriotCategory `json:"compatriot_category" gorm:"foreignkey:id_compatriot_category"`
	IdCompatriotCategory   uint               `json:"id_compatriot_category"`
	DocumentIdentification Identifications    `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint               `json:"id_ident_document"`
	Uid                    string             `json:"uid"`
	DocName                string             `json:"doc_name"`
	DocOrg                 string             `json:"doc_org"`
	Entrant                Entrants           `json:"entrant" gorm:"foreignkey:id_entrant"`
	IdEntrant              uint               `json:"id_entrant"`
	Created                time.Time          `json:"created"` // Дата создания
	PathFiles              string             `json:"path_files"`
	Checked                bool               `json:"checked"`
}
type Composition struct {
	Id                     uint              `json:"id"` // Идентификатор
	Entrant                Entrants          `json:"entrant" gorm:"foreignkey:id_entrant"`
	IdEntrant              uint              `json:"id_entrant"`
	DocumentIdentification Identifications   `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint              `json:"id_ident_document"`
	DocumentType           DocumentTypes     `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint              `json:"id_document_type"`
	Uid                    string            `json:"uid"`
	DocName                string            `json:"doc_name"`
	DocOrg                 string            `json:"doc_org"`
	DocYear                int64             `json:"doc_year"`
	CompositionThemes      CompositionThemes `json:"composition_theme" gorm:"foreignkey:id_composition_theme"`
	IdCompositionTheme     uint              `json:"id_composition_theme"`
	AppealStatuses         AppealStatuses    `json:"appeal_statuses" gorm:"foreignkey:id_appeal_status"`
	IdAppealStatus         uint              `json:"id_appeal_status"`
	HasAppealed            bool              `json:"has_appealed"`
	Checked                bool              `json:"checked"`
	Result                 bool              `json:"result"`
	IssueDate              time.Time         `json:"issue_date"`
	PathFiles              string            `json:"path_files"`
	Created                time.Time         `json:"created"` // Дата создания
}
type Ege struct {
	Id                     uint            `json:"id"` // Идентификатор
	Entrant                Entrants        `json:"entrant" gorm:"foreignkey:id_entrant"`
	IdEntrant              uint            `json:"id_entrant"`
	DocumentIdentification Identifications `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint            `json:"id_ident_document"`
	DocumentType           DocumentTypes   `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint            `json:"id_document_type"`
	DocName                string          `json:"doc_name"`
	DocOrg                 string          `json:"doc_org"`
	Uid                    string          `json:"uid"`
	RegisterNumber         string          `json:"register_number"`
	DocNumber              string          `json:"doc_number"`
	Region                 Region          `json:"region" gorm:"foreignkey:id_region"`
	IdRegion               uint            `json:"id_region"`
	Subject                Subject         `json:"subject" gorm:"foreignkey:id_subject"`
	IdSubject              uint            `json:"id_subject"`
	Mark                   int64           `json:"mark"`
	Checked                bool            `json:"checked"`
	ResultDate             time.Time       `json:"result_date"`
	IssueDate              time.Time       `json:"issue_date"`
	Created                time.Time       `json:"created"` // Дата создания
}
type Educations struct {
	Id                     uint            `json:"id"` // Идентификатор
	Entrant                Entrants        `json:"entrant" gorm:"foreignkey:id_entrant"`
	IdEntrant              uint            `json:"id_entrant"`
	DocumentIdentification Identifications `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint            `json:"id_ident_document"`
	DocumentType           DocumentTypes   `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint            `json:"id_document_type"`
	DocName                string          `json:"doc_name"`
	DocOrg                 string          `json:"doc_org"`
	Uid                    string          `json:"uid"`
	DocSeries              string          `json:"doc_series"`
	DocNumber              string          `json:"doc_number"`
	RegisterNumber         string          `json:"register_number"`
	Direction              Direction       `json:"direction" gorm:"foreignkey:id_direction"`
	IdDirection            uint            `json:"id_direction"`
	EducationLevel         EducationLevel  `json:"education_level" gorm:"foreignkey:id_education_level"`
	IdEducationLevel       uint            `json:"id_education_level"`
	IssueDate              time.Time       `json:"issue_date"`
	Created                time.Time       `json:"created"` // Дата создания
	Checked                bool            `json:"checked"`
}
type Disability struct {
	Id                     uint            `json:"id"` // Идентификатор
	Entrant                Entrants        `json:"entrant" gorm:"foreignkey:id_entrant"`
	IdEntrant              uint            `json:"id_entrant"`
	DocumentIdentification Identifications `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint            `json:"id_ident_document"`
	DocumentType           DocumentTypes   `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint            `json:"id_document_type"`
	DocName                string          `json:"doc_name"`
	DocOrg                 string          `json:"doc_org"`
	Uid                    string          `json:"uid"`
	DocNumber              string          `json:"doc_number"`
	Checked                bool            `json:"checked"`
	IssueDate              time.Time       `json:"issue_date"`
	DisabilityType         DisabilityTypes `json:"disability_type" gorm:"foreignkey:id_disability_type"`
	IdDisabilityType       uint            `json:"id_disability_type"`
	PathFiles              string          `json:"path_files"`
	Created                time.Time       `json:"created"` // Дата создания
}
type Identifications struct {
	Id              uint          `json:"id"` // Идентификатор
	Entrants        Entrants      `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId      uint          `json:"id_entrant" gorm:"column:id_entrant"`
	DocumentType    DocumentTypes `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType  uint          `json:"id_document_type"`
	Surname         string        `json:"surname"`
	Name            string        `json:"name"`
	Patronymic      string        `json:"patronymic"`
	DocSeries       string        `json:"doc_series"`
	DocNumber       string        `json:"doc_number"`
	DocOrganization string        `json:"doc_organization"`
	IssueDate       time.Time     `gorm:"type:time" json:"issue_date"`
	SubdivisionCode string        `json:"subdivision_code"`
	Okcm            Okcm          `json:"okcm" gorm:"foreignkey:id_okcm"`
	IdOkcm          uint          `json:"id_okcm"`
	Checked         bool          `json:"checked"`
	Created         time.Time     `json:"created"` // Дата создания
}
type Militaries struct {
	Id                     uint               `json:"id"` // Идентификатор
	Entrants               Entrants           `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint               `json:"id_entrant" gorm:"column:id_entrant"`
	DocumentType           DocumentTypes      `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint               `json:"id_document_type"`
	DocumentIdentification Identifications    `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint               `json:"id_ident_document"`
	DocName                string             `json:"doc_name"`
	DocOrg                 string             `json:"doc_org"`
	Uid                    string             `json:"uid"`
	DocNumber              string             `json:"doc_number"`
	DocSeries              string             `json:"doc_series"`
	Checked                bool               `json:"checked"`
	IssueDate              time.Time          `json:"issue_date"`
	MilitaryCategories     MilitaryCategories `json:"military_categories" gorm:"foreignkey:id_category"`
	IdCategory             uint               `json:"id_category"`
	PathFiles              string             `json:"path_files"`
	Created                time.Time          `json:"created"` // Дата создания
}
type OlympicsDocs struct {
	Id                     uint            `json:"id"` // Идентификатор
	Entrants               Entrants        `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint            `json:"id_entrant" gorm:"column:id_entrant"`
	DocumentType           DocumentTypes   `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint            `json:"id_document_type"`
	DocumentIdentification Identifications `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint            `json:"id_ident_document"`
	DocName                string          `json:"doc_name"`
	DocOrg                 string          `json:"doc_org"`
	Uid                    string          `json:"uid"`
	DocNumber              string          `json:"doc_number"`
	DocSeries              string          `json:"doc_series"`
	Checked                bool            `json:"checked"`
	IssueDate              time.Time       `json:"issue_date"`
	Olympics               Olympics        `json:"olympics" gorm:"foreignkey:id_olympic"`
	IdOlympic              uint            `json:"id_olympic"`
	PathFiles              string          `json:"path_files"`
	Created                time.Time       `json:"created"` // Дата создания
}
type Orphans struct {
	Id                     uint             `json:"id"` // Идентификатор
	Entrants               Entrants         `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint             `json:"id_entrant" gorm:"column:id_entrant"`
	DocumentType           DocumentTypes    `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint             `json:"id_document_type"`
	DocumentIdentification Identifications  `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint             `json:"id_ident_document"`
	DocName                string           `json:"doc_name"`
	DocOrg                 string           `json:"doc_org"`
	Uid                    string           `json:"uid"`
	DocNumber              string           `json:"doc_number"`
	DocSeries              string           `json:"doc_series"`
	Checked                bool             `json:"checked"`
	IssueDate              time.Time        `json:"issue_date"`
	OrphanCategories       OrphanCategories `json:"olympics" gorm:"foreignkey:id_category"`
	IdCategory             uint             `json:"id_category"`
	PathFiles              string           `json:"path_files"`
	Created                time.Time        `json:"created"` // Дата создания
}
type Other struct {
	Id                     uint            `json:"id"` // Идентификатор
	Entrants               Entrants        `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint            `json:"id_entrant" gorm:"column:id_entrant"`
	DocumentType           DocumentTypes   `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint            `json:"id_document_type"`
	DocumentIdentification Identifications `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint            `json:"id_ident_document"`
	DocName                string          `json:"doc_name"`
	DocOrg                 string          `json:"doc_org"`
	Uid                    string          `json:"uid"`
	DocNumber              string          `json:"doc_number"`
	DocSeries              string          `json:"doc_series"`
	Checked                bool            `json:"checked"`
	IssueDate              time.Time       `json:"issue_date"`
	PathFiles              string          `json:"path_files"`
	Created                time.Time       `json:"created"` // Дата создания
}
type ParentsLost struct {
	Id                     uint                  `json:"id"` // Идентификатор
	Entrants               Entrants              `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint                  `json:"id_entrant" gorm:"column:id_entrant"`
	DocumentType           DocumentTypes         `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint                  `json:"id_document_type"`
	DocumentIdentification Identifications       `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint                  `json:"id_ident_document"`
	DocName                string                `json:"doc_name"`
	DocOrg                 string                `json:"doc_org"`
	Uid                    string                `json:"uid"`
	DocNumber              string                `json:"doc_number"`
	DocSeries              string                `json:"doc_series"`
	Checked                bool                  `json:"checked"`
	IssueDate              time.Time             `json:"issue_date"`
	ParentsLostCategory    ParentsLostCategories `json:"parents_lost_category" gorm:"foreignkey:id_category"`
	IdCategory             uint                  `json:"id_category"`
	PathFiles              string                `json:"path_files"`
	Created                time.Time             `json:"created"` // Дата создания
}
type RadiationWork struct {
	Id                     uint                    `json:"id"` // Идентификатор
	Entrants               Entrants                `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint                    `json:"id_entrant" gorm:"column:id_entrant"`
	DocumentType           DocumentTypes           `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint                    `json:"id_document_type"`
	DocumentIdentification Identifications         `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint                    `json:"id_ident_document"`
	DocName                string                  `json:"doc_name"`
	DocOrg                 string                  `json:"doc_org"`
	Uid                    string                  `json:"uid"`
	DocNumber              string                  `json:"doc_number"`
	DocSeries              string                  `json:"doc_series"`
	Checked                bool                    `json:"checked"`
	IssueDate              time.Time               `json:"issue_date"`
	RadiationWorkCategory  RadiationWorkCategories `json:"radiation_work_category" gorm:"foreignkey:id_category"`
	IdCategory             uint                    `json:"id_category"`
	PathFiles              string                  `json:"path_files"`
	Created                time.Time               `json:"created"` // Дата создания
}
type Veteran struct {
	Id                     uint              `json:"id"` // Идентификатор
	Entrants               Entrants          `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint              `json:"id_entrant" gorm:"column:id_entrant"`
	DocumentType           DocumentTypes     `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint              `json:"id_document_type"`
	DocumentIdentification Identifications   `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint              `json:"id_ident_document"`
	DocName                string            `json:"doc_name"`
	DocOrg                 string            `json:"doc_org"`
	Uid                    string            `json:"uid"`
	DocNumber              string            `json:"doc_number"`
	DocSeries              string            `json:"doc_series"`
	Checked                bool              `json:"checked"`
	IssueDate              time.Time         `json:"issue_date"`
	VeteranCategory        VeteranCategories `json:"veteran_category" gorm:"foreignkey:id_category"`
	IdCategory             uint              `json:"id_category"`
	PathFiles              string            `json:"path_files"`
	Created                time.Time         `json:"created"` // Дата создания
}

type VDocuments struct {
	IdDocument        uint          `json:"id"` // Идентификатор
	Entrant           Entrants      `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId        uint          `json:"id_entrant" gorm:"column:id_entrant"`
	Snils             string        `json:"snils"`
	Surname           string        `json:"surname"`
	Name              string        `json:"name"`
	Patronymic        string        `json:"patronymic"`
	Gender            Gender        `json:"gender" gorm:"foreignkey:IdGender"`
	IdGender          uint          `json:"id_gender"`
	Birthday          time.Time     `json:"birthday"`
	IdSysCategories   uint          `json:"id_sys_categories"`
	DocumentName      string        `json:"document_name"`
	NameSysCategories string        `json:"name_sys_categories"`
	NameTable         string        `json:"name_table"`
	DocumentType      DocumentTypes `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType    uint          `json:"id_document_type"`
	Created           time.Time     `json:"created"` // Дата создания
}

func (Compatriot) TableName() string {
	return "documents.compatriot"
}
func (Composition) TableName() string {
	return "documents.composition"
}
func (Disability) TableName() string {
	return "documents.disability"
}
func (Educations) TableName() string {
	return "documents.educations"
}
func (Ege) TableName() string {
	return "documents.ege"
}
func (Identifications) TableName() string {
	return "documents.identification"
}
func (Militaries) TableName() string {
	return "documents.militaries"
}
func (OlympicsDocs) TableName() string {
	return "documents.olympics"
}
func (Orphans) TableName() string {
	return "documents.orphans"
}
func (Other) TableName() string {
	return "documents.other"
}
func (ParentsLost) TableName() string {
	return "documents.parents_lost"
}
func (RadiationWork) TableName() string {
	return "documents.radiation_work"
}
func (Veteran) TableName() string {
	return "documents.veteran"
}
func (VDocuments) TableName() string {
	return "persons.v_documents"
}
