package digest

import (
	"errors"
	"mime/multipart"
	"persons/config"
	"time"
)

type File struct {
	MultFile multipart.File
	Header   multipart.FileHeader
}

type FileC struct {
	Content []byte `json:"content,omitempty"`
	Title   string `json:"title,omitempty"`
	Size    int64  `json:"size,omitempty"`
	Type    string `json:"type, omitempty"`
}

type Compatriot struct {
	Id                     uint               `json:"id"` // Идентификатор
	DocumentType           DocumentTypes      `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint               `json:"id_document_type" schema:"id_document_type"`
	CompatriotCategory     CompatriotCategory `json:"compatriot_category" gorm:"foreignkey:id_compatriot_category"`
	IdCompatriotCategory   uint               `json:"id_compatriot_category" schema:"id_compatriot_category"`
	DocumentIdentification Identifications    `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint               `json:"id_ident_document" schema:"id_ident_document"`
	Uid                    *string            `json:"uid" schema:"uid"`
	DocName                string             `json:"doc_name" schema:"doc_name"`
	DocOrg                 string             `json:"doc_org" schema:"doc_org"`
	Entrant                Entrants           `json:"entrant" gorm:"foreignkey:id_entrant"`
	IdEntrant              uint               `json:"id_entrant"`
	Created                time.Time          `json:"created"` // Дата создания
	PathFile               *string            `json:"path_file" schema:"file" schema:"file"`
	Checked                bool               `json:"checked"`
	IdOrganization         uint               `json:"id_organization"`
}
type Composition struct {
	Id                     uint              `json:"id"` // Идентификатор
	Entrant                Entrants          `json:"entrant" gorm:"foreignkey:id_entrant"`
	IdEntrant              uint              `json:"id_entrant" schema:"id_entrant"`
	DocumentIdentification Identifications   `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint              `json:"id_ident_document" schema:"id_ident_document"`
	DocumentType           DocumentTypes     `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint              `json:"id_document_type" schema:"id_document_type"`
	Uid                    *string           `json:"uid" schema:"uid"`
	DocName                string            `json:"doc_name" schema:"doc_name"`
	DocOrg                 string            `json:"doc_org" schema:"doc_org"`
	DocYear                int64             `json:"doc_year" schema:"doc_year"`
	CompositionThemes      CompositionThemes `json:"composition_theme" gorm:"foreignkey:id_composition_theme"`
	IdCompositionTheme     uint              `json:"id_composition_theme" schema:"id_composition_theme"`
	AppealStatuses         AppealStatuses    `json:"appeal_statuses" gorm:"foreignkey:id_appeal_status"`
	IdAppealStatus         uint              `json:"id_appeal_status" schema:"id_appeal_status"`
	HasAppealed            bool              `json:"has_appealed" schema:"has_appealed"`
	Checked                bool              `json:"checked"`
	Result                 bool              `json:"result" schema:"result"`
	IssueDate              time.Time         `json:"issue_date" schema:"issue_date"`
	PathFile               *string           `json:"path_file" schema:"file"`
	Created                time.Time         `json:"created"` // Дата создания
	IdOrganization         uint              `json:"id_organization"`
}

type Disability struct {
	Id                     uint            `json:"id" schema:"id"` // Идентификатор
	Entrant                Entrants        `json:"entrant" gorm:"foreignkey:id_entrant"`
	IdEntrant              uint            `json:"id_entrant" schema:"id_entrant"`
	DocumentIdentification Identifications `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint            `json:"id_ident_document" schema:"id_ident_document"`
	DocumentType           DocumentTypes   `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint            `json:"id_document_type" schema:"id_document_type"`
	DocName                string          `json:"doc_name" schema:"doc_name"`
	DocOrg                 string          `json:"doc_org" schema:"doc_org"`
	Uid                    *string         `json:"uid" schema:"uid"`
	DocNumber              string          `json:"doc_number" schema:"doc_number"`
	Checked                bool            `json:"checked"`
	IssueDate              time.Time       `json:"issue_date" schema:"issue_date"`
	DisabilityType         DisabilityTypes `json:"disability_type" gorm:"foreignkey:id_disability_type"`
	IdDisabilityType       uint            `json:"id_disability_type" schema:"id_disability_type"`
	PathFile               *string         `json:"path_file" schema:"file"`
	Created                time.Time       `json:"created"` // Дата создания
	IdOrganization         uint            `json:"id_organization"`
}

type Ege struct {
	Id                     uint            `json:"id" schema:"id"` // Идентификатор
	Entrant                Entrants        `json:"entrant" gorm:"foreignkey:id_entrant"`
	IdEntrant              uint            `json:"id_entrant" schema:"id_entrant"`
	DocumentIdentification Identifications `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint            `json:"id_ident_document" schema:"id_ident_document"`
	DocumentType           DocumentTypes   `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint            `json:"id_document_type" schema:"id_document_type"`
	DocName                string          `json:"doc_name" schema:"doc_name"`
	DocOrg                 string          `json:"doc_org" schema:"doc_org"`
	Uid                    *string         `json:"uid" schema:"uid"`
	RegisterNumber         string          `json:"register_number" schema:"register_number"`
	DocNumber              string          `json:"doc_number" schema:"doc_number"`
	Region                 Region          `json:"region" gorm:"foreignkey:id_region"`
	IdRegion               uint            `json:"id_region" schema:"id_region"`
	Subject                Subject         `json:"subject" gorm:"foreignkey:id_subject"`
	IdSubject              uint            `json:"id_subject" schema:"id_subject"`
	Mark                   int64           `json:"mark" schema:"mark"`
	Checked                bool            `json:"checked"`
	ResultDate             time.Time       `json:"result_date" schema:"result_date"`
	IssueDate              time.Time       `json:"issue_date" schema:"issue_date"`
	Created                time.Time       `json:"created"` // Дата создания
	PathFile               *string         `json:"path_file" schema:"file"`
	IdOrganization         uint            `json:"id_organization"`
}
type Educations struct {
	Id                     uint                   `json:"id" schema:"id"` // Идентификатор
	Entrant                Entrants               `json:"entrant" gorm:"foreignkey:id_entrant"`
	IdEntrant              uint                   `json:"id_entrant" schema:"id_entrant"`
	DocumentIdentification Identifications        `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint                   `json:"id_ident_document" schema:"id_ident_document"`
	DocumentType           DocumentTypes          `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint                   `json:"id_document_type" schema:"id_document_type"`
	DocName                string                 `json:"doc_name" schema:"doc_name"`
	DocOrg                 string                 `json:"doc_org" schema:"doc_org"`
	Uid                    *string                `json:"uid" schema:"uid"`
	DocSeries              string                 `json:"doc_series" schema:"doc_series"`
	DocNumber              string                 `json:"doc_number" schema:"doc_number"`
	RegisterNumber         string                 `json:"register_number" schema:"register_number"`
	Direction              Direction              `json:"direction" gorm:"foreignkey:id_direction"`
	IdDirection            *uint                  `json:"id_direction" schema:"id_direction"`
	EducationLevel         DocumentEducationLevel `json:"education_level" gorm:"foreignkey:id_education_level"`
	IdEducationLevel       uint                   `json:"id_education_level" schema:"id_education_level"`
	IssueDate              time.Time              `json:"issue_date" schema:"issue_date"`
	Created                time.Time              `json:"created"` // Дата создания
	Checked                bool                   `json:"checked"`
	PathFile               *string                `json:"path_file" schema:"file"`
	IdOrganization         uint                   `json:"id_organization"`
}
type General struct {
	Id                     uint            `json:"id" schema:"id"` // Идентификатор
	Entrant                Entrants        `json:"entrant" gorm:"foreignkey:id_entrant"`
	IdEntrant              uint            `json:"id_entrant" schema:"id_entrant"`
	DocumentIdentification Identifications `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint            `json:"id_ident_document" schema:"id_ident_document"`
	DocumentType           DocumentTypes   `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint            `json:"id_document_type" schema:"id_document_type"`
	DocName                string          `json:"doc_name" schema:"doc_name"`
	DocOrg                 string          `json:"doc_org" schema:"doc_org"`
	Uid                    *string         `json:"uid" schema:"uid"`
	Created                time.Time       `json:"created"` // Дата создания
	PathFile               *string         `json:"path_file" schema:"file"`
}
type Identifications struct {
	Id              uint          `json:"id" schema:"id"` // Идентификатор
	Entrants        Entrants      `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId      uint          `json:"id_entrant" gorm:"column:id_entrant" schema:"id_entrant"`
	DocumentType    DocumentTypes `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType  uint          `json:"id_document_type" schema:"id_document_type"`
	Surname         string        `json:"surname" schema:"surname"`
	Name            string        `json:"name" schema:"name"`
	Uid             *string       `json:"uid" schema:"uid"`
	Patronymic      string        `json:"patronymic" schema:"patronymic"`
	DocSeries       string        `json:"doc_series" schema:"doc_series"`
	DocNumber       string        `json:"doc_number" schema:"doc_number"`
	DocOrganization string        `json:"doc_organization" schema:"doc_organization"`
	IssueDate       time.Time     `gorm:"type:time" json:"issue_date" schema:"issue_date"`
	SubdivisionCode string        `json:"subdivision_code" schema:"subdivision_code"`
	Okcm            Okcm          `json:"okcm" gorm:"foreignkey:id_okcm"`
	IdOkcm          uint          `json:"id_okcm" schema:"id_okcm"`
	Checked         bool          `json:"checked"`
	Created         time.Time     `json:"created"` // Дата создания
	PathFile        *string       `json:"path_file" schema:"file"`
	IdOrganization  uint          `json:"id_organization"`
}
type Militaries struct {
	Id                     uint               `json:"id" schema:"id"` // Идентификатор
	Entrants               Entrants           `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint               `json:"id_entrant" gorm:"column:id_entrant" schema:"id_entrant"`
	DocumentType           DocumentTypes      `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint               `json:"id_document_type" schema:"id_document_type"`
	DocumentIdentification Identifications    `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint               `json:"id_ident_document" schema:"id_ident_document"`
	DocName                string             `json:"doc_name" schema:"doc_name"`
	DocOrg                 string             `json:"doc_org" schema:"doc_org"`
	Uid                    *string            `json:"uid" schema:"uid"`
	DocNumber              string             `json:"doc_number" schema:"doc_number"`
	DocSeries              string             `json:"doc_series" schema:"doc_series"`
	Checked                bool               `json:"checked"`
	IssueDate              time.Time          `json:"issue_date" schema:"issue_date"`
	MilitaryCategories     MilitaryCategories `json:"military_categories" gorm:"foreignkey:id_category"`
	IdCategory             uint               `json:"id_category" schema:"id_category"`
	Created                time.Time          `json:"created"` // Дата создания
	PathFile               *string            `json:"path_file" schema:"file"`
	IdOrganization         uint               `json:"id_organization"`
}
type OlympicsDocs struct {
	Id                     uint            `json:"id" schema:"id"` // Идентификатор
	Entrants               Entrants        `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint            `json:"id_entrant" gorm:"column:id_entrant" schema:"id_entrant"`
	DocumentType           DocumentTypes   `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint            `json:"id_document_type" schema:"id_document_type"`
	DocumentIdentification Identifications `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint            `json:"id_ident_document" schema:"id_ident_document"`
	DocName                string          `json:"doc_name" schema:"doc_name"`
	DocOrg                 string          `json:"doc_org" schema:"doc_org"`
	Uid                    *string         `json:"uid" schema:"uid"`
	DocNumber              string          `json:"doc_number" schema:"doc_number"`
	DocSeries              string          `json:"doc_series" schema:"doc_series"`
	Checked                bool            `json:"checked"`
	IssueDate              time.Time       `json:"issue_date" schema:"issue_date"`
	Olympics               Olympics        `json:"olympics" gorm:"foreignkey:id_olympic"`
	IdOlympic              uint            `json:"id_olympic" schema:"id_olympic"`
	Created                time.Time       `json:"created"` // Дата создания
	PathFile               *string         `json:"path_file" schema:"file"`
	IdOrganization         uint            `json:"id_organization"`
}
type Orphans struct {
	Id                     uint             `json:"id" schema:"id"` // Идентификатор
	Entrants               Entrants         `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint             `json:"id_entrant" gorm:"column:id_entrant" schema:"id_entrant"`
	DocumentType           DocumentTypes    `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint             `json:"id_document_type" schema:"id_document_type"`
	DocumentIdentification Identifications  `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint             `json:"id_ident_document" schema:"id_ident_document"`
	DocName                string           `json:"doc_name" schema:"doc_name"`
	DocOrg                 string           `json:"doc_org" schema:"doc_org"`
	Uid                    *string          `json:"uid" schema:"uid"`
	DocNumber              string           `json:"doc_number" schema:"doc_number"`
	DocSeries              string           `json:"doc_series" schema:"doc_series"`
	Checked                bool             `json:"checked"`
	IssueDate              time.Time        `json:"issue_date" schema:"issue_date"`
	OrphanCategories       OrphanCategories `json:"olympics" gorm:"foreignkey:id_category"`
	IdCategory             uint             `json:"id_category" schema:"id_category"`
	Created                time.Time        `json:"created"` // Дата создания
	PathFile               *string          `json:"path_file" schema:"file"`
	IdOrganization         uint             `json:"id_organization"`
}
type Other struct {
	Id                     uint            `json:"id" schema:"id"` // Идентификатор
	Entrants               Entrants        `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint            `json:"id_entrant" gorm:"column:id_entrant" schema:"id_entrant"`
	DocumentType           DocumentTypes   `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint            `json:"id_document_type" schema:"id_document_type"`
	DocumentIdentification Identifications `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint            `json:"id_ident_document" schema:"id_ident_document"`
	DocName                string          `json:"doc_name" schema:"doc_name"`
	DocOrg                 string          `json:"doc_org" schema:"doc_org"`
	Uid                    *string         `json:"uid" schema:"uid"`
	DocNumber              string          `json:"doc_number" schema:"doc_number"`
	DocSeries              string          `json:"doc_series" schema:"doc_series"`
	Checked                bool            `json:"checked"`
	IssueDate              time.Time       `json:"issue_date" schema:"issue_date"`
	Created                time.Time       `json:"created"` // Дата создания
	PathFile               *string         `json:"path_file" schema:"file"`
	IdOrganization         uint            `json:"id_organization"`
}
type ParentsLost struct {
	Id                     uint                  `json:"id" schema:"id"` // Идентификатор
	Entrants               Entrants              `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint                  `json:"id_entrant" gorm:"column:id_entrant" schema:"id_entrant"`
	DocumentType           DocumentTypes         `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint                  `json:"id_document_type" schema:"id_document_type"`
	DocumentIdentification Identifications       `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint                  `json:"id_ident_document" schema:"id_ident_document"`
	DocName                string                `json:"doc_name" schema:"doc_name"`
	DocOrg                 string                `json:"doc_org" schema:"doc_org"`
	Uid                    *string               `json:"uid" schema:"uid"`
	DocNumber              string                `json:"doc_number" schema:"doc_number"`
	DocSeries              string                `json:"doc_series" schema:"doc_series"`
	Checked                bool                  `json:"checked"`
	IssueDate              time.Time             `json:"issue_date" schema:"issue_date"`
	ParentsLostCategory    ParentsLostCategories `json:"parents_lost_category" gorm:"foreignkey:id_category"`
	IdCategory             uint                  `json:"id_category" schema:"id_category"`
	Created                time.Time             `json:"created"` // Дата создания
	PathFile               *string               `json:"path_file" schema:"file"`
	IdOrganization         uint                  `json:"id_organization"`
}
type RadiationWork struct {
	Id                     uint                    `json:"id" schema:"id"` // Идентификатор
	Entrants               Entrants                `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint                    `json:"id_entrant" gorm:"column:id_entrant" schema:"id_entrant"`
	DocumentType           DocumentTypes           `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint                    `json:"id_document_type" schema:"id_document_type"`
	DocumentIdentification Identifications         `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint                    `json:"id_ident_document" schema:"id_ident_document"`
	DocName                string                  `json:"doc_name" schema:"doc_name"`
	DocOrg                 string                  `json:"doc_org" schema:"doc_org"`
	Uid                    *string                 `json:"uid" schema:"uid"`
	DocNumber              string                  `json:"doc_number" schema:"doc_number"`
	DocSeries              string                  `json:"doc_series" schema:"doc_series"`
	Checked                bool                    `json:"checked"`
	IssueDate              time.Time               `json:"issue_date" schema:"issue_date"`
	RadiationWorkCategory  RadiationWorkCategories `json:"radiation_work_category" gorm:"foreignkey:id_category"`
	IdCategory             uint                    `json:"id_category" schema:"id_category"`
	Created                time.Time               `json:"created"` // Дата создания
	PathFile               *string                 `json:"path_file" schema:"file"`
	IdOrganization         uint                    `json:"id_organization"`
}
type Veteran struct {
	Id                     uint              `json:"id" schema:"id"` // Идентификатор
	Entrants               Entrants          `json:"entrant" gorm:"association_foreignkey:Id"`
	EntrantsId             uint              `json:"id_entrant" gorm:"column:id_entrant" schema:"id_entrant"`
	DocumentType           DocumentTypes     `json:"document_type" gorm:"foreignkey:id_document_type"`
	IdDocumentType         uint              `json:"id_document_type" schema:"id_document_type"`
	DocumentIdentification Identifications   `json:"document_identification" gorm:"foreignkey:id_ident_document"`
	IdIdentDocument        uint              `json:"id_ident_document" schema:"id_ident_document"`
	DocName                string            `json:"doc_name" schema:"doc_name"`
	DocOrg                 string            `json:"doc_org" schema:"doc_org"`
	Uid                    *string           `json:"uid" schema:"uid"`
	DocNumber              string            `json:"doc_number" schema:"doc_number"`
	DocSeries              string            `json:"doc_series" schema:"doc_series"`
	Checked                bool              `json:"checked"`
	IssueDate              time.Time         `json:"issue_date" schema:"issue_date"`
	VeteranCategory        VeteranCategories `json:"veteran_category" gorm:"foreignkey:id_category"`
	IdCategory             uint              `json:"id_category" schema:"id_category"`
	Created                time.Time         `json:"created"` // Дата создания
	PathFile               *string           `json:"path_file" schema:"file"`
	IdOrganization         uint              `json:"id_organization"`
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
	Created           time.Time     `json:"created"`                    // Дата создания
	PathFile          *string       `json:"path_file" schema:"file"`    // Дата создания
	DocName           *string       `json:"doc_name" schema:"doc_name"` // Дата создания
}

type AllDocuments struct {
	Id                uint       `json:"id"` // Идентификатор
	DocNumber         *string    `json:"doc_number"`
	DocSeries         *string    `json:"doc_series"`
	IssueDate         *time.Time `json:"issue_date" schema:"issue_date"`
	IdDocumentType    uint       `json:"id_document_type"`
	NameDocumentType  string     `json:"name_document_type"`
	IdSysCategories   uint       `json:"id_sys_categories"`
	NameSysCategories string     `json:"name_sys_categories"`
	NameTable         string     `json:"name_table"`
	Mark              *int64     `json:"mark"`
	NameSubject       *string    `json:"name_subject"`
	Checked           bool       `json:"checked"`
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
func (General) TableName() string {
	return "documents.general"
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

func GetVDocuments(id uint) (*VDocuments, error) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item VDocuments
	db := conn.Where(`id_document=?`, id).Find(&item)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			return nil, errors.New(`Документ не найден. `)
		}
		return nil, errors.New(`Ошибка подключения к БД. `)
	}
	if item.IdDocument <= 0 {
		return nil, errors.New(`Документ не найден. `)
	}
	return &item, nil
}
